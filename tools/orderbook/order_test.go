package orderbook

import (
	"testing"
)

func TestNewOrder(t *testing.T) {
	var orderList OrderList
	order := &Order{
		Timestamp: testTimestamp,
		Quantity:  testQuanity,
		Price:     testPrice,
		Id:        testOrderId,
		TradeId:   testTradeId,
		OrderList: &orderList,
	}

	if !(order.Timestamp == testTimestamp) {
		t.Errorf("Timesmape incorrect, got: %d, want: %d.", order.Timestamp, testTimestamp)
	}

	if !(order.Quantity.Equal(testQuanity)) {
		t.Errorf("quantity incorrect, got: %d, want: %d.", order.Quantity, testQuanity)
	}

	if !(order.Price.Equal(testPrice)) {
		t.Errorf("price incorrect, got: %d, want: %d.", order.Price, testPrice)
	}

	if !(order.Id == testOrderId) {
		t.Errorf("order id incorrect, got: %d, want: %d.", order.Id, testOrderId)
	}

	if !(order.TradeId == testTradeId) {
		t.Errorf("trade id incorrect, got: %d, want: %d.", order.TradeId, testTradeId)
	}
}

func TestOrder(t *testing.T) {
	orderList := NewOrderList(testPrice)

	order := &Order{
		Timestamp: testTimestamp,
		Quantity:  testQuanity,
		Price:     testPrice,
		Id:        testOrderId,
		TradeId:   testTradeId,
		OrderList: orderList,
	}

	orderList.AppendOrder(order)

	order.UpdateQuantity(testQuanity1, testTimestamp1)

	if !(order.Quantity.Equal(testQuanity1)) {
		t.Errorf("order id incorrect, got: %d, want: %d.", order.Id, testOrderId)
	}

	if !(order.Timestamp == testTimestamp1) {
		t.Errorf("trade id incorrect, got: %d, want: %d.", order.TradeId, testTradeId)
	}
}
