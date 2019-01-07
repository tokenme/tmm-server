package common

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/utils"
	"strings"
	"time"
)

type Redpacket struct {
	Id           uint64          `json:"id"`
	Creator      uint64          `json:"creator"`
	Message      string          `json:"message"`
	Tmm          decimal.Decimal `json:"tmm"`
	Recipients   uint            `json:"recipients"`
	FundTx       string          `json:"fund_tx,omitempty"`
	FundTxStatus uint            `json:"fund_tx_status,omitempty"`
	ExpireTime   time.Time       `json:"expire_time,omitempty"`
	InsertedAt   time.Time       `json:"inserted_at,omitempty"`
	UpdatedAt    time.Time       `json:"updated_at,omitempty"`
}

func NewRedpacket(service *Service, creator uint64, tmm decimal.Decimal, recipients uint) (*Redpacket, error) {
	db := service.Db
	_, ret, err := db.Query(`INSERT INTO tmm.redpackets (creator, tmm, recipients) VALUES (%d, %s, %d)`, creator, tmm.String(), recipients)
	if err != nil {
		return nil, err
	}
	redpacket := &Redpacket{
		Id:         ret.InsertId(),
		Tmm:        tmm,
		Recipients: recipients,
	}
	redpacket.PrepareRecipients(service)
	return redpacket, nil
}

func (this Redpacket) GetLink(config Config, userId uint64) (string, error) {
	encryptedId, err := utils.EncryptUint64(this.Id, []byte(config.LinkSalt))
	if err != nil {
		return "", err
	}
	encryptedUserId, err := utils.EncryptUint64(userId, []byte(config.LinkSalt))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", config.RedpacketUrl, encryptedId, encryptedUserId), nil
}

func DecryptRedpacketLink(encryptedId string, encryptedUserId string, config Config) (id uint64, userId uint64, err error) {
	id, err = utils.DecryptUint64(encryptedId, []byte(config.LinkSalt))
	if err != nil {
		return
	}
	userId, err = utils.DecryptUint64(encryptedUserId, []byte(config.LinkSalt))
	if err != nil {
		return
	}
	return id, userId, nil
}

func (this Redpacket) PrepareRecipients(service *Service) error {
	db := service.Db
	rows, _, err := db.Query(`SELECT 1 FROM tmm.redpacket_recipients WHERE redpacket_id=%d LIMIT 1`, this.Id)
	if err != nil {
		return err
	}
	if len(rows) > 0 {
		return nil
	}
	packs := GenerateRecipients(this.Tmm, uint64(this.Recipients), 1)
	var val []string
	for _, p := range packs {
		val = append(val, fmt.Sprintf("(%d, %s)", this.Id, p.String()))
	}
	_, _, err = db.Query(`INSERT INTO tmm.redpacket_recipients (redpacket_id, tmm) VALUES %s`, strings.Join(val, ","))
	if err != nil {
		return err
	}
	return nil
}

func GenerateRecipients(tokens decimal.Decimal, num uint64, min uint64) []decimal.Decimal {
	var (
		i        uint64 = 1
		resp     []decimal.Decimal
		decimals = decimal.New(1, 4)
		total    = uint64(tokens.Mul(decimals).IntPart())
	)
	for i < num {
		safeTotal := (total - (num-i)*min) / (num - i) //随机安全上限
		money := utils.RangeRandUint64(min, safeTotal)
		resp = append(resp, decimal.New(int64(money), 0).Div(decimals))
		total = total - money
		i += 1
	}
	resp = append(resp, decimal.New(int64(total), 0).Div(decimals))
	//spew.Dump(resp)
	return resp
}

type RedpacketRecipient struct {
	UserId  uint64          `json:"-"`
	Nick    string          `json:"nick"`
	Avatar  string          `json:"avatar"`
	UnionId string          `json:"-"`
	Tmm     decimal.Decimal `json:"tmm"`
	Cash    decimal.Decimal `json:"cash"`
}
