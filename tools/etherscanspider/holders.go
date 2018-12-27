package etherscanspider

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	//"github.com/davecgh/go-spew/spew"
	"github.com/levigross/grequests"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	"strings"
)

type Holder struct {
	Address string
	Amount  decimal.Decimal
}

func GetHolders(service *common.Service) error {
	var (
		page    int
		holders []Holder
	)
	for {
		hs, err := getHoldersPage(page)
		if err != nil {
			log.Error(err.Error())
		} else {
			holders = append(holders, hs...)
		}
		if len(hs) < 50 {
			break
		}
		page += 1
	}
	db := service.Db
	var val []string
	for _, h := range holders {
		log.Info("W: %s, B: %s", h.Address, h.Amount.String())
		val = append(val, fmt.Sprintf("('%s', %s)", db.Escape(strings.ToLower(h.Address)), h.Amount.String()))
	}
	if len(val) > 0 {
		_, _, err := db.Query("TRUNCATE TABLE ucoin.holders")
		if err != nil {
			log.Error(err.Error())
		}
		_, _, err = db.Query(`INSERT IGNORE INTO ucoin.holders (wallet, balance) VALUES %s`, strings.Join(val, ","))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func getHoldersPage(page int) (holders []Holder, err error) {
	baseUrl := "https://etherscan.io/token/generic-tokenholders2?a=0x5aeba72b15e4ef814460e49beca6d176caec228b&p=%d"
	ro := &grequests.RequestOptions{
		UserAgent:    "Mozilla/5.0 (Linux; U; Android 4.3; en-us; SM-N900T Build/JSS15J) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
		UseCookieJar: false,
	}
	httpClient := grequests.NewSession(ro)
	resp, err := httpClient.Get(fmt.Sprintf(baseUrl, page), ro)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	reader := bytes.NewBuffer(resp.Bytes())
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	doc.Find("#maintable table.table tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		var holder Holder
		s.Find("td").Each(func(ii int, ss *goquery.Selection) {
			if ii == 1 {
				holder.Address = ss.Find("a").Text()
			} else if ii == 2 {
				holder.Amount, _ = decimal.NewFromString(ss.Text())
			}
		})
		if holder.Address != "" {
			holders = append(holders, holder)
		}
	})
	return holders, nil
}
