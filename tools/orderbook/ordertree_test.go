package orderbook

import (
	"github.com/shopspring/decimal"
	"testing"
)

var testTimestamp = 123452342343
var testQuanity, _ = decimal.NewFromString("0.1")
var testPrice, _ = decimal.NewFromString("0.1")
var testOrderId uint64 = 1
var testTradeId uint64 = 1

var testTimestamp1 = 123452342345
var testQuanity1, _ = decimal.NewFromString("0.2")
var testPrice1, _ = decimal.NewFromString("0.1")
var testOrderId1 uint64 = 2
var testTradeId1 uint64 = 2

var testTimestamp2 = 123452342340
var testQuanity2, _ = decimal.NewFromString("0.2")
var testPrice2, _ = decimal.NewFromString("0.3")
var testOrderId2 uint64 = 3
var testTradeId2 uint64 = 3

var testTimestamp3 = 1234523
var testQuanity3, _ = decimal.NewFromString("200.0")
var testPrice3, _ = decimal.NewFromString("1.3")
var testOrderId3 uint64 = 3
var testTradeId3 uint64 = 3

func TestNewOrderTree(t *testing.T) {
	orderTree := NewOrderTree()

	dummyOrder := &Order{
		Timestamp: testTimestamp,
		Quantity:  testQuanity,
		Price:     testPrice,
		Id:        testOrderId,
		TradeId:   testTradeId,
	}

	dummyOrder1 := &Order{
		Timestamp: testTimestamp1,
		Quantity:  testQuanity1,
		Price:     testPrice1,
		Id:        testOrderId1,
		TradeId:   testTradeId1,
	}

	dummyOrder2 := &Order{
		Timestamp: testTimestamp2,
		Quantity:  testQuanity2,
		Price:     testPrice2,
		Id:        testOrderId2,
		TradeId:   testTradeId2,
	}

	dummyOrder3 := &Order{
		Timestamp: testTimestamp3,
		Quantity:  testQuanity3,
		Price:     testPrice3,
		Id:        testOrderId3,
		TradeId:   testTradeId3,
	}

	if !(orderTree.volume.Equal(decimal.Zero)) {
		t.Errorf("orderTree.volume incorrect, got: %d, want: %d.", orderTree.volume, decimal.Zero)
	}

	if !(orderTree.Length() == 0) {
		t.Errorf("orderTree.Length() incorrect, got: %d, want: %d.", orderTree.Length(), 0)
	}

	orderTree.InsertOrder(dummyOrder)
	orderTree.InsertOrder(dummyOrder1)

	if !(orderTree.PriceExist(testPrice)) {
		t.Errorf("orderTree.numOrders incorrect, got: %d, want: %d.", orderTree.numOrders, 2)
	}

	if !(orderTree.PriceExist(testPrice1)) {
		t.Errorf("orderTree.numOrders incorrect, got: %d, want: %d.", orderTree.numOrders, 2)
	}

	if !(orderTree.Length() == 2) {
		t.Errorf("orderTree.numOrders incorrect, got: %d, want: %d.", orderTree.numOrders, 2)
	}

	orderTree.RemoveOrderById(dummyOrder1.Id)
	orderTree.RemoveOrderById(dummyOrder.Id)

	if !(orderTree.Length() == 0) {
		t.Errorf("orderTree.numOrders incorrect, got: %d, want: %d.", orderTree.numOrders, 2)
	}

	orderTree.InsertOrder(dummyOrder)
	orderTree.InsertOrder(dummyOrder1)
	orderTree.InsertOrder(dummyOrder2)
	orderTree.InsertOrder(dummyOrder3)

	if !(orderTree.MaxPrice().Equal(testPrice3)) {
		t.Errorf("orderTree.MaxPrice incorrect, got: %d, want: %d.", orderTree.MaxPrice(), testPrice3)
	}

	if !(orderTree.MinPrice().Equal(testPrice)) {
		t.Errorf("orderTree.MinPrice incorrect, got: %d, want: %d.", orderTree.MinPrice(), testPrice)
	}

	orderTree.RemovePrice(testPrice)

	if orderTree.PriceExist(testPrice) {
		t.Errorf("orderTree.MinPrice incorrect, got: %d, want: %d.", orderTree.MinPrice(), testPrice)
	}

	// TODO Check PriceList as well and verify with the orders
}
