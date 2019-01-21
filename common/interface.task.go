package common

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/utils"
	"github.com/tokenme/tmm/utils/binary"
	"net/url"
	"strings"
)

type ShareTaskTrackSource = uint

const (
	TrackFromUCoin   ShareTaskTrackSource = 0
	TrackFromWechat  ShareTaskTrackSource = 1
	TrackFromUnknown ShareTaskTrackSource = 2
)

type GeneralTask struct {
	Id                 uint64          `json:"id"`
	Creator            uint64          `json:"creator,omitempty"`
	Title              string          `json:"title,omitempty"`
	Summary            string          `json:"summary,omitempty"`
	Image              string          `json:"image,omitempty"`
	Points             decimal.Decimal `json:"points,omitempty"`
	PointsLeft         decimal.Decimal `json:"points_left,omitempty"`
	Bonus              decimal.Decimal `json:"bonus,omitempty"`
	InsertedAt         string          `json:"inserted_at,omitempty"`
	UpdatedAt          string          `json:"updated_at,omitempty"`
	OnlineStatus       int8            `json:"online_status,omitempty"`
	Details            string          `json:"details,omitempty"`
	CertificateStatus  int8            `json:"certificate_status,omitempty"`
	CertificateImages  string          `json:"certificate_images,omitempty"`
	CertificateComment string          `json:"certificate_comment,omitempty"`
}

type ShareTask struct {
	Id            uint64          `json:"id"`
	Creator       uint64          `json:"creator,omitempty"`
	Title         string          `json:"title,omitempty"`
	Summary       string          `json:"summary,omitempty"`
	Link          string          `json:"link,omitempty"`
	ShareLink     string          `json:"share_link,omitempty"`
	VideoLink     string          `json:"video_link,omitempty"`
	Image         string          `json:"image,omitempty"`
	Images        []string        `json:"images,omitempty"`
	Points        decimal.Decimal `json:"points,omitempty"`
	PointsLeft    decimal.Decimal `json:"points_left,omitempty"`
	Bonus         decimal.Decimal `json:"bonus,omitempty"`
	MaxViewers    uint            `json:"max_viewers,omitempty"`
	Viewers       uint            `json:"viewers,omitempty"`
	InsertedAt    string          `json:"inserted_at,omitempty"`
	UpdatedAt     string          `json:"updated_at,omitempty"`
	OnlineStatus  int8            `json:"online_status,omitempty"`
	IsVideo       uint8           `json:"is_video,omitempty"`
	IsTask        bool            `json:"is_task,omitempty"`
	InIframe      bool            `json:"-"`
	TimelineOnly  bool            `json:"-"`
	ShowBonusHint bool            `json:"show_bonus_hint,omitempty"`
	Creative      *Creative       `json:"creative,omitempty"`
	Cid           []int           `json:"cid,omitempty"`
	TotalReadUser int             `json:"total_read_user,omitempty"`
	ReadDuration  int             `json:"read_duration,omitempty"`
}

func (this ShareTask) ShouldUseIframe() bool {
	return strings.HasPrefix(this.Link, "https://static.tianxi100.com") || strings.HasPrefix(this.Link, "https://tmm.tokenmama.io")
}

func (this ShareTask) TrackLink(link string, userId uint64, config Config) (string, error) {
	encrypted, err := utils.EncryptUint64(this.Id, []byte(config.LinkSalt))
	if err != nil {
		return "", err
	}
	var encryptedUserId string
	if userId > 0 {
		encryptedUserId, err = utils.EncryptUint64(userId, []byte(config.LinkSalt))
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s/%s?uid=%s&url=%s", config.ShareTrackUrl, encrypted, url.QueryEscape(encryptedUserId), url.QueryEscape(link)), nil
}

func (this ShareTask) GetShareLink(deviceId string, config Config) (string, error) {
	encrypted, err := utils.EncryptUint64(this.Id, []byte(config.LinkSalt))
	if err != nil {
		return "", err
	}
	encryptedDeviceId, err := utils.AESEncrypt([]byte(config.LinkSalt), deviceId)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", config.ShareUrl, encrypted, encryptedDeviceId), nil
}

func (this ShareTask) GetShareImpLink(deviceId string, config Config) (string, error) {
	encrypted, err := utils.EncryptUint64(this.Id, []byte(config.LinkSalt))
	if err != nil {
		return "", err
	}
	encryptedDeviceId, err := utils.AESEncrypt([]byte(config.LinkSalt), deviceId)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", config.ShareImpUrl, encrypted, encryptedDeviceId), nil
}

func DecryptShareTaskLink(encryptedTaskId string, encryptedDeviceId string, config Config) (taskId uint64, deviceId string, err error) {
	taskId, err = utils.DecryptUint64(encryptedTaskId, []byte(config.LinkSalt))
	if err != nil {
		return
	}
	deviceId, err = utils.AESDecrypt([]byte(config.LinkSalt), encryptedDeviceId)
	if err != nil {
		return
	}
	return taskId, deviceId, nil
}

func (this ShareTask) Imp(deviceId string, bonus decimal.Decimal, service *Service, config Config) error {
	pointsPerTs, err := GetPointsPerTs(service)
	if err != nil {
		return err
	}
	db := service.Db
	var bonusRate float64 = 1
	{ // Getting bonus rate
		rows, _, _ := db.Query(`SELECT ul.task_bonus_rate FROM tmm.user_settings AS us INNER JOIN tmm.user_levels AS ul ON (ul.id=us.level) INNER JOIN tmm.devices AS d ON (d.user_id=us.user_id) WHERE d.id='%s' LIMIT 1`, db.Escape(deviceId))
		if len(rows) > 0 {
			bonusRate = rows[0].ForceFloat(0) / 100
		}
	}

	{ // Update device points
		query := `UPDATE tmm.devices AS d, tmm.device_share_tasks AS dst, tmm.share_tasks AS st
            SET
                d.points = d.points + IF(st.points_left > st.bonus, st.bonus, st.points_left) * %.2f,
                d.total_ts = d.total_ts + CEIL(IF(st.points_left > st.bonus, st.bonus, st.points_left) / %s),
                dst.points = dst.points + IF(st.points_left > st.bonus, st.bonus, st.points_left) * %.2f,
                dst.viewers = dst.viewers + 1,
                st.points_left = IF(st.points_left > st.bonus, st.points_left - st.bonus, 0),
                st.viewers = st.viewers + 1
            WHERE
                d.id='%s'
            AND dst.device_id = d.id
            AND dst.task_id = %d
            AND st.id = dst.task_id`
		_, _, err = db.Query(query, bonusRate, pointsPerTs.String(), bonusRate, db.Escape(deviceId), this.Id)
		if err != nil {
			return err
		}
	}

	{ // Update inviter bonus
		query := `SELECT t.id, t.inviter_id, t.user_id, t.is_grand FROM
(SELECT id, inviter_id, user_id, false AS is_grand FROM
    (SELECT
    d.id,
    ic.parent_id AS inviter_id,
    ic.user_id
    FROM tmm.invite_codes AS ic
    INNER JOIN tmm.wx AS wx ON (wx.user_id=ic.user_id)
    LEFT JOIN tmm.devices AS d ON (d.user_id=ic.parent_id)
    LEFT JOIN tmm.devices AS d2 ON (d2.user_id=ic.user_id)
    WHERE d2.id='%s' AND ic.parent_id > 0
    ORDER BY d.lastping_at DESC LIMIT 1) AS t1
    UNION
    SELECT id, inviter_id, user_id, true AS is_grand FROM
    (SELECT
    d.id,
    ic.grand_id AS inviter_id,
    ic.user_id
    FROM tmm.invite_codes AS ic
    INNER JOIN tmm.wx AS wx ON (wx.user_id=ic.user_id)
    LEFT JOIN tmm.devices AS d ON (d.user_id=ic.grand_id)
    LEFT JOIN tmm.devices AS d2 ON (d2.user_id=ic.user_id)
    WHERE d2.id='%s' AND ic.grand_id > 0
    ORDER BY d.lastping_at DESC LIMIT 1) AS t2
) AS t
LEFT JOIN tmm.wx AS wx ON (wx.user_id=t.inviter_id)
LEFT JOIN tmm.wx AS wx2 ON (wx2.user_id=t.user_id)
WHERE ISNULL(wx.open_id) OR ISNULL(wx2.open_id) OR wx.open_id!=wx2.open_id`
		rows, _, err := db.Query(query, db.Escape(deviceId), db.Escape(deviceId))
		if err != nil {
			return err
		}
		for _, row := range rows {
			dId := row.Str(0)
			inviterId := row.Uint64(1)
			userId := row.Uint64(2)
			isGrand := row.Bool(3)
			var inviterBonus decimal.Decimal
			if isGrand {
				inviterBonus = bonus.Mul(decimal.NewFromFloat(bonusRate * config.InviteBonusRate * config.InviteBonusRate))
			} else {
				inviterBonus = bonus.Mul(decimal.NewFromFloat(bonusRate * config.InviteBonusRate))
			}
			_, ret, err := db.Query(`UPDATE tmm.devices SET points = points + %s, total_ts = total_ts + %d WHERE id='%s'`, inviterBonus.String(), inviterBonus.Div(pointsPerTs).IntPart(), db.Escape(dId))
			if err != nil {
				continue
			}
			if ret.AffectedRows() > 0 {
				_, _, err = db.Query(`INSERT INTO tmm.invite_bonus (user_id, from_user_id, bonus, task_type, task_id) VALUES (%d, %d, %s, 1, %d)`, inviterId, userId, inviterBonus.String(), this.Id)
				if err != nil {
					continue
				}
			}
		}
	}
	return nil
}

func (this ShareTask) CookieKey() string {
	return fmt.Sprintf("share-task-%d", this.Id)
}

func (this ShareTask) IpKey(ip string) string {
	return fmt.Sprintf("share-task-%d-ip-%s", this.Id, ip)
}

func (this ShareTask) OpenidKey(openId string) string {
	return fmt.Sprintf("share-task-%d-openid-%s", this.Id, openId)
}

func (this ShareTask) UidKey(uid uint64) string {
	return fmt.Sprintf("st-%d-uid-%s", this.Id, uid)
}

func (this ShareTask) UserRateLimitSecondKey(uid uint64) string {
	return fmt.Sprintf("st-%d-rate-sec", uid)
}

func (this ShareTask) UserRateLimitMinuteKey(uid uint64) string {
	return fmt.Sprintf("st-%d-rate-min", uid)
}

func (this ShareTask) UserRateLimitSecondBlockKey(uid uint64) string {
	return fmt.Sprintf("st-%d-rb-sec", uid)
}

func (this ShareTask) UserRateLimitBlockKey(uid uint64) string {
	return fmt.Sprintf("st-%d-rb-min", uid)
}

func WxCodeKey(code string) string {
	return fmt.Sprintf("wxcode-%s", code)
}

type AppTask struct {
	Id                 uint64          `json:"id"`
	Creator            uint64          `json:"creator",omitempty`
	Name               string          `json:"name,omitempty"`
	Platform           Platform        `json:"platform,omitempty"`
	SchemeId           uint64          `json:"scheme_id,omitempty"`
	BundleId           string          `json:"bundle_id,omitempty"`
	StoreId            uint64          `json:"store_id,omitempty"`
	Icon               string          `json:"icon,omitempty"`
	Points             decimal.Decimal `json:"points,omitempty"`
	PointsLeft         decimal.Decimal `json:"points_left,omitempty"`
	Bonus              decimal.Decimal `json:"bonus,omitempty"`
	DownloadUrl        string          `json:"download_url,omitempty"`
	Downloads          uint            `json:"downloads,omitempty"`
	InsertedAt         string          `json:"inserted_at,omitempty"`
	UpdatedAt          string          `json:"updated_at,omitempty"`
	Size               uint            `json:"size,omitempty"`
	OnlineStatus       int8            `json:"online_status,omitempty"`
	InstallStatus      int8            `json:"install_status,omitempty"`
	Details            string          `json:"details,omitempty"`
	CertificateStatus  int8            `json:"certificate_status,omitempty"`
	CertificateImages  string          `json:"certificate_images,omitempty"`
	CertificateComment string          `json:"certificate_comment,omitempty"`
}

func (this AppTask) Install(user User, deviceId string, service *Service, config Config) (bonus decimal.Decimal, err error) {
	db := service.Db
	{ // Check App installed
		rows, _, err := db.Query(`SELECT 1 FROM tmm.device_app_tasks WHERE device_id='%s' AND task_id=%d AND status=-1 LIMIT 1`, db.Escape(deviceId), this.Id)
		if err != nil {
			return bonus, err
		}
		if len(rows) > 0 {
			return bonus, errors.New("You have been finished the task")
		}
	}

	pointsPerTs, err := GetPointsPerTs(service)
	if err != nil {
		return bonus, err
	}
	var bonusRate float64 = 1
	{ // Getting bonus rate
		rows, _, _ := db.Query(`SELECT ul.task_bonus_rate FROM tmm.user_settings AS us INNER JOIN tmm.user_levels AS ul ON (ul.id=us.level) INNER JOIN tmm.devices AS d ON (d.user_id=us.user_id) WHERE d.id='%s' LIMIT 1`, db.Escape(deviceId))
		if len(rows) > 0 {
			bonusRate = rows[0].ForceFloat(0) / 100
		}
	}
	{ // Update device bonus
		query := `UPDATE tmm.devices AS d, tmm.device_app_tasks AS dat, tmm.app_tasks AS appt
SET d.points = d.points + IF(appt.points_left > appt.bonus, appt.bonus, appt.points_left) * %.2f,
    d.total_ts = d.total_ts + CEIL(IF(appt.points_left > appt.bonus, appt.bonus, appt.points_left) / %s),
    appt.points_left = IF(appt.points_left > appt.bonus, appt.points_left - appt.bonus, 0),
    appt.downloads = appt.downloads + 1,
    dat.points = IF(appt.points_left > appt.bonus, appt.bonus, appt.points_left) * %.2f,
    dat.status = 1
WHERE
    d.id = '%s'
    AND appt.id = dat.task_id
    AND dat.device_id = d.id
    AND dat.task_id = %d
    AND dat.status != 1`
		_, _, err = db.Query(query, bonusRate, pointsPerTs.String(), bonusRate, db.Escape(deviceId), this.Id)
		if err != nil {
			return bonus, err
		}
	}
	{ // Check device bonus
		rows, _, err := db.Query(`SELECT points FROM tmm.device_app_tasks WHERE device_id='%s' AND task_id=%d LIMIT 1`, db.Escape(deviceId), this.Id)
		if err != nil {
			return bonus, nil
		}
		if len(rows) == 0 {
			return bonus, errors.New("Task not finished")
		}
		bonus, _ = decimal.NewFromString(rows[0].Str(0))
	}
	{ // Give bonus to inviters
		query := `SELECT t.id, t.inviter_id, t.is_grand FROM
(SELECT id, inviter_id,user_id, false AS is_grand FROM
    (SELECT
    d.id,
    d.user_id AS inviter_id,
	ic.user_id
    FROM tmm.devices AS d
    LEFT JOIN tmm.invite_codes AS ic ON (ic.parent_id=d.user_id)
    INNER JOIN tmm.wx AS wx ON (wx.user_id=ic.user_id)
    WHERE ic.user_id = %d AND d.user_id > 0
    AND NOT EXISTS (SELECT 1 FROM tmm.invite_bonus AS ib WHERE ib.user_id=d.user_id AND ib.from_user_id=ic.user_id AND task_type=2 AND task_id=%d LIMIT 1)
    ORDER BY d.lastping_at DESC LIMIT 1) AS t1
    UNION
    SELECT id,inviter_id, user_id, true AS is_grand FROM
    (SELECT
    d.id,
    d.user_id AS inviter_id,
	ic.user_id
    FROM tmm.devices AS d
    LEFT JOIN tmm.invite_codes AS ic ON (ic.grand_id=d.user_id)
    INNER JOIN tmm.wx AS wx ON (wx.user_id=ic.user_id)
    WHERE ic.user_id = %d AND d.user_id > 0
    AND NOT EXISTS (SELECT 1 FROM tmm.invite_bonus AS ib WHERE ib.user_id=d.user_id AND ib.from_user_id=ic.user_id AND task_type=2 AND task_id=%d LIMIT 1)
    ORDER BY d.lastping_at DESC LIMIT 1) AS t2
) AS t
LEFT JOIN tmm.wx AS wx ON (wx.user_id=t.inviter_id)
LEFT JOIN tmm.wx AS wx2 ON (wx2.user_id=t.user_id)
WHERE ISNULL(wx.open_id) OR ISNULL(wx2.open_id) OR wx.open_id!=wx2.open_id`
		rows, _, err := db.Query(query, user.Id, this.Id, user.Id, this.Id)
		if err != nil {
			return bonus, err
		}
		for _, row := range rows {
			dId := row.Str(0)
			inviterId := row.Uint64(1)
			isGrand := row.Bool(2)
			var inviterBonus decimal.Decimal
			if isGrand {
				inviterBonus = bonus.Mul(decimal.NewFromFloat(bonusRate * config.InviteBonusRate * config.InviteBonusRate))
			} else {
				inviterBonus = bonus.Mul(decimal.NewFromFloat(bonusRate * config.InviteBonusRate))
			}
			_, ret, err := db.Query(`UPDATE tmm.devices SET points = points + %s, total_ts = total_ts + %d WHERE id='%s'`, inviterBonus.String(), inviterBonus.Div(pointsPerTs).IntPart(), db.Escape(dId))
			if err != nil {
				continue
			}
			if ret.AffectedRows() > 0 {
				_, _, err = db.Query(`INSERT INTO tmm.invite_bonus (user_id, from_user_id, bonus, task_type, task_id) VALUES (%d, %d, %s, 2, %d)`, inviterId, user.Id, inviterBonus.String(), this.Id)
				if err != nil {
					continue
				}
			}
		}
	}
	return bonus, nil
}

func (this AppTask) Uninstall(user User, deviceId string, service *Service, config Config) (bonus decimal.Decimal, err error) {
	db := service.Db
	{ // Check device bonus
		rows, _, err := db.Query(`SELECT points FROM tmm.device_app_tasks WHERE device_id='%s' AND task_id=%d LIMIT 1`, db.Escape(deviceId), this.Id)
		if err != nil {
			return bonus, err
		}
		if len(rows) > 0 {
			return bonus, errors.New("You have been finished the task")
		}
		bonus, _ = decimal.NewFromString(rows[0].Str(0))
	}

	pointsPerTs, err := GetPointsPerTs(service)
	if err != nil {
		return bonus, err
	}
	{ // Update device bonus
		query := `UPDATE tmm.devices AS d, tmm.device_app_tasks AS dat, tmm.app_tasks AS appt
SET d.points = IF(d.points > dat.points, d.points - dat.points, 0),
    d.consumed_ts = d.consumed_ts + CEIL(IF(d.points > dat.points, d.points - dat.points, 0) / %s),
    appt.points_left = appt.points_left + IF(d.points > dat.points, dat.points, 0),
    appt.downloads = appt.downloads - IF(d.points > dat.points, 1, 0),
    dat.status = -1
WHERE
    d.id = '%s'
    AND appt.id = dat.task_id
    AND dat.device_id = d.id
    AND dat.task_id = %d
    AND dat.status = 1`
		_, _, err = db.Query(query, pointsPerTs, db.Escape(deviceId), this.Id)
		if err != nil {
			return bonus, err
		}
	}
	{ // Update invite bonus
		query := `SELECT id, user_id FROM
(SELECT
d.id,
d.user_id
FROM tmm.devices AS d
LEFT JOIN tmm.invite_codes AS ic ON (ic.parent_id=d.user_id)
WHERE ic.user_id = %d AND d.user_id > 0
ORDER BY d.points DESC LIMIT 1) AS t1
UNION
SELECT id, user_id FROM
(SELECT
d.id,
d.user_id
FROM tmm.devices AS d
LEFT JOIN tmm.invite_codes AS ic ON (ic.grand_id=d.user_id)
WHERE ic.user_id = %d AND d.user_id > 0
ORDER BY d.points DESC LIMIT 1) AS t2`
		rows, _, err := db.Query(query, user.Id, user.Id)
		if err != nil {
			return bonus, err
		}
		var (
			deviceIds  []string
			inviterIds []string
		)
		for _, row := range rows {
			deviceIds = append(deviceIds, fmt.Sprintf("'%s'", db.Escape(row.Str(0))))
			inviterIds = append(inviterIds, fmt.Sprintf("%d", row.Uint64(1)))
		}
		if len(deviceIds) > 0 {
			_, ret, err := db.Query(`UPDATE tmm.devices AS d, tmm.invite_bonus AS ib SET d.points = IF(d.points>ib.bonus, d.points-ib.bonus, 0) WHERE ib.user_id IN (%s) AND d.id IN (%s) AND ib.task_type=2 AND ib.task_id=%d`, strings.Join(inviterIds, ","), strings.Join(deviceIds, ","), this.Id)
			if err != nil {
				return bonus, err
			}
			if ret.AffectedRows() > 0 {
				_, _, err = db.Query(`DELETE FROM tmm.invite_bonus WHERE user_id IN (%s) AND from_user_id=%d AND task_type=2 AND task_id=%d`, strings.Join(inviterIds, ","), user.Id, this.Id)
				if err != nil {
					return bonus, err
				}
			}
		}
	}
	return bonus, nil
}

type TaskType = uint

const (
	AppTaskType   TaskType = 1
	ShareTaskType TaskType = 2
)

type TaskRecord struct {
	Type      TaskType        `json:"type"`
	Title     string          `json:"title"`
	Points    decimal.Decimal `json:"points"`
	Image     string          `json:"image,omitempty"`
	Viewers   uint            `json:"viewers,omitempty"`
	UpdatedAt string          `json:"updated_at,omitempty"`
}

type CryptOpenid struct {
	Openid string `json:"openid"`
	Ts     int64  `json:"ts"`
}

func (this CryptOpenid) Encode(key []byte) (string, error) {
	enc := binary.NewEncoder()
	enc.Encode(this)
	return utils.AESEncryptBytes(key, enc.Buffer())
}

func DecodeCryptOpenid(key []byte, cryptoText string) (openid CryptOpenid, err error) {
	data, err := utils.AESDecryptBytes(key, cryptoText)
	if err != nil {
		return openid, err
	}
	dec := binary.NewDecoder()
	dec.SetBuffer(data)
	dec.Decode(&openid)
	return openid, nil
}
