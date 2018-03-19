package exchange

import (
	"orders"
)

type CoreExchangeService interface {
	SubmitAskOrder(order *orders.Order) int64
	SubmitBidOrder(order *orders.Order) int64
	CancelAskOrder(orderID int64) (bool, *orders.Order)
	CancelBidOrder(orderID int64) (bool, *orders.Order)
	AskDepth() *orders.OrderList
	BidDepth() *orders.OrderList
	Fills() []*Fill

	fillOrder(price float64, volume float64, creditAccount int64, debitAccount int64)
}
