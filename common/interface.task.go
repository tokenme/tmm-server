package common

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/utils"
	"github.com/tokenme/tmm/utils/binary"
	"strings"
)

type ShareTask struct {
	Id            uint64          `json:"id"`
	Creator       uint64          `json:"creator,omitempty"`
	Title         string          `json:"title,omitempty"`
	Summary       string          `json:"summary,omitempty"`
	Link          string          `json:"link,omitempty"`
	ShareLink     string          `json:"share_link,omitempty"`
	VideoLink     string          `json:"video_link,omitempty"`
	Image         string          `json:"image,omitempty"`
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

func (this ShareTask) CookieKey() string {
	return fmt.Sprintf("share-task-%d", this.Id)
}

func (this ShareTask) IpKey(ip string) string {
	return fmt.Sprintf("share-task-%d-ip-%s", this.Id, ip)
}

func (this ShareTask) OpenidKey(openId string) string {
	return fmt.Sprintf("share-task-%d-openid-%s", this.Id, openId)
}

type AppTask struct {
	Id            uint64          `json:"id"`
	Creator       uint64          `json:"creator",omitempty`
	Name          string          `json:"name,omitempty"`
	Platform      Platform        `json:"platform,omitempty"`
	SchemeId      uint64          `json:"scheme_id,omitempty"`
	BundleId      string          `json:"bundle_id,omitempty"`
	StoreId       uint64          `json:"store_id,omitempty"`
	Icon          string          `json:"icon,omitempty"`
	Points        decimal.Decimal `json:"points,omitempty"`
	PointsLeft    decimal.Decimal `json:"points_left,omitempty"`
	Bonus         decimal.Decimal `json:"bonus,omitempty"`
	DownloadUrl   string          `json:"download_url,omitempty"`
	Downloads     uint            `json:"downloads,omitempty"`
	InsertedAt    string          `json:"inserted_at,omitempty"`
	UpdatedAt     string          `json:"updated_at,omitempty"`
	OnlineStatus  int8            `json:"online_status,omitempty"`
	InstallStatus int8            `json:"install_status,omitempty"`
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
