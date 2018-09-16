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

// TokenABI is the input ABI used to generate the binding from.
const TokenABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowanceProxy\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"mintingFinished\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isAgent\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approveProxy\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"removeMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"initialSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalTransfers\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isPauser\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferProxy\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"addresses\",\"type\":\"address[]\"},{\"name\":\"tokenAmount\",\"type\":\"uint256[]\"}],\"name\":\"batchTransferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalHolders\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"removePauser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renouncePauser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finishMinting\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowanceProxy\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addPauser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addAgent\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addresses\",\"type\":\"address[]\"},{\"name\":\"tokenAmount\",\"type\":\"uint256[]\"}],\"name\":\"batchTransfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"circulatingSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"removeAgent\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isMinter\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceAgent\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"meta\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"uint8\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\"},{\"name\":\"decimals\",\"type\":\"uint8\"},{\"name\":\"initialSupply\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AgentAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AgentRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"PauserAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"PauserRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"MintingFinished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// TokenBin is the compiled bytecode used for deploying new contracts.
const TokenBin = `0x60806040526008805460ff19908116909155600a805490911690553480156200002757600080fd5b5060405162002419380380620024198339810160409081528151602083015191830151606084015160008054600160a060020a0319163390811790915592850194939093019290916200008b9060079064010000000062001d366200015782021704565b620000a660093364010000000062001d366200015782021704565b620000c1600b3364010000000062001d366200015782021704565b8351620000d690600c90602087019062000192565b508251620000ec90600d90602086019062000192565b50600e805460ff191660ff841617905560068190556003819055336000818152600160209081526040808320859055805185815290517fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a35050505062000237565b600160a060020a03811615156200016d57600080fd5b600160a060020a0316600090815260209190915260409020805460ff19166001179055565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620001d557805160ff191683800117855562000205565b8280016001018555821562000205579182015b8281111562000205578251825591602001919060010190620001e8565b506200021392915062000217565b5090565b6200023491905b808211156200021357600081556001016200021e565b90565b6121d280620002476000396000f30060806040526004361061022f5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166302a4e488811461023457806305d2035b1461027257806306fdde0314610287578063095ea7b31461031157806318160ddd146103355780631ffbb0641461035c5780632376fe701461037d57806323b872dd146103a75780633092afd5146103d1578063313ce567146103f4578063378dc3dc1461041f57806339509351146104345780633bf47720146104585780633f4ba83a1461046d57806340c10f191461048257806342966c68146104a657806346fbf68e146104be5780634733dc8f146104df5780634885b2541461050957806353d74fdf146105a55780635c975abb146105ba5780636b2c0f55146105cf5780636ef8d66d146105f057806370a0823114610605578063715018a61461062657806379cc67901461063b5780637d64bcb41461065f57806382bcef791461067457806382dc1ec41461069e5780638456cb59146106bf57806384e79842146106d457806388d695b2146106f55780638da5cb5b146107835780638f32d59b146107b45780639358928b146107c957806395d89b41146107de57806397a6278e146107f3578063983b2d56146108145780639865027514610835578063a457c2d71461084a578063a9059cbb1461086e578063aa271e1a14610892578063c692f4cf146108b3578063dd62ed3e146108c8578063e021deff146108ef578063f2fde38b14610a25575b600080fd5b34801561024057600080fd5b5061025e600160a060020a0360043581169060243516604435610a46565b604080519115158252519081900360200190f35b34801561027e57600080fd5b5061025e610b0f565b34801561029357600080fd5b5061029c610b18565b6040805160208082528351818301528351919283929083019185019080838360005b838110156102d65781810151838201526020016102be565b50505050905090810190601f1680156103035780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561031d57600080fd5b5061025e600160a060020a0360043516602435610bae565b34801561034157600080fd5b5061034a610bcb565b60408051918252519081900360200190f35b34801561036857600080fd5b5061025e600160a060020a0360043516610bd1565b34801561038957600080fd5b5061025e600160a060020a0360043581169060243516604435610bea565b3480156103b357600080fd5b5061025e600160a060020a0360043581169060243516604435610c78565b3480156103dd57600080fd5b506103f2600160a060020a0360043516610c9e565b005b34801561040057600080fd5b50610409610ccc565b6040805160ff9092168252519081900360200190f35b34801561042b57600080fd5b5061034a610cd5565b34801561044057600080fd5b5061025e600160a060020a0360043516602435610cdb565b34801561046457600080fd5b5061034a610cf8565b34801561047957600080fd5b506103f2610cfe565b34801561048e57600080fd5b5061025e600160a060020a0360043516602435610d66565b3480156104b257600080fd5b506103f2600435610dad565b3480156104ca57600080fd5b5061025e600160a060020a0360043516610db7565b3480156104eb57600080fd5b5061025e600160a060020a0360043581169060243516604435610dca565b34801561051557600080fd5b5060408051602060046024803582810135848102808701860190975280865261025e968435600160a060020a031696369660449591949091019291829185019084908082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750949750610f899650505050505050565b3480156105b157600080fd5b5061034a611012565b3480156105c657600080fd5b5061025e611018565b3480156105db57600080fd5b506103f2600160a060020a0360043516611021565b3480156105fc57600080fd5b506103f261104c565b34801561061157600080fd5b5061034a600160a060020a0360043516611081565b34801561063257600080fd5b506103f261109c565b34801561064757600080fd5b506103f2600160a060020a0360043516602435611104565b34801561066b57600080fd5b5061025e611112565b34801561068057600080fd5b5061025e600160a060020a0360043581169060243516604435611182565b3480156106aa57600080fd5b506103f2600160a060020a03600435166111f1565b3480156106cb57600080fd5b506103f261125b565b3480156106e057600080fd5b506103f2600160a060020a03600435166112c5565b34801561070157600080fd5b506040805160206004803580820135838102808601850190965280855261025e95369593946024949385019291829185019084908082843750506040805187358901803560208181028481018201909552818452989b9a99890198929750908201955093508392508501908490808284375094975061132f9650505050505050565b34801561078f57600080fd5b506107986113b6565b60408051600160a060020a039092168252519081900360200190f35b3480156107c057600080fd5b5061025e6113c5565b3480156107d557600080fd5b5061034a6113d6565b3480156107ea57600080fd5b5061029c611449565b3480156107ff57600080fd5b506103f2600160a060020a03600435166114aa565b34801561082057600080fd5b506103f2600160a060020a03600435166114d5565b34801561084157600080fd5b506103f261153f565b34801561085657600080fd5b5061025e600160a060020a0360043516602435611572565b34801561087a57600080fd5b5061025e600160a060020a036004351660243561158f565b34801561089e57600080fd5b5061025e600160a060020a03600435166115ac565b3480156108bf57600080fd5b506103f26115bf565b3480156108d457600080fd5b5061034a600160a060020a03600435811690602435166115f2565b3480156108fb57600080fd5b50610910600160a060020a036004351661161d565b6040518080602001806020018a60ff1660ff16815260200189815260200188815260200187815260200186815260200185815260200184815260200183810383528c818151815260200191508051906020019080838360005b83811015610981578181015183820152602001610969565b50505050905090810190601f1680156109ae5780820380516001836020036101000a031916815260200191505b5083810382528b5181528b516020918201918d019080838360005b838110156109e15781810151838201526020016109c9565b50505050905090810190601f168015610a0e5780820380516001836020036101000a031916815260200191505b509b50505050505050505050505060405180910390f35b348015610a3157600080fd5b506103f2600160a060020a0360043516611827565b6000610a506113c5565b80610a5f5750610a5f33610bd1565b1515610a6a57600080fd5b600160a060020a0383161515610a7f57600080fd5b600160a060020a03808516600090815260026020908152604080832093871683529290522054610ab5908363ffffffff61184316565b600160a060020a0385811660008181526002602090815260408083209489168084529482529182902085905581519485529051929391926000805160206121878339815191529281900390910190a35060015b9392505050565b60085460ff1690565b600c8054604080516020601f6002600019610100600188161502019095169490940493840181900481028201810190925282815260609390929091830182828015610ba45780601f10610b7957610100808354040283529160200191610ba4565b820191906000526020600020905b815481529060010190602001808311610b8757829003601f168201915b5050505050905090565b600a5460009060ff1615610bc157600080fd5b610b08838361185a565b60035490565b6000610be4600b8363ffffffff6118c616565b92915050565b6000610bf46113c5565b80610c035750610c0333610bd1565b1515610c0e57600080fd5b600160a060020a0383161515610c2357600080fd5b600160a060020a03808516600081815260026020908152604080832094881680845294825291829020869055815186815291516000805160206121878339815191529281900390910190a35060019392505050565b600a5460009060ff1615610c8b57600080fd5b610c968484846118fd565b949350505050565b610ca66113c5565b80610cb55750610cb5336115ac565b1515610cc057600080fd5b610cc981611ad0565b50565b600e5460ff1690565b60065490565b600a5460009060ff1615610cee57600080fd5b610b088383611ad9565b60055490565b610d066113c5565b80610d155750610d1533610db7565b1515610d2057600080fd5b600a5460ff161515610d3157600080fd5b600a805460ff191690556040517fa45f47fdea8a1efdd9029a5691c7f759c32b7c698632b563573e155625d1693390600090a1565b6000610d706113c5565b80610d7f5750610d7f336115ac565b1515610d8a57600080fd5b60085460ff1615610d9a57600080fd5b610da48383611b77565b50600192915050565b610cc93382611c48565b6000610be460098363ffffffff6118c616565b6000610dd46113c5565b80610de35750610de333610bd1565b1515610dee57600080fd5b33600160a060020a0385161415610e1057610e09838361158f565b9050610b08565b600160a060020a038416600090815260016020526040902054821115610e3557600080fd5b600160a060020a0383161515610e4a57600080fd5b600160a060020a038416600090815260016020526040902054610e73908363ffffffff61184316565b600160a060020a0385166000908152600160205260409020819055158015610e9d57506000600454115b15610eba57600454610eb690600163ffffffff61184316565b6004555b600160a060020a0383166000908152600160205260409020541515610ef157600454610eed90600163ffffffff611c5216565b6004555b600160a060020a038316600090815260016020526040902054610f1a908363ffffffff611c5216565b600160a060020a038416600090815260016020819052604090912091909155600554610f4b9163ffffffff611c5216565b600555604080518381529051600160a060020a0380861692908716916000805160206121678339815191529181900360200190a35060019392505050565b60008060008060008651118015610fa1575084518651145b1515610fac57600080fd5b600092505b8551831015611005578583815181101515610fc857fe5b9060200190602002015191508483815181101515610fe257fe5b906020019060200201519050610ff9878383610c78565b50600190920191610fb1565b5060019695505050505050565b60045490565b600a5460ff1690565b6110296113c5565b80611038575061103833610db7565b151561104357600080fd5b610cc981611c64565b6110546113c5565b80611063575061106333610db7565b151561106e57600080fd5b61107f60093363ffffffff611c6d16565b565b600160a060020a031660009081526001602052604090205490565b6110a46113c5565b15156110af57600080fd5b60008054604051600160a060020a03909116917ff8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c6482091a26000805473ffffffffffffffffffffffffffffffffffffffff19169055565b61110e8282611ca4565b5050565b600061111c6113c5565b8061112b575061112b336115ac565b151561113657600080fd5b60085460ff161561114657600080fd5b6008805460ff191660011790556040517fb828d9b5c78095deeeeff2eca2e5d4fe046ce3feb4c99702624a3fd384ad2dbc90600090a150600190565b600061118c6113c5565b8061119b575061119b33610bd1565b15156111a657600080fd5b600160a060020a03831615156111bb57600080fd5b600160a060020a03808516600090815260026020908152604080832093871683529290522054610ab5908363ffffffff611c5216565b6111f96113c5565b80611208575061120833610db7565b151561121357600080fd5b61122460098263ffffffff611d3616565b604051600160a060020a038216907f6719d08c1888103bea251a4ed56406bd0c3e69723c8a1686e017e7bbe159b6f890600090a250565b6112636113c5565b80611272575061127233610db7565b151561127d57600080fd5b600a5460ff161561128d57600080fd5b600a805460ff191660011790556040517f9e87fac88ff661f02d44f95383c817fece4bce600a3dab7a54406878b965e75290600090a1565b6112cd6113c5565b806112dc57506112dc33610bd1565b15156112e757600080fd5b6112f8600b8263ffffffff611d3616565b604051600160a060020a038216907ff68e73cec97f2d70aa641fb26e87a4383686e2efacb648f2165aeb02ac562ec590600090a250565b60008060008060008651118015611347575084518651145b151561135257600080fd5b600092505b85518310156113aa57858381518110151561136e57fe5b906020019060200201519150848381518110151561138857fe5b90602001906020020151905061139e828261158f565b50600190920191611357565b50600195945050505050565b600054600160a060020a031690565b600054600160a060020a0316331490565b6000600160006113e46113b6565b600160a060020a03168152602081019190915260400160002054600354101561140c57600080fd5b6114446001600061141b6113b6565b600160a060020a031681526020810191909152604001600020546003549063ffffffff61184316565b905090565b600d8054604080516020601f6002600019610100600188161502019095169490940493840181900481028201810190925282815260609390929091830182828015610ba45780601f10610b7957610100808354040283529160200191610ba4565b6114b26113c5565b806114c157506114c133610bd1565b15156114cc57600080fd5b610cc981611d70565b6114dd6113c5565b806114ec57506114ec336115ac565b15156114f757600080fd5b61150860078263ffffffff611d3616565b604051600160a060020a038216907f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f690600090a250565b6115476113c5565b806115565750611556336115ac565b151561156157600080fd5b61107f60073363ffffffff611c6d16565b600a5460009060ff161561158557600080fd5b610b088383611d79565b600a5460009060ff16156115a257600080fd5b610b088383611dc4565b6000610be460078363ffffffff6118c616565b6115c76113c5565b806115d657506115d633610bd1565b15156115e157600080fd5b61107f600b3363ffffffff611c6d16565b600160a060020a03918216600090815260026020908152604080832093909416825291909152205490565b606080600080808080808080806001816116356113b6565b600160a060020a0316600160a060020a031681526020019081526020016000205460035411156116715761166e6001600061141b6113b6565b91505b506000600160a060020a038c16156116a25750600160a060020a038b166000908152600160205260409020546116ba565b33156116ba5750336000908152600160205260409020545b600c600d600e60009054906101000a900460ff166006546003546005546004548888888054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156117715780601f1061174657610100808354040283529160200191611771565b820191906000526020600020905b81548152906001019060200180831161175457829003601f168201915b50508b5460408051602060026001851615610100026000190190941693909304601f8101849004840282018401909252818152959e508d9450925084019050828280156117ff5780601f106117d4576101008083540402835291602001916117ff565b820191906000526020600020905b8154815290600101906020018083116117e257829003601f168201915b505050505097509a509a509a509a509a509a509a509a509a5050509193959799909294969850565b61182f6113c5565b151561183a57600080fd5b610cc981611f1e565b6000808383111561185357600080fd5b5050900390565b6000600160a060020a038316151561187157600080fd5b336000818152600260209081526040808320600160a060020a0388168085529083529281902086905580518681529051929392600080516020612187833981519152929181900390910190a350600192915050565b6000600160a060020a03821615156118dd57600080fd5b50600160a060020a03166000908152602091909152604090205460ff1690565b600033600160a060020a038516141561191a57610e09838361158f565b600160a060020a03841660009081526001602052604090205482111561193f57600080fd5b600160a060020a038416600090815260026020908152604080832033845290915290205482111561196f57600080fd5b600160a060020a038316151561198457600080fd5b600160a060020a0384166000908152600160205260409020546119ad908363ffffffff61184316565b600160a060020a03851660009081526001602052604090208190551580156119d757506000600454115b156119f4576004546119f090600163ffffffff61184316565b6004555b600160a060020a0383166000908152600160205260409020541515611a2b57600454611a2790600163ffffffff611c5216565b6004555b600160a060020a038316600090815260016020526040902054611a54908363ffffffff611c5216565b600160a060020a038085166000908152600160209081526040808320949094559187168152600282528281203382529091522054611a98908363ffffffff61184316565b600160a060020a0385166000908152600260209081526040808320338452909152902055600554610f4b90600163ffffffff611c5216565b610cc981611f9b565b6000600160a060020a0383161515611af057600080fd5b336000908152600260209081526040808320600160a060020a0387168452909152902054611b24908363ffffffff611c5216565b336000818152600260209081526040808320600160a060020a038916808552908352928190208590558051948552519193600080516020612187833981519152929081900390910190a350600192915050565b600160a060020a0382161515611b8c57600080fd5b600354611b9f908263ffffffff611c5216565b600355600160a060020a0382166000908152600160205260409020541515611bd957600454611bd590600163ffffffff611c5216565b6004555b600160a060020a038216600090815260016020526040902054611c02908263ffffffff611c5216565b600160a060020a03831660008181526001602090815260408083209490945583518581529351929391926000805160206121678339815191529281900390910190a35050565b61110e8282611fe3565b600082820183811015610b0857600080fd5b610cc9816120d6565b600160a060020a0381161515611c8257600080fd5b600160a060020a0316600090815260209190915260409020805460ff19169055565b600160a060020a0382166000908152600260209081526040808320338452909152902054811115611cd457600080fd5b600160a060020a0382166000908152600260209081526040808320338452909152902054611d08908263ffffffff61184316565b600160a060020a038316600090815260026020908152604080832033845290915290205561110e8282611c48565b600160a060020a0381161515611d4b57600080fd5b600160a060020a0316600090815260209190915260409020805460ff19166001179055565b610cc98161211e565b6000600160a060020a0383161515611d9057600080fd5b336000908152600260209081526040808320600160a060020a0387168452909152902054611b24908363ffffffff61184316565b33600090815260016020526040812054821115611de057600080fd5b600160a060020a0383161515611df557600080fd5b33600090815260016020526040902054611e15908363ffffffff61184316565b336000908152600160205260409020819055158015611e3657506000600454115b15611e5357600454611e4f90600163ffffffff61184316565b6004555b600160a060020a0383166000908152600160205260409020541515611e8a57600454611e8690600163ffffffff611c5216565b6004555b600160a060020a038316600090815260016020526040902054611eb3908363ffffffff611c5216565b600160a060020a038416600090815260016020819052604090912091909155600554611ee49163ffffffff611c5216565b600555604080518381529051600160a060020a0385169133916000805160206121678339815191529181900360200190a350600192915050565b600160a060020a0381161515611f3357600080fd5b60008054604051600160a060020a03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b611fac60078263ffffffff611c6d16565b604051600160a060020a038216907fe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb6669290600090a250565b600160a060020a0382161515611ff857600080fd5b600160a060020a03821660009081526001602052604090205481111561201d57600080fd5b600354612030908263ffffffff61184316565b600355600160a060020a03821660009081526001602052604090205461205c908263ffffffff61184316565b600160a060020a038316600090815260016020526040902081905515801561208657506000600454115b156120a35760045461209f90600163ffffffff61184316565b6004555b604080518281529051600091600160a060020a038516916000805160206121678339815191529181900360200190a35050565b6120e760098263ffffffff611c6d16565b604051600160a060020a038216907fcd265ebaf09df2871cc7bd4133404a235ba12eff2041bb89d9c714a2621c7c7e90600090a250565b61212f600b8263ffffffff611c6d16565b604051600160a060020a038216907fed9c8ad8d5a0a66898ea49d2956929c93ae2e8bd50281b2ed897c5d1a6737e0b90600090a2505600ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925a165627a7a7230582046c77d9f8f3d60cde60bf945d3d9877cd65709a499c195c164c83c2c022498d30029`

// DeployToken deploys a new Ethereum contract, binding an instance of Token to it.
func DeployToken(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, decimals uint8, initialSupply *big.Int) (common.Address, *types.Transaction, *Token, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TokenBin), backend, name, symbol, decimals, initialSupply)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Token{TokenCaller: TokenCaller{contract: contract}, TokenTransactor: TokenTransactor{contract: contract}, TokenFilterer: TokenFilterer{contract: contract}}, nil
}

// Token is an auto generated Go binding around an Ethereum contract.
type Token struct {
	TokenCaller     // Read-only binding to the contract
	TokenTransactor // Write-only binding to the contract
	TokenFilterer   // Log filterer for contract events
}

// TokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenSession struct {
	Contract     *Token            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenCallerSession struct {
	Contract *TokenCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenTransactorSession struct {
	Contract     *TokenTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenRaw struct {
	Contract *Token // Generic contract binding to access the raw methods on
}

// TokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenCallerRaw struct {
	Contract *TokenCaller // Generic read-only contract binding to access the raw methods on
}

// TokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenTransactorRaw struct {
	Contract *TokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewToken creates a new instance of Token, bound to a specific deployed contract.
func NewToken(address common.Address, backend bind.ContractBackend) (*Token, error) {
	contract, err := bindToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Token{TokenCaller: TokenCaller{contract: contract}, TokenTransactor: TokenTransactor{contract: contract}, TokenFilterer: TokenFilterer{contract: contract}}, nil
}

// NewTokenCaller creates a new read-only instance of Token, bound to a specific deployed contract.
func NewTokenCaller(address common.Address, caller bind.ContractCaller) (*TokenCaller, error) {
	contract, err := bindToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCaller{contract: contract}, nil
}

// NewTokenTransactor creates a new write-only instance of Token, bound to a specific deployed contract.
func NewTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenTransactor, error) {
	contract, err := bindToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenTransactor{contract: contract}, nil
}

// NewTokenFilterer creates a new log filterer instance of Token, bound to a specific deployed contract.
func NewTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenFilterer, error) {
	contract, err := bindToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenFilterer{contract: contract}, nil
}

// bindToken binds a generic wrapper to an already deployed contract.
func bindToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token *TokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Token.Contract.TokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token *TokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.Contract.TokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token *TokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token.Contract.TokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token *TokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Token.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token *TokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token *TokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_Token *TokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_Token *TokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Token.Contract.Allowance(&_Token.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_Token *TokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Token.Contract.Allowance(&_Token.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_Token *TokenCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "balanceOf", owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_Token *TokenSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _Token.Contract.BalanceOf(&_Token.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_Token *TokenCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _Token.Contract.BalanceOf(&_Token.CallOpts, owner)
}

// CirculatingSupply is a free data retrieval call binding the contract method 0x9358928b.
//
// Solidity: function circulatingSupply() constant returns(uint256)
func (_Token *TokenCaller) CirculatingSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "circulatingSupply")
	return *ret0, err
}

// CirculatingSupply is a free data retrieval call binding the contract method 0x9358928b.
//
// Solidity: function circulatingSupply() constant returns(uint256)
func (_Token *TokenSession) CirculatingSupply() (*big.Int, error) {
	return _Token.Contract.CirculatingSupply(&_Token.CallOpts)
}

// CirculatingSupply is a free data retrieval call binding the contract method 0x9358928b.
//
// Solidity: function circulatingSupply() constant returns(uint256)
func (_Token *TokenCallerSession) CirculatingSupply() (*big.Int, error) {
	return _Token.Contract.CirculatingSupply(&_Token.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_Token *TokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_Token *TokenSession) Decimals() (uint8, error) {
	return _Token.Contract.Decimals(&_Token.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_Token *TokenCallerSession) Decimals() (uint8, error) {
	return _Token.Contract.Decimals(&_Token.CallOpts)
}

// InitialSupply is a free data retrieval call binding the contract method 0x378dc3dc.
//
// Solidity: function initialSupply() constant returns(uint256)
func (_Token *TokenCaller) InitialSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "initialSupply")
	return *ret0, err
}

// InitialSupply is a free data retrieval call binding the contract method 0x378dc3dc.
//
// Solidity: function initialSupply() constant returns(uint256)
func (_Token *TokenSession) InitialSupply() (*big.Int, error) {
	return _Token.Contract.InitialSupply(&_Token.CallOpts)
}

// InitialSupply is a free data retrieval call binding the contract method 0x378dc3dc.
//
// Solidity: function initialSupply() constant returns(uint256)
func (_Token *TokenCallerSession) InitialSupply() (*big.Int, error) {
	return _Token.Contract.InitialSupply(&_Token.CallOpts)
}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(account address) constant returns(bool)
func (_Token *TokenCaller) IsAgent(opts *bind.CallOpts, account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "isAgent", account)
	return *ret0, err
}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(account address) constant returns(bool)
func (_Token *TokenSession) IsAgent(account common.Address) (bool, error) {
	return _Token.Contract.IsAgent(&_Token.CallOpts, account)
}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(account address) constant returns(bool)
func (_Token *TokenCallerSession) IsAgent(account common.Address) (bool, error) {
	return _Token.Contract.IsAgent(&_Token.CallOpts, account)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(account address) constant returns(bool)
func (_Token *TokenCaller) IsMinter(opts *bind.CallOpts, account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "isMinter", account)
	return *ret0, err
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(account address) constant returns(bool)
func (_Token *TokenSession) IsMinter(account common.Address) (bool, error) {
	return _Token.Contract.IsMinter(&_Token.CallOpts, account)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(account address) constant returns(bool)
func (_Token *TokenCallerSession) IsMinter(account common.Address) (bool, error) {
	return _Token.Contract.IsMinter(&_Token.CallOpts, account)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Token *TokenCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Token *TokenSession) IsOwner() (bool, error) {
	return _Token.Contract.IsOwner(&_Token.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Token *TokenCallerSession) IsOwner() (bool, error) {
	return _Token.Contract.IsOwner(&_Token.CallOpts)
}

// IsPauser is a free data retrieval call binding the contract method 0x46fbf68e.
//
// Solidity: function isPauser(account address) constant returns(bool)
func (_Token *TokenCaller) IsPauser(opts *bind.CallOpts, account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "isPauser", account)
	return *ret0, err
}

// IsPauser is a free data retrieval call binding the contract method 0x46fbf68e.
//
// Solidity: function isPauser(account address) constant returns(bool)
func (_Token *TokenSession) IsPauser(account common.Address) (bool, error) {
	return _Token.Contract.IsPauser(&_Token.CallOpts, account)
}

// IsPauser is a free data retrieval call binding the contract method 0x46fbf68e.
//
// Solidity: function isPauser(account address) constant returns(bool)
func (_Token *TokenCallerSession) IsPauser(account common.Address) (bool, error) {
	return _Token.Contract.IsPauser(&_Token.CallOpts, account)
}

// Meta is a free data retrieval call binding the contract method 0xe021deff.
//
// Solidity: function meta(account address) constant returns(string, string, uint8, uint256, uint256, uint256, uint256, uint256, uint256)
func (_Token *TokenCaller) Meta(opts *bind.CallOpts, account common.Address) (string, string, uint8, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	var (
		ret0 = new(string)
		ret1 = new(string)
		ret2 = new(uint8)
		ret3 = new(*big.Int)
		ret4 = new(*big.Int)
		ret5 = new(*big.Int)
		ret6 = new(*big.Int)
		ret7 = new(*big.Int)
		ret8 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
		ret5,
		ret6,
		ret7,
		ret8,
	}
	err := _Token.contract.Call(opts, out, "meta", account)
	return *ret0, *ret1, *ret2, *ret3, *ret4, *ret5, *ret6, *ret7, *ret8, err
}

// Meta is a free data retrieval call binding the contract method 0xe021deff.
//
// Solidity: function meta(account address) constant returns(string, string, uint8, uint256, uint256, uint256, uint256, uint256, uint256)
func (_Token *TokenSession) Meta(account common.Address) (string, string, uint8, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _Token.Contract.Meta(&_Token.CallOpts, account)
}

// Meta is a free data retrieval call binding the contract method 0xe021deff.
//
// Solidity: function meta(account address) constant returns(string, string, uint8, uint256, uint256, uint256, uint256, uint256, uint256)
func (_Token *TokenCallerSession) Meta(account common.Address) (string, string, uint8, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _Token.Contract.Meta(&_Token.CallOpts, account)
}

// MintingFinished is a free data retrieval call binding the contract method 0x05d2035b.
//
// Solidity: function mintingFinished() constant returns(bool)
func (_Token *TokenCaller) MintingFinished(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "mintingFinished")
	return *ret0, err
}

// MintingFinished is a free data retrieval call binding the contract method 0x05d2035b.
//
// Solidity: function mintingFinished() constant returns(bool)
func (_Token *TokenSession) MintingFinished() (bool, error) {
	return _Token.Contract.MintingFinished(&_Token.CallOpts)
}

// MintingFinished is a free data retrieval call binding the contract method 0x05d2035b.
//
// Solidity: function mintingFinished() constant returns(bool)
func (_Token *TokenCallerSession) MintingFinished() (bool, error) {
	return _Token.Contract.MintingFinished(&_Token.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_Token *TokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_Token *TokenSession) Name() (string, error) {
	return _Token.Contract.Name(&_Token.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_Token *TokenCallerSession) Name() (string, error) {
	return _Token.Contract.Name(&_Token.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Token *TokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Token *TokenSession) Owner() (common.Address, error) {
	return _Token.Contract.Owner(&_Token.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Token *TokenCallerSession) Owner() (common.Address, error) {
	return _Token.Contract.Owner(&_Token.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Token *TokenCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Token *TokenSession) Paused() (bool, error) {
	return _Token.Contract.Paused(&_Token.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Token *TokenCallerSession) Paused() (bool, error) {
	return _Token.Contract.Paused(&_Token.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_Token *TokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_Token *TokenSession) Symbol() (string, error) {
	return _Token.Contract.Symbol(&_Token.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_Token *TokenCallerSession) Symbol() (string, error) {
	return _Token.Contract.Symbol(&_Token.CallOpts)
}

// TotalHolders is a free data retrieval call binding the contract method 0x53d74fdf.
//
// Solidity: function totalHolders() constant returns(uint256)
func (_Token *TokenCaller) TotalHolders(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "totalHolders")
	return *ret0, err
}

// TotalHolders is a free data retrieval call binding the contract method 0x53d74fdf.
//
// Solidity: function totalHolders() constant returns(uint256)
func (_Token *TokenSession) TotalHolders() (*big.Int, error) {
	return _Token.Contract.TotalHolders(&_Token.CallOpts)
}

// TotalHolders is a free data retrieval call binding the contract method 0x53d74fdf.
//
// Solidity: function totalHolders() constant returns(uint256)
func (_Token *TokenCallerSession) TotalHolders() (*big.Int, error) {
	return _Token.Contract.TotalHolders(&_Token.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Token *TokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Token *TokenSession) TotalSupply() (*big.Int, error) {
	return _Token.Contract.TotalSupply(&_Token.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Token *TokenCallerSession) TotalSupply() (*big.Int, error) {
	return _Token.Contract.TotalSupply(&_Token.CallOpts)
}

// TotalTransfers is a free data retrieval call binding the contract method 0x3bf47720.
//
// Solidity: function totalTransfers() constant returns(uint256)
func (_Token *TokenCaller) TotalTransfers(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "totalTransfers")
	return *ret0, err
}

// TotalTransfers is a free data retrieval call binding the contract method 0x3bf47720.
//
// Solidity: function totalTransfers() constant returns(uint256)
func (_Token *TokenSession) TotalTransfers() (*big.Int, error) {
	return _Token.Contract.TotalTransfers(&_Token.CallOpts)
}

// TotalTransfers is a free data retrieval call binding the contract method 0x3bf47720.
//
// Solidity: function totalTransfers() constant returns(uint256)
func (_Token *TokenCallerSession) TotalTransfers() (*big.Int, error) {
	return _Token.Contract.TotalTransfers(&_Token.CallOpts)
}

// AddAgent is a paid mutator transaction binding the contract method 0x84e79842.
//
// Solidity: function addAgent(account address) returns()
func (_Token *TokenTransactor) AddAgent(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "addAgent", account)
}

// AddAgent is a paid mutator transaction binding the contract method 0x84e79842.
//
// Solidity: function addAgent(account address) returns()
func (_Token *TokenSession) AddAgent(account common.Address) (*types.Transaction, error) {
	return _Token.Contract.AddAgent(&_Token.TransactOpts, account)
}

// AddAgent is a paid mutator transaction binding the contract method 0x84e79842.
//
// Solidity: function addAgent(account address) returns()
func (_Token *TokenTransactorSession) AddAgent(account common.Address) (*types.Transaction, error) {
	return _Token.Contract.AddAgent(&_Token.TransactOpts, account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(account address) returns()
func (_Token *TokenTransactor) AddMinter(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "addMinter", account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(account address) returns()
func (_Token *TokenSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _Token.Contract.AddMinter(&_Token.TransactOpts, account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(account address) returns()
func (_Token *TokenTransactorSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _Token.Contract.AddMinter(&_Token.TransactOpts, account)
}

// AddPauser is a paid mutator transaction binding the contract method 0x82dc1ec4.
//
// Solidity: function addPauser(account address) returns()
func (_Token *TokenTransactor) AddPauser(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "addPauser", account)
}

// AddPauser is a paid mutator transaction binding the contract method 0x82dc1ec4.
//
// Solidity: function addPauser(account address) returns()
func (_Token *TokenSession) AddPauser(account common.Address) (*types.Transaction, error) {
	return _Token.Contract.AddPauser(&_Token.TransactOpts, account)
}

// AddPauser is a paid mutator transaction binding the contract method 0x82dc1ec4.
//
// Solidity: function addPauser(account address) returns()
func (_Token *TokenTransactorSession) AddPauser(account common.Address) (*types.Transaction, error) {
	return _Token.Contract.AddPauser(&_Token.TransactOpts, account)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_Token *TokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_Token *TokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Approve(&_Token.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_Token *TokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Approve(&_Token.TransactOpts, spender, value)
}

// ApproveProxy is a paid mutator transaction binding the contract method 0x2376fe70.
//
// Solidity: function approveProxy(from address, spender address, value uint256) returns(bool)
func (_Token *TokenTransactor) ApproveProxy(opts *bind.TransactOpts, from common.Address, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "approveProxy", from, spender, value)
}

// ApproveProxy is a paid mutator transaction binding the contract method 0x2376fe70.
//
// Solidity: function approveProxy(from address, spender address, value uint256) returns(bool)
func (_Token *TokenSession) ApproveProxy(from common.Address, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.ApproveProxy(&_Token.TransactOpts, from, spender, value)
}

// ApproveProxy is a paid mutator transaction binding the contract method 0x2376fe70.
//
// Solidity: function approveProxy(from address, spender address, value uint256) returns(bool)
func (_Token *TokenTransactorSession) ApproveProxy(from common.Address, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.ApproveProxy(&_Token.TransactOpts, from, spender, value)
}

// BatchTransfer is a paid mutator transaction binding the contract method 0x88d695b2.
//
// Solidity: function batchTransfer(addresses address[], tokenAmount uint256[]) returns(bool)
func (_Token *TokenTransactor) BatchTransfer(opts *bind.TransactOpts, addresses []common.Address, tokenAmount []*big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "batchTransfer", addresses, tokenAmount)
}

// BatchTransfer is a paid mutator transaction binding the contract method 0x88d695b2.
//
// Solidity: function batchTransfer(addresses address[], tokenAmount uint256[]) returns(bool)
func (_Token *TokenSession) BatchTransfer(addresses []common.Address, tokenAmount []*big.Int) (*types.Transaction, error) {
	return _Token.Contract.BatchTransfer(&_Token.TransactOpts, addresses, tokenAmount)
}

// BatchTransfer is a paid mutator transaction binding the contract method 0x88d695b2.
//
// Solidity: function batchTransfer(addresses address[], tokenAmount uint256[]) returns(bool)
func (_Token *TokenTransactorSession) BatchTransfer(addresses []common.Address, tokenAmount []*big.Int) (*types.Transaction, error) {
	return _Token.Contract.BatchTransfer(&_Token.TransactOpts, addresses, tokenAmount)
}

// BatchTransferFrom is a paid mutator transaction binding the contract method 0x4885b254.
//
// Solidity: function batchTransferFrom(_from address, addresses address[], tokenAmount uint256[]) returns(bool)
func (_Token *TokenTransactor) BatchTransferFrom(opts *bind.TransactOpts, _from common.Address, addresses []common.Address, tokenAmount []*big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "batchTransferFrom", _from, addresses, tokenAmount)
}

// BatchTransferFrom is a paid mutator transaction binding the contract method 0x4885b254.
//
// Solidity: function batchTransferFrom(_from address, addresses address[], tokenAmount uint256[]) returns(bool)
func (_Token *TokenSession) BatchTransferFrom(_from common.Address, addresses []common.Address, tokenAmount []*big.Int) (*types.Transaction, error) {
	return _Token.Contract.BatchTransferFrom(&_Token.TransactOpts, _from, addresses, tokenAmount)
}

// BatchTransferFrom is a paid mutator transaction binding the contract method 0x4885b254.
//
// Solidity: function batchTransferFrom(_from address, addresses address[], tokenAmount uint256[]) returns(bool)
func (_Token *TokenTransactorSession) BatchTransferFrom(_from common.Address, addresses []common.Address, tokenAmount []*big.Int) (*types.Transaction, error) {
	return _Token.Contract.BatchTransferFrom(&_Token.TransactOpts, _from, addresses, tokenAmount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(value uint256) returns()
func (_Token *TokenTransactor) Burn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "burn", value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(value uint256) returns()
func (_Token *TokenSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Burn(&_Token.TransactOpts, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(value uint256) returns()
func (_Token *TokenTransactorSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Burn(&_Token.TransactOpts, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(from address, value uint256) returns()
func (_Token *TokenTransactor) BurnFrom(opts *bind.TransactOpts, from common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "burnFrom", from, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(from address, value uint256) returns()
func (_Token *TokenSession) BurnFrom(from common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.BurnFrom(&_Token.TransactOpts, from, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(from address, value uint256) returns()
func (_Token *TokenTransactorSession) BurnFrom(from common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.BurnFrom(&_Token.TransactOpts, from, value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(spender address, subtractedValue uint256) returns(success bool)
func (_Token *TokenTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(spender address, subtractedValue uint256) returns(success bool)
func (_Token *TokenSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Token.Contract.DecreaseAllowance(&_Token.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(spender address, subtractedValue uint256) returns(success bool)
func (_Token *TokenTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Token.Contract.DecreaseAllowance(&_Token.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowanceProxy is a paid mutator transaction binding the contract method 0x02a4e488.
//
// Solidity: function decreaseAllowanceProxy(from address, spender address, subtractedValue uint256) returns(success bool)
func (_Token *TokenTransactor) DecreaseAllowanceProxy(opts *bind.TransactOpts, from common.Address, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "decreaseAllowanceProxy", from, spender, subtractedValue)
}

// DecreaseAllowanceProxy is a paid mutator transaction binding the contract method 0x02a4e488.
//
// Solidity: function decreaseAllowanceProxy(from address, spender address, subtractedValue uint256) returns(success bool)
func (_Token *TokenSession) DecreaseAllowanceProxy(from common.Address, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Token.Contract.DecreaseAllowanceProxy(&_Token.TransactOpts, from, spender, subtractedValue)
}

// DecreaseAllowanceProxy is a paid mutator transaction binding the contract method 0x02a4e488.
//
// Solidity: function decreaseAllowanceProxy(from address, spender address, subtractedValue uint256) returns(success bool)
func (_Token *TokenTransactorSession) DecreaseAllowanceProxy(from common.Address, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Token.Contract.DecreaseAllowanceProxy(&_Token.TransactOpts, from, spender, subtractedValue)
}

// FinishMinting is a paid mutator transaction binding the contract method 0x7d64bcb4.
//
// Solidity: function finishMinting() returns(bool)
func (_Token *TokenTransactor) FinishMinting(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "finishMinting")
}

// FinishMinting is a paid mutator transaction binding the contract method 0x7d64bcb4.
//
// Solidity: function finishMinting() returns(bool)
func (_Token *TokenSession) FinishMinting() (*types.Transaction, error) {
	return _Token.Contract.FinishMinting(&_Token.TransactOpts)
}

// FinishMinting is a paid mutator transaction binding the contract method 0x7d64bcb4.
//
// Solidity: function finishMinting() returns(bool)
func (_Token *TokenTransactorSession) FinishMinting() (*types.Transaction, error) {
	return _Token.Contract.FinishMinting(&_Token.TransactOpts)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(spender address, addedValue uint256) returns(success bool)
func (_Token *TokenTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(spender address, addedValue uint256) returns(success bool)
func (_Token *TokenSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Token.Contract.IncreaseAllowance(&_Token.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(spender address, addedValue uint256) returns(success bool)
func (_Token *TokenTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Token.Contract.IncreaseAllowance(&_Token.TransactOpts, spender, addedValue)
}

// IncreaseAllowanceProxy is a paid mutator transaction binding the contract method 0x82bcef79.
//
// Solidity: function increaseAllowanceProxy(from address, spender address, addedValue uint256) returns(success bool)
func (_Token *TokenTransactor) IncreaseAllowanceProxy(opts *bind.TransactOpts, from common.Address, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "increaseAllowanceProxy", from, spender, addedValue)
}

// IncreaseAllowanceProxy is a paid mutator transaction binding the contract method 0x82bcef79.
//
// Solidity: function increaseAllowanceProxy(from address, spender address, addedValue uint256) returns(success bool)
func (_Token *TokenSession) IncreaseAllowanceProxy(from common.Address, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Token.Contract.IncreaseAllowanceProxy(&_Token.TransactOpts, from, spender, addedValue)
}

// IncreaseAllowanceProxy is a paid mutator transaction binding the contract method 0x82bcef79.
//
// Solidity: function increaseAllowanceProxy(from address, spender address, addedValue uint256) returns(success bool)
func (_Token *TokenTransactorSession) IncreaseAllowanceProxy(from common.Address, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Token.Contract.IncreaseAllowanceProxy(&_Token.TransactOpts, from, spender, addedValue)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(to address, amount uint256) returns(bool)
func (_Token *TokenTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(to address, amount uint256) returns(bool)
func (_Token *TokenSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Mint(&_Token.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(to address, amount uint256) returns(bool)
func (_Token *TokenTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Mint(&_Token.TransactOpts, to, amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Token *TokenTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Token *TokenSession) Pause() (*types.Transaction, error) {
	return _Token.Contract.Pause(&_Token.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Token *TokenTransactorSession) Pause() (*types.Transaction, error) {
	return _Token.Contract.Pause(&_Token.TransactOpts)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(account address) returns()
func (_Token *TokenTransactor) RemoveAgent(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "removeAgent", account)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(account address) returns()
func (_Token *TokenSession) RemoveAgent(account common.Address) (*types.Transaction, error) {
	return _Token.Contract.RemoveAgent(&_Token.TransactOpts, account)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(account address) returns()
func (_Token *TokenTransactorSession) RemoveAgent(account common.Address) (*types.Transaction, error) {
	return _Token.Contract.RemoveAgent(&_Token.TransactOpts, account)
}

// RemoveMinter is a paid mutator transaction binding the contract method 0x3092afd5.
//
// Solidity: function removeMinter(account address) returns()
func (_Token *TokenTransactor) RemoveMinter(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "removeMinter", account)
}

// RemoveMinter is a paid mutator transaction binding the contract method 0x3092afd5.
//
// Solidity: function removeMinter(account address) returns()
func (_Token *TokenSession) RemoveMinter(account common.Address) (*types.Transaction, error) {
	return _Token.Contract.RemoveMinter(&_Token.TransactOpts, account)
}

// RemoveMinter is a paid mutator transaction binding the contract method 0x3092afd5.
//
// Solidity: function removeMinter(account address) returns()
func (_Token *TokenTransactorSession) RemoveMinter(account common.Address) (*types.Transaction, error) {
	return _Token.Contract.RemoveMinter(&_Token.TransactOpts, account)
}

// RemovePauser is a paid mutator transaction binding the contract method 0x6b2c0f55.
//
// Solidity: function removePauser(account address) returns()
func (_Token *TokenTransactor) RemovePauser(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "removePauser", account)
}

// RemovePauser is a paid mutator transaction binding the contract method 0x6b2c0f55.
//
// Solidity: function removePauser(account address) returns()
func (_Token *TokenSession) RemovePauser(account common.Address) (*types.Transaction, error) {
	return _Token.Contract.RemovePauser(&_Token.TransactOpts, account)
}

// RemovePauser is a paid mutator transaction binding the contract method 0x6b2c0f55.
//
// Solidity: function removePauser(account address) returns()
func (_Token *TokenTransactorSession) RemovePauser(account common.Address) (*types.Transaction, error) {
	return _Token.Contract.RemovePauser(&_Token.TransactOpts, account)
}

// RenounceAgent is a paid mutator transaction binding the contract method 0xc692f4cf.
//
// Solidity: function renounceAgent() returns()
func (_Token *TokenTransactor) RenounceAgent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "renounceAgent")
}

// RenounceAgent is a paid mutator transaction binding the contract method 0xc692f4cf.
//
// Solidity: function renounceAgent() returns()
func (_Token *TokenSession) RenounceAgent() (*types.Transaction, error) {
	return _Token.Contract.RenounceAgent(&_Token.TransactOpts)
}

// RenounceAgent is a paid mutator transaction binding the contract method 0xc692f4cf.
//
// Solidity: function renounceAgent() returns()
func (_Token *TokenTransactorSession) RenounceAgent() (*types.Transaction, error) {
	return _Token.Contract.RenounceAgent(&_Token.TransactOpts)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_Token *TokenTransactor) RenounceMinter(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "renounceMinter")
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_Token *TokenSession) RenounceMinter() (*types.Transaction, error) {
	return _Token.Contract.RenounceMinter(&_Token.TransactOpts)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_Token *TokenTransactorSession) RenounceMinter() (*types.Transaction, error) {
	return _Token.Contract.RenounceMinter(&_Token.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Token *TokenTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Token *TokenSession) RenounceOwnership() (*types.Transaction, error) {
	return _Token.Contract.RenounceOwnership(&_Token.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Token *TokenTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Token.Contract.RenounceOwnership(&_Token.TransactOpts)
}

// RenouncePauser is a paid mutator transaction binding the contract method 0x6ef8d66d.
//
// Solidity: function renouncePauser() returns()
func (_Token *TokenTransactor) RenouncePauser(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "renouncePauser")
}

// RenouncePauser is a paid mutator transaction binding the contract method 0x6ef8d66d.
//
// Solidity: function renouncePauser() returns()
func (_Token *TokenSession) RenouncePauser() (*types.Transaction, error) {
	return _Token.Contract.RenouncePauser(&_Token.TransactOpts)
}

// RenouncePauser is a paid mutator transaction binding the contract method 0x6ef8d66d.
//
// Solidity: function renouncePauser() returns()
func (_Token *TokenTransactorSession) RenouncePauser() (*types.Transaction, error) {
	return _Token.Contract.RenouncePauser(&_Token.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_Token *TokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_Token *TokenSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Transfer(&_Token.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_Token *TokenTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Transfer(&_Token.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_Token *TokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_Token *TokenSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.TransferFrom(&_Token.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_Token *TokenTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.TransferFrom(&_Token.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Token *TokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Token *TokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Token.Contract.TransferOwnership(&_Token.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Token *TokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Token.Contract.TransferOwnership(&_Token.TransactOpts, newOwner)
}

// TransferProxy is a paid mutator transaction binding the contract method 0x4733dc8f.
//
// Solidity: function transferProxy(from address, to address, value uint256) returns(bool)
func (_Token *TokenTransactor) TransferProxy(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "transferProxy", from, to, value)
}

// TransferProxy is a paid mutator transaction binding the contract method 0x4733dc8f.
//
// Solidity: function transferProxy(from address, to address, value uint256) returns(bool)
func (_Token *TokenSession) TransferProxy(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.TransferProxy(&_Token.TransactOpts, from, to, value)
}

// TransferProxy is a paid mutator transaction binding the contract method 0x4733dc8f.
//
// Solidity: function transferProxy(from address, to address, value uint256) returns(bool)
func (_Token *TokenTransactorSession) TransferProxy(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Token.Contract.TransferProxy(&_Token.TransactOpts, from, to, value)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Token *TokenTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Token *TokenSession) Unpause() (*types.Transaction, error) {
	return _Token.Contract.Unpause(&_Token.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Token *TokenTransactorSession) Unpause() (*types.Transaction, error) {
	return _Token.Contract.Unpause(&_Token.TransactOpts)
}

// TokenAgentAddedIterator is returned from FilterAgentAdded and is used to iterate over the raw logs and unpacked data for AgentAdded events raised by the Token contract.
type TokenAgentAddedIterator struct {
	Event *TokenAgentAdded // Event containing the contract specifics and raw log

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
func (it *TokenAgentAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAgentAdded)
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
		it.Event = new(TokenAgentAdded)
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
func (it *TokenAgentAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenAgentAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenAgentAdded represents a AgentAdded event raised by the Token contract.
type TokenAgentAdded struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAdded is a free log retrieval operation binding the contract event 0xf68e73cec97f2d70aa641fb26e87a4383686e2efacb648f2165aeb02ac562ec5.
//
// Solidity: e AgentAdded(account indexed address)
func (_Token *TokenFilterer) FilterAgentAdded(opts *bind.FilterOpts, account []common.Address) (*TokenAgentAddedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "AgentAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return &TokenAgentAddedIterator{contract: _Token.contract, event: "AgentAdded", logs: logs, sub: sub}, nil
}

// WatchAgentAdded is a free log subscription operation binding the contract event 0xf68e73cec97f2d70aa641fb26e87a4383686e2efacb648f2165aeb02ac562ec5.
//
// Solidity: e AgentAdded(account indexed address)
func (_Token *TokenFilterer) WatchAgentAdded(opts *bind.WatchOpts, sink chan<- *TokenAgentAdded, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "AgentAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenAgentAdded)
				if err := _Token.contract.UnpackLog(event, "AgentAdded", log); err != nil {
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

// TokenAgentRemovedIterator is returned from FilterAgentRemoved and is used to iterate over the raw logs and unpacked data for AgentRemoved events raised by the Token contract.
type TokenAgentRemovedIterator struct {
	Event *TokenAgentRemoved // Event containing the contract specifics and raw log

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
func (it *TokenAgentRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenAgentRemoved)
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
		it.Event = new(TokenAgentRemoved)
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
func (it *TokenAgentRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenAgentRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenAgentRemoved represents a AgentRemoved event raised by the Token contract.
type TokenAgentRemoved struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentRemoved is a free log retrieval operation binding the contract event 0xed9c8ad8d5a0a66898ea49d2956929c93ae2e8bd50281b2ed897c5d1a6737e0b.
//
// Solidity: e AgentRemoved(account indexed address)
func (_Token *TokenFilterer) FilterAgentRemoved(opts *bind.FilterOpts, account []common.Address) (*TokenAgentRemovedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "AgentRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return &TokenAgentRemovedIterator{contract: _Token.contract, event: "AgentRemoved", logs: logs, sub: sub}, nil
}

// WatchAgentRemoved is a free log subscription operation binding the contract event 0xed9c8ad8d5a0a66898ea49d2956929c93ae2e8bd50281b2ed897c5d1a6737e0b.
//
// Solidity: e AgentRemoved(account indexed address)
func (_Token *TokenFilterer) WatchAgentRemoved(opts *bind.WatchOpts, sink chan<- *TokenAgentRemoved, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "AgentRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenAgentRemoved)
				if err := _Token.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
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

// TokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Token contract.
type TokenApprovalIterator struct {
	Event *TokenApproval // Event containing the contract specifics and raw log

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
func (it *TokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenApproval)
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
		it.Event = new(TokenApproval)
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
func (it *TokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenApproval represents a Approval event raised by the Token contract.
type TokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_Token *TokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TokenApprovalIterator{contract: _Token.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_Token *TokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenApproval)
				if err := _Token.contract.UnpackLog(event, "Approval", log); err != nil {
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

// TokenMinterAddedIterator is returned from FilterMinterAdded and is used to iterate over the raw logs and unpacked data for MinterAdded events raised by the Token contract.
type TokenMinterAddedIterator struct {
	Event *TokenMinterAdded // Event containing the contract specifics and raw log

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
func (it *TokenMinterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenMinterAdded)
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
		it.Event = new(TokenMinterAdded)
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
func (it *TokenMinterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenMinterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenMinterAdded represents a MinterAdded event raised by the Token contract.
type TokenMinterAdded struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinterAdded is a free log retrieval operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: e MinterAdded(account indexed address)
func (_Token *TokenFilterer) FilterMinterAdded(opts *bind.FilterOpts, account []common.Address) (*TokenMinterAddedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "MinterAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return &TokenMinterAddedIterator{contract: _Token.contract, event: "MinterAdded", logs: logs, sub: sub}, nil
}

// WatchMinterAdded is a free log subscription operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: e MinterAdded(account indexed address)
func (_Token *TokenFilterer) WatchMinterAdded(opts *bind.WatchOpts, sink chan<- *TokenMinterAdded, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "MinterAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenMinterAdded)
				if err := _Token.contract.UnpackLog(event, "MinterAdded", log); err != nil {
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

// TokenMinterRemovedIterator is returned from FilterMinterRemoved and is used to iterate over the raw logs and unpacked data for MinterRemoved events raised by the Token contract.
type TokenMinterRemovedIterator struct {
	Event *TokenMinterRemoved // Event containing the contract specifics and raw log

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
func (it *TokenMinterRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenMinterRemoved)
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
		it.Event = new(TokenMinterRemoved)
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
func (it *TokenMinterRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenMinterRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenMinterRemoved represents a MinterRemoved event raised by the Token contract.
type TokenMinterRemoved struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinterRemoved is a free log retrieval operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: e MinterRemoved(account indexed address)
func (_Token *TokenFilterer) FilterMinterRemoved(opts *bind.FilterOpts, account []common.Address) (*TokenMinterRemovedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "MinterRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return &TokenMinterRemovedIterator{contract: _Token.contract, event: "MinterRemoved", logs: logs, sub: sub}, nil
}

// WatchMinterRemoved is a free log subscription operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: e MinterRemoved(account indexed address)
func (_Token *TokenFilterer) WatchMinterRemoved(opts *bind.WatchOpts, sink chan<- *TokenMinterRemoved, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "MinterRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenMinterRemoved)
				if err := _Token.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
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

// TokenMintingFinishedIterator is returned from FilterMintingFinished and is used to iterate over the raw logs and unpacked data for MintingFinished events raised by the Token contract.
type TokenMintingFinishedIterator struct {
	Event *TokenMintingFinished // Event containing the contract specifics and raw log

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
func (it *TokenMintingFinishedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenMintingFinished)
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
		it.Event = new(TokenMintingFinished)
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
func (it *TokenMintingFinishedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenMintingFinishedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenMintingFinished represents a MintingFinished event raised by the Token contract.
type TokenMintingFinished struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterMintingFinished is a free log retrieval operation binding the contract event 0xb828d9b5c78095deeeeff2eca2e5d4fe046ce3feb4c99702624a3fd384ad2dbc.
//
// Solidity: e MintingFinished()
func (_Token *TokenFilterer) FilterMintingFinished(opts *bind.FilterOpts) (*TokenMintingFinishedIterator, error) {

	logs, sub, err := _Token.contract.FilterLogs(opts, "MintingFinished")
	if err != nil {
		return nil, err
	}
	return &TokenMintingFinishedIterator{contract: _Token.contract, event: "MintingFinished", logs: logs, sub: sub}, nil
}

// WatchMintingFinished is a free log subscription operation binding the contract event 0xb828d9b5c78095deeeeff2eca2e5d4fe046ce3feb4c99702624a3fd384ad2dbc.
//
// Solidity: e MintingFinished()
func (_Token *TokenFilterer) WatchMintingFinished(opts *bind.WatchOpts, sink chan<- *TokenMintingFinished) (event.Subscription, error) {

	logs, sub, err := _Token.contract.WatchLogs(opts, "MintingFinished")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenMintingFinished)
				if err := _Token.contract.UnpackLog(event, "MintingFinished", log); err != nil {
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

// TokenOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the Token contract.
type TokenOwnershipRenouncedIterator struct {
	Event *TokenOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *TokenOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenOwnershipRenounced)
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
		it.Event = new(TokenOwnershipRenounced)
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
func (it *TokenOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenOwnershipRenounced represents a OwnershipRenounced event raised by the Token contract.
type TokenOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Token *TokenFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*TokenOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TokenOwnershipRenouncedIterator{contract: _Token.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Token *TokenFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *TokenOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenOwnershipRenounced)
				if err := _Token.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// TokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Token contract.
type TokenOwnershipTransferredIterator struct {
	Event *TokenOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenOwnershipTransferred)
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
		it.Event = new(TokenOwnershipTransferred)
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
func (it *TokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenOwnershipTransferred represents a OwnershipTransferred event raised by the Token contract.
type TokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Token *TokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TokenOwnershipTransferredIterator{contract: _Token.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Token *TokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenOwnershipTransferred)
				if err := _Token.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// TokenPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Token contract.
type TokenPausedIterator struct {
	Event *TokenPaused // Event containing the contract specifics and raw log

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
func (it *TokenPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPaused)
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
		it.Event = new(TokenPaused)
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
func (it *TokenPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenPaused represents a Paused event raised by the Token contract.
type TokenPaused struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x9e87fac88ff661f02d44f95383c817fece4bce600a3dab7a54406878b965e752.
//
// Solidity: e Paused()
func (_Token *TokenFilterer) FilterPaused(opts *bind.FilterOpts) (*TokenPausedIterator, error) {

	logs, sub, err := _Token.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &TokenPausedIterator{contract: _Token.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x9e87fac88ff661f02d44f95383c817fece4bce600a3dab7a54406878b965e752.
//
// Solidity: e Paused()
func (_Token *TokenFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *TokenPaused) (event.Subscription, error) {

	logs, sub, err := _Token.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenPaused)
				if err := _Token.contract.UnpackLog(event, "Paused", log); err != nil {
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

// TokenPauserAddedIterator is returned from FilterPauserAdded and is used to iterate over the raw logs and unpacked data for PauserAdded events raised by the Token contract.
type TokenPauserAddedIterator struct {
	Event *TokenPauserAdded // Event containing the contract specifics and raw log

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
func (it *TokenPauserAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPauserAdded)
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
		it.Event = new(TokenPauserAdded)
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
func (it *TokenPauserAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenPauserAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenPauserAdded represents a PauserAdded event raised by the Token contract.
type TokenPauserAdded struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPauserAdded is a free log retrieval operation binding the contract event 0x6719d08c1888103bea251a4ed56406bd0c3e69723c8a1686e017e7bbe159b6f8.
//
// Solidity: e PauserAdded(account indexed address)
func (_Token *TokenFilterer) FilterPauserAdded(opts *bind.FilterOpts, account []common.Address) (*TokenPauserAddedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "PauserAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return &TokenPauserAddedIterator{contract: _Token.contract, event: "PauserAdded", logs: logs, sub: sub}, nil
}

// WatchPauserAdded is a free log subscription operation binding the contract event 0x6719d08c1888103bea251a4ed56406bd0c3e69723c8a1686e017e7bbe159b6f8.
//
// Solidity: e PauserAdded(account indexed address)
func (_Token *TokenFilterer) WatchPauserAdded(opts *bind.WatchOpts, sink chan<- *TokenPauserAdded, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "PauserAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenPauserAdded)
				if err := _Token.contract.UnpackLog(event, "PauserAdded", log); err != nil {
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

// TokenPauserRemovedIterator is returned from FilterPauserRemoved and is used to iterate over the raw logs and unpacked data for PauserRemoved events raised by the Token contract.
type TokenPauserRemovedIterator struct {
	Event *TokenPauserRemoved // Event containing the contract specifics and raw log

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
func (it *TokenPauserRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenPauserRemoved)
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
		it.Event = new(TokenPauserRemoved)
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
func (it *TokenPauserRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenPauserRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenPauserRemoved represents a PauserRemoved event raised by the Token contract.
type TokenPauserRemoved struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPauserRemoved is a free log retrieval operation binding the contract event 0xcd265ebaf09df2871cc7bd4133404a235ba12eff2041bb89d9c714a2621c7c7e.
//
// Solidity: e PauserRemoved(account indexed address)
func (_Token *TokenFilterer) FilterPauserRemoved(opts *bind.FilterOpts, account []common.Address) (*TokenPauserRemovedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "PauserRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return &TokenPauserRemovedIterator{contract: _Token.contract, event: "PauserRemoved", logs: logs, sub: sub}, nil
}

// WatchPauserRemoved is a free log subscription operation binding the contract event 0xcd265ebaf09df2871cc7bd4133404a235ba12eff2041bb89d9c714a2621c7c7e.
//
// Solidity: e PauserRemoved(account indexed address)
func (_Token *TokenFilterer) WatchPauserRemoved(opts *bind.WatchOpts, sink chan<- *TokenPauserRemoved, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "PauserRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenPauserRemoved)
				if err := _Token.contract.UnpackLog(event, "PauserRemoved", log); err != nil {
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

// TokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Token contract.
type TokenTransferIterator struct {
	Event *TokenTransfer // Event containing the contract specifics and raw log

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
func (it *TokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenTransfer)
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
		it.Event = new(TokenTransfer)
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
func (it *TokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenTransfer represents a Transfer event raised by the Token contract.
type TokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_Token *TokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TokenTransferIterator{contract: _Token.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_Token *TokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenTransfer)
				if err := _Token.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// TokenUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Token contract.
type TokenUnpausedIterator struct {
	Event *TokenUnpaused // Event containing the contract specifics and raw log

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
func (it *TokenUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenUnpaused)
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
		it.Event = new(TokenUnpaused)
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
func (it *TokenUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenUnpaused represents a Unpaused event raised by the Token contract.
type TokenUnpaused struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0xa45f47fdea8a1efdd9029a5691c7f759c32b7c698632b563573e155625d16933.
//
// Solidity: e Unpaused()
func (_Token *TokenFilterer) FilterUnpaused(opts *bind.FilterOpts) (*TokenUnpausedIterator, error) {

	logs, sub, err := _Token.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &TokenUnpausedIterator{contract: _Token.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0xa45f47fdea8a1efdd9029a5691c7f759c32b7c698632b563573e155625d16933.
//
// Solidity: e Unpaused()
func (_Token *TokenFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *TokenUnpaused) (event.Subscription, error) {

	logs, sub, err := _Token.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenUnpaused)
				if err := _Token.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
