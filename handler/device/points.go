package device

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	//"github.com/davecgh/go-spew/spew"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/utils"
	"net/http"
	"strings"
)

const (
	DEVICE_POINTS_CACHE_KEY = "DevicePoints:%s"
)

func PointsHandler(c *gin.Context) {
	req := c.MustGet("Request").(APIRequest)
	secret := c.MustGet("Secret").(string)
	decrepted, err := utils.DesDecrypt(req.Payload, []byte(secret))
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	var deviceRequest common.DeviceRequest
	err = json.Unmarshal(decrepted, &deviceRequest)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if Check(deviceRequest.IsEmulator || deviceRequest.IsJailBrojen, "invalid request", c) {
		return
	}
	deviceId := deviceRequest.DeviceId()
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	cacheKey := fmt.Sprintf(DEVICE_POINTS_CACHE_KEY, deviceId)
	points, err := getPoints(Service, deviceId)
	if CheckErr(err, c) {
		return
	}
	lastPointsStr, err := redis.String(redisConn.Do("GET", cacheKey))
	redisConn.Do("SETX", cacheKey, 60*2, points.String())
	var lastPoints decimal.Decimal
	if err == nil {
		lastPoints, err = decimal.NewFromString(lastPointsStr)
		if err == nil && points.GreaterThan(lastPoints.Add(decimal.New(1, -3))) {
			title := "New UCoin Points"
			desc := fmt.Sprintf("You just earn %s UCoin points, check UCoin Wallet for more information.", points.StringFixed(4))
			icon := "https://static.tianxi100.com/ucoin/icon-ios-marketing-1024@1x.png"
			if strings.HasPrefix(deviceRequest.Language, "zh") {
				title = "UCoin 积分提醒"
				desc = fmt.Sprintf("您刚刚获得 %s UCoin积分，打开 UCoin 钱包查看详情。", points.StringFixed(4))
			}
			c.JSON(http.StatusOK, gin.H{"title": title, "desc": desc, "icon": icon})
		}
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}

func getPoints(service *common.Service, deviceId string) (points decimal.Decimal, err error) {
	db := service.Db

	query := `SELECT
                points
            FROM tmm.devices AS d
            WHERE
                d.id='%s'`
	rows, _, err := db.Query(query, db.Escape(deviceId))
	if err != nil {
		raven.CaptureError(err, nil)
		return points, err
	}
	if len(rows) == 0 {
		return points, nil
	}
	row := rows[0]
	return decimal.NewFromString(row.Str(0))
}
