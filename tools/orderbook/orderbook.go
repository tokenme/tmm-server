package orderbook

import (
	"fmt"
	"github.com/shopspring/decimal"
	lane "gopkg.in/oleiade/lane.v1"
)

type OrderBook struct {
	deque       *lane.Deque
	bids        *OrderTree
	asks        *OrderTree
	time        int
	nextOrderID uint64
}

func NewOrderBook() *OrderBook {
	deque := lane.NewDeque()
	bids := NewOrderTree()
	asks := NewOrderTree()
	return &OrderBook{deque, bids, asks, 0, 0}
}

func (orderBook *OrderBook) UpdateTime() {
	orderBook.time++
}

func (orderBook *OrderBook) BestBid() (value decimal.Decimal) {
	value = orderBook.bids.MaxPrice()
	return
}

func (orderBook *OrderBook) BestAsk() (value decimal.Decimal) {
	value = orderBook.asks.MinPrice()
	return
}

func (orderBook *OrderBook) WorstBid() (value decimal.Decimal) {
	value = orderBook.bids.MinPrice()
	return
}

func (orderBook *OrderBook) WorstAsk() (value decimal.Decimal) {
	value = orderBook.asks.MaxPrice()
	return
}

func (orderBook *OrderBook) ProcessMarketOrder(order *Order, verbose bool) []*Trade {
	var (
		trades    []*Trade
		newTrades []*Trade
	)
	quantityToTrade := order.Quantity
	if order.Side == Bid {
		for quantityToTrade.GreaterThan(decimal.Zero) && orderBook.asks.Length() > 0 {
			bestPriceAsks := orderBook.asks.MinPriceList()
			quantityToTrade, newTrades = orderBook.ProcessOrderList(Ask, bestPriceAsks, quantityToTrade, order, verbose)
			trades = append(trades, newTrades...)
		}
	} else if order.Side == Ask {
		for quantityToTrade.GreaterThan(decimal.Zero) && orderBook.bids.Length() > 0 {
			bestPriceBids := orderBook.bids.MaxPriceList()
			quantityToTrade, newTrades = orderBook.ProcessOrderList(Bid, bestPriceBids, quantityToTrade, order, verbose)
			trades = append(trades, newTrades...)
		}
	}
	return trades
}

func (orderBook *OrderBook) ProcessLimitOrder(order *Order, verbose bool) ([]*Trade, *Order) {
	var (
		trades      []*Trade
		newTrades   []*Trade
		orderInBook *Order
	)

	quantityToTrade := order.Quantity
	price := order.Price

	if order.Side == Bid {
		minPrice := orderBook.asks.MinPrice()
		for quantityToTrade.GreaterThan(decimal.Zero) && orderBook.asks.Length() > 0 && price.GreaterThanOrEqual(minPrice) {
			bestPriceAsks := orderBook.asks.MinPriceList()
			quantityToTrade, newTrades = orderBook.ProcessOrderList(Ask, bestPriceAsks, quantityToTrade, order, verbose)
			trades = append(trades, newTrades...)
			minPrice = orderBook.asks.MinPrice()
		}

		if quantityToTrade.GreaterThan(decimal.Zero) {
			newOrder := order.Copy()
			newOrder.Id = orderBook.nextOrderID
			newOrder.Quantity = quantityToTrade
			orderBook.bids.InsertOrder(newOrder)
			orderInBook = newOrder
		}

	} else if order.Side == Ask {
		maxPrice := orderBook.bids.MaxPrice()
		for quantityToTrade.GreaterThan(decimal.Zero) && orderBook.bids.Length() > 0 && price.LessThanOrEqual(maxPrice) {
			bestPriceBids := orderBook.bids.MaxPriceList()
			quantityToTrade, newTrades = orderBook.ProcessOrderList(Bid, bestPriceBids, quantityToTrade, order, verbose)
			trades = append(trades, newTrades...)
			maxPrice = orderBook.bids.MaxPrice()
		}

		if quantityToTrade.GreaterThan(decimal.Zero) {
			newOrder := order.Copy()
			newOrder.Id = orderBook.nextOrderID
			newOrder.Quantity = quantityToTrade
			orderBook.asks.InsertOrder(newOrder)
			orderInBook = newOrder
		}
	}
	return trades, orderInBook
}

func (orderBook *OrderBook) ProcessOrder(order *Order, verbose bool) ([]*Trade, *Order) {
	var orderInBook *Order
	var trades []*Trade

	orderBook.UpdateTime()
	order.Timestamp = orderBook.time
	orderBook.nextOrderID++

	if order.ProcessType == Market {
		trades = orderBook.ProcessMarketOrder(order, verbose)
	} else {
		trades, orderInBook = orderBook.ProcessLimitOrder(order, verbose)
	}
	return trades, orderInBook
}

func (orderBook *OrderBook) ProcessOrderList(side Side, orderList *OrderList, quantityStillToTrade decimal.Decimal, order *Order, verbose bool) (decimal.Decimal, []*Trade) {
	quantityToTrade := quantityStillToTrade
	var trades []*Trade

	for orderList.Length() > 0 && quantityToTrade.GreaterThan(decimal.Zero) {
		headOrder := orderList.HeadOrder()
		tradedPrice := headOrder.Price
		// counterParty := headOrder.trade_id
		var newBookQuantity decimal.Decimal
		var tradedQuantity decimal.Decimal

		if quantityToTrade.LessThan(headOrder.Quantity) {
			tradedQuantity = quantityToTrade
			// Do the transaction
			newBookQuantity = headOrder.Quantity.Sub(quantityToTrade)
			headOrder.UpdateQuantity(newBookQuantity, headOrder.Timestamp)
			quantityToTrade = decimal.Zero
		} else if quantityToTrade.Equal(headOrder.Quantity) {
			tradedQuantity = quantityToTrade
			if side == Bid {
				orderBook.bids.RemoveOrderById(headOrder.Id)
			} else {
				orderBook.asks.RemoveOrderById(headOrder.Id)
			}
			quantityToTrade = decimal.Zero
		} else {
			tradedQuantity = headOrder.Quantity
			if side == Bid {
				orderBook.bids.RemoveOrderById(headOrder.Id)
			} else {
				orderBook.asks.RemoveOrderById(headOrder.Id)
			}
		}

		if verbose {
			fmt.Printf("TRADE: Time - %v, Price - %v, Quantity - %v, TradeID - %v, Matching TradeID - %v\n", orderBook.time, tradedPrice.String(), tradedQuantity.String(), headOrder.TradeId, order.TradeId)
		}

		transactionRecord := &Trade{
			Id:                 order.TradeId,
			Wallet:             order.Wallet,
			CounterParty:       headOrder.TradeId,
			CounterPartyWallet: headOrder.Wallet,
			Timestamp:          orderBook.time,
			Price:              tradedPrice,
			Quantity:           tradedQuantity,
			Side:               order.Side,
		}

		orderBook.deque.Append(transactionRecord)
		trades = append(trades, transactionRecord)
	}
	return quantityToTrade, trades
}

func (orderBook *OrderBook) CancelOrder(side Side, orderId uint64) {
	orderBook.UpdateTime()

	if side == Bid {
		if orderBook.bids.OrderExist(orderId) {
			orderBook.bids.RemoveOrderById(orderId)
		}
	} else {
		if orderBook.asks.OrderExist(orderId) {
			orderBook.asks.RemoveOrderById(orderId)
		}
	}
}

func (orderBook *OrderBook) ModifyOrder(newOrder *Order) {
	orderBook.UpdateTime()
	newOrder.Timestamp = orderBook.time

	if newOrder.Side == Bid {
		if orderBook.bids.OrderExist(newOrder.Id) {
			orderBook.bids.UpdateOrder(newOrder)
		}
	} else {
		if orderBook.asks.OrderExist(newOrder.Id) {
			orderBook.asks.UpdateOrder(newOrder)
		}
	}
}

func (orderBook *OrderBook) VolumeAtPrice(side Side, price decimal.Decimal) decimal.Decimal {
	if side == Bid {
		volume := decimal.Zero
		if orderBook.bids.PriceExist(price) {
			volume = orderBook.bids.PriceList(price).volume
		}

		return volume

	} else {
		volume := decimal.Zero
		if orderBook.asks.PriceExist(price) {
			volume = orderBook.asks.PriceList(price).volume
		}
		return volume
	}
}
