package token

import (
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/ethereum/go-ethereum/params"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	//"github.com/tokenme/tmm/tools/ethgasstation-api"
	commonutils "github.com/tokenme/tmm/utils"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type TransferRequest struct {
	Token  string          `json:"token" form:"token" binding:"required"`
	Amount decimal.Decimal `json:"amount" form:"amount" binding:"required"`
	To     string          `json:"to" form:"to" binding:"required"`
}

func TransferHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req TransferRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	db := Service.Db
	if CheckErr(user.IsBlocked(Service), c) {
		return
	}
	/*
		var gasPrice *big.Int
		gas, err := ethgasstation.Gas()
		if err != nil {
			gasPrice = nil
		} else {
			gasPrice = new(big.Int).Mul(big.NewInt(gas.SafeLow.Div(decimal.New(10, 0)).IntPart()), big.NewInt(params.GWei))
		}
	*/
	gasPrice, err := Service.Geth.SuggestGasPrice(c)
	if err == nil && gasPrice.Cmp(eth.MinGas) == -1 {
		gasPrice = eth.MinGas
	}
	minGas := new(big.Int).Mul(big.NewInt(60000), gasPrice)
	if req.Token == "" {
		decimals := decimal.New(params.Ether, 0)
		amountInt := req.Amount.Mul(decimals)
		amount, ok := new(big.Int).SetString(amountInt.Floor().String(), 10)
		if Check(!ok, fmt.Sprintf("Internal Error: %s", amountInt.Floor().String()), c) {
			return
		}
		ethBalance, err := eth.BalanceOf(Service.Geth, c, user.Wallet)
		if CheckErr(err, c) {
			return
		}
		minRequired := new(big.Int).Add(amount, minGas)
		minRequiredDecimal := decimal.NewFromBigInt(minRequired, 0).Div(decimals)
		if ethBalance.Cmp(minRequired) < 1 {
			c.JSON(NOT_ENOUGH_ETH_ERROR, gin.H{"min_eth": minRequiredDecimal})
			return
		}
		db := Service.Db
		rows, _, err := db.Query(`SELECT wallet, wallet_salt FROM ucoin.users WHERE id=%d`, user.Id)
		if CheckErr(err, c) {
			return
		}
		if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
			return
		}
		walletData := rows[0].Str(0)
		walletSalt := rows[0].Str(1)
		userPrivateKey, err := commonutils.AddressDecrypt(walletData, walletSalt, Config.TokenSalt)
		if CheckErr(err, c) {
			return
		}
		transactor := eth.TransactorAccount(userPrivateKey)
		GlobalLock.Lock()
		defer GlobalLock.Unlock()
		nonce, err := eth.Nonce(c, Service.Geth, Service.Redis.Master, user.Wallet, Config.Geth)
		if CheckErr(err, c) {
			return
		}
		transactorOpts := eth.TransactorOptions{
			Nonce:    nonce,
			GasPrice: gasPrice,
			GasLimit: 60000,
			Value:    amount,
		}
		eth.TransactorUpdate(transactor, transactorOpts, c)
		tx, err := eth.Transfer(transactor, Service.Geth, c, req.To)
		if CheckErr(err, c) {
			return
		}

		err = eth.NonceIncr(c, Service.Geth, Service.Redis.Master, user.Wallet, Config.Geth)
		if err != nil {
			log.Error(err.Error())
		}
		receipt := common.Transaction{
			Receipt:    tx.Hash().Hex(),
			Value:      req.Amount,
			Status:     2,
			InsertedAt: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusOK, receipt)
	}

	token, err := utils.NewToken(req.Token, Service.Geth)
	if CheckErr(err, c) {
		return
	}
	tokenDecimal, err := utils.TokenDecimal(token)
	if CheckErr(err, c) {
		return
	}
	decimals := decimal.New(1, int32(tokenDecimal))
	amountInt := req.Amount.Mul(decimals)
	amount, ok := new(big.Int).SetString(amountInt.Floor().String(), 10)
	if Check(!ok, fmt.Sprintf("Internal Error: %s", amountInt.Floor().String()), c) {
		return
	}

	tokenBalance, err := utils.TokenBalanceOf(token, user.Wallet)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(tokenBalance.Cmp(amount) == -1, NOT_ENOUGH_TOKEN_ERROR, "not enough token", c) {
		return
	}

	if strings.ToLower(req.Token) == Config.TMMTokenAddress {
		minTransferRequired := decimal.New(1, 3)
		if Check(req.Amount.LessThan(minTransferRequired), fmt.Sprintf("最低转账要求%s", minTransferRequired.String()), c) {
			return
		}
		agentPrivKey, err := commonutils.AddressDecrypt(Config.TMMAgentWallet.Data, Config.TMMAgentWallet.Salt, Config.TMMAgentWallet.Key)
		if CheckErr(err, c) {
			return
		}
		agentPubKey, err := eth.AddressFromHexPrivateKey(agentPrivKey)
		if CheckErr(err, c) {
			return
		}

		transactor := eth.TransactorAccount(agentPrivKey)
		GlobalLock.Lock()
		defer GlobalLock.Unlock()
		nonce, err := eth.Nonce(c, Service.Geth, Service.Redis.Master, agentPubKey, Config.Geth)
		if CheckErr(err, c) {
			return
		}
		transactorOpts := eth.TransactorOptions{
			Nonce:    nonce,
			GasPrice: gasPrice,
			GasLimit: 540000,
		}
		eth.TransactorUpdate(transactor, transactorOpts, c)
		tx, err := utils.TransferProxy(token, transactor, user.Wallet, req.To, amount)
		if CheckErr(err, c) {
			return
		}

		err = eth.NonceIncr(c, Service.Geth, Service.Redis.Master, agentPubKey, Config.Geth)
		if err != nil {
			log.Error(err.Error())
		}
		receipt := common.Transaction{
			Receipt:    tx.Hash().Hex(),
			Value:      req.Amount,
			Status:     2,
			InsertedAt: time.Now().Format(time.RFC3339),
		}

		_, _, err = db.Query(`INSERT INTO tmm.transfer_txs (tx, from_addr, to_addr, amount) VALUES ('%s', '%s', '%s', %s)`, db.Escape(receipt.Receipt), db.Escape(user.Wallet), db.Escape(strings.ToLower(req.To)), receipt.Value.String())
		if err != nil {
			log.Error(err.Error())
		}
		c.JSON(http.StatusOK, receipt)
	} else {
		ethBalance, err := eth.BalanceOf(Service.Geth, c, user.Wallet)
		if CheckErr(err, c) {
			return
		}
		minRequiredDecimal := decimal.NewFromBigInt(minGas, 0).Div(decimals)

		if ethBalance.Cmp(minGas) < 1 {
			c.JSON(NOT_ENOUGH_ETH_ERROR, gin.H{"eth": minRequiredDecimal})
			return
		}
		rows, _, err := db.Query(`SELECT wallet, wallet_salt FROM ucoin.users WHERE id=%d`, user.Id)
		if CheckErr(err, c) {
			return
		}
		if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
			return
		}
		walletData := rows[0].Str(0)
		walletSalt := rows[0].Str(1)
		userPrivateKey, err := commonutils.AddressDecrypt(walletData, walletSalt, Config.TokenSalt)
		if CheckErr(err, c) {
			return
		}
		transactor := eth.TransactorAccount(userPrivateKey)
		GlobalLock.Lock()
		defer GlobalLock.Unlock()
		nonce, err := eth.Nonce(c, Service.Geth, Service.Redis.Master, user.Wallet, Config.Geth)
		if CheckErr(err, c) {
			return
		}
		transactorOpts := eth.TransactorOptions{
			Nonce:    nonce,
			GasPrice: gasPrice,
			GasLimit: 540000,
		}
		eth.TransactorUpdate(transactor, transactorOpts, c)
		tx, err := utils.TransferProxy(token, transactor, user.Wallet, req.To, amount)
		if CheckErr(err, c) {
			return
		}

		err = eth.NonceIncr(c, Service.Geth, Service.Redis.Master, user.Wallet, Config.Geth)
		if err != nil {
			log.Error(err.Error())
		}
		receipt := common.Transaction{
			Receipt:    tx.Hash().Hex(),
			Value:      req.Amount,
			Status:     2,
			InsertedAt: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusOK, receipt)
	}
}
