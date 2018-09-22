package orderbook

import (
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/tools/orderbook/extend"
)

type Comparator func(a, b interface{}) int

func decimalComparator(a, b interface{}) int {
	aAsserted := a.(decimal.Decimal)
	bAsserted := b.(decimal.Decimal)
	switch {
	case aAsserted.GreaterThan(bAsserted):
		return 1
	case aAsserted.LessThan(bAsserted):
		return -1
	default:
		return 0
	}
}

type OrderTree struct {
	priceTree *redblacktreeextended.RedBlackTreeExtended
	priceMap  map[string]*OrderList // Dictionary containing price : OrderList object
	orderMap  map[uint64]*Order     // Dictionary containing order_id : Order object
	volume    decimal.Decimal       // Contains total quantity from all Orders in tree
	numOrders int                   // Contains count of Orders in tree
	depth     int                   // Number of different prices in tree (http://en.wikipedia.org/wiki/Order_book_(trading)#Book_depth)
}

func NewOrderTree() *OrderTree {
	priceTree := &redblacktreeextended.RedBlackTreeExtended{rbt.NewWith(decimalComparator)}
	priceMap := make(map[string]*OrderList)
	orderMap := make(map[uint64]*Order)
	return &OrderTree{priceTree, priceMap, orderMap, decimal.Zero, 0, 0}
}

func (ordertree *OrderTree) Length() int {
	return len(ordertree.orderMap)
}

func (ordertree *OrderTree) Order(orderId uint64) *Order {
	return ordertree.orderMap[orderId]
}

func (ordertree *OrderTree) PriceList(price decimal.Decimal) *OrderList {
	return ordertree.priceMap[price.String()]
}

func (ordertree *OrderTree) CreatePrice(price decimal.Decimal) {
	ordertree.depth = ordertree.depth + 1
	newList := NewOrderList(price)
	ordertree.priceTree.Put(price, newList)
	ordertree.priceMap[price.String()] = newList
}

func (ordertree *OrderTree) RemovePrice(price decimal.Decimal) {
	ordertree.depth = ordertree.depth - 1
	ordertree.priceTree.Remove(price)
	delete(ordertree.priceMap, price.String())
}

func (ordertree *OrderTree) PriceExist(price decimal.Decimal) bool {
	if _, ok := ordertree.priceMap[price.String()]; ok {
		return true
	}
	return false
}

func (ordertree *OrderTree) OrderExist(orderId uint64) bool {
	if _, ok := ordertree.orderMap[orderId]; ok {
		return true
	}
	return false
}

func (ordertree *OrderTree) RemoveOrderById(orderId uint64) {
	ordertree.numOrders = ordertree.numOrders - 1
	order := ordertree.orderMap[orderId]
	ordertree.volume = ordertree.volume.Sub(order.Quantity)
	order.OrderList.RemoveOrder(order)
	if order.OrderList.Length() == 0 {
		ordertree.RemovePrice(order.Price)
	}
	delete(ordertree.orderMap, orderId)
}

func (ordertree *OrderTree) MaxPrice() decimal.Decimal {
	if ordertree.depth > 0 {
		value, found := ordertree.priceTree.GetMax()
		if found {
			return value.(*OrderList).price
		}
		return decimal.Zero

	} else {
		return decimal.Zero
	}
}

func (ordertree *OrderTree) MinPrice() decimal.Decimal {
	if ordertree.depth > 0 {
		value, found := ordertree.priceTree.GetMin()
		if found {
			return value.(*OrderList).price
		} else {
			return decimal.Zero
		}

	} else {
		return decimal.Zero
	}
}

func (ordertree *OrderTree) MaxPriceList() *OrderList {
	if ordertree.depth > 0 {
		price := ordertree.MaxPrice()
		return ordertree.priceMap[price.String()]
	}
	return nil

}

func (ordertree *OrderTree) MinPriceList() *OrderList {
	if ordertree.depth > 0 {
		price := ordertree.MinPrice()
		return ordertree.priceMap[price.String()]
	}
	return nil
}

func (ordertree *OrderTree) InsertOrder(order *Order) {

	if ordertree.OrderExist(order.Id) {
		ordertree.RemoveOrderById(order.Id)
	}
	ordertree.numOrders++

	if !ordertree.PriceExist(order.Price) {
		ordertree.CreatePrice(order.Price)
	}
	order.OrderList = ordertree.priceMap[order.Price.String()]
	ordertree.priceMap[order.Price.String()].AppendOrder(order)
	ordertree.orderMap[order.Id] = order
	ordertree.volume = ordertree.volume.Add(order.Quantity)
}

func (ordertree *OrderTree) UpdateOrder(newOrder *Order) {
	order := ordertree.orderMap[newOrder.Id]
	originalQuantity := order.Quantity

	if !newOrder.Price.Equal(order.Price) {
		// Price changed. Remove order and update tree.
		orderList := ordertree.priceMap[order.Price.String()]
		orderList.RemoveOrder(order)
		if orderList.Length() == 0 {
			ordertree.RemovePrice(newOrder.Price)
		}
		ordertree.InsertOrder(newOrder)
	} else {
		order.UpdateQuantity(newOrder.Quantity, newOrder.Timestamp)
	}
	ordertree.volume = ordertree.volume.Add(order.Quantity.Sub(originalQuantity))
}
