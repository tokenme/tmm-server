pragma solidity ^0.4.24;

import "Roles.sol";
import "Ownable.sol";


contract AgentRole is Ownable {
  using Roles for Roles.Role;

  event AgentAdded(address indexed account);
  event AgentRemoved(address indexed account);

  Roles.Role private agencies;

  constructor() public {
    agencies.add(msg.sender);
  }

  modifier onlyAgent() {
    require(isOwner() || isAgent(msg.sender));
    _;
  }

  function isAgent(address account) public view returns (bool) {
    return agencies.has(account);
  }

  function addAgent(address account) public onlyAgent {
    agencies.add(account);
    emit AgentAdded(account);
  }

  function renounceAgent() public onlyAgent {
    agencies.remove(msg.sender);
  }

  function _removeAgent(address account) internal {
    agencies.remove(account);
    emit AgentRemoved(account);
  }
}