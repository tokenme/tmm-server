pragma solidity ^0.4.24;

import "IERC20.sol";
import "SafeMath.sol";
import "SafeERC20.sol";
import "Ownable.sol";
import "AgentRole.sol";

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