package exchange

import (
	"Asset-Exchange/orders"
)

type CoreExchangeService interface {
	SubmitOrder(order *orders.Order) int64
	CancelOrder(orderID int64, orderType orders.OrderType) (bool, *orders.Order)
	FillOrder(price float64, volume float64, creditAccount int64, debitAccount int64)
	AskDepth() *orders.OrderList
	BidDepth() *orders.OrderList
	Fills() []*Fill
}
