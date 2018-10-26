package common

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/utils"
	"strings"
)

type ShareTask struct {
	Id            uint64          `json:"id"`
	Creator       uint64          `json:"creator",omitempty`
	Title         string          `json:"title"`
	Summary       string          `json:"summary"`
	Link          string          `json:"link"`
	ShareLink     string          `json:"share_link"`
	Image         string          `json:"image,omitempty"`
	Points        decimal.Decimal `json:"points,omitempty"`
	PointsLeft    decimal.Decimal `json:"points_left,omitempty"`
	Bonus         decimal.Decimal `json:"bonus,omitempty"`
	MaxViewers    uint            `json:"max_viewers,omitempty"`
	Viewers       uint            `json:"viewers,omitempty"`
	InsertedAt    string          `json:"inserted_at,omitempty"`
	UpdatedAt     string          `json:"updated_at,omitempty"`
	OnlineStatus  int8            `json:"online_status,omitempty"`
	InIframe      bool            `json:"-"`
	ShowBonusHint bool            `json:"show_bonus_hint,omitempty"`
}

func (this ShareTask) ShouldUseIframe() bool {
	return !strings.HasPrefix(this.Link, "https://mp.weixin.qq.com") && !strings.HasPrefix(this.Link, "https://www.taobao.com")
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
