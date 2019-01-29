package account_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"strconv"
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
	dgt.points AS dgt_points,
    IFNULL(pu.id, 0) AS parent_id,
    IF(IFNULL(pus.blocked, 0)=0 OR pus.block_whitelist=1 ,0, 1) AS parent_blocked,
    IFNULL(pwx.nick, pu.nickname) AS parent_nick,
    IFNULL(ru.id, 0) AS root_id,
    IF(IFNULL(rus.blocked, 0)=0 OR rus.block_whitelist=1 ,0, 1) AS root_blocked,
    IFNULL(rwx.nick, ru.nickname) AS root_nick,
	IF(IFNULL(us.blocked, 0)=0 OR us.block_whitelist=1, 0, 1) AS blocked,
	us.comments AS message,
	us.level AS _level,
    dev.id AS d_id,
    dev.platform AS d_platform,
    dev.idfa AS d_idfa,
    dev.imei AS d_imei,
    dev.mac AS d_mac,
    dev.device_name AS d_device_name,
    dev.system_version AS d_system_version,
    dev.os_version AS d_os_version,
    dev.language AS d_language,
    dev.model AS d_model,
    dev.timezone AS d_timezone,
    dev.country AS d_country,
    dev.is_emulator AS d_is_emulator,
    dev.is_jailbrojen AS d_is_jailbrojen,
    dev.is_tablet AS d_is_tablet,
    dev.points AS d_points,
    dev.total_ts AS d_total_ts,
    dev.tmp_ts AS d_tmp_ts,
    dev.consumed_ts AS d_consumed_ts,
    dev.lastping_at AS d_lastping_at,
    dev.inserted_at AS d_inserted_at,
    dev.updated_at AS d_updated_at,
	IF(
        (
            EXISTS (SELECT 1 FROM tmm.device_share_tasks AS dst WHERE dst.device_id=dev.id AND dst.inserted_at>=DATE_ADD(NOW(), INTERVAL 3 DAY) LIMIT 1) OR
            EXISTS (SELECT 1 FROM tmm.device_app_tasks AS dat WHERE dat.device_id=dev.id AND dat.inserted_at>=DATE_ADD(NOW(), INTERVAL 3 DAY) LIMIT 1) OR
            EXISTS (SELECT 1 FROM tmm.reading_logs AS rl WHERE rl.user_id=u.id AND (rl.inserted_at>=DATE_ADD(NOW(), INTERVAL 3 DAY) OR rl.updated_at>=DATE_ADD(NOW(), INTERVAL 3 DAY)) LIMIT 1) OR
            EXISTS (SELECT 1 FROM tmm.daily_bonus_logs AS dbl WHERE dbl.user_id=u.id AND dbl.updated_on>=DATE_ADD(NOW(), INTERVAL 3 DAY) LIMIT 1)),
        TRUE, FALSE
    ) AS _active,
	IF(COUNT(dev_app.app_id)>0, TRUE, FALSE) AS app_id,
	three_active.total AS total
FROM
	ucoin.users AS u
LEFT JOIN tmm.devices AS dev ON (dev.user_id=u.id)
LEFT JOIN tmm.device_apps AS dev_app ON (dev_app.device_id=dev.id)
LEFT JOIN tmm.user_settings AS us ON (us.user_id=u.id )
LEFT JOIN tmm.wx AS wx ON (wx.user_id =u.id)
LEFT JOIN tmm.invite_codes AS ic ON (ic.user_id=u.id)
LEFT JOIN ucoin.users AS pu ON (pu.id=ic.parent_id)
LEFT JOIN tmm.user_settings AS pus ON (pus.user_id=pu.id)
LEFT JOIN tmm.wx AS pwx ON (pwx.user_id=pu.id)
LEFT JOIN ucoin.users AS ru ON (ru.id=ic.root_id)
LEFT JOIN tmm.user_settings AS rus ON (rus.user_id=ru.id)
LEFT JOIN tmm.wx AS rwx ON (rwx.user_id=ru.id)
LEFT JOIN (
	SELECT SUM(cny) AS cny
	FROM tmm.withdraw_txs
	WHERE user_id=%d AND tx_status = 1
) AS uc ON 1 = 1
LEFT JOIN (
	SELECT SUM(cny) AS cny
	FROM tmm.point_withdraws
    WHERE user_id=%d AND verified!=-1
) AS point ON 1 = 1
LEFT JOIN (
	SELECT SUM(dgt.points) AS points
	FROM tmm.device_general_tasks AS dgt 
	INNER JOIN tmm.devices AS dev  ON (dev.id = dgt.device_id)
	WHERE  dev.user_id = %d  AND dgt.status = 1
) AS dgt ON 1 = 1
LEFT JOIN (
	SELECT
        COUNT(DISTINCT IF(inv.parent_id=%d, inv.user_id, NULL)) AS direct,
        COUNT(DISTINCT IF(inv.grand_id=%d, inv.user_id, NULL)) AS indirect,
        COUNT(DISTINCT IF(inv.parent_id=%d AND us.blocked=1 AND us.block_whitelist=0, inv.user_id, NULL)) AS direct_blocked,
        COUNT(DISTINCT IF(inv.grand_id=%d AND us.blocked=1 AND us.block_whitelist=0, inv.user_id, NULL)) AS indirect_blocked,
        COUNT(
            DISTINCT IF(
                    (EXISTS (SELECT 1 FROM tmm.device_share_tasks AS dst WHERE dst.device_id=d.id AND dst.inserted_at>=DATE_ADD(NOW(), INTERVAL 3 DAY) LIMIT 1) OR
                    EXISTS (SELECT 1 FROM tmm.device_app_tasks AS dat WHERE dat.device_id=d.id AND dat.inserted_at>=DATE_ADD(NOW(), INTERVAL 3 DAY) LIMIT 1) OR
                    EXISTS (SELECT 1 FROM tmm.reading_logs AS rl WHERE rl.user_id=d.user_id AND (rl.inserted_at>=DATE_ADD(NOW(), INTERVAL 3 DAY) OR rl.updated_at>=DATE_ADD(NOW(), INTERVAL 3 DAY)) LIMIT 1) OR
                    EXISTS (SELECT 1 FROM tmm.daily_bonus_logs AS dbl WHERE dbl.user_id=d.user_id AND dbl.updated_on>=DATE_ADD(NOW(), INTERVAL 3 DAY) LIMIT 1)),
            inv.user_id, NULL)
        ) AS active,
        COUNT(
            DISTINCT IF(
                u.created > DATE_SUB(NOW(),INTERVAL 3 DAY)
                AND
                    (EXISTS (SELECT 1 FROM tmm.device_share_tasks AS dst WHERE dst.device_id=d.id AND dst.inserted_at>=DATE_ADD(NOW(), INTERVAL 3 DAY) LIMIT 1) OR
                    EXISTS (SELECT 1 FROM tmm.device_app_tasks AS dat WHERE dat.device_id=d.id AND dat.inserted_at>=DATE_ADD(NOW(), INTERVAL 3 DAY) LIMIT 1) OR
                    EXISTS (SELECT 1 FROM tmm.reading_logs AS rl WHERE rl.user_id=d.user_id AND (rl.inserted_at>=DATE_ADD(NOW(), INTERVAL 3 DAY) OR rl.updated_at>=DATE_ADD(NOW(), INTERVAL 3 DAY)) LIMIT 1) OR
                    EXISTS (SELECT 1 FROM tmm.daily_bonus_logs AS dbl WHERE dbl.user_id=d.user_id AND dbl.updated_on>=DATE_ADD(NOW(), INTERVAL 3 DAY) LIMIT 1)),
            inv.user_id, NULL)
        ) AS invite_firend_active,
        COUNT(DISTINCT u.id) AS invite_By_Number
    FROM tmm.invite_codes  AS inv
    INNER JOIN ucoin.users AS u ON (u.id=inv.user_id)
    LEFT JOIN tmm.devices AS d ON (d.user_id=u.id)
    LEFT JOIN tmm.user_settings AS us ON (us.user_id=u.id)
    WHERE inv.parent_id=%d OR inv.grand_id=%d
) AS inv ON 1 = 1
LEFT JOIN (
	SELECT SUM(points) AS points FROM
    (
        SELECT d.user_id, SUM(dst.points) AS points
    	FROM tmm.device_share_tasks AS dst
    	INNER JOIN tmm.devices AS d ON (d.id = dst.device_id)
        WHERE d.user_id=%d
        UNION ALL
    	SELECT user_id, SUM(bonus) AS points
    	FROM tmm.invite_bonus
    	WHERE task_type!=0 AND user_id=%d
	) AS tmp
) AS sha ON 1 = 1
LEFT JOIN (
	SELECT SUM(IF(app.status=1, app.points, 0)) AS points
	FROM tmm.device_app_tasks AS app
	INNER JOIN tmm.devices AS d ON (d.id = app.device_id)
	WHERE d.user_id=%d AND app.status=1
) AS app ON 1 = 1
LEFT JOIN (
	SELECT SUM(bonus) AS inv_bonus
	FROM tmm.invite_bonus
	WHERE user_id=%d AND task_type=0
) AS bonus ON 1 = 1
LEFT JOIN (
	SELECT SUM(point) AS point
	FROM tmm.reading_logs
	WHERE user_id=%d
) AS reading ON 1 = 1
LEFT JOIN (
	SELECT COUNT(DISTINCT inv.user_id) AS total
	FROM tmm.invite_codes AS inv
	INNER JOIN ucoin.users AS u ON u.id = inv.user_id
	INNER JOIN tmm.devices AS d ON d.user_id = u.id
	WHERE inv.parent_id=%d OR inv.grand_id=%d
    AND
        (EXISTS (SELECT 1 FROM tmm.device_share_tasks AS dst WHERE dst.device_id=d.id AND dst.inserted_at<DATE_ADD(u.created, INTERVAL 3 DAY) LIMIT 1) OR
        EXISTS (SELECT 1 FROM tmm.device_app_tasks AS dat WHERE dat.device_id=d.id AND dat.inserted_at<DATE_ADD(u.created, INTERVAL 3 DAY) LIMIT 1) OR
        EXISTS (SELECT 1 FROM tmm.reading_logs AS rl WHERE rl.user_id=u.id AND rl.inserted_at<DATE_ADD(u.created, INTERVAL 3 DAY) LIMIT 1))
) AS three_active ON 1 = 1
WHERE
	u.id = %d
LIMIT 1
	`

	db := Service.Db
	rows, res, err := db.Query(query, id, id, id, id, id, id, id, id, id, id, id, id, id, id, id, id, id)
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
		IsActive:       row.Bool(res.Map(`_active`)),
		IsHaveAppId:    row.Bool(res.Map(`app_id`)),
		BlockedMessage: row.Str(res.Map(`message`)),
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
	user.PointByGeneralTask = int(row.Float(res.Map(`dgt_points`)))

	user.Id = row.Uint64(res.Map(`id`))
	user.Mobile = row.Str(res.Map(`mobile`))
	user.Nick = row.Str(res.Map(`nick`))
	user.Created = row.Str(res.Map(`created`))
	user.CountryCode = row.Uint(res.Map(`country_code`))
	user.Level.Id = row.Uint(res.Map(`_level`))
	user.Wallet = row.Str(res.Map(`addr`))
	user.Tmm = tmm.StringFixed(2)
	user.Point = point.StringFixed(0)
	user.Blocked = row.Int(res.Map(`blocked`))

	user.WxInsertedAt = row.Str(res.Map(`wx_inserted_at`))
	user.OpenId = row.Str(res.Map(`open_id`))
	user.WxUnionId = row.Str(res.Map(`union_id`))

	user.TotalMakePoint = user.PointByShare + user.PointByReading + user.PointByInvite + user.PointByDownLoadApp + user.PointByGeneralTask
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

	user.Parent = &admin.User{
		User: common.User{
			Id:   row.Uint64(res.Map(`parent_id`)),
			Nick: row.Str(res.Map(`parent_nick`)),
		},
		Blocked: row.Int(res.Map(`parent_blocked`)),
	}

	user.Root = &admin.User{
		User: common.User{
			Id:   row.Uint64(res.Map(`root_id`)),
			Nick: row.Str(res.Map(`root_nick`)),
		},
		Blocked: row.Int(res.Map(`root_blocked`)),
	}

	for _, row := range rows {
		deviceId := row.Str(res.Map("d_id"))
		if deviceId == "" {
			continue
		}
		device := &admin.Device{
			Device: common.Device{
				Id:         deviceId,
				Platform:   row.Str(res.Map("d_platform")),
				Idfa:       row.Str(res.Map("d_idfa")),
				Imei:       row.Str(res.Map("d_imei")),
				Mac:        row.Str(res.Map("d_mac")),
				Name:       row.Str(res.Map("d_device_name")),
				Model:      row.Str(res.Map("d_model")),
				IsEmulator: row.Bool(res.Map("d_is_emulator")),
				IsTablet:   row.Bool(res.Map("d_is_tablet")),
				Points:     decimal.NewFromFloat(row.Float(res.Map("d_points"))).Ceil(),
				TotalTs:    row.Uint64(res.Map("d_total_ts")),
				LastPingAt: row.Str(res.Map("d_lastping_at")),
				InsertedAt: row.Str(res.Map("d_inserted_at")),
				UpdatedAt:  row.Str(res.Map("d_updated_at")),
			},
			SystemVersion: row.Str(res.Map("d_system_version")),
			OsVersion:     row.Str(res.Map("d_os_version")),
			Language:      row.Str(res.Map("d_language")),
			Timezone:      row.Str(res.Map("d_timezone")),
			Country:       row.Str(res.Map("d_country")),
			IsJailbrojen:  row.Bool(res.Map("d_is_jailbrojen")),
			TmpTs:         row.Uint64(res.Map("d_tmp_ts")),
			ConsumedTs:    row.Float(res.Map("d_consumed_ts")),
		}

		user.DeviceList = append(user.DeviceList, device)
		if device.IsEmulator {
			user.IsHaveEmulatorDevices = true
		}
	}

	if user.OpenId != "" {
		rows, _, err := db.Query(`SELECT user_id FROM tmm.wx WHERE open_id='%s' AND user_id!=%d`, user.OpenId, id)
		if CheckErr(err, c) {
			return
		}
		for _, row := range rows {
			user.OtherAccount = append(user.OtherAccount, row.Str(0))
		}
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    user,
	})
}
