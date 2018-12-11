package device

import (
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	ADZONE_TOP_ID    = 3
	ADZONE_BOTTOM_ID = 4
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
    d.idfa,
    d.is_tablet,
    d.total_ts,
    COUNT(da.app_id),
    SUM(a.tmm),
    d.points,
    d.lastping_at,
    d.inserted_at,
    d.updated_at,
    d.mac,
    d.imei
FROM devices AS d
LEFT JOIN device_apps AS da ON (da.device_id=d.id)
INNER JOIN apps AS a ON (a.id=da.app_id)
WHERE d.user_id=%d
GROUP BY d.id`
	rows, _, err := db.Query(query, user.Id)
	if CheckErr(err, c) {
		return
	}
	platform := c.GetString("tmm-platform")
	buildVersionStr := c.GetString("tmm-build")
	buildVersion, _ := strconv.ParseUint(buildVersionStr, 10, 64)
	adsMap := make(map[int][]*common.Adgroup)
	if platform == "ios" && buildVersion > 42 || platform == "android" && buildVersion > 211 {
		adsMap, err = getCreatives(platform)
		if err != nil {
			log.Error(err.Error())
		}
	}
	var devices []common.Device
	if adgroups, found := adsMap[ADZONE_TOP_ID]; found {
		adgroupIdx := rand.Intn(len(adgroups))
		creatives := adgroups[adgroupIdx].Creatives
		creativeIdx := rand.Intn(len(creatives))
		d := common.Device{
			Creative: creatives[creativeIdx],
		}
		devices = append(devices, d)
	}
	for _, row := range rows {
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
			Mac:        row.Str(13),
			Imei:       row.Str(14),
		}
		device.GrowthFactor, _ = device.GetGrowthFactor(Service)
		devices = append(devices, device)
	}

	if adgroups, found := adsMap[ADZONE_BOTTOM_ID]; found {
		adgroupIdx := rand.Intn(len(adgroups))
		creatives := adgroups[adgroupIdx].Creatives
		creativeIdx := rand.Intn(len(creatives))
		d := common.Device{
			Creative: creatives[creativeIdx],
		}
		devices = append(devices, d)
	}
	c.JSON(http.StatusOK, devices)
}

func getCreatives(platform string) (map[int][]*common.Adgroup, error) {
	adsMap := make(map[int][]*common.Adgroup)
	adgroupsMap := make(map[uint64]*common.Adgroup)
	var constraint string
	if platform == "ios" {
		constraint = " AND c.platform IN (0, 1)"
	} else {
		constraint = " AND c.platform IN (0, 2)"
	}
	db := Service.Db
	query := `SELECT
        c.id,
        c.adgroup_id,
        c.image,
        c.link,
        c.width,
        c.height,
        z.id,
        c.share_image,
        c.title
    FROM tmm.creatives AS c
    INNER JOIN tmm.adgroups AS a ON (a.id=c.adgroup_id)
    INNER JOIN tmm.adzones AS z ON (z.id=a.adzone_id)
    WHERE z.id IN (%d, %d) AND a.online_status=1 AND c.online_status=1%s`
	rows, _, err := db.Query(query, ADZONE_TOP_ID, ADZONE_BOTTOM_ID, constraint)
	if err != nil {
		return nil, err
	} else if len(rows) > 0 {
		for _, row := range rows {
			adgroupId := row.Uint64(1)
			adzoneId := row.Int(6)
			creative := &common.Creative{
				Id:         row.Uint64(0),
				AdgroupId:  adgroupId,
				Image:      row.Str(2),
				Link:       row.Str(3),
				Width:      row.Uint(4),
				Height:     row.Uint(5),
                ShareImage: row.Str(7),
                Title:      row.Str(8),
			}
			creativeCode, err := creative.Code([]byte(Config.LinkSalt))
			if err != nil {
				continue
			}
			creative.Image = fmt.Sprintf("%s/%s", Config.AdImpUrl, creativeCode)
			creative.Link = fmt.Sprintf("%s/%s", Config.AdClkUrl, creativeCode)
			if ad, found := adgroupsMap[adgroupId]; found {
				ad.Creatives = append(ad.Creatives, creative)
			} else {
				ad := &common.Adgroup{
					Id:        adgroupId,
					Creatives: []*common.Creative{creative},
				}
				adsMap[adzoneId] = append(adsMap[adzoneId], ad)
			}
		}
	}
	return adsMap, nil
}
