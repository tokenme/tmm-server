package orderbook

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

var limitOrders []*Order

func TestNewOrderBook(t *testing.T) {
	orderBook := NewOrderBook()

	if !(orderBook.VolumeAtPrice(Bid, decimal.Zero).Equal(decimal.Zero)) {
		t.Errorf("orderBook.VolumeAtPrice incorrect, got: %d, want: %d.", orderBook.VolumeAtPrice(Bid, decimal.Zero), decimal.Zero)
	}

	if !(orderBook.BestAsk().Equal(decimal.Zero)) {
		t.Errorf("orderBook.BestAsk incorrect, got: %d, want: %d.", orderBook.BestAsk(), decimal.Zero)
	}

	if !(orderBook.WorstBid().Equal(decimal.Zero)) {
		t.Errorf("orderBook.WorstBid incorrect, got: %d, want: %d.", orderBook.WorstBid(), decimal.Zero)
	}

	if !(orderBook.WorstAsk().Equal(decimal.Zero)) {
		t.Errorf("orderBook.WorstAsk incorrect, got: %d, want: %d.", orderBook.WorstAsk(), decimal.Zero)
	}

	if !(orderBook.BestBid().Equal(decimal.Zero)) {
		t.Errorf("orderBook.BestBid incorrect, got: %d, want: %d.", orderBook.BestBid(), decimal.Zero)
	}
}

func TestOrderBook(t *testing.T) {
	orderBook := NewOrderBook()

	fmt.Println(orderBook.BestAsk())

	dummyOrder := &Order{
		ProcessType: Limit,
		Side:        Ask,
		Quantity:    decimal.New(5, 0),
		Price:       decimal.New(101, 0),
		TradeId:     100,
	}

	limitOrders = append(limitOrders, dummyOrder)

	dummyOrder1 := &Order{
		ProcessType: Limit,
		Side:        Ask,
		Quantity:    decimal.New(5, 0),
		Price:       decimal.New(103, 0),
		TradeId:     101,
	}

	limitOrders = append(limitOrders, dummyOrder1)

	dummyOrder2 := &Order{
		ProcessType: Limit,
		Side:        Ask,
		Quantity:    decimal.New(5, 0),
		Price:       decimal.New(101, 0),
		TradeId:     102,
	}

	limitOrders = append(limitOrders, dummyOrder2)

	dummyOrder7 := &Order{
		ProcessType: Limit,
		Side:        Ask,
		Quantity:    decimal.New(5, 0),
		Price:       decimal.New(101, 0),
		TradeId:     103,
	}

	limitOrders = append(limitOrders, dummyOrder7)

	dummyOrder3 := &Order{
		ProcessType: Limit,
		Side:        Bid,
		Quantity:    decimal.New(5, 0),
		Price:       decimal.New(99, 0),
		TradeId:     100,
	}

	limitOrders = append(limitOrders, dummyOrder3)

	dummyOrder4 := &Order{
		ProcessType: Limit,
		Side:        Bid,
		Quantity:    decimal.New(5, 0),
		Price:       decimal.New(98, 0),
		TradeId:     101,
	}

	limitOrders = append(limitOrders, dummyOrder4)

	dummyOrder5 := &Order{
		ProcessType: Limit,
		Side:        Bid,
		Quantity:    decimal.New(5, 0),
		Price:       decimal.New(99, 0),
		TradeId:     102,
	}

	limitOrders = append(limitOrders, dummyOrder5)

	dummyOrder6 := &Order{
		ProcessType: Limit,
		Side:        Bid,
		Quantity:    decimal.New(5, 0),
		Price:       decimal.New(97, 0),
		TradeId:     103,
	}

	limitOrders = append(limitOrders, dummyOrder6)

	for _, order := range limitOrders {
		orderBook.ProcessOrder(order, true)
	}

	value, _ := decimal.NewFromString("101")
	if !(orderBook.BestAsk().Equal(value)) {
		t.Errorf("orderBook.BestAsk incorrect, got: %d, want: %d.", orderBook.BestAsk(), value)
	}

	value, _ = decimal.NewFromString("103")
	if !(orderBook.WorstAsk().Equal(value)) {
		t.Errorf("orderBook.WorstBid incorrect, got: %d, want: %d.", orderBook.WorstAsk(), value)
	}

	value, _ = decimal.NewFromString("99")
	if !(orderBook.BestBid().Equal(value)) {
		t.Errorf("orderBook.BestBid incorrect, got: %d, want: %d.", orderBook.BestBid(), value)
	}

	value, _ = decimal.NewFromString("97")
	if !(orderBook.WorstBid().Equal(value)) {
		t.Errorf("orderBook.BestBid incorrect, got: %d, want: %d.", orderBook.WorstBid(), value)
	}

	value, _ = decimal.NewFromString("15")
	pricePoint, _ := decimal.NewFromString("101")
	if !(orderBook.VolumeAtPrice(Ask, pricePoint).Equal(value)) {
		t.Errorf("orderBook.VolumeAtPrice incorrect, got: %d, want: %d.", orderBook.VolumeAtPrice(Bid, decimal.Zero), decimal.Zero)
	}

	//Submitting a limit order that crosses the opposing best price will result in a trade
	marketOrder := &Order{
		ProcessType: Limit,
		Side:        Bid,
		Quantity:    decimal.New(2, 0),
		Price:       decimal.New(102, 0),
		TradeId:     109,
	}

	trades, orderInBook := orderBook.ProcessOrder(marketOrder, true)

	tradedPrice := trades[0].Price
	tradedQuantity := trades[0].Quantity

	if !(tradedPrice.IntPart() == 101 && tradedQuantity.IntPart() == 2 && orderInBook == nil) {
		t.Errorf("orderBook.ProcessOrder incorrect, tradePrice:%s, tradedQuantity: %s", tradedPrice.String(), tradedQuantity.String())
	}

	// If a limit crosses but is only partially matched, the remaning volume will
	// be placed in the book as an outstanding order
	bigOrder := &Order{
		ProcessType: Limit,
		Side:        Bid,
		Quantity:    decimal.New(50, 0),
		Price:       decimal.New(102, 0),
		TradeId:     111,
	}

	trades, orderInBook = orderBook.ProcessOrder(bigOrder, true)

	fmt.Println(trades)
	fmt.Println(orderInBook)

	if orderInBook == nil {
		t.Errorf("orderBook.ProcessOrder incorrect")
	}

	// Market orders only require that a user specifies a side (bid or ask), a quantity, and their unique trade id
	marketOrder = &Order{
		ProcessType: Market,
		Side:        Ask,
		Quantity:    decimal.New(20, 0),
		TradeId:     111,
	}
	trades, orderInBook = orderBook.ProcessOrder(marketOrder, true)

}
