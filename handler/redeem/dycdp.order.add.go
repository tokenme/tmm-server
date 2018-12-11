package redeem

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/nlopes/slack"
	//"github.com/davecgh/go-spew/spew"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/nu7hatch/gouuid"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/dycdp"
	"github.com/tokenme/tmm/tools/forex"
	"github.com/tokenme/tmm/utils"
	"github.com/xluohome/phonedata"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DycdpOrderAddRequest struct {
	OfferId  uint64 `json:"offer_id" form:"offer_id" binding:"required"`
	DeviceId string `json:"device_id" form:"device_id" binding:"required"`
}

func DycdpOrderAddHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req DycdpOrderAddRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	if CheckWithCode(user.CountryCode != 86, INVALID_CDP_VENDOR_ERROR, "the cdp vendor not supported", c) {
		return
	}

	db := Service.Db
	{
		rows, _, err := db.Query(`SELECT 1 FROM tmm.user_settings WHERE user_id=%d AND blocked=1 AND block_whitelist=0`, user.Id)
		if CheckErr(err, c) {
			return
		}
		if Check(len(rows) > 0, "您的账户存在异常操作，疑似恶意邀请用户行为，不能执行提现操作。如有疑问请联系客服。", c) {
			return
		}
	}
	userMobile := strings.TrimSpace(user.Mobile)
	phone, err := phonedata.Find(userMobile)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(phone.CardType != "中国移动" && phone.CardType != "中国联通" && phone.CardType != "中国电信", INVALID_CDP_VENDOR_ERROR, "the cdp vendor not supported", c) {
		return
	}
	cdpClient, err := dycdp.NewClientWithAccessKey(Config.Aliyun.RegionId, Config.Aliyun.AK, Config.Aliyun.AS)
	if CheckErr(err, c) {
		return
	}
	cdpRequest := dycdp.CreateQueryCdpOfferByIdRequest()
	cdpRequest.OfferId = strconv.FormatUint(req.OfferId, 10)
	cdpResponse, err := cdpClient.QueryCdpOfferById(cdpRequest)
	if CheckErr(err, c) {
		return
	}
	offers := cdpResponse.FlowOffers.FlowOffer
	if CheckWithCode(len(offers) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	offer := offers[0]
	price := decimal.New(int64(offer.Price), 0).Div(decimal.New(100, 0))
	pointPrice := common.GetPointPrice(Service, Config)
	currencyRate := forex.Rate(Service, "USD", "CNY")
	pointPrice = pointPrice.Mul(currencyRate)
	if CheckWithCode(!pointPrice.GreaterThan(decimal.Zero), INTERNAL_ERROR, "system error", c) {
		return
	}
	points := price.Div(pointPrice)
	rows, _, err := db.Query(`SELECT points FROM tmm.devices WHERE id='%s' AND user_id=%d`, db.Escape(req.DeviceId), user.Id)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	userPoints, _ := decimal.NewFromString(rows[0].Str(0))
	if CheckWithCode(userPoints.LessThan(points), NOT_ENOUGH_POINTS_ERROR, "not enough points", c) {
		return
	}

	token, err := uuid.NewV4()
	if CheckErr(err, c) {
		return
	}
	outerOrderId := utils.Sha1(token.String())

	orderRequest := dycdp.CreateCreateCdpOrderRequest()
	orderRequest.PhoneNumber = userMobile
	orderRequest.OfferId = requests.NewInteger64(int64(req.OfferId))
	orderRequest.OutOrderId = outerOrderId
	orderResponse, err := cdpClient.CreateCdpOrder(orderRequest)
	if CheckErr(err, c) {
		return
	}
	if orderResponse.Code == "OK" {
		pointsPerTs, _ := common.GetTMMPerTs(Config, Service)
		ts := points.Div(pointsPerTs)
		_, _, err := db.Query(`UPDATE tmm.devices SET points = IF(points > %s, points - %s, 0), consumed_ts = consumed_ts + %d WHERE id='%s'`, points.String(), points.String(), ts.IntPart(), db.Escape(req.DeviceId))
		if err != nil {
			log.Error(err.Error())
		} else {
			_, _, err = db.Query(`INSERT INTO tmm.cdp_orders (outer_order_id, device_id, order_id, points, offer_id, grade) VALUES ('%s', '%s', '%s', %s, %d, '%s')`, db.Escape(outerOrderId), db.Escape(req.DeviceId), db.Escape(orderResponse.Data.OrderId), points.String(), req.OfferId, db.Escape(offer.Grade))
			if err != nil {
				log.Error(err.Error())
			}
		}
	} else if Check(orderResponse.Message != "", orderResponse.Message, c) {
		return
	}

	params := slack.PostMessageParameters{Parse: "full", UnfurlMedia: true, Markdown: true}
	attachment := slack.Attachment{
		Color:      "success",
		AuthorName: user.ShowName,
		AuthorIcon: user.Avatar,
		Title:      "Redeem Mobile Data",
		Fallback:   "Fallback message",
		Fields: []slack.AttachmentField{
			{
				Title: "CountryCode",
				Value: strconv.FormatUint(uint64(user.CountryCode), 10),
				Short: true,
			},
			{
				Title: "UserID",
				Value: strconv.FormatUint(user.Id, 10),
				Short: true,
			},
			{
				Title: "Points",
				Value: points.StringFixed(4),
				Short: true,
			},
			{
				Title: "Grade",
				Value: offer.Grade,
				Short: true,
			},
			{
				Title: "Price",
				Value: strconv.FormatUint(offer.Price, 10),
				Short: true,
			},
			{
				Title: "Discount",
				Value: offer.Discount.String(),
				Short: true,
			},
		},
		Ts: json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	params.Attachments = []slack.Attachment{attachment}
	_, _, err = Service.Slack.PostMessage(Config.Slack.OpsChannel, "Redeem Mobile Data", params)
	if err != nil {
		log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
