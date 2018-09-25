package utils

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	//"github.com/mkideal/log"
	"github.com/tokenme/tmm/coins/eth"
	"math/big"
)

func DeployEscrow(auth *bind.TransactOpts, backend bind.ContractBackend, tokenAddress string) (common.Address, *types.Transaction, *eth.Escrow, error) {
	return eth.DeployEscrow(auth, backend, common.HexToAddress(tokenAddress))
}

func NewEscrow(addr string, backend bind.ContractBackend) (*eth.Escrow, error) {
	return eth.NewEscrow(common.HexToAddress(addr), backend)
}

func EscrowBuy(escrow *eth.Escrow, opts *bind.TransactOpts) (*types.Transaction, error) {
	return escrow.Buy(opts)
}

func EscrowSellFrom(escrow *eth.Escrow, opts *bind.TransactOpts, seller string, amount *big.Int) (*types.Transaction, error) {
	return escrow.SellFrom(opts, common.HexToAddress(seller), amount)
}

func EscrowWithdrawBid(escrow *eth.Escrow, opts *bind.TransactOpts, buyer string, value *big.Int) (*types.Transaction, error) {
	return escrow.WithdrawBid(opts, common.HexToAddress(buyer), value)
}

func EscrowWithdrawAsk(escrow *eth.Escrow, opts *bind.TransactOpts, seller string, value *big.Int) (*types.Transaction, error) {
	return escrow.WithdrawAsk(opts, common.HexToAddress(seller), value)
}
