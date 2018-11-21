package common

import (
	"fmt"
	"github.com/tokenme/tmm/utils"
	tokenUtils "github.com/tokenme/tmm/utils/token"
	"strings"
	"time"
)

type User struct {
	Id              uint64           `json:"id,omitempty"`
	CountryCode     uint             `json:"country_code,omitempty"`
	Mobile          string           `json:"mobile,omitempty"`
	Nick            string           `json:"nick,omitempty"`
	Name            string           `json:"realname,omitempty"`
	ShowName        string           `json:"showname,omitempty"`
	Avatar          string           `json:"avatar,omitempty"`
	Salt            string           `json:"-"`
	Password        string           `json:"-"`
	Wallet          string           `json:"wallet"`
	WalletPK        string           `json:"-"`
	InviteCode      tokenUtils.Token `json:"invite_code,omitempty"`
	InviterCode     tokenUtils.Token `json:"inviter_code,omitempty"`
	CanPay          uint             `json:"can_pay,omitempty"`
	ExchangeEnabled bool             `json:"exchange_enabled,omitempty"`
	WxBinded        bool             `json:"wx_binded,omitempty"`
	Wechat          *Wechat          `json:"-"`
}

type Wechat struct {
	UnionId     string    `json:"union_id,omitempty"`
	Nick        string    `json:"nick,omitempty"`
	Gender      uint      `json:"gender,omitempty"`
	Avatar      string    `json:"avatar,omitempty"`
	AccessToken string    `json:"access_token,omitempty"`
	Expires     time.Time `json:"expires,omitempty"`
}

func (this User) GetShowName() string {
	if this.Wechat != nil && this.Wechat.Nick != "" {
		return this.Wechat.Nick
	}
	/*
		if this.Name != "" {
			return this.Name
		}
		if this.Nick != "" {
			return this.Nick
		}
	*/
	return fmt.Sprintf("+%d%s", this.CountryCode, this.Mobile)
}

func (this User) GetAvatar(cdn string) string {
	if this.Wechat != nil && this.Wechat.Avatar != "" {
		if strings.HasPrefix(this.Wechat.Avatar, "http://") {
			return strings.Replace(this.Wechat.Avatar, "http://", "https://", 1)
		}
		return this.Wechat.Avatar
	}
	if this.Avatar != "" {
		return this.Avatar
	}
	key := utils.Md5(fmt.Sprintf("+%d%s", this.CountryCode, this.Mobile))
	return fmt.Sprintf("%suser/avatar/%s", cdn, key)
}
