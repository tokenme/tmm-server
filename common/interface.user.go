package common

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	commonutils "github.com/tokenme/tmm/utils"
	tokenUtils "github.com/tokenme/tmm/utils/token"
	"math/big"
	"strings"
	"sync"
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
	Level           CreditLevel      `json:"level,omitempty"`
	Role            uint8            `json:"role,omitempty"`
	OpenId          string           `json:"openid,omitempty"`
	WxBinded        bool             `json:"wx_binded,omitempty"`
	DirectFriend    bool             `json:"direct_friend,omitempty"`
	Contribute      decimal.Decimal  `json:"contribute,omitempty"`
	Wechat          *Wechat          `json:"-"`
}

type Wechat struct {
	UnionId     string    `json:"union_id,omitempty"`
	OpenId      string    `json:"open_id,omitempty"`
	Nick        string    `json:"nick,omitempty"`
	Gender      uint      `json:"gender,omitempty"`
	Avatar      string    `json:"avatar,omitempty"`
	AccessToken string    `json:"access_token,omitempty"`
	Expires     time.Time `json:"expires,omitempty"`
}

type CreditLevel struct {
	Id            uint            `json:"id,omitemty"`
	Name          string          `json:"name,omitempty"`
	Enname        string          `json:"enname,omitempty"`
	Desc          string          `json:"desc,omitempty"`
	Endesc        string          `json:"endesc,omitempty"`
	Invites       uint            `json:"invites,omitempty"`
	TaskBonusRate decimal.Decimal `json:"task_bonus_rate"`
}

func (this User) IsAdmin() bool {
	return this.Role == 1
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
	key := commonutils.Md5(fmt.Sprintf("+%d%s", this.CountryCode, this.Mobile))
	return fmt.Sprintf("%suser/avatar/%s", cdn, key)
}

type UserBalance struct {
	Points decimal.Decimal `json:"points"`
	TMM    decimal.Decimal `json:"tmm"`
	Cash   decimal.Decimal `json:"cash"`
}

type UserWithdraw struct {
	Points decimal.Decimal `json:"points"`
	TMM    decimal.Decimal `json:"tmm"`
}

func (this User) Reset(ctx context.Context, service *Service, config Config, locker *sync.Mutex) (string, error) {
	db := service.Db
	_, _, err := db.Query(`UPDATE tmm.devices SET points=0, total_ts=0, consumed_ts=0 WHERE user_id=%d`, this.Id)
	if err != nil {
		return "", err
	}
	_, _, err = db.Query(`UPDATE tmm.devices AS d, tmm.user_devices AS ud SET d.points=0, d.total_ts=0, d.consumed_ts=0 WHERE d.id=ud.device_id AND ud.user_id=%d`, this.Id)
	if err != nil {
		return "", err
	}
	_, _, err = db.Query(`UPDATE tmm.user_settings SET level=0 WHERE user_id=%d`, this.Id)
	if err != nil {
		return "", err
	}
	agentPrivKey, err := commonutils.AddressDecrypt(config.TMMAgentWallet.Data, config.TMMAgentWallet.Salt, config.TMMAgentWallet.Key)
	if err != nil {
		return "", err
	}
	agentPubKey, err := eth.AddressFromHexPrivateKey(agentPrivKey)
	if err != nil {
		return "", err
	}

	token, err := utils.NewToken(config.TMMTokenAddress, service.Geth)
	if err != nil {
		return "", err
	}
	rows, _, err := db.Query(`SELECT wallet_addr FROM ucoin.users WHERE id=%d LIMIT 1`, this.Id)
	if err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return "", errors.New("not found user")
	}
	userWallet := rows[0].Str(0)
	tokenBalance, err := utils.TokenBalanceOf(token, userWallet)
	if err != nil {
		return "", err
	}
	if tokenBalance.Cmp(big.NewInt(0)) == 1 {
		locker.Lock()
		defer locker.Unlock()
		transactor := eth.TransactorAccount(agentPrivKey)
		nonce, err := eth.Nonce(ctx, service.Geth, service.Redis.Master, agentPubKey, config.Geth)
		if err != nil {
			return "", err
		}
		gasPrice, err := service.Geth.SuggestGasPrice(ctx)
		if err == nil && gasPrice.Cmp(eth.MinGas) == -1 {
			gasPrice = eth.MinGas
		} else {
			gasPrice = nil
		}
		transactorOpts := eth.TransactorOptions{
			Nonce:    nonce,
			GasPrice: gasPrice,
			GasLimit: 210000,
		}
		eth.TransactorUpdate(transactor, transactorOpts, ctx)
		tx, err := utils.TransferProxy(token, transactor, userWallet, agentPubKey, tokenBalance)
		if err != nil {
			return "", err
		}
		err = eth.NonceIncr(ctx, service.Geth, service.Redis.Master, agentPubKey, config.Geth)
		if err != nil {
			return "", err
		}
		return tx.Hash().Hex(), nil
	}
	return "", nil
}

func (this User) IsBlocked(service *Service) error {
	db := service.Db
	rows, _, err := db.Query(`SELECT 1 FROM tmm.user_settings WHERE user_id=%d AND blocked=1 AND block_whitelist=0 LIMIT 1`, this.Id)
	if err != nil {
		return err
	}
	if len(rows) > 0 {
		return errors.New("您的账户存在异常操作（异常行为包括但不限于：疑似恶意邀请用户行为、恶意刷分享文章、异常阅读文章等行为），不能执行提现及兑换操作。如有疑问请联系客服。客服微信搜索 \"jjxseven\"")
	}
	return nil
}

func (this User) BlockReason(service *Service) error {
	db := service.Db
	query := `SELECT
    COUNT( DISTINCT ib.from_user_id ) AS invites,
    SUM(IF(da.app_id IS NULL, 0, 1)) AS apps,
    SUM(ib.bonus) AS bonus,
    SUM(IFNULL(da.total_ts, 0)) AS ts
    FROM
        tmm.invite_bonus AS ib
        LEFT JOIN tmm.devices AS d ON (d.user_id=ib.from_user_id)
        LEFT JOIN tmm.device_apps AS da ON ( da.device_id = d.id )
WHERE ib.task_type=0 AND ib.user_id=%d
HAVING invites>=10 AND apps<invites/2
UNION
SELECT 0, 0, 0, 0
FROM tmm.wx AS ws
WHERE EXISTS (
    SELECT
        1
    FROM tmm.wx AS wx
    INNER JOIN tmm.user_settings AS us ON (us.user_id=wx.user_id)
    WHERE us.blocked=1 AND wx.open_id=ws.open_id AND wx.user_id!=ws.user_id LIMIT 1
) AND ws.user_id=%d`
	rows, _, err := db.Query(query, this.Id, this.Id)
	if err != nil {
		return err
	}
	if len(rows) > 0 {
		var reasons []string
		for _, row := range rows {
			invites := row.Uint(0)
			appsInUse := row.Uint(1)
			if invites > 0 {
				reasons = append(reasons, fmt.Sprintf("邀请人数 %d, 使用APP人数 %d", invites, appsInUse))
			} else {
				reasons = append(reasons, fmt.Sprintf("与屏蔽用户使用相同微信账号"))
			}
		}
		return errors.New(strings.Join(reasons, "\n"))
	}
	return nil
}
