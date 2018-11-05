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

type OrderCancelRequest struct {
	Id uint64 `json:"id" form:"id" binding:"required"`
}

func OrderCancelHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req OrderCancelRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	escrow, err := utils.NewEscrow(Config.TMMEscrowAddress, Service.Geth)
	if CheckErr(err, c) {
		return
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT side, quantity, quantity-deal_quantity, deal_eth, price FROM tmm.orderbooks WHERE id=%d AND online_status=0 LIMIT 1`, req.Id)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	side := orderbook.Side(rows[0].Uint(0))
	quantity, _ := decimal.NewFromString(rows[0].Str(1))
	quantityLeft, _ := decimal.NewFromString(rows[0].Str(2))
	dealEth, _ := decimal.NewFromString(rows[0].Str(3))
	price, _ := decimal.NewFromString(rows[0].Str(4))
	var txHash string
	var gasPrice *big.Int
	gas, err := ethgasstation.Gas()
	if err != nil {
		gasPrice = new(big.Int).Mul(big.NewInt(2), big.NewInt(params.GWei))
	} else {
		gasPrice = new(big.Int).Mul(big.NewInt(gas.SafeLow.Div(decimal.New(10, 0)).IntPart()), big.NewInt(params.GWei))
	}
	var gasLimit uint64 = 540000

	agentPrivKey, err := commonutils.AddressDecrypt(Config.TMMAgentWallet.Data, Config.TMMAgentWallet.Salt, Config.TMMAgentWallet.Key)
	if CheckErr(err, c) {
		return
	}
	agentPubKey, err := eth.AddressFromHexPrivateKey(agentPrivKey)
	if CheckErr(err, c) {
		return
	}

	transactor := eth.TransactorAccount(agentPrivKey)

	if side == orderbook.Ask {
		token, err := utils.NewToken(Config.TMMTokenAddress, Service.Geth)
		if CheckErr(err, c) {
			return
		}
		decimals, err := utils.TokenDecimal(token)
		if CheckErr(err, c) {
			return
		}
		amountInt := quantityLeft.Mul(decimal.New(1, int32(decimals)))
		amount, ok := new(big.Int).SetString(amountInt.Floor().String(), 10)
		if Check(!ok, fmt.Sprintf("Internal Error: %s", amountInt.Floor().String()), c) {
			return
		}
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
		tx, err := utils.EscrowWithdrawAsk(escrow, transactor, user.Wallet, amount)
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
		amountInt := quantity.Mul(price).Mul(decimals).Sub(dealEth)
		amount, ok := new(big.Int).SetString(amountInt.Floor().String(), 10)
		if Check(!ok, fmt.Sprintf("Internal Error: %s", amountInt.Floor().String()), c) {
			return
		}
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
		tx, err := utils.EscrowWithdrawBid(escrow, transactor, user.Wallet, amount)
		if CheckErr(err, c) {
			return
		}
		err = eth.NonceIncr(c, Service.Geth, Service.Redis.Master, GlobalLock, agentPubKey, Config.Geth)
		if err != nil {
			log.Error(err.Error())
		}
		txHash = tx.Hash().Hex()
	}
	query := `UPDATE tmm.orderbooks SET withdraw_tx='%s', withdraw_tx_status=2, online_status=-1 WHERE id=%d`
	_, _, err = db.Query(query, db.Escape(txHash), req.Id)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
