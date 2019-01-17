package account_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth/utils"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"strconv"
	"strings"
)

func UserInfoHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.DefaultQuery(`id`, `-1`))
	if CheckErr(err, c) {
		return
	}
	if Check(id < 1, `错误参数`, c) {
		return
	}

	query := `
SELECT 
	u.id AS id ,
	u.mobile AS mobile ,
	u.country_code AS country_code,
    DATE_ADD(u.created,INTERVAL 8 HOUR) created,
	wx.nick AS nick,
	DATE_ADD(wx.inserted_at,INTERVAL 8 HOUR) AS wx_inserted_at,
	wx.union_id AS union_id,
	wx.expires AS expires,
	wx.open_id AS open_id,
	u.wallet_addr AS addr,
	uc.cny AS uc_cny,
	point.cny AS point_cny,
	IFNULL(uc.cny,0)+IFNULL(point.cny,0) AS cny,
	IFNULL(SUM(dev.points),0) AS points,
	inv.direct AS direct,
	inv.indirect AS indirect,
	inv.active AS active,
	inv.invite_firend_active AS invite_firend_active,
	inv.invite_By_Number AS invite_By_Number,
	inv.direct_blocked AS direct_blocked,
	inv.indirect_blocked AS indirect_blocked,
	bonus.inv_bonus AS inv_bonus,
	sha.points AS sha_points,
	app.points AS app_points,
	reading.point AS reading_point,	
	IFNULL(parent.id,0) AS parent_id,
	IFNULL(parent.blocked,0) AS parent_blocked,
	IFNULL(parent.nick,NULL) AS parent_nick,	
	IFNULL(root.id,0) AS root_id,
	IFNULL(root.blocked,0) AS root_blocked,
	IFNULL(root.nick,NULL) AS root_nick,
	IF(us_set.user_id > 0,IF(us_set.blocked = us_set.block_whitelist,0,1),0) AS blocked,
	us_set.comments AS message,
	us_set.level AS _level,
	IF(EXISTS(
		SELECT 
		1
		FROM tmm.devices AS dev 
		LEFT JOIN tmm.device_share_tasks AS sha ON (sha.device_id = dev.id AND sha.inserted_at > DATE_SUB(NOW(),INTERVAL 3 DAY))
		LEFT JOIN tmm.device_app_tasks AS app ON (app.device_id = dev.id  AND  app.inserted_at > DATE_SUB(NOW(),INTERVAL 3 DAY))
		LEFT JOIN reading_logs AS reading ON (reading.user_id = dev.user_id  AND (reading.updated_at > DATE_SUB(NOW(),INTERVAL 3 DAY) OR  reading.inserted_at > DATE_SUB(NOW(),INTERVAL 3 DAY)))
		LEFT JOIN tmm.daily_bonus_logs AS daily ON (daily.user_id = dev.user_id AND daily.updated_on >= DATE_SUB(NOW(),INTERVAL 3 DAY))  
		WHERE dev.user_id = u.id AND ( 
		sha.task_id > 0  
		OR app.task_id > 0   OR reading.user_id > 0
		OR daily.user_id > 0)
		LIMIT 1
	),TRUE,FALSE) AS _active,
	IF(dev_app.app_id IS NOT NULL,TRUE,FALSE) AS  app_id,
	three_active.total AS total
FROM 
	ucoin.users AS u
LEFT JOIN  
	tmm.devices AS dev ON dev.user_id = u.id
LEFT JOIN
	tmm.device_apps AS dev_app ON dev_app.device_id = dev.id
LEFT JOIN 
	tmm.user_settings AS us_set ON (us_set.user_id = u.id )
LEFT JOIN 
	tmm.wx AS wx ON (wx.user_id =u.id)
LEFT JOIN (
	SELECT 
		IFNULL(wx.nick,us.nickname) AS nick,
		us.id AS id ,
		IF(us_set.user_id > 0,IF(us_set.blocked = us_set.block_whitelist,0,1),0) AS blocked
	FROM 
		tmm.invite_codes  AS ic
	INNER JOIN 
		ucoin.users AS us  ON us.id = ic.parent_id
	LEFT JOIN 
		tmm.user_settings AS us_set ON (us_set.user_id = us.id )
	LEFT JOIN 
		tmm.wx AS wx ON wx.user_id = us.id 
	WHERE 
		ic.user_id = %d
) AS parent ON 1 = 1 
LEFT JOIN (
	SELECT 
		IFNULL(wx.nick,us.nickname) AS nick,
		us.id AS id ,
		IF(us_set.user_id > 0,IF(us_set.blocked = us_set.block_whitelist,0,1),0) AS blocked
	FROM 
		tmm.invite_codes  AS ic
	INNER JOIN 
		ucoin.users AS us  ON us.id = ic.root_id
	LEFT JOIN 
		tmm.user_settings AS us_set ON (us_set.user_id = us.id )
	LEFT JOIN 
		tmm.wx AS wx ON wx.user_id = us.id 
	WHERE 
		ic.user_id = %d
) AS root ON 1 = 1 
LEFT JOIN (
	SELECT 
		SUM(cny) AS cny
	FROM 
		tmm.withdraw_txs
	WHERE 
		user_id = %d AND tx_status = 1
) AS uc ON 1 = 1
LEFT JOIN (
	SELECT 
		SUM(cny) AS cny
	FROM 
		tmm.point_withdraws 
		WHERE user_id = %d
) AS point ON 1 = 1
LEFT JOIN (
	SELECT
		COUNT(IF(inv.parent_id = %d,inv.user_id,NULL)) AS direct,
		COUNT(IF(inv.grand_id = %d,inv.user_id,NULL)) AS indirect,
		COUNT(IF(inv.parent_id = %d AND EXISTS(
		SELECT 1 FROM user_settings AS us  WHERE us.blocked= 1 AND us.user_id=inv.user_id AND us.block_whitelist=0  LIMIT 1
		),1,NULL)) AS direct_blocked,
		COUNT(IF(inv.grand_id = %d  AND EXISTS(
		SELECT 1 FROM user_settings AS us  WHERE us.blocked= 1 AND us.user_id=inv.user_id AND us.block_whitelist=0  LIMIT 1
		),1,NULL)) AS indirect_blocked,
		COUNT( 
			IF(EXISTS(
			SELECT 
				1
			FROM 
				tmm.devices AS dev 
			LEFT JOIN 
				tmm.device_share_tasks AS sha ON (sha.device_id = dev.id AND sha.inserted_at > DATE_SUB(NOW(),INTERVAL 3 DAY))
			LEFT JOIN 
				tmm.device_app_tasks AS app ON (app.device_id = dev.id  AND  app.inserted_at > DATE_SUB(NOW(),INTERVAL 3 DAY))
			LEFT JOIN 
				reading_logs AS reading ON (reading.user_id = dev.user_id  AND (reading.updated_at > DATE_SUB(NOW(),INTERVAL 3 DAY) OR  reading.inserted_at > DATE_SUB(NOW(),INTERVAL 3 DAY)))
			LEFT JOIN 
				tmm.daily_bonus_logs AS daily ON (daily.user_id = dev.user_id AND daily.updated_on >= DATE_SUB(NOW(),INTERVAL 3 DAY))  
			WHERE 
				dev.user_id = inv.user_id AND ( 
				sha.task_id > 0	OR app.task_id > 0 OR 
				reading.user_id > 0 OR daily.user_id > 0)
			LIMIT 1
			),inv.user_id,NULL)
		) AS active,
		COUNT(
		IF(u.id > 0  AND EXISTS(
			SELECT 
				1
			FROM 
				tmm.devices AS dev 
			LEFT JOIN 
				tmm.device_share_tasks AS sha ON (sha.device_id = dev.id AND sha.inserted_at > DATE_SUB(NOW(),INTERVAL 3 DAY))
			LEFT JOIN 
				tmm.device_app_tasks AS app ON (app.device_id = dev.id  AND  app.inserted_at > DATE_SUB(NOW(),INTERVAL 3 DAY))
			LEFT JOIN 
				reading_logs AS reading ON (reading.user_id = dev.user_id  AND (reading.updated_at > DATE_SUB(NOW(),INTERVAL 3 DAY) OR  reading.inserted_at > DATE_SUB(NOW(),INTERVAL 3 DAY)))
			LEFT JOIN 
				tmm.daily_bonus_logs AS daily ON (daily.user_id = dev.user_id AND daily.updated_on >= DATE_SUB(NOW(),INTERVAL 3 DAY))  
			WHERE 
				dev.user_id = u.id AND ( 
				sha.task_id > 0 OR app.task_id > 0  OR
				reading.user_id > 0 OR daily.user_id > 0)
			LIMIT 1
			),inv.user_id,NULL) 
		) AS invite_firend_active,
		COUNT(IF(u.id>0,1,NULL)) AS invite_By_Number
	FROM
		tmm.invite_codes  AS inv
	LEFT JOIN 
		ucoin.users AS u ON (u.id = inv.user_id AND u.created > DATE_SUB(NOW(),INTERVAL 3 DAY))
	WHERE
		inv.parent_id = %d OR inv.grand_id = %d
) AS inv ON 1 = 1
LEFT JOIN (
	SELECT 
		IFNULL(SUM(sha.points),0)+IFNULL(bonus.points,0) AS points
	FROM
		tmm.device_share_tasks AS sha
	INNER JOIN
		tmm.devices AS dev ON  (dev.id = sha.device_id)
	LEFT JOIN (
	SELECT 
		  user_id, 
		  SUM(bonus) AS points
	FROM 
		tmm.invite_bonus 
	WHERE  task_type != 0  AND user_id = %d
	) AS bonus ON bonus.user_id = dev.user_id
	WHERE
		dev.user_id = %d
) AS sha ON 1 = 1 
LEFT JOIN (
	SELECT
		SUM(IF(app.status = 1,app.points,0)) AS points
	FROM
		tmm.device_app_tasks AS app
	INNER JOIN
		tmm.devices AS dev ON  (dev.id = app.device_id)
	WHERE
		dev.user_id = %d AND app.status = 1
) AS app ON 1 = 1
LEFT JOIN (
	SELECT
		SUM(bonus) AS inv_bonus
	FROM 
		tmm.invite_bonus 
	WHERE 
		user_id = %d  AND task_type  = 0
) AS bonus ON 1 = 1
LEFT JOIN (
	SELECT 
		SUM(point) AS point
	FROM
		tmm.reading_logs
	WHERE user_id = %d
) AS reading ON 1 = 1
LEFT JOIN (
	SELECT
		COUNT(DISTINCT IF(sha.task_id > 0  OR app.task_id > 0  OR reading.user_id  > 0   ,inv.user_id,NULL)) AS total
	FROM
		tmm.invite_codes AS inv
	INNER JOIN 
		ucoin.users AS u ON u.id = inv.user_id
	LEFT JOIN
		tmm.devices AS dev ON dev.user_id = u.id
	LEFT JOIN 
		tmm.device_share_tasks AS sha ON (sha.device_id = dev.id AND sha.inserted_at < DATE_ADD(u.created,INTERVAL 3 DAY))
	LEFT JOIN 
		tmm.device_app_tasks AS app ON (app.device_id = dev.id  AND  app.inserted_at < DATE_ADD(u.created,INTERVAL 3 DAY))
	LEFT JOIN 
		reading_logs AS reading ON (reading.user_id = dev.user_id  AND reading.inserted_at < DATE_ADD(u.created,INTERVAL 3 DAY))
	WHERE   
		inv.parent_id = %d OR inv.grand_id = %d
) AS three_active ON 1 = 1
WHERE 
	u.id = %d
LIMIT 1 
	`

	db := Service.Db
	rows, res, err := db.Query(query, id, id, id, id, id, id, id, id, id, id, id, id, id, id, id, id, id, id)
	if CheckErr(err, c) {
		return
	}

	if Check(len(rows) == 0, admin.Not_Found, c) {
		return
	}

	row := rows[0]
	point, err := decimal.NewFromString(row.Str(res.Map(`points`)))
	if CheckErr(err, c) {
		return
	}

	tokenABI, err := utils.NewToken(Config.TMMTokenAddress, Service.Geth)
	if CheckErr(err, c) {
		return
	}
	_, _, decimals, _, _, _, _, _, balance, err := utils.TokenMeta(tokenABI, row.Str(res.Map(`addr`)))
	balanceDecimal, err := decimal.NewFromString(balance.String())
	tmm := balanceDecimal.Div(decimal.New(1, int32(decimals)))

	user := &admin.UserStats{
		IsActive:           row.Bool(res.Map(`_active`)),
		IsHaveAppId:              row.Bool(res.Map(`app_id`)),
		BlockedMessage:           row.Str(res.Map(`message`)),
	}
	user.DrawCashByPoint = fmt.Sprintf("%.2f", row.Float(res.Map(`point_cny`)))
	user.DrawCashByUc = fmt.Sprintf("%.2f", row.Float(res.Map(`uc_cny`)))
	user.DrawCash = fmt.Sprintf("%.2f", row.Float(res.Map(`cny`)))

	user.DirectFriends = row.Int(res.Map(`direct`))
	user.IndirectFriends = row.Int(res.Map(`indirect`))
	user.ActiveFriends = row.Int(res.Map(`active`))
	user.DirectBlockedCount = row.Int(res.Map(`direct_blocked`))
	user.InDirectBlockedCount = row.Int(res.Map(`indirect_blocked`))
	user.InviteNewUserByThreeDays = row.Int(res.Map(`invite_By_Number`))
	user.InviteNewUserActiveCount = row.Int(res.Map(`invite_firend_active`))

	user.PointByShare = int(row.Float(res.Map(`sha_points`)))
	user.PointByReading = int(row.Float(res.Map(`reading_point`)))
	user.PointByInvite = int(row.Float(res.Map(`inv_bonus`)))
	user.PointByDownLoadApp = int(row.Float(res.Map(`app_points`)))

	user.Id = row.Uint64(res.Map(`id`))
	user.Mobile = row.Str(res.Map(`mobile`))
	user.Nick = row.Str(res.Map(`nick`))
	user.InsertedAt = row.Str(res.Map(`created`))
	user.CountryCode = row.Uint(res.Map(`country_code`))
	user.Level.Id = row.Uint(res.Map(`_level`))
	user.Wallet = row.Str(res.Map(`addr`))
	user.Tmm = tmm.StringFixed(2)
	user.Point = point.StringFixed(0)
	user.Blocked = row.Int(res.Map(`blocked`))

	user.WxInsertedAt = row.Str(res.Map(`wx_inserted_at`))
	user.OpenId = row.Str(res.Map(`open_id`))
	user.WxUnionId = row.Str(res.Map(`union_id`))

	user.TotalMakePoint = user.PointByShare + user.PointByReading +
		user.PointByInvite + user.PointByDownLoadApp
	threeActiveCount := row.Float(res.Map(`total`))
	if user.DirectFriends+user.IndirectFriends > 0 && threeActiveCount > 0 {
		user.NotActive = fmt.Sprintf("%.2f", 100-threeActiveCount/float64(user.DirectFriends+user.IndirectFriends)*100) + "%"
	} else {
		if user.DirectFriends+user.InDirectBlockedCount > 0 && threeActiveCount == 0 {
			user.NotActive = fmt.Sprint("100%")
		} else {
			user.NotActive = fmt.Sprint("0%")
		}
	}

	parent := admin.User{}
	parent.Id = row.Uint64(res.Map(`parent_id`))
	parent.Blocked = row.Int(res.Map(`parent_blocked`))
	parent.Nick = row.Str(res.Map(`parent_nick`))
	user.Parent = parent

	root := admin.User{}
	root.Id = row.Uint64(res.Map(`root_id`))
	root.Blocked = row.Int(res.Map(`root_blocked`))
	root.Nick = row.Str(res.Map(`root_nick`))
	user.Root = root

	rows, _, err = db.Query(`SELECT 
		id,
		platform,
		idfa,
		imei,
		mac,
		device_name,
		system_version,
		os_version,
		language,
		model,
		timezone,
		country,
		is_emulator,
		is_jailbrojen,
		is_tablet,
		points,
		total_ts,
		tmp_ts,
		consumed_ts,
		lastping_at,
		inserted_at,
		updated_at
FROM 
	tmm.devices
WHERE 
	user_id = %d`, id)
	if CheckErr(err, c) {
		return
	}

	for _, row := range rows {
		device := &admin.Device{}
		device.Id = row.Str(0)
		device.Platform = row.Str(1)
		device.Idfa = row.Str(2)
		device.Imei = row.Str(3)
		device.Mac = row.Str(4)
		device.Name = row.Str(5)
		device.SystemVersion = row.Str(6)
		device.OsVersion = row.Str(7)
		device.Language = row.Str(8)
		device.Model = row.Str(9)
		device.Timezone = row.Str(10)
		device.Country = row.Str(11)
		device.IsEmulator = row.Bool(12)
		device.IsJailbrojen = row.Bool(13)
		device.IsTablet = row.Bool(14)
		device.Points = decimal.NewFromFloat(row.Float(15)).Ceil()
		device.TotalTs = row.Uint64(16)
		device.TmpTs = row.Uint64(17)
		device.ConsumedTs = row.Float(18)
		device.LastPingAt = row.Str(19)
		device.InsertedAt = row.Str(20)
		device.UpdatedAt = row.Str(21)
		user.DeviceList = append(user.DeviceList, device)
	}

	if user.OpenId != "" {
		rows, _, err = db.Query(`
	select 
		group_concat(user_id), 
		count(1) num, open_id, nick
	from 
		tmm.wx
	WHERE 
		open_id = '%s'
	group by 
		open_id
	having num > 1`, user.OpenId)
		if CheckErr(err, c) {
			return
		}
		if len(rows) > 0 {
			var accountList []string
			for _, account := range strings.Split(rows[0].Str(0), ",") {
				if account == strconv.Itoa(id) {
					continue
				}
				accountList = append(accountList, account)
			}
			user.OtherAccount = accountList
		}
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    user,
	})
}
