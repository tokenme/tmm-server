package common

import (
	"fmt"
	"github.com/tokenme/tmm/utils"
	tokenUtils "github.com/tokenme/tmm/utils/token"
)

type User struct {
	Id          uint64           `json:"id,omitempty"`
	CountryCode uint             `json:"country_code,omitempty"`
	Mobile      string           `json:"mobile,omitempty"`
	Nick        string           `json:"nick,omitempty"`
	Name        string           `json:"realname,omitempty"`
	ShowName    string           `json:"showname,omitempty"`
	Avatar      string           `json:"avatar,omitempty"`
	Salt        string           `json:"-"`
	Password    string           `json:"-"`
	Wallet      string           `json:"wallet"`
	WalletPK    string           `json:"-"`
	InviteCode  tokenUtils.Token `json:"invite_code,omitempty"`
	InviterCode tokenUtils.Token `json:"inviter_code,omitempty"`
	CanPay      uint             `json:"can_pay,omitempty"`
}

func (this User) GetShowName() string {
	if this.Name != "" {
		return this.Name
	}
	if this.Nick != "" {
		return this.Nick
	}
	return fmt.Sprintf("+%d%s", this.CountryCode, this.Mobile)
}

func (this User) GetAvatar(cdn string) string {
	if this.Avatar != "" {
		return this.Avatar
	}
	key := utils.Md5(fmt.Sprintf("+%d%s", this.CountryCode, this.Mobile))
	return fmt.Sprintf("%suser/avatar/%s", cdn, key)
}
