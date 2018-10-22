package token

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/params"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/etherscan-api"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strconv"
	"time"
)

func TransactionsHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	client := etherscan.New(etherscan.Mainnet, Config.EtherscanAPIKey)
	tokenAddress := c.Param("address")
	page, _ := strconv.ParseUint(c.Param("page"), 10, 64)
	if page == 0 {
		page = 1
	}
	offset, _ := strconv.ParseUint(c.Param("pageSize"), 10, 64)
	if offset == 0 || offset > 10 {
		offset = 10
	}
	var transactions []common.Transaction
	etherValue := decimal.New(params.Ether, 0)
	txs, err := client.ERC20Transfers(&tokenAddress, &user.Wallet, nil, nil, int(page), int(offset), true)
	if err != nil || len(txs) == 0 {
		if err != nil {
			log.Error(err.Error())
		}
		c.JSON(http.StatusOK, transactions)
		return
	}
	token, err := utils.NewToken(tokenAddress, Service.Geth)
	if CheckErr(err, c) {
		return
	}
	tokenDecimal, err := utils.TokenDecimal(token)
	if CheckErr(err, c) {
		return
	}
	for _, tx := range txs {
		transaction := common.Transaction{
			Receipt:           tx.Hash,
			From:              tx.From,
			To:                tx.To,
			Value:             decimal.NewFromBigInt(tx.Value.Int(), -1*int32(tokenDecimal)),
			Gas:               decimal.New(int64(tx.Gas), 0).Div(etherValue),
			GasPrice:          decimal.NewFromBigInt(tx.GasPrice.Int(), 0).Div(etherValue),
			GasUsed:           decimal.New(int64(tx.GasUsed), 0).Div(etherValue),
			CumulativeGasUsed: decimal.New(int64(tx.CumulativeGasUsed), 0).Div(etherValue),
			Confirmations:     tx.Confirmations,
			InsertedAt:        tx.TimeStamp.Time().Format(time.RFC3339),
		}
		transactions = append(transactions, transaction)
	}
	c.JSON(http.StatusOK, transactions)
}
