pragma solidity ^0.4.24;

import "ERC20Burnable.sol";
import "ERC20Mintable.sol";
import "ERC20Pausable.sol";
import "ERC20Agentable.sol";

contract Token is ERC20Burnable, ERC20Mintable, ERC20Pausable, ERC20Agentable {

  string private _name;
  string private _symbol;
  uint8 private _decimals;

  constructor(string name, string symbol, uint8 decimals, uint256 initialSupply) public {
    _name = name;
    _symbol = symbol;
    _decimals = decimals;
    _initialSupply = initialSupply;
    _totalSupply = _initialSupply;
    _balances[msg.sender] = _initialSupply;
    emit Transfer(0x0, msg.sender, _initialSupply);
  }

  /**
   * @return the name of the token.
   */
  function name() public view returns(string) {
    return _name;
  }

  /**
   * @return the symbol of the token.
   */
  function symbol() public view returns(string) {
    return _symbol;
  }

  /**
   * @return the number of decimals of the token.
   */
  function decimals() public view returns(uint8) {
    return _decimals;
  }

  function meta(address account) public view returns (string, string, uint8, uint256, uint256, uint256, uint256, uint256, uint256) {
    uint256 circulating = 0;
    if (_totalSupply > _balances[owner()]) {
      circulating = _totalSupply.sub(_balances[owner()]);
    }
    uint256 balance = 0;
    if (account != address(0)) {
      balance = _balances[account];
    } else if (msg.sender != address(0)) {
      balance = _balances[msg.sender];
    }
    return (_name, _symbol, _decimals, _initialSupply, _totalSupply, _totalTransfers, _totalHolders, circulating, balance);
  }

  function batchTransfer(address[] addresses, uint256[] tokenAmount) public returns (bool) {
    require(addresses.length > 0 && addresses.length == tokenAmount.length);
    for (uint i = 0; i < addresses.length; i++) {
        address _to = addresses[i];
        uint256 _value = tokenAmount[i];
        super.transfer(_to, _value);
    }
    return true;
  }

  function batchTransferFrom(address _from, address[] addresses, uint256[] tokenAmount) public returns (bool) {
    require(addresses.length > 0 && addresses.length == tokenAmount.length);
    for (uint i = 0; i < addresses.length; i++) {
        address _to = addresses[i];
        uint256 _value = tokenAmount[i];
        super.transferFrom(_from, _to, _value);
    }
    return true;
  }


}