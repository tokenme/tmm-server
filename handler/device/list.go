package device

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"time"
)

func ListHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db
	query := `SELECT
    d.id,
    d.platform,
    d.device_name,
    d.model,
    d.is_tablet,
    d.total_ts,
    COUNT(da.app_id),
    SUM(a.tmm),
    d.points,
    d.lastping_at,
    d.inserted_at,
    d.updated_at
FROM user_devices AS ud
INNER JOIN devices AS d ON (d.id=ud.device_id)
LEFT JOIN device_apps AS da ON (da.device_id=d.id)
INNER JOIN apps AS a ON (a.id=da.app_id)
WHERE ud.user_id=%d
GROUP BY d.id`
	rows, _, err := db.Query(query, user.Id)
	if CheckErr(err, c) {
		return
	}
	var devices []common.Device
	for _, row := range rows {
		var isTablet bool
		if row.Uint(4) == 1 {
			isTablet = true
		}
		tmmBalance, _ := decimal.NewFromString(row.Str(7))
		points, _ := decimal.NewFromString(row.Str(8))
		device := common.Device{
			Id:         row.Str(0),
			Platform:   row.Str(1),
			Name:       row.Str(2),
			Model:      row.Str(3),
			IsTablet:   isTablet,
			TotalTs:    row.Uint64(5),
			TotalApps:  row.Uint(6),
			TMMBalance: tmmBalance,
			Points:     points,
			LastPingAt: row.ForceLocaltime(9).Format(time.RFC3339),
			InsertedAt: row.ForceLocaltime(10).Format(time.RFC3339),
			UpdatedAt:  row.ForceLocaltime(11).Format(time.RFC3339),
		}
		device.GrowthFactor, _ = device.GetGrowthFactor(Service)
		devices = append(devices, device)
	}
	c.JSON(http.StatusOK, devices)
}
