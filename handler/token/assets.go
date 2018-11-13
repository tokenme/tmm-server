package token

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/panjf2000/ants"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/ethgasstation-api"
	"github.com/tokenme/tmm/tools/ethplorer-api"
	"github.com/tokenme/tmm/tools/forex"
	"math/big"
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
	var (
		tokens   []*common.Token
		tokenMap = make(map[string]*common.Token)
		ethPrice = common.GetETHPrice(Service, Config)
		currency = c.Query("currency")
	)

	if currency == "" {
		currency = "USD"
	}

	ethBalance, _ := eth.BalanceOf(Service.Geth, c, user.Wallet)
	if ethBalance != nil && ethBalance.Cmp(big.NewInt(0)) == 1 {
		token := &common.Token{
			Name:     "Ethereum",
			Symbol:   "ETH",
			Decimals: 18,
			Icon:     "https://www.ethereum.org/images/logos/ETHEREUM-ICON_Black_small.png",
			Balance:  decimal.NewFromBigInt(ethBalance, -18),
		}
		token.Price = ethPrice
		tokens = append(tokens, token)
	}

	gasPrice := decimal.New(8, 0)
	gas, err := ethgasstation.Gas()
	if err == nil {
		gas.SafeLow.Div(decimal.New(10, 0))
	}
	gasLimit := decimal.New(210000, -9)
	generalMinGas := gasPrice.Mul(gasLimit)

	db := Service.Db
	tokenAddresses := make(map[string]struct{})
	var escapedAddresses []string
	client := ethplorer.NewClient(Config.EthplorerAPIKey)
	addressInfo, _ := client.GetAddressInfo(user.Wallet, "")
	if len(addressInfo.Tokens) > 0 {
		for _, addrToken := range addressInfo.Tokens {
			if addrToken.Token.Address == "" && addrToken.Token.Symbol == "" {
				continue
			}
			addr := strings.ToLower(addrToken.Token.Address)
			if _, found := tokenAddresses[addr]; !found {
				escapedAddresses = append(escapedAddresses, fmt.Sprintf("'%s'", db.Escape(addr)))
				var minGas decimal.Decimal
				if addr != Config.TMMTokenAddress {
					minGas = generalMinGas
				}
				token := &common.Token{
					Address:  addr,
					Name:     addrToken.Token.Name,
					Symbol:   addrToken.Token.Symbol,
					Decimals: uint(addrToken.Token.Decimals.IntPart()),
					MinGas:   minGas,
					Price:    addrToken.Token.Price.Rate,
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
		if addr == Config.TMMTokenAddress {
			price = common.GetTMMPrice(Service, Config, common.MarketPrice)
		}
		if token, found := tokenMap[addr]; found {
			if token.Price.Equal(decimal.New(0, 0)) {
				token.Price = price
			}
			token.Icon = icon
		} else {
			var minGas decimal.Decimal
			if addr != Config.TMMTokenAddress {
				minGas = generalMinGas
			}
			token := &common.Token{
				Address:  addr,
				Name:     row.Str(3),
				Symbol:   row.Str(4),
				Decimals: row.Uint(5),
				Icon:     icon,
				Price:    price,
				MinGas:   minGas,
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
			log.Error(err.Error())
			return err
		}
		balance, err := utils.TokenBalanceOf(tokenABI, user.Wallet)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		balanceDecimal, err := decimal.NewFromString(balance.String())
		if err != nil {
			log.Error(err.Error())
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
	var (
		val       []string
		updateVal []string
	)
	for _, token := range tokenMap {
		if _, found := knownAddressMap[token.Address]; found && token.Address != Config.TMMTokenAddress {
			updateVal = append(updateVal, fmt.Sprintf("('%s', %s)", token.Address, token.Price.String()))
			continue
		}
		val = append(val, fmt.Sprintf("('%s', '%s', '%s', %d, %s)", token.Address, db.Escape(token.Name), db.Escape(token.Symbol), token.Decimals, token.Price.String()))
	}
	if len(val) > 0 {
		_, _, err := db.Query(`INSERT IGNORE INTO tmm.erc20 (address, name, symbol, decimals, price) VALUES %s`, strings.Join(val, ","))
		if err != nil {
			log.Error(err.Error())
		}
	}
	if len(updateVal) > 0 {
		_, _, err := db.Query(`INSERT INTO tmm.erc20 (address, price) VALUES %s ON DUPLICATE KEY UPDATE price=VALUES(price)`, strings.Join(updateVal, ","))
		if err != nil {
			log.Error(err.Error())
		}
	}
	if currency != "USD" {
		rate := forex.Rate(Service, "USD", currency)
		for _, token := range tokens {
			token.Price = token.Price.Mul(rate)
		}
	}
	c.JSON(http.StatusOK, tokens)
}
