package account_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
)

type types int

const (
	Direct types = iota
	Indirect
	Children
	Active
)

func FriendsHandler(c *gin.Context) {
	var req PageOptions
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if Check(req.Id < 0, admin.Error_Param, c) {
		return
	}

	var offset int
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Page > 0 {
		offset = (req.Page - 1) * req.Limit
	}

	query := `
	SELECT
		inv.user_id,
		u.mobile,
		wx.nick,
		DATE_ADD(u.created,INTERVAL 8 HOUR),
		IF(IFNULL(us.blocked, 0)=0 OR us.block_whitelist=1, 0, 1) AS blocked,
		IF(COUNT(IF(
			(sha.inserted_at > u.created AND sha.inserted_at < DATE_ADD(u.created,INTERVAL 1 day) )  OR
			(app.inserted_at > u.created AND app.inserted_at < DATE_ADD(u.created,INTERVAL 1 day) )  OR
			(reading.inserted_at > u.created AND reading.inserted_at < DATE_ADD(u.created,INTERVAL 1 day))
		,1,NULL)) > 0,TRUE,FALSE
		) AS flrst_day_active,
		IF(COUNT(IF(
			(sha.inserted_at > DATE_ADD(u.created,INTERVAL 1 day) AND sha.inserted_at < DATE_ADD(u.created,INTERVAL 2 day)) OR
			(app.inserted_at > DATE_ADD(u.created,INTERVAL 1 day) AND app.inserted_at < DATE_ADD(u.created,INTERVAL 2 day)) OR
			(reading.inserted_at > DATE_ADD(u.created,INTERVAL 1 day) AND reading.inserted_at < DATE_ADD(u.created,INTERVAL 2 day)) OR
			(reading.updated_at > DATE_ADD(u.created,INTERVAL 1 day) AND reading.updated_at < DATE_ADD(u.created,INTERVAL 2 day))
			,1,NULL)) > 0,TRUE,FALSE
		) AS second_day_active,
		IF(COUNT(IF(
			sha.inserted_at > DATE_ADD(u.created,INTERVAL 2 day) OR
			app.inserted_at > DATE_ADD(u.created,INTERVAL 2 day) OR
			reading.updated_at > DATE_ADD(u.created,INTERVAL 2 day) OR
			reading.inserted_at > DATE_ADD(u.created,INTERVAL 2 day)
			,1,NULL)) > 0,TRUE,FALSE
		) AS three_day_active,
		bonus.bonus,
		IF(inv.parent_id=%d, 0, 1 ) AS fiends_types,
		IF(dev_app.app_id IS NOT NULL,TRUE,FALSE) AS app_id
	FROM tmm.invite_codes AS inv
	INNER JOIN ucoin.users AS u ON (u.id=inv.user_id)
	LEFT JOIN tmm.wx AS wx ON (wx.user_id=u.id)
	LEFT JOIN tmm.user_settings AS us ON (us.user_id=u.id)
    LEFT JOIN tmm.devices AS dev ON (dev.user_id=u.id)
    LEFT JOIN tmm.device_apps AS dev_app ON (dev_app.device_id=dev.id)
	LEFT JOIN(
		SELECT
			from_user_id,
			SUM(bonus) AS bonus
		FROM tmm.invite_bonus
		WHERE user_id=%d
		GROUP BY from_user_id
		) AS bonus ON (bonus.from_user_id=inv.user_id)
	LEFT JOIN tmm.device_share_tasks AS sha ON (sha.device_id = dev.id AND sha.inserted_at < DATE_ADD(u.created,INTERVAL 3 day ) )
	LEFT JOIN tmm.device_app_tasks AS app ON (app.device_id = dev.id AND app.inserted_at < DATE_ADD(u.created,INTERVAL 3 day ))
	LEFT JOIN reading_logs AS reading ON (reading.user_id = dev.user_id AND (reading.inserted_at <  DATE_ADD(u.created,INTERVAL 3 day ) OR reading.updated_at <  DATE_ADD(u.created,INTERVAL 3 day ) ) )
 	WHERE %s
	GROUP BY inv.user_id
	ORDER BY inv.user_id DESC
	LIMIT %d OFFSET %d`
	totalquery := `SELECT COUNT(DISTINCT inv.user_id) FROM tmm.invite_codes AS inv LEFT JOIN tmm.devices AS dev ON (dev.user_id=inv.user_id) WHERE %s`

	switch types(req.Types) {
	case Direct:
		direct := fmt.Sprintf("inv.parent_id=%d", req.Id)
		query = fmt.Sprintf(query, req.Id, req.Id, direct, req.Limit, offset)
		totalquery = fmt.Sprintf(totalquery, direct)

	case Indirect:
		indirect := fmt.Sprintf("inv.grand_id=%d", req.Id)
		query = fmt.Sprintf(query, req.Id, req.Id, indirect, req.Limit, offset)
		totalquery = fmt.Sprintf(totalquery, indirect)

	case Children:
		online := fmt.Sprintf("inv.parent_id=%d OR inv.grand_id=%d", req.Id, req.Id)
		query = fmt.Sprintf(query, req.Id, req.Id, online, req.Limit, offset)
		totalquery = fmt.Sprintf(totalquery, online)

	case Active:
		active := fmt.Sprintf(`(inv.parent_id=%d OR inv.grand_id=%d) AND (
            EXISTS (SELECT 1 FROM tmm.device_share_tasks AS dst WHERE dst.device_id=dev.id AND dst.inserted_at>=DATE_SUB(NOW(), INTERVAL 3 DAY) LIMIT 1) OR
            EXISTS (SELECT 1 FROM tmm.device_app_tasks AS dat WHERE dat.device_id=dev.id AND dat.inserted_at>=DATE_SUB(NOW(), INTERVAL 3 DAY) LIMIT 1) OR
            EXISTS (SELECT 1 FROM tmm.reading_logs AS rl WHERE rl.user_id=inv.user_id AND (rl.inserted_at>=DATE_SUB(NOW(), INTERVAL 3 DAY) OR rl.updated_at>=DATE_ADD(NOW(), INTERVAL 3 DAY)) LIMIT 1) OR
            EXISTS (SELECT 1 FROM tmm.daily_bonus_logs AS dbl WHERE dbl.user_id=inv.user_id AND dbl.updated_on>=DATE_SUB(NOW(), INTERVAL 3 DAY) LIMIT 1))`, req.Id, req.Id)
		query = fmt.Sprintf(query, req.Id, req.Id, active, req.Limit, offset)
		totalquery = fmt.Sprintf(totalquery, active)
	default:
		c.JSON(http.StatusOK, nil)
		return
	}

	var List []*admin.UserStats
	db := Service.Db
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}

	if Check(len(rows) == 0, admin.Not_Found, c) {
		return
	}

	for _, row := range rows {
		user := &admin.UserStats{}
		user.Id = row.Uint64(0)
		user.Mobile = row.Str(1)
		user.Nick = row.Str(2)
		user.Created = row.Str(3)
		user.Blocked = row.Int(4)
		user.FirstDayActive = row.Bool(5)
		user.SecondDayActive = row.Bool(6)
		user.ThreeDayActive = row.Bool(7)
		user.FriendBonus = fmt.Sprintf("%.1f", row.Float(8))
		user.FriendType = FirendMap[row.Int(9)]
		user.IsHaveAppId = row.Bool(10)
		List = append(List, user)
	}

	var total int
	rows, _, err = db.Query(totalquery)
	if CheckErr(err, c) {
		return
	}
	if len(rows) > 0 {
		total = rows[0].Int(0)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			"data":  List,
			"total": total,
		},
	})
}
