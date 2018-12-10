package mobilecode

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/tools/zz253"
	"github.com/tokenme/tmm/utils/twilio"
	"math/rand"
	"time"
)

const (
	BaseCode    = "0123456789"
	RegisterMsg = "【友币】您的验证码是: %s"
)

type Client struct {
	service *common.Service
	config  common.Config
}

func NewClient(service *common.Service, config common.Config) *Client {
	return &Client{
		service: service,
		config:  config,
	}
}

func (this *Client) Send(mobile string) (string, error) {
	db := this.service.Db
	rows, _, err := db.Query(`SELECT 1 FROM tmm.mobile_codes WHERE mobile='%s' AND updated>=DATE_SUB(NOW(), INTERVAL 1 MINUTE) LIMIT 1`, mobile)
	if len(rows) > 0 || err != nil {
		if err != nil {
			log.Error(err.Error())
		}
		return "", errors.New("请稍后1分钟后重试")
	}
	code := genCode()
	apiReq := zz253.SendRequest{}
	apiReq.Phones = []string{mobile}
	content := fmt.Sprintf(RegisterMsg, code)
	apiReq.Content = content
	apiReq.BaseRequest = zz253.BaseRequest{
		Account:  this.config.ZZ253.Account,
		Password: this.config.ZZ253.Password,
	}
	_, err = zz253.Send(&apiReq)
	if err != nil {
		log.Error(err.Error())
		return "", errors.New("短信发送失败")
	}
	_, _, err = db.Query(`INSERT INTO tmm.mobile_codes(mobile, code) VALUES ('%s', '%s') ON DUPLICATE KEY UPDATE code=VALUES(code)`, db.Escape(mobile), db.Escape(code))
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	return code, nil
}

func (this *Client) Verify(mobile string, code string) twilio.AuthVerificationResponse {
	db := this.service.Db
	res := twilio.AuthVerificationResponse{Success: true, Message: "OK"}
	row, _, err := db.Query(`SELECT 1 FROM tmm.mobile_codes WHERE mobile='%s' AND code='%s' LIMIT 1`, mobile, code)
	if err != nil || row == nil {
		if err != nil {
			log.Error(err.Error())
		}
		res.Success = false
		res.Message = "验证码错误,请重新输入"
		return res
	}
	return res

}

func genCode() string {
	code := shuffle(BaseCode)
	code = shuffle(code)
	return code[0:4]
}

func shuffle(s string) (ret string) {
	var buffer bytes.Buffer
	r := rand.New(rand.NewSource(time.Now().Unix()))
	runes := []rune(s)
	perm := r.Perm(len(runes))
	for _, randIndex := range perm {
		buffer.WriteString(string(runes[randIndex]))
	}
	return buffer.String()

}
