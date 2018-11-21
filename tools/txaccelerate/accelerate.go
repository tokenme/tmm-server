package txaccelerate

import (
	"context"
	"errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/common"
	commonutils "github.com/tokenme/tmm/utils"
	"log"
	"math/big"
)

func Accelerate(service *common.Service, config common.Config, tx string, gas int64) error {
	ctx := context.Background()
	geth := service.Geth
	transaction, isPending, err := geth.TransactionByHash(ctx, ethcommon.HexToHash(tx))
	if err != nil {
		return err
	}
	if !isPending {
		return nil
	}
	gasPrice := new(big.Int).Mul(big.NewInt(gas), big.NewInt(params.GWei))
	rawTx := types.NewTransaction(transaction.Nonce(), *transaction.To(), transaction.Value(), transaction.Gas(), gasPrice, transaction.Data())
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
	log.Println("Tx: " + newTransaction.Hash().Hex())
	return nil
}