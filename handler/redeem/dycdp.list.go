package redeem

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/dycdp"
	"github.com/tokenme/tmm/tools/forex"
	"github.com/xluohome/phonedata"
	"net/http"
	"sort"
)

func DycdpListHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	if CheckWithCode(user.CountryCode != 86, INVALID_CDP_VENDOR_ERROR, "the cdp vendor not supported", c) {
		return
	}
	phone, err := phonedata.Find(user.Mobile)
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
	cdpRequest := dycdp.CreateQueryCdpOfferRequest()
	cdpRequest.ChannelType = "分省"
	cdpRequest.Vendor = phone.CardType
	cdpRequest.Province = phone.Province
	cdpResponse, err := cdpClient.QueryCdpOffer(cdpRequest)
	if CheckErr(err, c) {
		return
	}
	offers := cdpResponse.FlowOffers.FlowOffer
	{
		cdpRequest := dycdp.CreateQueryCdpOfferRequest()
		cdpRequest.ChannelType = "全国"
		cdpRequest.Vendor = phone.CardType
		cdpResponse, err := cdpClient.QueryCdpOffer(cdpRequest)
		if CheckErr(err, c) {
			return
		}
		newOffers := cdpResponse.FlowOffers.FlowOffer
		offers = append(offers, newOffers...)
	}
	var redeemCdps common.RedeemCdpSlice
	gradeMap := make(map[string]*common.RedeemCdp)
	for _, offer := range offers {
		if offer.Discount.GreaterThan(decimal.NewFromFloat(0.7)) {
			continue
		}
		price := decimal.New(int64(offer.Price), 0).Div(decimal.New(100, 0))
		if grade, found := gradeMap[offer.Grade]; found {
			if grade.Price.GreaterThan(price) {
				grade.OfferId = offer.OfferId
				grade.Price = price
			}
		} else {
			gradeMap[offer.Grade] = &common.RedeemCdp{
				OfferId: offer.OfferId,
				Grade:   offer.Grade,
				Price:   price,
			}
		}
	}
	pointPrice := common.GetPointPrice(Service, Config)
	currencyRate := forex.Rate(Service, "USD", "CNY")
	pointPrice = pointPrice.Mul(currencyRate)
	if pointPrice.GreaterThan(decimal.Zero) {
		for _, grade := range gradeMap {
			grade.Points = grade.Price.Div(pointPrice)
			redeemCdps = append(redeemCdps, grade)
		}
	}
	sort.Sort(redeemCdps)
	c.JSON(http.StatusOK, redeemCdps)
}
