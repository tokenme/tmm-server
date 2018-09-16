pragma solidity ^0.4.24;

import "ERC20.sol";
import "AgentRole.sol";


/**
 * @title Agent token
 * @dev ERC20 modified with agentable transfers.
 **/
contract ERC20Agentable is ERC20, AgentRole {

  function removeAgent(address account) public onlyAgent {
    _removeAgent(account);
  }

  function _removeAgent(address account) internal {
    super._removeAgent(account);
  }

  function transferProxy(
    address from,
    address to,
    uint256 value
  )
    public
    onlyAgent
    returns (bool)
  {
    if (msg.sender == from) {
      return transfer(to, value);
    }

    require(value <= _balances[from]);
    require(to != address(0));

    _balances[from] = _balances[from].sub(value);

    if (_balances[from] == 0 && _totalHolders > 0) {
      _totalHolders = _totalHolders.sub(1);
    }
    if (_balances[to] == 0) {
      _totalHolders = _totalHolders.add(1);
    }

    _balances[to] = _balances[to].add(value);
    _totalTransfers = _totalTransfers.add(1);
    emit Transfer(from, to, value);
    return true;
  }

  function approveProxy(
    address from,
    address spender,
    uint256 value
  )
    public
    onlyAgent
    returns (bool)
  {
    require(spender != address(0));

    _allowed[from][spender] = value;
    emit Approval(from, spender, value);
    return true;
  }

  function increaseAllowanceProxy(
    address from,
    address spender,
    uint addedValue
  )
    public
    onlyAgent
    returns (bool success)
  {
    require(spender != address(0));

    _allowed[from][spender] = (
      _allowed[from][spender].add(addedValue));
    emit Approval(from, spender, _allowed[from][spender]);
    return true;
  }

  function decreaseAllowanceProxy(
    address from,
    address spender,
    uint subtractedValue
  )
    public
    onlyAgent
    returns (bool success)
  {
    require(spender != address(0));

    _allowed[from][spender] = (
      _allowed[from][spender].sub(subtractedValue));
    emit Approval(from, spender, _allowed[from][spender]);
    return true;
  }
}