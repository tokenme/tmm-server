package utils

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	//"github.com/mkideal/log"
	"github.com/tokenme/tmm/coins/eth"
	"math"
	"math/big"
)

func DeployToken(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, decimals uint, totalSupply uint64) (common.Address, *types.Transaction, *eth.Token, error) {
	initialSupply := new(big.Int).Mul(new(big.Int).SetUint64(totalSupply), new(big.Int).SetUint64(uint64(math.Pow(10, float64(decimals)))))
	return eth.DeployToken(auth, backend, name, symbol, uint8(decimals), initialSupply)
}

func NewToken(addr string, backend bind.ContractBackend) (*eth.Token, error) {
	return eth.NewToken(common.HexToAddress(addr), backend)
}

func TokenAllowance(token *eth.Token, wallet string, contract string) (*big.Int, error) {
	return token.Allowance(nil, common.HexToAddress(wallet), common.HexToAddress(contract))
}

func TokenApprove(token *eth.Token, opts *bind.TransactOpts, spender string, value *big.Int) (*types.Transaction, error) {
	return token.Approve(opts, common.HexToAddress(spender), value)
}

func BalanceOfToken(geth *ethclient.Client, tokenAddress string, wallet string) (*big.Int, error) {
	token, err := eth.NewToken(common.HexToAddress(tokenAddress), geth)
	if err != nil {
		return nil, err
	}
	return TokenBalanceOf(token, wallet)
}

func TokenBalanceOf(token *eth.Token, wallet string) (*big.Int, error) {
	return token.BalanceOf(nil, common.HexToAddress(wallet))
}

func Transfer(token *eth.Token, opts *bind.TransactOpts, _to string, _value *big.Int) (*types.Transaction, error) {
	return token.Transfer(opts, common.HexToAddress(_to), _value)
}

func TransferProxy(token *eth.Token, opts *bind.TransactOpts, _from string, _to string, _value *big.Int) (*types.Transaction, error) {
	return token.TransferProxy(opts, common.HexToAddress(_from), common.HexToAddress(_to), _value)
}

func TokenMeta(token *eth.Token, owner string) (name string, symbol string, decimals uint8, initialSupply *big.Int, totalSupply *big.Int, totalTransfers *big.Int, totalHolers *big.Int, circulatingSupply *big.Int, balance *big.Int, err error) {
	return token.Meta(nil, common.HexToAddress(owner))
}

func TokenDecimal(token *eth.Token) (int, error) {
	d, err := token.Decimals(nil)
	if err != nil {
		return 0, err
	}
	return int(d), nil
}
