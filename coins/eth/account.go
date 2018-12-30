package eth

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/garyburd/redigo/redis"
	"github.com/mkideal/log"
	"math/big"
	"strings"
)

var (
	InitialGas = new(big.Int).Mul(big.NewInt(10), big.NewInt(params.Ether))
	MinGas     = new(big.Int).Mul(big.NewInt(2), big.NewInt(params.Ether))
)

const (
	MAIN_CHAIN = "main"
	UC_CHAIN   = "uc"
)

func GenerateAccount() (string, string, error) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	prvKey := hex.EncodeToString(crypto.FromECDSA(key))
	pubKey := "0x" + hex.EncodeToString(addr[:])
	return prvKey, pubKey, nil
}

func AddressFromHexPrivateKey(hexkey string) (string, error) {
	key, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return "", err
	}
	addr := crypto.PubkeyToAddress(key.PublicKey)
	pubKey := "0x" + hex.EncodeToString(addr[:])
	return pubKey, nil
}

func PendingNonce(client *ethclient.Client, ctx context.Context, wallet string) (uint64, error) {
	return client.PendingNonceAt(ctx, common.HexToAddress(wallet))
}

func TransactorAccount(hexkey string) *bind.TransactOpts {
	key, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return nil
	}
	return bind.NewKeyedTransactor(key)
}

type TransactorOptions struct {
	Nonce    uint64
	Value    *big.Int
	GasPrice *big.Int
	GasLimit uint64
}

func TransactorUpdate(transactor *bind.TransactOpts, opt TransactorOptions, ctx context.Context) {
	if opt.Nonce > 0 {
		transactor.Nonce = new(big.Int).SetUint64(opt.Nonce)
	}
	transactor.Value = opt.Value
	transactor.GasPrice = opt.GasPrice
	transactor.GasLimit = opt.GasLimit
	transactor.Context = ctx
}

func Transfer(transactor *bind.TransactOpts, client *ethclient.Client, ctx context.Context, _to string) (tx *types.Transaction, err error) {
	var nonce uint64
	if transactor.Nonce == nil {
		nonce, err = client.PendingNonceAt(ctx, transactor.From)
		if err != nil {
			return nil, err
		}
	} else {
		nonce = transactor.Nonce.Uint64()
	}
	if transactor.GasPrice == nil {
		transactor.GasPrice, err = client.SuggestGasPrice(ctx)
		if err != nil {
			return nil, err
		}
	}
	toAddr := common.HexToAddress(_to)
	if transactor.GasLimit == 0 {
		msg := ethereum.CallMsg{From: transactor.From, To: &toAddr, Value: transactor.Value, Data: nil}
		transactor.GasLimit, err = client.EstimateGas(ctx, msg)
	}
	rawTx := types.NewTransaction(nonce, toAddr, transactor.Value, transactor.GasLimit, transactor.GasPrice, nil)
	if transactor.Signer == nil {
		return nil, errors.New("no signer to authorize the transaction with")
	}
	tx, err = transactor.Signer(types.HomesteadSigner{}, transactor.From, rawTx)
	if err != nil {
		return nil, err
	}
	err = client.SendTransaction(ctx, tx)
	return tx, err
}

func BalanceOf(client *ethclient.Client, ctx context.Context, addr string) (*big.Int, error) {
	return client.BalanceAt(ctx, common.HexToAddress(addr), nil)
}

func Nonce(ctx context.Context, client *ethclient.Client, redisConn *redis.Pool, addr string, chain string) (uint64, error) {
	conn := redisConn.Get()
	defer conn.Close()
	chainType := UC_CHAIN
	if strings.Contains(chain, "mainnet.infura.io") {
		chainType = MAIN_CHAIN
	}
	key := fmt.Sprintf("%s-%s", addr, chainType)
	nonceSaved, _ := redis.Uint64(conn.Do("GET", key))
	nonce, err := client.PendingNonceAt(ctx, common.HexToAddress(addr))
	if err != nil {
		return 0, err
	}
	if nonceSaved < nonce {
		log.Warn("UPDATE nonce: %d, saved: %d", nonce, nonceSaved)
		conn.Do("SET", key, nonce)
		return nonce, nil
	}
	log.Warn("nonce: %d, saved: %d", nonce, nonceSaved)
	return nonceSaved, nil
}

func NonceIncr(ctx context.Context, client *ethclient.Client, redisConn *redis.Pool, addr string, chain string) error {
	conn := redisConn.Get()
	defer conn.Close()
	chainType := UC_CHAIN
	if strings.Contains(chain, "mainnet.infura.io") {
		chainType = MAIN_CHAIN
	}
	key := fmt.Sprintf("%s-%s", addr, chainType)
	_, err := conn.Do("INCR", key)
	return err
}
