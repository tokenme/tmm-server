package common

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/utils"
)

type ShareTask struct {
	Id         uint64          `json:"id"`
	Title      string          `json:"title"`
	Summary    string          `json:"summary"`
	Link       string          `json:"link"`
	ShareLink  string          `json:"share_link"`
	Image      string          `json:"image,omitempty"`
	Points     decimal.Decimal `json:"points,omitempty"`
	PointsLeft decimal.Decimal `json:"points_left,omitempty"`
	Bonus      decimal.Decimal `json:"bonus,omitempty"`
	MaxViewers uint            `json:"max_viewers,omitempty"`
	Viewers    uint            `json:"viewers,omitempty"`
	InsertedAt string          `json:"inserted_at,omitempty"`
	UpdatedAt  string          `json:"updated_at,omitempty"`
}

type ShareTaskProto struct {
	Id   uint64 `json:"id"`
	Link string `json:"link"`
}

func (this ShareTask) GetShareLink(userId uint64, config Config) (string, error) {
	encrypted, err := utils.EncryptUint64(this.Id, []byte(config.LinkSalt))
	if err != nil {
		return "", err
	}
	encryptedUserId, err := utils.EncryptUint64(userId, []byte(config.LinkSalt))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", config.ShareUrl, encrypted, encryptedUserId), nil
}

func (this ShareTask) CookieKey() string {
	return fmt.Sprintf("share-task-%d", this.Id)
}
