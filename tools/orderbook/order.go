package orderbook

import (
	"github.com/shopspring/decimal"
)

type Side = uint8

const (
	Ask Side = 1
	Bid Side = 2
)

type ProcessType = uint8

const (
	Market ProcessType = 0
	Limit  ProcessType = 1
)

type Order struct {
	Timestamp   int
	Quantity    decimal.Decimal
	Price       decimal.Decimal
	Id          uint64
	TradeId     uint64
	Side        Side
	ProcessType ProcessType
	Wallet      string
	nextOrder   *Order
	prevOrder   *Order
	OrderList   *OrderList
}

func (o *Order) Copy() *Order {
	return &Order{
		Timestamp:   o.Timestamp,
		Quantity:    o.Quantity,
		Price:       o.Price,
		Id:          o.Id,
		TradeId:     o.TradeId,
		Side:        o.Side,
		Wallet:      o.Wallet,
		ProcessType: o.ProcessType,
		nextOrder:   o.nextOrder,
		prevOrder:   o.prevOrder,
		OrderList:   o.OrderList,
	}
}

func (o *Order) NextOrder() *Order {
	return o.nextOrder
}

func (o *Order) PrevOrder() *Order {
	return o.prevOrder
}

func (o *Order) UpdateQuantity(newQuantity decimal.Decimal, newTimestamp int) {
	if newQuantity.GreaterThan(o.Quantity) && o.OrderList.tail_order != o {
		o.OrderList.MoveToTail(o)
	}
	o.OrderList.volume = o.OrderList.volume.Sub(o.Quantity.Sub(newQuantity))
	o.Timestamp = newTimestamp
	o.Quantity = newQuantity
}
