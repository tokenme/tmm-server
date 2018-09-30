package device

import (
	"encoding/json"
	//"github.com/davecgh/go-spew/spew"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/utils"
	"net/http"
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
	//if Check(pingRequest.Device.IsEmulator || pingRequest.Device.IsJailBrojen, "invalid request", c) {
	//  return
	//}
	points, err := getPoints(Service, deviceRequest)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"points": points})
}

func getPoints(service *common.Service, deviceRequest common.DeviceRequest) (points decimal.Decimal, err error) {
	db := service.Db
	deviceId := deviceRequest.DeviceId()
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
