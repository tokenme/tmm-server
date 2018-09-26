package app

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strconv"
)

const (
	DEFAULT_PAGE_SIZE = 10
)

func SdksHandler(c *gin.Context) {
	_, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}

	platform := c.Param("platform")
	page, _ := strconv.ParseUint(c.Param("page"), 10, 64)
	pageSize, _ := strconv.ParseUint(c.Param("pageSize"), 10, 64)
	if page == 0 {
		page = 1
	}
	if pageSize == 0 || pageSize > DEFAULT_PAGE_SIZE {
		pageSize = DEFAULT_PAGE_SIZE
	}

	db := Service.Db
	query := `SELECT
    a.id,
    a.platform,
    a.name,
    a.bundle_id,
    a.store_id,
    a.app_version,
    a.build_version,
    a.total_ts,
    a.tmm
FROM tmm.apps AS a
WHERE a.platform='%s' AND a.is_active=1 ORDER BY a.total_ts DESC LIMIT %d, %d`
	rows, _, err := db.Query(query, platform, (page-1)*pageSize, pageSize)
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
