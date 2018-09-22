package orderbook

import (
	"github.com/shopspring/decimal"
	"testing"
)

func TestNewOrderList(t *testing.T) {
	orderList := NewOrderList(testPrice)

	if !(orderList.length == 0) {
		t.Errorf("Orderlist Length incorrect, got: %d, want: %d.", orderList.length, 0)
	}

	if !(orderList.price.Equal(testPrice)) {
		t.Errorf("Orderlist price incorrect, got: %d, want: %d.", orderList.length, 0)
	}

	if !(orderList.volume.Equal(decimal.Zero)) {
		t.Errorf("Orderlist Length incorrect, got: %d, want: %d.", orderList.length, 0)
	}
}

func TestOrderList(t *testing.T) {
	orderList := NewOrderList(testPrice)

	var order_list OrderList
	dummyOrder := &Order{
		Id:        testOrderId,
		Timestamp: testTimestamp,
		Quantity:  testQuanity,
		Price:     testPrice,
		TradeId:   testTradeId,
		OrderList: &order_list,
	}

	orderList.AppendOrder(dummyOrder)

	if !(orderList.Length() == 1) {
		t.Errorf("Orderlist Length incorrect, got: %d, want: %d.", orderList.length, 0)
	}

	if !(orderList.price.Equal(testPrice)) {
		t.Errorf("Orderlist price incorrect, got: %s, want: %s.", orderList.price.String(), testPrice.String())
	}

	if !(orderList.volume.Equal(dummyOrder.Quantity)) {
		t.Errorf("Orderlist Length incorrect, got: %d, want: %d.", orderList.length, 0)
	}

	if !(orderList.volume.Equal(dummyOrder.Quantity)) {
		t.Errorf("Orderlist Length incorrect, got: %d, want: %d.", orderList.length, 0)
	}

	dummyOrder1 := &Order{
		Id:        testOrderId1,
		Timestamp: testTimestamp1,
		Quantity:  testQuanity1,
		Price:     testPrice1,
		TradeId:   testTradeId1,
		OrderList: &order_list,
	}

	orderList.AppendOrder(dummyOrder1)

	if !(orderList.Length() == 2) {
		t.Errorf("Orderlist Length incorrect, got: %d, want: %d.", orderList.length, 0)
	}

	if !(orderList.volume.Equal(dummyOrder.Quantity.Add(dummyOrder1.Quantity))) {
		t.Errorf("Orderlist Length incorrect, got: %d, want: %d.", orderList.length, 0)
	}

	headOrder := orderList.HeadOrder()
	if !(headOrder.Id == 1) {
		t.Errorf("headorder id incorrect, got: %d want: %d.", headOrder.Id, 0)
	}

	nextOrder := headOrder.NextOrder()

	if !(nextOrder.Id == 2) {
		t.Errorf("Next headorder id incorrect, got: %d, want: %d.", headOrder.NextOrder().Id, 2)
	}
}
