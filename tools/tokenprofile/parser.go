package tokenprofile

import (
	//"bytes"
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/json-iterator/go"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/tools/qiniu"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Update(service *common.Service, config common.Config) {
	ctx := context.Background()
	filepath.Walk(filepath.Join(config.TokenProfilePath, "erc20"), func(path string, info os.FileInfo, err error) error {
		if strings.Contains(info.Name(), "$template") {
			return nil
		}
		if !strings.Contains(info.Name(), ".json") {
			return nil
		}
		fd, err := os.Open(path)
		if err != nil {
			return err
		}
		content, err := ioutil.ReadAll(fd)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		var token common.Token
		err = json.Unmarshal(content, &token)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		c, _ := context.WithCancel(ctx)
		token.Icon, err = uploadImage(c, token.Address, config)
		spew.Dump(token)
		saveToken(token, service)
		return nil
	})
}

func uploadImage(ctx context.Context, address string, config common.Config) (string, error) {
	imgFd, err := os.Open(filepath.Join(config.TokenProfilePath, "images", fmt.Sprintf("%s.png", address)))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	imageData, err := ioutil.ReadAll(imgFd)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	link, _, err := qiniu.Upload(ctx, config.Qiniu, config.Qiniu.LogoPath, fmt.Sprintf("%s.png", address), imageData)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return link, err
}

func saveToken(token common.Token, service *common.Service) error {
	db := service.Db
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var (
		icon         = "NULL"
		email        = "NULL"
		website      = "NULL"
		state        = "NULL"
		publishOn    = "NULL"
		overview     = "NULL"
		initialPrice = "NULL"
		links        = "NULL"
	)
	if token.Icon != "" {
		icon = fmt.Sprintf("'%s'", db.Escape(token.Icon))
	}
	if token.Email != "" {
		email = fmt.Sprintf("'%s'", db.Escape(token.Email))
	}
	if token.Website != "" {
		website = fmt.Sprintf("'%s'", db.Escape(token.Website))
	}
	if token.State != "" {
		state = fmt.Sprintf("'%s'", db.Escape(token.State))
	}
	if token.PublishOn != "" {
		publishOn = fmt.Sprintf("'%s'", db.Escape(token.PublishOn))
	}
	if token.Overview != nil {
		js, err := json.Marshal(token.Overview)
		if err == nil {
			overview = fmt.Sprintf("'%s'", db.Escape(string(js)))
		}
	}
	if token.InitialPrice != nil {
		js, err := json.Marshal(token.InitialPrice)
		if err == nil {
			initialPrice = fmt.Sprintf("'%s'", db.Escape(string(js)))
		}
	}
	if token.Links != nil {
		js, err := json.Marshal(token.Links)
		if err == nil {
			links = fmt.Sprintf("'%s'", db.Escape(string(js)))
		}
	}
	_, _, err := db.Query(`INSERT INTO tmm.erc20 (address, name, symbol, icon, email, website, state, publish_on, overview, initial_price, links) VALUES ('%s', '%s', '%s', %s, %s, %s, %s, %s, %s, %s, %s) ON DUPLICATE KEY UPDATE symbol=VALUES(symbol), icon=VALUES(icon), email=VALUES(email), website=VALUES(website), state=VALUES(state), publish_on=VALUES(publish_on), overview=VALUES(overview), initial_price=VALUES(initial_price), links=VALUES(links)`, db.Escape(strings.ToLower(token.Address)), db.Escape(token.Name), db.Escape(token.Symbol), icon, email, website, state, publishOn, overview, initialPrice, links)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}
