package common

import (
	"github.com/tokenme/tmm/utils"
	"github.com/tokenme/tmm/utils/binary"
)

type Adgroup struct {
	Id           uint64      `json:"id"`
	OnlineStatus uint8       `json:"online_status"`
	Adzone       *Adzone     `json:"-"`
	Creatives    []*Creative `json:"creatives"`
}

type Creative struct {
	Id        uint64 `json:"id"`
	AdgroupId uint64 `json:"adgroup_id,omitempty"`
	Image     string `json:"image,omitempty"`
	Link      string `json:"link,omitempty"`
	Width     uint   `json:"width,omitempty"`
	Height    uint   `json:"height,omitempty"`
}

func (this Creative) Code(key []byte) (string, error) {
	enc := binary.NewEncoder()
	enc.Encode(this)
	return utils.AESEncryptBytes(key, enc.Buffer())
}

func DecodeCreative(key []byte, cryptoText string) (creative Creative, err error) {
	data, err := utils.AESDecryptBytes(key, cryptoText)
	if err != nil {
		return creative, err
	}
	dec := binary.NewDecoder()
	dec.SetBuffer(data)
	dec.Decode(&creative)
	return creative, nil
}

type Adzone struct {
	Id   uint64 `json:"id"`
	Cid  uint   `json:"cid"`
	Page uint   `json:"page"`
	Idx  int    `json:"idx"`
}
