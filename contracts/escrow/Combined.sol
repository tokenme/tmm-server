pragma solidity ^0.4.24;

/**
 * @title SafeMath
 * @dev Math operations with safety checks that revert on error
 */
library SafeMath {

  /**
  * @dev Multiplies two numbers, reverts on overflow.
  */
  function mul(uint256 a, uint256 b) internal pure returns (uint256) {
    // Gas optimization: this is cheaper than requiring 'a' not being zero, but the
    // benefit is lost if 'b' is also tested.
    // See: https://github.com/OpenZeppelin/openzeppelin-solidity/pull/522
    if (a == 0) {
      return 0;
    }

    uint256 c = a * b;
    require(c / a == b);

    return c;
  }

  /**
  * @dev Integer division of two numbers truncating the quotient, reverts on division by zero.
  */
  function div(uint256 a, uint256 b) internal pure returns (uint256) {
    require(b > 0); // Solidity only automatically asserts when dividing by 0
    uint256 c = a / b;
    // assert(a == b * c + a % b); // There is no case in which this doesn't hold

    return c;
  }

  /**
  * @dev Subtracts two numbers, reverts on overflow (i.e. if subtrahend is greater than minuend).
  */
  function sub(uint256 a, uint256 b) internal pure returns (uint256) {
    require(b <= a);
    uint256 c = a - b;

    return c;
  }

  /**
  * @dev Adds two numbers, reverts on overflow.
  */
  function add(uint256 a, uint256 b) internal pure returns (uint256) {
    uint256 c = a + b;
    require(c >= a);

    return c;
  }

  /**
  * @dev Divides two numbers and returns the remainder (unsigned integer modulo),
  * reverts when dividing by zero.
  */
  function mod(uint256 a, uint256 b) internal pure returns (uint256) {
    require(b != 0);
    return a % b;
  }
}

/**
 * @title Roles
 * @dev Library for managing addresses assigned to a Role.
 */
library Roles {
  struct Role {
    mapping (address => bool) bearer;
  }

  /**
   * @dev give an account access to this role
   */
  function add(Role storage role, address account) internal {
    require(account != address(0));
    role.bearer[account] = true;
  }

  /**
   * @dev remove an account's access to this role
   */
  function remove(Role storage role, address account) internal {
    require(account != address(0));
    role.bearer[account] = false;
  }

  /**
   * @dev check if an account has this role
   * @return bool
   */
  function has(Role storage role, address account)
    internal
    view
    returns (bool)
  {
    require(account != address(0));
    return role.bearer[account];
  }
}

/**
 * @title Ownable
 * @dev The Ownable contract has an owner address, and provides basic authorization control
 * functions, this simplifies the implementation of "user permissions".
 */
contract Ownable {
  address private _owner;


  event OwnershipRenounced(address indexed previousOwner);
  event OwnershipTransferred(
    address indexed previousOwner,
    address indexed newOwner
  );


  /**
   * @dev The Ownable constructor sets the original `owner` of the contract to the sender
   * account.
   */
  constructor() public {
    _owner = msg.sender;
  }

  /**
   * @return the address of the owner.
   */
  function owner() public view returns(address) {
    return _owner;
  }

  /**
   * @dev Throws if called by any account other than the owner.
   */
  modifier onlyOwner() {
    require(isOwner());
    _;
  }

  /**
   * @return true if `msg.sender` is the owner of the contract.
   */
  function isOwner() public view returns(bool) {
    return msg.sender == _owner;
  }

  /**
   * @dev Allows the current owner to relinquish control of the contract.
   * @notice Renouncing to ownership will leave the contract without an owner.
   * It will not be possible to call the functions with the `onlyOwner`
   * modifier anymore.
   */
  function renounceOwnership() public onlyOwner {
    emit OwnershipRenounced(_owner);
    _owner = address(0);
  }

  /**
   * @dev Allows the current owner to transfer control of the contract to a newOwner.
   * @param newOwner The address to transfer ownership to.
   */
  function transferOwnership(address newOwner) public onlyOwner {
    _transferOwnership(newOwner);
  }

  /**
   * @dev Transfers control of the contract to a newOwner.
   * @param newOwner The address to transfer ownership to.
   */
  function _transferOwnership(address newOwner) internal {
    require(newOwner != address(0));
    emit OwnershipTransferred(_owner, newOwner);
    _owner = newOwner;
  }
}

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

/**
 * @title ERC20 interface
 * @dev see https://github.com/ethereum/EIPs/issues/20
 */
interface IERC20 {

  function totalSupply() external view returns (uint256);

  function balanceOf(address who) external view returns (uint256);

  function allowance(address owner, address spender)
    external view returns (uint256);

  function transfer(address to, uint256 value) external returns (bool);

  function approve(address spender, uint256 value)
    external returns (bool);

  function transferFrom(address from, address to, uint256 value)
    external returns (bool);

  function transferProxy(address from, address to, uint256 value)
    external returns (bool);

  event Transfer(
    address indexed from,
    address indexed to,
    uint256 value
  );

  event Approval(
    address indexed owner,
    address indexed spender,
    uint256 value
  );
}

/**
 * @title Standard ERC20 token
 *
 * @dev Implementation of the basic standard token.
 * https://github.com/ethereum/EIPs/blob/master/EIPS/eip-20.md
 * Originally based on code by FirstBlood: https://github.com/Firstbloodio/token/blob/master/smart_contract/FirstBloodToken.sol
 */
contract ERC20 is IERC20, Ownable {
  using SafeMath for uint256;

  mapping (address => uint256) internal _balances;

  mapping (address => mapping (address => uint256)) internal _allowed;

  uint256 internal _totalSupply;

  uint256 internal _totalHolders;

  uint256 internal _totalTransfers;

  uint256 internal _initialSupply;

  function initialSupply() public view returns (uint256) {
    return _initialSupply;
  }

  /**
  * @dev Total number of tokens in existence
  */
  function totalSupply() public view returns (uint256) {
    return _totalSupply;
  }

  function circulatingSupply() public view returns (uint256) {
    require(_totalSupply >= _balances[owner()]);
    return _totalSupply.sub(_balances[owner()]);
  }

  /**
  * @dev total number of token holders in existence
  */
  function totalHolders() public view returns (uint256) {
    return _totalHolders;
  }

  /**
  * @dev total number of token transfers
  */
  function totalTransfers() public view returns (uint256) {
    return _totalTransfers;
  }

  /**
  * @dev Gets the balance of the specified address.
  * @param owner The address to query the balance of.
  * @return An uint256 representing the amount owned by the passed address.
  */
  function balanceOf(address owner) public view returns (uint256) {
    return _balances[owner];
  }

  /**
   * @dev Function to check the amount of tokens that an owner allowed to a spender.
   * @param owner address The address which owns the funds.
   * @param spender address The address which will spend the funds.
   * @return A uint256 specifying the amount of tokens still available for the spender.
   */
  function allowance(
    address owner,
    address spender
   )
    public
    view
    returns (uint256)
  {
    return _allowed[owner][spender];
  }

  /**
  * @dev Transfer token for a specified address
  * @param to The address to transfer to.
  * @param value The amount to be transferred.
  */
  function transfer(address to, uint256 value) public returns (bool) {
    require(value <= _balances[msg.sender]);
    require(to != address(0));

    _balances[msg.sender] = _balances[msg.sender].sub(value);
    if (_balances[msg.sender] == 0 && _totalHolders > 0) {
      _totalHolders = _totalHolders.sub(1);
    }
    if (_balances[to] == 0) {
      _totalHolders = _totalHolders.add(1);
    }
    _balances[to] = _balances[to].add(value);
    _totalTransfers = _totalTransfers.add(1);
    emit Transfer(msg.sender, to, value);
    return true;
  }

  /**
   * @dev Approve the passed address to spend the specified amount of tokens on behalf of msg.sender.
   * Beware that changing an allowance with this method brings the risk that someone may use both the old
   * and the new allowance by unfortunate transaction ordering. One possible solution to mitigate this
   * race condition is to first reduce the spender's allowance to 0 and set the desired value afterwards:
   * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
   * @param spender The address which will spend the funds.
   * @param value The amount of tokens to be spent.
   */
  function approve(address spender, uint256 value) public returns (bool) {
    require(spender != address(0));

    _allowed[msg.sender][spender] = value;
    emit Approval(msg.sender, spender, value);
    return true;
  }

  /**
   * @dev Transfer tokens from one address to another
   * @param from address The address which you want to send tokens from
   * @param to address The address which you want to transfer to
   * @param value uint256 the amount of tokens to be transferred
   */
  function transferFrom(
    address from,
    address to,
    uint256 value
  )
    public
    returns (bool)
  {
    if (msg.sender == from) {
      return transfer(to, value);
    }

    require(value <= _balances[from]);
    require(value <= _allowed[from][msg.sender]);
    require(to != address(0));

    _balances[from] = _balances[from].sub(value);

    if (_balances[from] == 0 && _totalHolders > 0) {
      _totalHolders = _totalHolders.sub(1);
    }
    if (_balances[to] == 0) {
      _totalHolders = _totalHolders.add(1);
    }

    _balances[to] = _balances[to].add(value);
    _allowed[from][msg.sender] = _allowed[from][msg.sender].sub(value);
    _totalTransfers = _totalTransfers.add(1);
    emit Transfer(from, to, value);
    return true;
  }

  /**
   * @dev Increase the amount of tokens that an owner allowed to a spender.
   * approve should be called when allowed_[_spender] == 0. To increment
   * allowed value is better to use this function to avoid 2 calls (and wait until
   * the first transaction is mined)
   * From MonolithDAO Token.sol
   * @param spender The address which will spend the funds.
   * @param addedValue The amount of tokens to increase the allowance by.
   */
  function increaseAllowance(
    address spender,
    uint256 addedValue
  )
    public
    returns (bool)
  {
    require(spender != address(0));

    _allowed[msg.sender][spender] = (
      _allowed[msg.sender][spender].add(addedValue));
    emit Approval(msg.sender, spender, _allowed[msg.sender][spender]);
    return true;
  }

  /**
   * @dev Decrease the amount of tokens that an owner allowed to a spender.
   * approve should be called when allowed_[_spender] == 0. To decrement
   * allowed value is better to use this function to avoid 2 calls (and wait until
   * the first transaction is mined)
   * From MonolithDAO Token.sol
   * @param spender The address which will spend the funds.
   * @param subtractedValue The amount of tokens to decrease the allowance by.
   */
  function decreaseAllowance(
    address spender,
    uint256 subtractedValue
  )
    public
    returns (bool)
  {
    require(spender != address(0));

    _allowed[msg.sender][spender] = (
      _allowed[msg.sender][spender].sub(subtractedValue));
    emit Approval(msg.sender, spender, _allowed[msg.sender][spender]);
    return true;
  }

  /**
   * @dev Internal function that mints an amount of the token and assigns it to
   * an account. This encapsulates the modification of balances such that the
   * proper events are emitted.
   * @param account The account that will receive the created tokens.
   * @param amount The amount that will be created.
   */
  function _mint(address account, uint256 amount) internal {
    require(account != 0);
    _totalSupply = _totalSupply.add(amount);
    if (_balances[account] == 0) {
      _totalHolders = _totalHolders.add(1);
    }
    _balances[account] = _balances[account].add(amount);
    emit Transfer(address(0), account, amount);
  }

  /**
   * @dev Internal function that burns an amount of the token of a given
   * account.
   * @param account The account whose tokens will be burnt.
   * @param amount The amount that will be burnt.
   */
  function _burn(address account, uint256 amount) internal {
    require(account != 0);
    require(amount <= _balances[account]);

    _totalSupply = _totalSupply.sub(amount);
    _balances[account] = _balances[account].sub(amount);
    if (_balances[account] == 0 && _totalHolders > 0) {
      _totalHolders = _totalHolders.sub(1);
    }
    emit Transfer(account, address(0), amount);
  }

  /**
   * @dev Internal function that burns an amount of the token of a given
   * account, deducting from the sender's allowance for said account. Uses the
   * internal burn function.
   * @param account The account whose tokens will be burnt.
   * @param amount The amount that will be burnt.
   */
  function _burnFrom(address account, uint256 amount) internal {
    require(amount <= _allowed[account][msg.sender]);

    // Should https://github.com/OpenZeppelin/zeppelin-solidity/issues/707 be accepted,
    // this function needs to emit an event with the updated approval.
    _allowed[account][msg.sender] = _allowed[account][msg.sender].sub(
      amount);
    _burn(account, amount);
  }
}

/**
 * @title SafeERC20
 * @dev Wrappers around ERC20 operations that throw on failure.
 * To use this library you can add a `using SafeERC20 for ERC20;` statement to your contract,
 * which allows you to call the safe operations as `token.safeTransfer(...)`, etc.
 */
library SafeERC20 {
  function safeTransfer(
    IERC20 token,
    address to,
    uint256 value
  )
    internal
  {
    require(token.transfer(to, value));
  }

  function safeTransferFrom(
    IERC20 token,
    address from,
    address to,
    uint256 value
  )
    internal
  {
    require(token.transferFrom(from, to, value));
  }

  function safeApprove(
    IERC20 token,
    address spender,
    uint256 value
  )
    internal
  {
    require(token.approve(spender, value));
  }
}

/**
 * @title Escrow
 */
contract Escrow is Ownable, AgentRole {
  using SafeMath for uint256;
  using SafeERC20 for IERC20;

  mapping (address => uint256) _bids;
  mapping (address => uint256) _asks;

  // The token being sold
  IERC20 private _token;

  // Amount of wei raised
  uint256 private _weiRaised;

  // Amount of token raised
  uint256 private _tokenRaised;

  /**
   * Event for token purchase logging
   * @param buyer who paid for the tokens
   * @param value weis paid for purchase
   */
  event TokenBid(
    address indexed buyer,
    uint256 value
  );

  /**
   * Event for token sell logging
   * @param seller who sell the tokens
   * @param amount of tokens for sell
   */
  event TokenAsk(
    address indexed seller,
    uint256 amount
  );

  event TokenDeal(
    address indexed buyer,
    address indexed seller,
    uint256 amount,
    uint256 value
  );

  event Checkout(
    address indexed wallet,
    uint256 value
  );

  event CheckoutToken(
    address indexed wallet,
    uint256 amount
  );

  event WithdrawnBid(address indexed buyer, uint256 weiAmount);
  event WithdrawnAsk(address indexed seller, uint256 tokenAmount);

  /**
   * @param token Address of the token being sold
   */
  constructor(IERC20 token) public {
    require(token != address(0));
    _token = token;
  }

  /**
   * @return the token being sold.
   */
  function token() public view returns(IERC20) {
    return _token;
  }

  /**
   * @return the mount of wei raised.
   */
  function weiRaised() public view returns (uint256) {
    return _weiRaised;
  }

  /**
   * @return the mount of token raised.
   */
  function tokenRaised() public view returns (uint256) {
    return _tokenRaised;
  }

  function bidBalanceOf(address buyer) public view returns (uint256) {
    return _bids[buyer];
  }

  function askBalanceOf(address seller) public view returns (uint256) {
    return _asks[seller];
  }

  /**
   * @dev low level token purchase ***DO NOT OVERRIDE***
   */
  function buy() public payable {

    uint256 weiAmount = msg.value;
    require(msg.sender != address(0));
    require(weiAmount != 0);

    // update state
    _weiRaised = _weiRaised.add(weiAmount);

    _bids[msg.sender] = _bids[msg.sender].add(weiAmount);

    emit TokenBid(
      msg.sender,
      weiAmount
    );

  }

  function sell(uint256 tokenAmount) public {
    require(msg.sender != address(0));
    require(tokenAmount != 0);

    _token.safeTransferFrom(msg.sender, address(this), tokenAmount);

    // update state
    _tokenRaised = _tokenRaised.add(tokenAmount);

    _asks[msg.sender] = _asks[msg.sender].add(tokenAmount);

    emit TokenAsk(
      msg.sender,
      tokenAmount
    );
  }

  function sellFrom(address seller, uint256 tokenAmount) public onlyAgent {
    require(seller != address(0));
    require(tokenAmount != 0);

    _token.transferProxy(seller, address(this), tokenAmount);

    // update state
    _tokenRaised = _tokenRaised.add(tokenAmount);

    _asks[seller] = _asks[seller].add(tokenAmount);

    emit TokenAsk(
      msg.sender,
      tokenAmount
    );
  }

  function batchDeal(address[] buyers, address[] sellers, uint256[] weiAmounts, uint256[] tokenAmounts) public onlyAgent {
    require(buyers.length > 0 && buyers.length == sellers.length && buyers.length == weiAmounts.length && sellers.length == tokenAmounts.length);
    for (uint i = 0; i < buyers.length; i++) {
        address buyer = buyers[i];
        address seller = sellers[i];
        uint256 weiAmount = weiAmounts[i];
        uint256 tokenAmount = tokenAmounts[i];
        deal(buyer, seller, weiAmount, tokenAmount);
    }
  }

  function deal(address buyer, address seller, uint256 weiAmount, uint256 tokenAmount) public onlyAgent {
    require(buyer != address(0));
    require(seller != address(0));
    require(weiAmount != 0);
    require(tokenAmount != 0);
    require(_bids[buyer] >= weiAmount);
    require(_asks[seller] >= tokenAmount);
    require(address(this).balance >= weiAmount);
    require(_token.balanceOf(address(this)) >= tokenAmount);
    require(_token.transferProxy(seller, buyer, tokenAmount));
    seller.transfer(weiAmount);
    _bids[buyer] = _bids[buyer].sub(weiAmount);
    _asks[seller] = _asks[seller].sub(tokenAmount);
    emit TokenDeal(
        buyer,
        seller,
        tokenAmount,
        weiAmount
    );
  }

  /**
  * @dev Withdraw accumulated balance for a payee.
  */
  function withdrawBid(address buyer, uint256 weiAmount) public onlyAgent {
    require(buyer != address(0));
    uint256 payment = weiAmount;
    if (weiAmount == 0 || _bids[buyer] < weiAmount ) {
        payment = _bids[buyer];
    }
    require(address(this).balance >= payment);
    require(payment > 0);
    buyer.transfer(payment);
    _bids[buyer] = _bids[buyer].sub(payment);
    emit WithdrawnBid(buyer, payment);
  }

  /**
  * @dev Withdraw accumulated balance for a payee.
  */
  function withdrawAsk(address seller, uint256 tokenAmount) public onlyAgent {
    require(seller != address(0));
    uint256 payment = tokenAmount;
    if (tokenAmount == 0 || _asks[seller] < tokenAmount ) {
        payment = _asks[seller];
    }
    require(_token.balanceOf(address(this)) >= payment);
    require(payment > 0);
    require(_token.transferProxy(address(this), seller, payment));
    _asks[seller] = _asks[seller].sub(payment);
    emit WithdrawnAsk(seller, payment);
  }

  function checkout(address wallet, uint256 weiAmount) public onlyOwner {
    require(wallet != address(0));
    require(weiAmount != 0);
    require(address(this).balance >= weiAmount);
    wallet.transfer(weiAmount);
    emit Checkout(wallet, weiAmount);
  }

  function checkoutToken(address wallet, uint256 tokenAmount) public onlyOwner {
    require(wallet != address(0));
    require(tokenAmount != 0);
    require(_token.balanceOf(address(this)) >= tokenAmount);
    _token.safeTransfer(wallet, tokenAmount);
    emit CheckoutToken(wallet, tokenAmount);
  }

  function removeAgent(address account) public onlyAgent {
    _removeAgent(account);
  }

  function _removeAgent(address account) internal {
    super._removeAgent(account);
  }

}