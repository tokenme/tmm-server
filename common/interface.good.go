package common

import (
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/utils"
	"github.com/tokenme/tmm/utils/binary"
)

type GoodInvestSummary struct {
	Invest decimal.Decimal `json:"invest"`
	Bonus  decimal.Decimal `json:"bonus"`
	Income decimal.Decimal `json:"income"`
}

type GoodInvest struct {
	Nick         string          `json:"user_name,omitempty"`
	Points       decimal.Decimal `json:"points"`
	Income       decimal.Decimal `json:"income,omitempty"`
	Avatar       string          `json:"avatar,omitempty"`
	RedeemStatus uint            `json:"redeem_status,omitempty"`
	InsertedAt   string          `json:"inserted_at,omitempty"`
	GoodId       uint64          `json:"good_id,omitempty"`
	GoodName     string          `json:"good_name,omitempty"`
	GoodPic      string          `json:"good_pic,omitempty"`
}

type GoodTx struct {
	OrderId   uint64 `json:"oid"`
	Uid       uint64 `json:"uid"`
	GoodId    uint64 `json:"good_id"`
	Amount    uint   `json:"amount"`
	Income    uint   `json:"income"`
	CreatedAt string `json:"created_at"`
}

func (this GoodTx) Encode(key []byte) (string, error) {
	enc := binary.NewEncoder()
	enc.Encode(this)
	return utils.AESEncryptBytes(key, enc.Buffer())
}

func DecodeGoodTx(key []byte, cryptoText string) (tx GoodTx, err error) {
	data, err := utils.AESDecryptBytes(key, cryptoText)
	if err != nil {
		return tx, err
	}
	dec := binary.NewDecoder()
	dec.SetBuffer(data)
	dec.Decode(&tx)
	return tx, nil
}

type GoodTxs []GoodTx

func (this GoodTxs) Encode(key []byte) (string, error) {
	enc := binary.NewEncoder()
	enc.Encode(this)
	return utils.AESEncryptBytes(key, enc.Buffer())
}

func DecodeGoodTxs(key []byte, cryptoText string) (txs GoodTxs, err error) {
	data, err := utils.AESDecryptBytes(key, cryptoText)
	if err != nil {
		return txs, err
	}
	dec := binary.NewDecoder()
	dec.SetBuffer(data)
	dec.Decode(&txs)
	return txs, nil
}
