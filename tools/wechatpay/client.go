package wechatpay

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strings"
)

func (this *Client) Pay(param *Request) (ret Response, err error) {
	param.AppId = this.AppId
	param.MchId = this.MchId
	if param.CheckName == "" {
		param.CheckName = NO_CHECK
	}
	param.Nonce = getNonceStr()
	err = param.GenSign(this.Key)
	if err != nil {
		return ret, err
	}
	xmlBytes, err := xml.Marshal(param)
	if err != nil {
		return ret, err
	}
	fmt.Println(string(xmlBytes))
	reader := bytes.NewReader(xmlBytes)
	request, err := http.NewRequest("POST", GATEWAY_PAY, reader)
	if err != nil {
		return ret, err
	}
	crtPool := x509.NewCertPool()
	caCrt, err := ioutil.ReadFile(this.CertCrt)
	if err != nil {
		return ret, err
	}
	crtPool.AppendCertsFromPEM(caCrt)
	cliCrt, err := tls.LoadX509KeyPair(this.CertCrt, this.CertKey)
	if err != nil {
		return ret, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            crtPool,
			Certificates:       []tls.Certificate{cliCrt},
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: tr}
	resp, err := client.Do(request)
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}
	err = xml.Unmarshal(respBytes, &ret)
	return ret, err
}

func getNonceStr() (nonceStr string) {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	for i := 0; i < 32; i++ {
		idx := rand.Intn(len(chars) - 1)
		nonceStr += chars[idx : idx+1]
	}
	return
}

func (param *Request) GenSign(key string) error {
	m := param.toMap()
	var signData []string
	for k, v := range m {
		if v != "" {
			signData = append(signData, fmt.Sprintf("%s=%s", k, v))
		}
	}

	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	signStr = signStr + "&key=" + key
	fmt.Println(signStr)
	c := md5.New()
	_, err := c.Write([]byte(signStr))
	if err != nil {
		return err
	}
	signByte := c.Sum(nil)
	param.Sign = strings.ToUpper(fmt.Sprintf("%x", signByte))
	return nil

}
