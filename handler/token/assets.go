package token

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/panjf2000/ants"
	"github.com/shopspring/decimal"
	"github.com/tokenme/etherscan-api"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strings"
	"sync"
)

func AssetsHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	user.Wallet = "0xb2D44E5eA830333aBD20b769fAD052bF13D40A41"

	var (
		tokens   []*common.Token
		tokenMap = make(map[string]*common.Token)
	)

	db := Service.Db

	client := etherscan.New(etherscan.Mainnet, Config.EtherscanAPIKey)
	offset := 1000
	page := 1
	tokenAddresses := make(map[string]struct{})
	tokenAddresses[Config.TMMTokenAddress] = struct{}{}
	var escapedAddresses []string
	for {
		txs, err := client.ERC20Transfers(nil, &user.Wallet, nil, nil, page, offset, false)
		if err != nil || len(txs) == 0 {
			break
		}
		page += 1
		for _, tx := range txs {
			addr := strings.ToLower(tx.ContractAddress)
			if tx.TokenName == "" {
				continue
			}
			if _, found := tokenAddresses[addr]; !found {
				escapedAddresses = append(escapedAddresses, fmt.Sprintf("'%s'", db.Escape(addr)))
				token := &common.Token{
					Address:  addr,
					Name:     tx.TokenName,
					Symbol:   tx.TokenSymbol,
					Decimals: uint(tx.TokenDecimal),
				}
				tokenMap[token.Address] = token
				tokens = append(tokens, token)
			}
			tokenAddresses[addr] = struct{}{}
		}
	}
	if _, found := tokenAddresses[Config.TMMTokenAddress]; !found {
		escapedAddresses = append(escapedAddresses, fmt.Sprintf("'%s'", db.Escape(Config.TMMTokenAddress)))
	}
	if len(escapedAddresses) == 0 {
		c.JSON(http.StatusOK, tokens)
		return
	}
	rows, _, err := db.Query(`SELECT address, icon, price, name, symbol, decimals FROM tmm.erc20 WHERE address IN (%s)`, strings.Join(escapedAddresses, ","))
	if CheckErr(err, c) {
		return
	}
	knownAddressMap := make(map[string]struct{})
	for _, row := range rows {
		addr := strings.ToLower(row.Str(0))
		icon := row.Str(1)
		price, _ := decimal.NewFromString(row.Str(2))
		if token, found := tokenMap[addr]; found {
			token.Price = price
			token.Icon = icon
		} else {
			token := &common.Token{
				Address:  addr,
				Name:     row.Str(3),
				Symbol:   row.Str(4),
				Decimals: row.Uint(5),
				Icon:     icon,
				Price:    price,
			}
			tokenMap[addr] = token
			tokens = append(tokens, token)
		}
		knownAddressMap[addr] = struct{}{}
	}
	var wg sync.WaitGroup
	balancePool, _ := ants.NewPoolWithFunc(10000, func(req interface{}) error {
		defer wg.Done()
		token := req.(*common.Token)
		tokenABI, err := utils.NewToken(token.Address, Service.Geth)
		if err != nil {
			return err
		}
		balance, err := utils.TokenBalanceOf(tokenABI, user.Wallet)
		if err != nil {
			return err
		}
		balanceDecimal, err := decimal.NewFromString(balance.String())
		if err != nil {
			return err
		}
		if token.Decimals > 0 {
			token.Balance = balanceDecimal.Div(decimal.New(1, int32(token.Decimals)))
		} else {
			token.Balance = balanceDecimal
		}

		return nil
	})
	for _, token := range tokenMap {
		wg.Add(1)
		balancePool.Serve(token)
	}
	wg.Wait()
	var val []string
	for _, token := range tokenMap {
		if _, found := knownAddressMap[token.Address]; found {
			continue
		}
		val = append(val, fmt.Sprintf("('%s', '%s', '%s', %d)", token.Address, db.Escape(token.Name), db.Escape(token.Symbol), token.Decimals))
	}
	if len(val) > 0 {
		_, _, err := db.Query(`INSERT IGNORE INTO tmm.erc20 (address, name, symbol, decimals) VALUES %s`, strings.Join(val, ","))
		if err != nil {
			log.Error(err.Error())
		}
	}
	c.JSON(http.StatusOK, tokens)
}
