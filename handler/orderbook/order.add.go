package orderbook

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
	"github.com/tokenme/tmm/tools/ethgasstation-api"
	"github.com/tokenme/tmm/tools/orderbook"
	commonutils "github.com/tokenme/tmm/utils"
	"math/big"
	"net/http"
)

type OrderAddRequest struct {
	ProcessType orderbook.ProcessType `json:"process_type" form:"process_type" binding:"required"`
	Side        orderbook.Side        `json:"side" form:"side" binding:"required"`
	Quantity    decimal.Decimal       `json:"quantity" form:"quantity" binding:"required"`
	Price       decimal.Decimal       `json:"price" form:"price" binding:"required"`
}

func OrderAddHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req OrderAddRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	escrow, err := utils.NewEscrow(Config.TMMEscrowAddress, Service.Geth)
	if CheckErr(err, c) {
		return
	}
	db := Service.Db
	var txHash string
	var gasPrice *big.Int
	gas, err := ethgasstation.Gas()
	if err != nil {
		gasPrice = new(big.Int).Mul(big.NewInt(2), big.NewInt(params.Shannon))
	} else {
		gasPrice = new(big.Int).Mul(big.NewInt(gas.SafeLow.Div(decimal.New(10, 0)).IntPart()), big.NewInt(params.Shannon))
	}
	var gasLimit uint64 = 540000
	if req.Side == orderbook.Ask {
		token, err := utils.NewToken(Config.TMMTokenAddress, Service.Geth)
		if CheckErr(err, c) {
			return
		}
		decimals, err := utils.TokenDecimal(token)
		if CheckErr(err, c) {
			return
		}
		amountInt := req.Quantity.Mul(decimal.New(1, int32(decimals)))
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
		agentPrivKey, err := commonutils.AddressDecrypt(Config.TMMAgentWallet.Data, Config.TMMAgentWallet.Salt, Config.TMMAgentWallet.Key)
		if CheckErr(err, c) {
			return
		}
		agentPubKey, err := eth.AddressFromHexPrivateKey(agentPrivKey)
		if CheckErr(err, c) {
			return
		}

		transactor := eth.TransactorAccount(agentPrivKey)
		nonce, err := eth.Nonce(c, Service.Geth, Service.Redis.Master, GlobalLock, agentPubKey, Config.Geth)
		if CheckErr(err, c) {
			return
		}
		transactorOpts := eth.TransactorOptions{
			Nonce:    nonce,
			GasPrice: gasPrice,
			GasLimit: gasLimit,
		}
		eth.TransactorUpdate(transactor, transactorOpts, c)
		tx, err := utils.EscrowSellFrom(escrow, transactor, user.Wallet, amount)
		if CheckErr(err, c) {
			return
		}
		err = eth.NonceIncr(c, Service.Geth, Service.Redis.Master, GlobalLock, agentPubKey, Config.Geth)
		if err != nil {
			log.Error(err.Error())
		}
		txHash = tx.Hash().Hex()
	} else {
		decimals := decimal.New(params.Ether, 0)
		amountInt := req.Quantity.Mul(req.Price).Mul(decimals)
		amount, ok := new(big.Int).SetString(amountInt.Floor().String(), 10)
		if Check(!ok, fmt.Sprintf("Internal Error: %s", amountInt.Floor().String()), c) {
			return
		}
		ethBalance, err := eth.BalanceOf(Service.Geth, c, user.Wallet)
		if CheckErr(err, c) {
			return
		}
		minGas := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gasLimit))
		minRequired := new(big.Int).Add(amount, minGas)
		minRequiredDecimal := decimal.NewFromBigInt(minRequired, 0).Div(decimals)
		if ethBalance.Cmp(minRequired) < 1 {
			c.JSON(NOT_ENOUGH_ETH_ERROR, gin.H{"min_eth": minRequiredDecimal})
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
		nonce, err := eth.Nonce(c, Service.Geth, Service.Redis.Master, GlobalLock, user.Wallet, Config.Geth)
		if CheckErr(err, c) {
			return
		}
		transactorOpts := eth.TransactorOptions{
			Nonce:    nonce,
			GasPrice: gasPrice,
			GasLimit: gasLimit,
			Value:    amount,
		}
		eth.TransactorUpdate(transactor, transactorOpts, c)
		tx, err := utils.EscrowBuy(escrow, transactor)
		if CheckErr(err, c) {
			return
		}
		err = eth.NonceIncr(c, Service.Geth, Service.Redis.Master, GlobalLock, user.Wallet, Config.Geth)
		if err != nil {
			log.Error(err.Error())
		}
		txHash = tx.Hash().Hex()
	}
	query := `INSERT INTO tmm.orderbooks (user_id, side, process_type, quantity, price, deposit_tx) VALUES (%d, %d, %d, %s, %s, '%s')`
	_, _, err = db.Query(query, user.Id, req.Side, req.ProcessType, req.Quantity, req.Price, db.Escape(txHash))
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
