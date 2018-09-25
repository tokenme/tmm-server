// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eth

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// EscrowABI is the input ABI used to generate the binding from.
const EscrowABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"buyers\",\"type\":\"address[]\"},{\"name\":\"sellers\",\"type\":\"address[]\"},{\"name\":\"weiAmounts\",\"type\":\"uint256[]\"},{\"name\":\"tokenAmounts\",\"type\":\"uint256[]\"}],\"name\":\"batchDeal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isAgent\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenRaised\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"buyer\",\"type\":\"address\"}],\"name\":\"bidBalanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"wallet\",\"type\":\"address\"},{\"name\":\"weiAmount\",\"type\":\"uint256\"}],\"name\":\"checkout\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"seller\",\"type\":\"address\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\"}],\"name\":\"withdrawAsk\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"weiRaised\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"buyer\",\"type\":\"address\"},{\"name\":\"weiAmount\",\"type\":\"uint256\"}],\"name\":\"withdrawBid\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addAgent\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"wallet\",\"type\":\"address\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\"}],\"name\":\"checkoutToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"buyer\",\"type\":\"address\"},{\"name\":\"seller\",\"type\":\"address\"},{\"name\":\"weiAmount\",\"type\":\"uint256\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\"}],\"name\":\"deal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"removeAgent\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"buy\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"seller\",\"type\":\"address\"}],\"name\":\"askBalanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceAgent\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"tokenAmount\",\"type\":\"uint256\"}],\"name\":\"sell\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"seller\",\"type\":\"address\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\"}],\"name\":\"sellFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"token\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TokenBid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"seller\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokenAsk\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"seller\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TokenDeal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"wallet\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Checkout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"wallet\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"CheckoutToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"weiAmount\",\"type\":\"uint256\"}],\"name\":\"WithdrawnBid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"seller\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenAmount\",\"type\":\"uint256\"}],\"name\":\"WithdrawnAsk\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AgentAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AgentRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// EscrowBin is the compiled bytecode used for deploying new contracts.
const EscrowBin = `60806040523480156200001157600080fd5b506040516020806200284383398101806040528101908080519060200190929190505050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620000993360016200011e6401000000000262002188179091906401000000009004565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614151515620000d657600080fd5b80600460006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050620001b9565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141515156200015b57600080fd5b60018260000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055505050565b61267a80620001c96000396000f300608060405260043610610128576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630a28ef331461012d5780631ffbb0641461025c57806326a21575146102b757806332a16bd5146102e257806332ec457c146103395780633d603698146103865780634042b66f146103d35780634cb8ef5b146103fe578063715018a61461044b57806384e798421461046257806388d3c62a146104a55780638da5cb5b146104f25780638f32d59b14610549578063904561811461057857806397a6278e146105ef578063a6f2ae3a14610632578063b517ba751461063c578063c692f4cf14610693578063e4849b32146106aa578063eb1715cc146106d7578063f2fde38b14610724578063fc0c546a14610767575b600080fd5b34801561013957600080fd5b5061025a600480360381019080803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192908035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437820191505050505050919291929080359060200190820180359060200190808060200260200160405190810160405280939291908181526020018383602002808284378201915050505050509192919290803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843782019150505050505091929192905050506107be565b005b34801561026857600080fd5b5061029d600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506108ba565b604051808215151515815260200191505060405180910390f35b3480156102c357600080fd5b506102cc6108d7565b6040518082815260200191505060405180910390f35b3480156102ee57600080fd5b50610323600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506108e1565b6040518082815260200191505060405180910390f35b34801561034557600080fd5b50610384600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061092a565b005b34801561039257600080fd5b506103d1600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610a48565b005b3480156103df57600080fd5b506103e8610e7d565b6040518082815260200191505060405180910390f35b34801561040a57600080fd5b50610449600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610e87565b005b34801561045757600080fd5b506104606110e4565b005b34801561046e57600080fd5b506104a3600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061119e565b005b3480156104b157600080fd5b506104f0600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061121b565b005b3480156104fe57600080fd5b5061050761141f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561055557600080fd5b5061055e611448565b604051808215151515815260200191505060405180910390f35b34801561058457600080fd5b506105ed600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291908035906020019092919050505061149f565b005b3480156105fb57600080fd5b50610630600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611a45565b005b61063a611a74565b005b34801561064857600080fd5b5061067d600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611bc6565b6040518082815260200191505060405180910390f35b34801561069f57600080fd5b506106a8611c0f565b005b3480156106b657600080fd5b506106d560048036038101908080359060200190929190505050611c48565b005b3480156106e357600080fd5b50610722600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050611de4565b005b34801561073057600080fd5b50610765600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061208a565b005b34801561077357600080fd5b5061077c6120a9565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b60008060008060006107ce611448565b806107de57506107dd336108ba565b5b15156107e957600080fd5b600089511180156107fb575087518951145b8015610808575086518951145b8015610815575085518851145b151561082057600080fd5b600094505b88518510156108af57888581518110151561083c57fe5b906020019060200201519350878581518110151561085657fe5b906020019060200201519250868581518110151561087057fe5b906020019060200201519150858581518110151561088a57fe5b9060200190602002015190506108a28484848461149f565b8480600101955050610825565b505050505050505050565b60006108d08260016120d390919063ffffffff16565b9050919050565b6000600654905090565b6000600260008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b610932611448565b151561093d57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415151561097957600080fd5b6000811415151561098957600080fd5b803073ffffffffffffffffffffffffffffffffffffffff1631101515156109af57600080fd5b8173ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f193505050501580156109f5573d6000803e3d6000fd5b508173ffffffffffffffffffffffffffffffffffffffff167f10190e7f174afbe5482194d509017e9df854671cd62fb4ef901fbd0080c1e6b6826040518082815260200191505060405180910390a25050565b6000610a52611448565b80610a625750610a61336108ba565b5b1515610a6d57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614151515610aa957600080fd5b8190506000821480610af9575081600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054105b15610b4157600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490505b80600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001915050602060405180830381600087803b158015610bff57600080fd5b505af1158015610c13573d6000803e3d6000fd5b505050506040513d6020811015610c2957600080fd5b810190808051906020019092919050505010151515610c4757600080fd5b600081111515610c5657600080fd5b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634733dc8f3085846040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b158015610d4f57600080fd5b505af1158015610d63573d6000803e3d6000fd5b505050506040513d6020811015610d7957600080fd5b81019080805190602001909291905050501515610d9557600080fd5b610de781600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461216790919063ffffffff16565b600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff167fc80bf1c92f6420fc18665ad14be6e5d7cc029baaf8fc547a6b2799a0f3bd012b826040518082815260200191505060405180910390a2505050565b6000600554905090565b6000610e91611448565b80610ea15750610ea0336108ba565b5b1515610eac57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614151515610ee857600080fd5b8190506000821480610f38575081600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054105b15610f8057600260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490505b803073ffffffffffffffffffffffffffffffffffffffff163110151515610fa657600080fd5b600081111515610fb557600080fd5b8273ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610ffb573d6000803e3d6000fd5b5061104e81600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461216790919063ffffffff16565b600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff167f064e09ad579edcb276d73f297cc24e5b0685e3b9574e9d1e5806e3d76cd949b8826040518082815260200191505060405180910390a2505050565b6110ec611448565b15156110f757600080fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167ff8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c6482060405160405180910390a260008060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b6111a6611448565b806111b657506111b5336108ba565b5b15156111c157600080fd5b6111d581600161218890919063ffffffff16565b8073ffffffffffffffffffffffffffffffffffffffff167ff68e73cec97f2d70aa641fb26e87a4383686e2efacb648f2165aeb02ac562ec560405160405180910390a250565b611223611448565b151561122e57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415151561126a57600080fd5b6000811415151561127a57600080fd5b80600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001915050602060405180830381600087803b15801561133857600080fd5b505af115801561134c573d6000803e3d6000fd5b505050506040513d602081101561136257600080fd5b81019080805190602001909291905050501015151561138057600080fd5b6113cd8282600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166122229092919063ffffffff16565b8173ffffffffffffffffffffffffffffffffffffffff167f1027be3db639879c20e5af001c80ef4b0621b3a7276e4668345ac6fb9ce5be58826040518082815260200191505060405180910390a25050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614905090565b6114a7611448565b806114b757506114b6336108ba565b5b15156114c257600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16141515156114fe57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415151561153a57600080fd5b6000821415151561154a57600080fd5b6000811415151561155a57600080fd5b81600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101515156115a857600080fd5b80600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101515156115f657600080fd5b813073ffffffffffffffffffffffffffffffffffffffff16311015151561161c57600080fd5b80600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001915050602060405180830381600087803b1580156116da57600080fd5b505af11580156116ee573d6000803e3d6000fd5b505050506040513d602081101561170457600080fd5b81019080805190602001909291905050501015151561172257600080fd5b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634733dc8f8486846040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b15801561181b57600080fd5b505af115801561182f573d6000803e3d6000fd5b505050506040513d602081101561184557600080fd5b8101908080519060200190929190505050151561186157600080fd5b8273ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f193505050501580156118a7573d6000803e3d6000fd5b506118fa82600260008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461216790919063ffffffff16565b600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555061198f81600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461216790919063ffffffff16565b600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f13b25d4eb0e3438ebba0fb869531f14c03b78b4c3e286a47cd3fd8627d4241658385604051808381526020018281526020019250505060405180910390a350505050565b611a4d611448565b80611a5d5750611a5c336108ba565b5b1515611a6857600080fd5b611a7181612310565b50565b6000349050600073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151515611ab557600080fd5b60008114151515611ac557600080fd5b611ada8160055461231c90919063ffffffff16565b600581905550611b3281600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461231c90919063ffffffff16565b600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055503373ffffffffffffffffffffffffffffffffffffffff167fd53242b5a0280628acfd89744b78a30e567c53d7c69397838d209f48dfe542a4826040518082815260200191505060405180910390a250565b6000600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b611c17611448565b80611c275750611c26336108ba565b5b1515611c3257600080fd5b611c4633600161233d90919063ffffffff16565b565b600073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151515611c8457600080fd5b60008114151515611c9457600080fd5b611ce3333083600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166123d7909392919063ffffffff16565b611cf88160065461231c90919063ffffffff16565b600681905550611d5081600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461231c90919063ffffffff16565b600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055503373ffffffffffffffffffffffffffffffffffffffff167f19bc6719e5b4b4a57c924567ee3c27073d21148f2801228976dc6361fd8190d6826040518082815260200191505060405180910390a250565b611dec611448565b80611dfc5750611dfb336108ba565b5b1515611e0757600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614151515611e4357600080fd5b60008114151515611e5357600080fd5b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634733dc8f8330846040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b158015611f4c57600080fd5b505af1158015611f60573d6000803e3d6000fd5b505050506040513d6020811015611f7657600080fd5b810190808051906020019092919050505050611f9d8160065461231c90919063ffffffff16565b600681905550611ff581600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461231c90919063ffffffff16565b600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055503373ffffffffffffffffffffffffffffffffffffffff167f19bc6719e5b4b4a57c924567ee3c27073d21148f2801228976dc6361fd8190d6826040518082815260200191505060405180910390a25050565b612092611448565b151561209d57600080fd5b6120a6816124fa565b50565b6000600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60008073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415151561211057600080fd5b8260000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905092915050565b60008083831115151561217957600080fd5b82840390508091505092915050565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141515156121c457600080fd5b60018260000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055505050565b8273ffffffffffffffffffffffffffffffffffffffff1663a9059cbb83836040518363ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b1580156122c557600080fd5b505af11580156122d9573d6000803e3d6000fd5b505050506040513d60208110156122ef57600080fd5b8101908080519060200190929190505050151561230b57600080fd5b505050565b612319816125f4565b50565b600080828401905083811015151561233357600080fd5b8091505092915050565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415151561237957600080fd5b60008260000160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055505050565b8373ffffffffffffffffffffffffffffffffffffffff166323b872dd8484846040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b1580156124ae57600080fd5b505af11580156124c2573d6000803e3d6000fd5b505050506040513d60208110156124d857600080fd5b810190808051906020019092919050505015156124f457600080fd5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415151561253657600080fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b61260881600161233d90919063ffffffff16565b8073ffffffffffffffffffffffffffffffffffffffff167fed9c8ad8d5a0a66898ea49d2956929c93ae2e8bd50281b2ed897c5d1a6737e0b60405160405180910390a2505600a165627a7a72305820d5aa34356d8e20587090c95cefa19a24b2468cf29829cd382154ec0a5109754f0029`

// DeployEscrow deploys a new Ethereum contract, binding an instance of Escrow to it.
func DeployEscrow(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address) (common.Address, *types.Transaction, *Escrow, error) {
	parsed, err := abi.JSON(strings.NewReader(EscrowABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(EscrowBin), backend, token)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Escrow{EscrowCaller: EscrowCaller{contract: contract}, EscrowTransactor: EscrowTransactor{contract: contract}, EscrowFilterer: EscrowFilterer{contract: contract}}, nil
}

// Escrow is an auto generated Go binding around an Ethereum contract.
type Escrow struct {
	EscrowCaller     // Read-only binding to the contract
	EscrowTransactor // Write-only binding to the contract
	EscrowFilterer   // Log filterer for contract events
}

// EscrowCaller is an auto generated read-only Go binding around an Ethereum contract.
type EscrowCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EscrowTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EscrowTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EscrowFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EscrowFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EscrowSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EscrowSession struct {
	Contract     *Escrow           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EscrowCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EscrowCallerSession struct {
	Contract *EscrowCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// EscrowTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EscrowTransactorSession struct {
	Contract     *EscrowTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EscrowRaw is an auto generated low-level Go binding around an Ethereum contract.
type EscrowRaw struct {
	Contract *Escrow // Generic contract binding to access the raw methods on
}

// EscrowCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EscrowCallerRaw struct {
	Contract *EscrowCaller // Generic read-only contract binding to access the raw methods on
}

// EscrowTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EscrowTransactorRaw struct {
	Contract *EscrowTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEscrow creates a new instance of Escrow, bound to a specific deployed contract.
func NewEscrow(address common.Address, backend bind.ContractBackend) (*Escrow, error) {
	contract, err := bindEscrow(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Escrow{EscrowCaller: EscrowCaller{contract: contract}, EscrowTransactor: EscrowTransactor{contract: contract}, EscrowFilterer: EscrowFilterer{contract: contract}}, nil
}

// NewEscrowCaller creates a new read-only instance of Escrow, bound to a specific deployed contract.
func NewEscrowCaller(address common.Address, caller bind.ContractCaller) (*EscrowCaller, error) {
	contract, err := bindEscrow(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EscrowCaller{contract: contract}, nil
}

// NewEscrowTransactor creates a new write-only instance of Escrow, bound to a specific deployed contract.
func NewEscrowTransactor(address common.Address, transactor bind.ContractTransactor) (*EscrowTransactor, error) {
	contract, err := bindEscrow(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EscrowTransactor{contract: contract}, nil
}

// NewEscrowFilterer creates a new log filterer instance of Escrow, bound to a specific deployed contract.
func NewEscrowFilterer(address common.Address, filterer bind.ContractFilterer) (*EscrowFilterer, error) {
	contract, err := bindEscrow(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EscrowFilterer{contract: contract}, nil
}

// bindEscrow binds a generic wrapper to an already deployed contract.
func bindEscrow(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EscrowABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Escrow *EscrowRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Escrow.Contract.EscrowCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Escrow *EscrowRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Escrow.Contract.EscrowTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Escrow *EscrowRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Escrow.Contract.EscrowTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Escrow *EscrowCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Escrow.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Escrow *EscrowTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Escrow.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Escrow *EscrowTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Escrow.Contract.contract.Transact(opts, method, params...)
}

// AskBalanceOf is a free data retrieval call binding the contract method 0xb517ba75.
//
// Solidity: function askBalanceOf(seller address) constant returns(uint256)
func (_Escrow *EscrowCaller) AskBalanceOf(opts *bind.CallOpts, seller common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Escrow.contract.Call(opts, out, "askBalanceOf", seller)
	return *ret0, err
}

// AskBalanceOf is a free data retrieval call binding the contract method 0xb517ba75.
//
// Solidity: function askBalanceOf(seller address) constant returns(uint256)
func (_Escrow *EscrowSession) AskBalanceOf(seller common.Address) (*big.Int, error) {
	return _Escrow.Contract.AskBalanceOf(&_Escrow.CallOpts, seller)
}

// AskBalanceOf is a free data retrieval call binding the contract method 0xb517ba75.
//
// Solidity: function askBalanceOf(seller address) constant returns(uint256)
func (_Escrow *EscrowCallerSession) AskBalanceOf(seller common.Address) (*big.Int, error) {
	return _Escrow.Contract.AskBalanceOf(&_Escrow.CallOpts, seller)
}

// BidBalanceOf is a free data retrieval call binding the contract method 0x32a16bd5.
//
// Solidity: function bidBalanceOf(buyer address) constant returns(uint256)
func (_Escrow *EscrowCaller) BidBalanceOf(opts *bind.CallOpts, buyer common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Escrow.contract.Call(opts, out, "bidBalanceOf", buyer)
	return *ret0, err
}

// BidBalanceOf is a free data retrieval call binding the contract method 0x32a16bd5.
//
// Solidity: function bidBalanceOf(buyer address) constant returns(uint256)
func (_Escrow *EscrowSession) BidBalanceOf(buyer common.Address) (*big.Int, error) {
	return _Escrow.Contract.BidBalanceOf(&_Escrow.CallOpts, buyer)
}

// BidBalanceOf is a free data retrieval call binding the contract method 0x32a16bd5.
//
// Solidity: function bidBalanceOf(buyer address) constant returns(uint256)
func (_Escrow *EscrowCallerSession) BidBalanceOf(buyer common.Address) (*big.Int, error) {
	return _Escrow.Contract.BidBalanceOf(&_Escrow.CallOpts, buyer)
}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(account address) constant returns(bool)
func (_Escrow *EscrowCaller) IsAgent(opts *bind.CallOpts, account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Escrow.contract.Call(opts, out, "isAgent", account)
	return *ret0, err
}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(account address) constant returns(bool)
func (_Escrow *EscrowSession) IsAgent(account common.Address) (bool, error) {
	return _Escrow.Contract.IsAgent(&_Escrow.CallOpts, account)
}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(account address) constant returns(bool)
func (_Escrow *EscrowCallerSession) IsAgent(account common.Address) (bool, error) {
	return _Escrow.Contract.IsAgent(&_Escrow.CallOpts, account)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Escrow *EscrowCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Escrow.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Escrow *EscrowSession) IsOwner() (bool, error) {
	return _Escrow.Contract.IsOwner(&_Escrow.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Escrow *EscrowCallerSession) IsOwner() (bool, error) {
	return _Escrow.Contract.IsOwner(&_Escrow.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Escrow *EscrowCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Escrow.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Escrow *EscrowSession) Owner() (common.Address, error) {
	return _Escrow.Contract.Owner(&_Escrow.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Escrow *EscrowCallerSession) Owner() (common.Address, error) {
	return _Escrow.Contract.Owner(&_Escrow.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Escrow *EscrowCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Escrow.contract.Call(opts, out, "token")
	return *ret0, err
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Escrow *EscrowSession) Token() (common.Address, error) {
	return _Escrow.Contract.Token(&_Escrow.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Escrow *EscrowCallerSession) Token() (common.Address, error) {
	return _Escrow.Contract.Token(&_Escrow.CallOpts)
}

// TokenRaised is a free data retrieval call binding the contract method 0x26a21575.
//
// Solidity: function tokenRaised() constant returns(uint256)
func (_Escrow *EscrowCaller) TokenRaised(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Escrow.contract.Call(opts, out, "tokenRaised")
	return *ret0, err
}

// TokenRaised is a free data retrieval call binding the contract method 0x26a21575.
//
// Solidity: function tokenRaised() constant returns(uint256)
func (_Escrow *EscrowSession) TokenRaised() (*big.Int, error) {
	return _Escrow.Contract.TokenRaised(&_Escrow.CallOpts)
}

// TokenRaised is a free data retrieval call binding the contract method 0x26a21575.
//
// Solidity: function tokenRaised() constant returns(uint256)
func (_Escrow *EscrowCallerSession) TokenRaised() (*big.Int, error) {
	return _Escrow.Contract.TokenRaised(&_Escrow.CallOpts)
}

// WeiRaised is a free data retrieval call binding the contract method 0x4042b66f.
//
// Solidity: function weiRaised() constant returns(uint256)
func (_Escrow *EscrowCaller) WeiRaised(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Escrow.contract.Call(opts, out, "weiRaised")
	return *ret0, err
}

// WeiRaised is a free data retrieval call binding the contract method 0x4042b66f.
//
// Solidity: function weiRaised() constant returns(uint256)
func (_Escrow *EscrowSession) WeiRaised() (*big.Int, error) {
	return _Escrow.Contract.WeiRaised(&_Escrow.CallOpts)
}

// WeiRaised is a free data retrieval call binding the contract method 0x4042b66f.
//
// Solidity: function weiRaised() constant returns(uint256)
func (_Escrow *EscrowCallerSession) WeiRaised() (*big.Int, error) {
	return _Escrow.Contract.WeiRaised(&_Escrow.CallOpts)
}

// AddAgent is a paid mutator transaction binding the contract method 0x84e79842.
//
// Solidity: function addAgent(account address) returns()
func (_Escrow *EscrowTransactor) AddAgent(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Escrow.contract.Transact(opts, "addAgent", account)
}

// AddAgent is a paid mutator transaction binding the contract method 0x84e79842.
//
// Solidity: function addAgent(account address) returns()
func (_Escrow *EscrowSession) AddAgent(account common.Address) (*types.Transaction, error) {
	return _Escrow.Contract.AddAgent(&_Escrow.TransactOpts, account)
}

// AddAgent is a paid mutator transaction binding the contract method 0x84e79842.
//
// Solidity: function addAgent(account address) returns()
func (_Escrow *EscrowTransactorSession) AddAgent(account common.Address) (*types.Transaction, error) {
	return _Escrow.Contract.AddAgent(&_Escrow.TransactOpts, account)
}

// BatchDeal is a paid mutator transaction binding the contract method 0x0a28ef33.
//
// Solidity: function batchDeal(buyers address[], sellers address[], weiAmounts uint256[], tokenAmounts uint256[]) returns()
func (_Escrow *EscrowTransactor) BatchDeal(opts *bind.TransactOpts, buyers []common.Address, sellers []common.Address, weiAmounts []*big.Int, tokenAmounts []*big.Int) (*types.Transaction, error) {
	return _Escrow.contract.Transact(opts, "batchDeal", buyers, sellers, weiAmounts, tokenAmounts)
}

// BatchDeal is a paid mutator transaction binding the contract method 0x0a28ef33.
//
// Solidity: function batchDeal(buyers address[], sellers address[], weiAmounts uint256[], tokenAmounts uint256[]) returns()
func (_Escrow *EscrowSession) BatchDeal(buyers []common.Address, sellers []common.Address, weiAmounts []*big.Int, tokenAmounts []*big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.BatchDeal(&_Escrow.TransactOpts, buyers, sellers, weiAmounts, tokenAmounts)
}

// BatchDeal is a paid mutator transaction binding the contract method 0x0a28ef33.
//
// Solidity: function batchDeal(buyers address[], sellers address[], weiAmounts uint256[], tokenAmounts uint256[]) returns()
func (_Escrow *EscrowTransactorSession) BatchDeal(buyers []common.Address, sellers []common.Address, weiAmounts []*big.Int, tokenAmounts []*big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.BatchDeal(&_Escrow.TransactOpts, buyers, sellers, weiAmounts, tokenAmounts)
}

// Buy is a paid mutator transaction binding the contract method 0xa6f2ae3a.
//
// Solidity: function buy() returns()
func (_Escrow *EscrowTransactor) Buy(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Escrow.contract.Transact(opts, "buy")
}

// Buy is a paid mutator transaction binding the contract method 0xa6f2ae3a.
//
// Solidity: function buy() returns()
func (_Escrow *EscrowSession) Buy() (*types.Transaction, error) {
	return _Escrow.Contract.Buy(&_Escrow.TransactOpts)
}

// Buy is a paid mutator transaction binding the contract method 0xa6f2ae3a.
//
// Solidity: function buy() returns()
func (_Escrow *EscrowTransactorSession) Buy() (*types.Transaction, error) {
	return _Escrow.Contract.Buy(&_Escrow.TransactOpts)
}

// Checkout is a paid mutator transaction binding the contract method 0x32ec457c.
//
// Solidity: function checkout(wallet address, weiAmount uint256) returns()
func (_Escrow *EscrowTransactor) Checkout(opts *bind.TransactOpts, wallet common.Address, weiAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.contract.Transact(opts, "checkout", wallet, weiAmount)
}

// Checkout is a paid mutator transaction binding the contract method 0x32ec457c.
//
// Solidity: function checkout(wallet address, weiAmount uint256) returns()
func (_Escrow *EscrowSession) Checkout(wallet common.Address, weiAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.Checkout(&_Escrow.TransactOpts, wallet, weiAmount)
}

// Checkout is a paid mutator transaction binding the contract method 0x32ec457c.
//
// Solidity: function checkout(wallet address, weiAmount uint256) returns()
func (_Escrow *EscrowTransactorSession) Checkout(wallet common.Address, weiAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.Checkout(&_Escrow.TransactOpts, wallet, weiAmount)
}

// CheckoutToken is a paid mutator transaction binding the contract method 0x88d3c62a.
//
// Solidity: function checkoutToken(wallet address, tokenAmount uint256) returns()
func (_Escrow *EscrowTransactor) CheckoutToken(opts *bind.TransactOpts, wallet common.Address, tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.contract.Transact(opts, "checkoutToken", wallet, tokenAmount)
}

// CheckoutToken is a paid mutator transaction binding the contract method 0x88d3c62a.
//
// Solidity: function checkoutToken(wallet address, tokenAmount uint256) returns()
func (_Escrow *EscrowSession) CheckoutToken(wallet common.Address, tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.CheckoutToken(&_Escrow.TransactOpts, wallet, tokenAmount)
}

// CheckoutToken is a paid mutator transaction binding the contract method 0x88d3c62a.
//
// Solidity: function checkoutToken(wallet address, tokenAmount uint256) returns()
func (_Escrow *EscrowTransactorSession) CheckoutToken(wallet common.Address, tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.CheckoutToken(&_Escrow.TransactOpts, wallet, tokenAmount)
}

// Deal is a paid mutator transaction binding the contract method 0x90456181.
//
// Solidity: function deal(buyer address, seller address, weiAmount uint256, tokenAmount uint256) returns()
func (_Escrow *EscrowTransactor) Deal(opts *bind.TransactOpts, buyer common.Address, seller common.Address, weiAmount *big.Int, tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.contract.Transact(opts, "deal", buyer, seller, weiAmount, tokenAmount)
}

// Deal is a paid mutator transaction binding the contract method 0x90456181.
//
// Solidity: function deal(buyer address, seller address, weiAmount uint256, tokenAmount uint256) returns()
func (_Escrow *EscrowSession) Deal(buyer common.Address, seller common.Address, weiAmount *big.Int, tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.Deal(&_Escrow.TransactOpts, buyer, seller, weiAmount, tokenAmount)
}

// Deal is a paid mutator transaction binding the contract method 0x90456181.
//
// Solidity: function deal(buyer address, seller address, weiAmount uint256, tokenAmount uint256) returns()
func (_Escrow *EscrowTransactorSession) Deal(buyer common.Address, seller common.Address, weiAmount *big.Int, tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.Deal(&_Escrow.TransactOpts, buyer, seller, weiAmount, tokenAmount)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(account address) returns()
func (_Escrow *EscrowTransactor) RemoveAgent(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Escrow.contract.Transact(opts, "removeAgent", account)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(account address) returns()
func (_Escrow *EscrowSession) RemoveAgent(account common.Address) (*types.Transaction, error) {
	return _Escrow.Contract.RemoveAgent(&_Escrow.TransactOpts, account)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(account address) returns()
func (_Escrow *EscrowTransactorSession) RemoveAgent(account common.Address) (*types.Transaction, error) {
	return _Escrow.Contract.RemoveAgent(&_Escrow.TransactOpts, account)
}

// RenounceAgent is a paid mutator transaction binding the contract method 0xc692f4cf.
//
// Solidity: function renounceAgent() returns()
func (_Escrow *EscrowTransactor) RenounceAgent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Escrow.contract.Transact(opts, "renounceAgent")
}

// RenounceAgent is a paid mutator transaction binding the contract method 0xc692f4cf.
//
// Solidity: function renounceAgent() returns()
func (_Escrow *EscrowSession) RenounceAgent() (*types.Transaction, error) {
	return _Escrow.Contract.RenounceAgent(&_Escrow.TransactOpts)
}

// RenounceAgent is a paid mutator transaction binding the contract method 0xc692f4cf.
//
// Solidity: function renounceAgent() returns()
func (_Escrow *EscrowTransactorSession) RenounceAgent() (*types.Transaction, error) {
	return _Escrow.Contract.RenounceAgent(&_Escrow.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Escrow *EscrowTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Escrow.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Escrow *EscrowSession) RenounceOwnership() (*types.Transaction, error) {
	return _Escrow.Contract.RenounceOwnership(&_Escrow.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Escrow *EscrowTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Escrow.Contract.RenounceOwnership(&_Escrow.TransactOpts)
}

// Sell is a paid mutator transaction binding the contract method 0xe4849b32.
//
// Solidity: function sell(tokenAmount uint256) returns()
func (_Escrow *EscrowTransactor) Sell(opts *bind.TransactOpts, tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.contract.Transact(opts, "sell", tokenAmount)
}

// Sell is a paid mutator transaction binding the contract method 0xe4849b32.
//
// Solidity: function sell(tokenAmount uint256) returns()
func (_Escrow *EscrowSession) Sell(tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.Sell(&_Escrow.TransactOpts, tokenAmount)
}

// Sell is a paid mutator transaction binding the contract method 0xe4849b32.
//
// Solidity: function sell(tokenAmount uint256) returns()
func (_Escrow *EscrowTransactorSession) Sell(tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.Sell(&_Escrow.TransactOpts, tokenAmount)
}

// SellFrom is a paid mutator transaction binding the contract method 0xeb1715cc.
//
// Solidity: function sellFrom(seller address, tokenAmount uint256) returns()
func (_Escrow *EscrowTransactor) SellFrom(opts *bind.TransactOpts, seller common.Address, tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.contract.Transact(opts, "sellFrom", seller, tokenAmount)
}

// SellFrom is a paid mutator transaction binding the contract method 0xeb1715cc.
//
// Solidity: function sellFrom(seller address, tokenAmount uint256) returns()
func (_Escrow *EscrowSession) SellFrom(seller common.Address, tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.SellFrom(&_Escrow.TransactOpts, seller, tokenAmount)
}

// SellFrom is a paid mutator transaction binding the contract method 0xeb1715cc.
//
// Solidity: function sellFrom(seller address, tokenAmount uint256) returns()
func (_Escrow *EscrowTransactorSession) SellFrom(seller common.Address, tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.SellFrom(&_Escrow.TransactOpts, seller, tokenAmount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Escrow *EscrowTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Escrow.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Escrow *EscrowSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Escrow.Contract.TransferOwnership(&_Escrow.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Escrow *EscrowTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Escrow.Contract.TransferOwnership(&_Escrow.TransactOpts, newOwner)
}

// WithdrawAsk is a paid mutator transaction binding the contract method 0x3d603698.
//
// Solidity: function withdrawAsk(seller address, tokenAmount uint256) returns()
func (_Escrow *EscrowTransactor) WithdrawAsk(opts *bind.TransactOpts, seller common.Address, tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.contract.Transact(opts, "withdrawAsk", seller, tokenAmount)
}

// WithdrawAsk is a paid mutator transaction binding the contract method 0x3d603698.
//
// Solidity: function withdrawAsk(seller address, tokenAmount uint256) returns()
func (_Escrow *EscrowSession) WithdrawAsk(seller common.Address, tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.WithdrawAsk(&_Escrow.TransactOpts, seller, tokenAmount)
}

// WithdrawAsk is a paid mutator transaction binding the contract method 0x3d603698.
//
// Solidity: function withdrawAsk(seller address, tokenAmount uint256) returns()
func (_Escrow *EscrowTransactorSession) WithdrawAsk(seller common.Address, tokenAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.WithdrawAsk(&_Escrow.TransactOpts, seller, tokenAmount)
}

// WithdrawBid is a paid mutator transaction binding the contract method 0x4cb8ef5b.
//
// Solidity: function withdrawBid(buyer address, weiAmount uint256) returns()
func (_Escrow *EscrowTransactor) WithdrawBid(opts *bind.TransactOpts, buyer common.Address, weiAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.contract.Transact(opts, "withdrawBid", buyer, weiAmount)
}

// WithdrawBid is a paid mutator transaction binding the contract method 0x4cb8ef5b.
//
// Solidity: function withdrawBid(buyer address, weiAmount uint256) returns()
func (_Escrow *EscrowSession) WithdrawBid(buyer common.Address, weiAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.WithdrawBid(&_Escrow.TransactOpts, buyer, weiAmount)
}

// WithdrawBid is a paid mutator transaction binding the contract method 0x4cb8ef5b.
//
// Solidity: function withdrawBid(buyer address, weiAmount uint256) returns()
func (_Escrow *EscrowTransactorSession) WithdrawBid(buyer common.Address, weiAmount *big.Int) (*types.Transaction, error) {
	return _Escrow.Contract.WithdrawBid(&_Escrow.TransactOpts, buyer, weiAmount)
}

// EscrowAgentAddedIterator is returned from FilterAgentAdded and is used to iterate over the raw logs and unpacked data for AgentAdded events raised by the Escrow contract.
type EscrowAgentAddedIterator struct {
	Event *EscrowAgentAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EscrowAgentAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EscrowAgentAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EscrowAgentAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EscrowAgentAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EscrowAgentAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EscrowAgentAdded represents a AgentAdded event raised by the Escrow contract.
type EscrowAgentAdded struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAdded is a free log retrieval operation binding the contract event 0xf68e73cec97f2d70aa641fb26e87a4383686e2efacb648f2165aeb02ac562ec5.
//
// Solidity: e AgentAdded(account indexed address)
func (_Escrow *EscrowFilterer) FilterAgentAdded(opts *bind.FilterOpts, account []common.Address) (*EscrowAgentAddedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Escrow.contract.FilterLogs(opts, "AgentAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return &EscrowAgentAddedIterator{contract: _Escrow.contract, event: "AgentAdded", logs: logs, sub: sub}, nil
}

// WatchAgentAdded is a free log subscription operation binding the contract event 0xf68e73cec97f2d70aa641fb26e87a4383686e2efacb648f2165aeb02ac562ec5.
//
// Solidity: e AgentAdded(account indexed address)
func (_Escrow *EscrowFilterer) WatchAgentAdded(opts *bind.WatchOpts, sink chan<- *EscrowAgentAdded, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Escrow.contract.WatchLogs(opts, "AgentAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EscrowAgentAdded)
				if err := _Escrow.contract.UnpackLog(event, "AgentAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// EscrowAgentRemovedIterator is returned from FilterAgentRemoved and is used to iterate over the raw logs and unpacked data for AgentRemoved events raised by the Escrow contract.
type EscrowAgentRemovedIterator struct {
	Event *EscrowAgentRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EscrowAgentRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EscrowAgentRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EscrowAgentRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EscrowAgentRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EscrowAgentRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EscrowAgentRemoved represents a AgentRemoved event raised by the Escrow contract.
type EscrowAgentRemoved struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentRemoved is a free log retrieval operation binding the contract event 0xed9c8ad8d5a0a66898ea49d2956929c93ae2e8bd50281b2ed897c5d1a6737e0b.
//
// Solidity: e AgentRemoved(account indexed address)
func (_Escrow *EscrowFilterer) FilterAgentRemoved(opts *bind.FilterOpts, account []common.Address) (*EscrowAgentRemovedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Escrow.contract.FilterLogs(opts, "AgentRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return &EscrowAgentRemovedIterator{contract: _Escrow.contract, event: "AgentRemoved", logs: logs, sub: sub}, nil
}

// WatchAgentRemoved is a free log subscription operation binding the contract event 0xed9c8ad8d5a0a66898ea49d2956929c93ae2e8bd50281b2ed897c5d1a6737e0b.
//
// Solidity: e AgentRemoved(account indexed address)
func (_Escrow *EscrowFilterer) WatchAgentRemoved(opts *bind.WatchOpts, sink chan<- *EscrowAgentRemoved, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Escrow.contract.WatchLogs(opts, "AgentRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EscrowAgentRemoved)
				if err := _Escrow.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// EscrowCheckoutIterator is returned from FilterCheckout and is used to iterate over the raw logs and unpacked data for Checkout events raised by the Escrow contract.
type EscrowCheckoutIterator struct {
	Event *EscrowCheckout // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EscrowCheckoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EscrowCheckout)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EscrowCheckout)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EscrowCheckoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EscrowCheckoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EscrowCheckout represents a Checkout event raised by the Escrow contract.
type EscrowCheckout struct {
	Wallet common.Address
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCheckout is a free log retrieval operation binding the contract event 0x10190e7f174afbe5482194d509017e9df854671cd62fb4ef901fbd0080c1e6b6.
//
// Solidity: e Checkout(wallet indexed address, value uint256)
func (_Escrow *EscrowFilterer) FilterCheckout(opts *bind.FilterOpts, wallet []common.Address) (*EscrowCheckoutIterator, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}

	logs, sub, err := _Escrow.contract.FilterLogs(opts, "Checkout", walletRule)
	if err != nil {
		return nil, err
	}
	return &EscrowCheckoutIterator{contract: _Escrow.contract, event: "Checkout", logs: logs, sub: sub}, nil
}

// WatchCheckout is a free log subscription operation binding the contract event 0x10190e7f174afbe5482194d509017e9df854671cd62fb4ef901fbd0080c1e6b6.
//
// Solidity: e Checkout(wallet indexed address, value uint256)
func (_Escrow *EscrowFilterer) WatchCheckout(opts *bind.WatchOpts, sink chan<- *EscrowCheckout, wallet []common.Address) (event.Subscription, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}

	logs, sub, err := _Escrow.contract.WatchLogs(opts, "Checkout", walletRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EscrowCheckout)
				if err := _Escrow.contract.UnpackLog(event, "Checkout", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// EscrowCheckoutTokenIterator is returned from FilterCheckoutToken and is used to iterate over the raw logs and unpacked data for CheckoutToken events raised by the Escrow contract.
type EscrowCheckoutTokenIterator struct {
	Event *EscrowCheckoutToken // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EscrowCheckoutTokenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EscrowCheckoutToken)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EscrowCheckoutToken)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EscrowCheckoutTokenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EscrowCheckoutTokenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EscrowCheckoutToken represents a CheckoutToken event raised by the Escrow contract.
type EscrowCheckoutToken struct {
	Wallet common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCheckoutToken is a free log retrieval operation binding the contract event 0x1027be3db639879c20e5af001c80ef4b0621b3a7276e4668345ac6fb9ce5be58.
//
// Solidity: e CheckoutToken(wallet indexed address, amount uint256)
func (_Escrow *EscrowFilterer) FilterCheckoutToken(opts *bind.FilterOpts, wallet []common.Address) (*EscrowCheckoutTokenIterator, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}

	logs, sub, err := _Escrow.contract.FilterLogs(opts, "CheckoutToken", walletRule)
	if err != nil {
		return nil, err
	}
	return &EscrowCheckoutTokenIterator{contract: _Escrow.contract, event: "CheckoutToken", logs: logs, sub: sub}, nil
}

// WatchCheckoutToken is a free log subscription operation binding the contract event 0x1027be3db639879c20e5af001c80ef4b0621b3a7276e4668345ac6fb9ce5be58.
//
// Solidity: e CheckoutToken(wallet indexed address, amount uint256)
func (_Escrow *EscrowFilterer) WatchCheckoutToken(opts *bind.WatchOpts, sink chan<- *EscrowCheckoutToken, wallet []common.Address) (event.Subscription, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}

	logs, sub, err := _Escrow.contract.WatchLogs(opts, "CheckoutToken", walletRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EscrowCheckoutToken)
				if err := _Escrow.contract.UnpackLog(event, "CheckoutToken", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// EscrowOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the Escrow contract.
type EscrowOwnershipRenouncedIterator struct {
	Event *EscrowOwnershipRenounced // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EscrowOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EscrowOwnershipRenounced)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EscrowOwnershipRenounced)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EscrowOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EscrowOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EscrowOwnershipRenounced represents a OwnershipRenounced event raised by the Escrow contract.
type EscrowOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Escrow *EscrowFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*EscrowOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Escrow.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &EscrowOwnershipRenouncedIterator{contract: _Escrow.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Escrow *EscrowFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *EscrowOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Escrow.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EscrowOwnershipRenounced)
				if err := _Escrow.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// EscrowOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Escrow contract.
type EscrowOwnershipTransferredIterator struct {
	Event *EscrowOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EscrowOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EscrowOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EscrowOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EscrowOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EscrowOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EscrowOwnershipTransferred represents a OwnershipTransferred event raised by the Escrow contract.
type EscrowOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Escrow *EscrowFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*EscrowOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Escrow.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &EscrowOwnershipTransferredIterator{contract: _Escrow.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Escrow *EscrowFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EscrowOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Escrow.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EscrowOwnershipTransferred)
				if err := _Escrow.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// EscrowTokenAskIterator is returned from FilterTokenAsk and is used to iterate over the raw logs and unpacked data for TokenAsk events raised by the Escrow contract.
type EscrowTokenAskIterator struct {
	Event *EscrowTokenAsk // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EscrowTokenAskIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EscrowTokenAsk)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EscrowTokenAsk)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EscrowTokenAskIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EscrowTokenAskIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EscrowTokenAsk represents a TokenAsk event raised by the Escrow contract.
type EscrowTokenAsk struct {
	Seller common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokenAsk is a free log retrieval operation binding the contract event 0x19bc6719e5b4b4a57c924567ee3c27073d21148f2801228976dc6361fd8190d6.
//
// Solidity: e TokenAsk(seller indexed address, amount uint256)
func (_Escrow *EscrowFilterer) FilterTokenAsk(opts *bind.FilterOpts, seller []common.Address) (*EscrowTokenAskIterator, error) {

	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _Escrow.contract.FilterLogs(opts, "TokenAsk", sellerRule)
	if err != nil {
		return nil, err
	}
	return &EscrowTokenAskIterator{contract: _Escrow.contract, event: "TokenAsk", logs: logs, sub: sub}, nil
}

// WatchTokenAsk is a free log subscription operation binding the contract event 0x19bc6719e5b4b4a57c924567ee3c27073d21148f2801228976dc6361fd8190d6.
//
// Solidity: e TokenAsk(seller indexed address, amount uint256)
func (_Escrow *EscrowFilterer) WatchTokenAsk(opts *bind.WatchOpts, sink chan<- *EscrowTokenAsk, seller []common.Address) (event.Subscription, error) {

	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _Escrow.contract.WatchLogs(opts, "TokenAsk", sellerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EscrowTokenAsk)
				if err := _Escrow.contract.UnpackLog(event, "TokenAsk", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// EscrowTokenBidIterator is returned from FilterTokenBid and is used to iterate over the raw logs and unpacked data for TokenBid events raised by the Escrow contract.
type EscrowTokenBidIterator struct {
	Event *EscrowTokenBid // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EscrowTokenBidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EscrowTokenBid)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EscrowTokenBid)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EscrowTokenBidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EscrowTokenBidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EscrowTokenBid represents a TokenBid event raised by the Escrow contract.
type EscrowTokenBid struct {
	Buyer common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTokenBid is a free log retrieval operation binding the contract event 0xd53242b5a0280628acfd89744b78a30e567c53d7c69397838d209f48dfe542a4.
//
// Solidity: e TokenBid(buyer indexed address, value uint256)
func (_Escrow *EscrowFilterer) FilterTokenBid(opts *bind.FilterOpts, buyer []common.Address) (*EscrowTokenBidIterator, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _Escrow.contract.FilterLogs(opts, "TokenBid", buyerRule)
	if err != nil {
		return nil, err
	}
	return &EscrowTokenBidIterator{contract: _Escrow.contract, event: "TokenBid", logs: logs, sub: sub}, nil
}

// WatchTokenBid is a free log subscription operation binding the contract event 0xd53242b5a0280628acfd89744b78a30e567c53d7c69397838d209f48dfe542a4.
//
// Solidity: e TokenBid(buyer indexed address, value uint256)
func (_Escrow *EscrowFilterer) WatchTokenBid(opts *bind.WatchOpts, sink chan<- *EscrowTokenBid, buyer []common.Address) (event.Subscription, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _Escrow.contract.WatchLogs(opts, "TokenBid", buyerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EscrowTokenBid)
				if err := _Escrow.contract.UnpackLog(event, "TokenBid", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// EscrowTokenDealIterator is returned from FilterTokenDeal and is used to iterate over the raw logs and unpacked data for TokenDeal events raised by the Escrow contract.
type EscrowTokenDealIterator struct {
	Event *EscrowTokenDeal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EscrowTokenDealIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EscrowTokenDeal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EscrowTokenDeal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EscrowTokenDealIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EscrowTokenDealIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EscrowTokenDeal represents a TokenDeal event raised by the Escrow contract.
type EscrowTokenDeal struct {
	Buyer  common.Address
	Seller common.Address
	Amount *big.Int
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokenDeal is a free log retrieval operation binding the contract event 0x13b25d4eb0e3438ebba0fb869531f14c03b78b4c3e286a47cd3fd8627d424165.
//
// Solidity: e TokenDeal(buyer indexed address, seller indexed address, amount uint256, value uint256)
func (_Escrow *EscrowFilterer) FilterTokenDeal(opts *bind.FilterOpts, buyer []common.Address, seller []common.Address) (*EscrowTokenDealIterator, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _Escrow.contract.FilterLogs(opts, "TokenDeal", buyerRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return &EscrowTokenDealIterator{contract: _Escrow.contract, event: "TokenDeal", logs: logs, sub: sub}, nil
}

// WatchTokenDeal is a free log subscription operation binding the contract event 0x13b25d4eb0e3438ebba0fb869531f14c03b78b4c3e286a47cd3fd8627d424165.
//
// Solidity: e TokenDeal(buyer indexed address, seller indexed address, amount uint256, value uint256)
func (_Escrow *EscrowFilterer) WatchTokenDeal(opts *bind.WatchOpts, sink chan<- *EscrowTokenDeal, buyer []common.Address, seller []common.Address) (event.Subscription, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _Escrow.contract.WatchLogs(opts, "TokenDeal", buyerRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EscrowTokenDeal)
				if err := _Escrow.contract.UnpackLog(event, "TokenDeal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// EscrowWithdrawnAskIterator is returned from FilterWithdrawnAsk and is used to iterate over the raw logs and unpacked data for WithdrawnAsk events raised by the Escrow contract.
type EscrowWithdrawnAskIterator struct {
	Event *EscrowWithdrawnAsk // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EscrowWithdrawnAskIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EscrowWithdrawnAsk)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EscrowWithdrawnAsk)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EscrowWithdrawnAskIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EscrowWithdrawnAskIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EscrowWithdrawnAsk represents a WithdrawnAsk event raised by the Escrow contract.
type EscrowWithdrawnAsk struct {
	Seller      common.Address
	TokenAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterWithdrawnAsk is a free log retrieval operation binding the contract event 0xc80bf1c92f6420fc18665ad14be6e5d7cc029baaf8fc547a6b2799a0f3bd012b.
//
// Solidity: e WithdrawnAsk(seller indexed address, tokenAmount uint256)
func (_Escrow *EscrowFilterer) FilterWithdrawnAsk(opts *bind.FilterOpts, seller []common.Address) (*EscrowWithdrawnAskIterator, error) {

	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _Escrow.contract.FilterLogs(opts, "WithdrawnAsk", sellerRule)
	if err != nil {
		return nil, err
	}
	return &EscrowWithdrawnAskIterator{contract: _Escrow.contract, event: "WithdrawnAsk", logs: logs, sub: sub}, nil
}

// WatchWithdrawnAsk is a free log subscription operation binding the contract event 0xc80bf1c92f6420fc18665ad14be6e5d7cc029baaf8fc547a6b2799a0f3bd012b.
//
// Solidity: e WithdrawnAsk(seller indexed address, tokenAmount uint256)
func (_Escrow *EscrowFilterer) WatchWithdrawnAsk(opts *bind.WatchOpts, sink chan<- *EscrowWithdrawnAsk, seller []common.Address) (event.Subscription, error) {

	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _Escrow.contract.WatchLogs(opts, "WithdrawnAsk", sellerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EscrowWithdrawnAsk)
				if err := _Escrow.contract.UnpackLog(event, "WithdrawnAsk", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// EscrowWithdrawnBidIterator is returned from FilterWithdrawnBid and is used to iterate over the raw logs and unpacked data for WithdrawnBid events raised by the Escrow contract.
type EscrowWithdrawnBidIterator struct {
	Event *EscrowWithdrawnBid // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EscrowWithdrawnBidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EscrowWithdrawnBid)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EscrowWithdrawnBid)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EscrowWithdrawnBidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EscrowWithdrawnBidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EscrowWithdrawnBid represents a WithdrawnBid event raised by the Escrow contract.
type EscrowWithdrawnBid struct {
	Buyer     common.Address
	WeiAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawnBid is a free log retrieval operation binding the contract event 0x064e09ad579edcb276d73f297cc24e5b0685e3b9574e9d1e5806e3d76cd949b8.
//
// Solidity: e WithdrawnBid(buyer indexed address, weiAmount uint256)
func (_Escrow *EscrowFilterer) FilterWithdrawnBid(opts *bind.FilterOpts, buyer []common.Address) (*EscrowWithdrawnBidIterator, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _Escrow.contract.FilterLogs(opts, "WithdrawnBid", buyerRule)
	if err != nil {
		return nil, err
	}
	return &EscrowWithdrawnBidIterator{contract: _Escrow.contract, event: "WithdrawnBid", logs: logs, sub: sub}, nil
}

// WatchWithdrawnBid is a free log subscription operation binding the contract event 0x064e09ad579edcb276d73f297cc24e5b0685e3b9574e9d1e5806e3d76cd949b8.
//
// Solidity: e WithdrawnBid(buyer indexed address, weiAmount uint256)
func (_Escrow *EscrowFilterer) WatchWithdrawnBid(opts *bind.WatchOpts, sink chan<- *EscrowWithdrawnBid, buyer []common.Address) (event.Subscription, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _Escrow.contract.WatchLogs(opts, "WithdrawnBid", buyerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EscrowWithdrawnBid)
				if err := _Escrow.contract.UnpackLog(event, "WithdrawnBid", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
