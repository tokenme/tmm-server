package device

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"time"
)

func AppsHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	deviceId := c.Param("deviceId")

	db := Service.Db
	query := `SELECT
    a.id,
    a.platform,
    a.name,
    a.bundle_id,
    a.store_id,
    a.app_version,
    a.build_version,
    da.total_ts,
    a.tmm,
    da.lastping_at,
    da.inserted_at,
    da.updated_at
FROM device_apps AS da
INNER JOIN devices AS d ON (d.id=da.device_id)
INNER JOIN apps AS a ON (a.id=da.app_id AND a.platform=d.platform)
INNER JOIN user_devices AS ud ON (ud.device_id=da.device_id)
WHERE ud.user_id=%d AND da.device_id='%s'`
	rows, _, err := db.Query(query, user.Id, deviceId)
	if CheckErr(err, c) {
		return
	}
	var apps []common.App
	for _, row := range rows {
		tmmBalance, _ := decimal.NewFromString(row.Str(8))
		app := common.App{
			Id:           row.Str(0),
			Platform:     row.Str(1),
			Name:         row.Str(2),
			BundleId:     row.Str(3),
			StoreId:      row.Uint64(4),
			Version:      row.Str(5),
			BuildVersion: row.Str(6),
			Ts:           row.Uint64(7),
			TMMBalance:   tmmBalance,
			LastPingAt:   row.ForceLocaltime(9).Format(time.RFC3339),
			InsertedAt:   row.ForceLocaltime(10).Format(time.RFC3339),
			UpdatedAt:    row.ForceLocaltime(11).Format(time.RFC3339),
		}
		app.GrowthFactor, _ = app.GetGrowthFactor(Service)
		if app.StoreId == 0 {
			lookup, err := app.LookUp(Service)
			if err == nil {
				app.StoreId = lookup.TrackId
				app.Icon = lookup.ArtworkUrl512
			}
		}
		apps = append(apps, app)
	}
	c.JSON(http.StatusOK, apps)
}
