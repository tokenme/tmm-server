package txaccelerate

import (
	"context"
	"errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/common"
	commonutils "github.com/tokenme/tmm/utils"
	"log"
	"math/big"
)

func Accelerate(service *common.Service, config common.Config, tx string, gas int64, data string, nonce uint64) error {
	ctx := context.Background()
	geth := service.Geth
	transaction, isPending, err := geth.TransactionByHash(ctx, ethcommon.HexToHash(tx))
	var notFound bool
	if err != nil {
		notFound = true
	}
	gasPrice := new(big.Int).Mul(big.NewInt(gas), big.NewInt(params.GWei))
	var rawTx *types.Transaction
	if data != "" && nonce > 0 && notFound {
		rawTx = types.NewTransaction(nonce, ethcommon.HexToAddress("0x5aeba72b15e4ef814460e49beca6d176caec228b"), big.NewInt(0), 540000, gasPrice, hexutil.MustDecode(data))
	} else if notFound {
		return errors.New("not found")
	} else if !isPending {
		return nil
	} else {
		rawTx = types.NewTransaction(transaction.Nonce(), *transaction.To(), transaction.Value(), transaction.Gas(), gasPrice, transaction.Data())
	}

	agentPrivKey, err := commonutils.AddressDecrypt(config.TMMAgentWallet.Data, config.TMMAgentWallet.Salt, config.TMMAgentWallet.Key)
	if err != nil {
		return err
	}

	transactor := eth.TransactorAccount(agentPrivKey)
	if transactor.Signer == nil {
		return errors.New("no signer to authorize the transaction with")
	}
	newTransaction, err := transactor.Signer(types.HomesteadSigner{}, transactor.From, rawTx)
	if err != nil {
		return err
	}
	err = geth.SendTransaction(ctx, newTransaction)
	if err != nil {
		return err
	}
	db := service.Db
	_, _, err = db.Query(`UPDATE tmm.withdraw_txs SET tx='%s' WHERE tx='%s' AND tx_status=2`, newTransaction.Hash().Hex(), tx)
	if err != nil {
		log.Println(err.Error())
	}
	_, _, err = db.Query(`UPDATE tmm.exchange_records AS er SET er.tx='%s' WHERE er.tx='%s' AND er.status=2`, newTransaction.Hash().Hex(), tx)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Tx: " + newTransaction.Hash().Hex())
	return nil
}
