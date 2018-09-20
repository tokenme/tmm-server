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

func GetHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	deviceId := c.Param("deviceId")

	db := Service.Db
	query := `SELECT
    d.id,
    d.platform,
    d.device_name,
    d.model,
    d.idfa,
    d.is_tablet,
    d.total_ts,
    COUNT(da.app_id),
    SUM(a.tmm),
    d.points,
    d.lastping_at,
    d.inserted_at,
    d.updated_at
FROM devices AS d
LEFT JOIN device_apps AS da ON (da.device_id=d.id)
INNER JOIN apps AS a ON (a.id=da.app_id)
WHERE d.user_id=%d AND d.id = '%s'
GROUP BY d.id`
	rows, _, err := db.Query(query, user.Id, deviceId)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	row := rows[0]
	var isTablet bool
	if row.Uint(5) == 1 {
		isTablet = true
	}
	tmmBalance, _ := decimal.NewFromString(row.Str(8))
	points, _ := decimal.NewFromString(row.Str(9))
	device := common.Device{
		Id:         row.Str(0),
		Platform:   row.Str(1),
		Name:       row.Str(2),
		Model:      row.Str(3),
		Idfa:       row.Str(4),
		IsTablet:   isTablet,
		TotalTs:    row.Uint64(6),
		TotalApps:  row.Uint(7),
		TMMBalance: tmmBalance,
		Points:     points,
		LastPingAt: row.ForceLocaltime(10).Format(time.RFC3339),
		InsertedAt: row.ForceLocaltime(11).Format(time.RFC3339),
		UpdatedAt:  row.ForceLocaltime(12).Format(time.RFC3339),
	}
	device.GrowthFactor, _ = device.GetGrowthFactor(Service)
	c.JSON(http.StatusOK, device)
}
