package orders

import "fmt"

type OrderType int

const (
	Ask OrderType = iota
	Bid
)

type Order struct {
	id        int64
	accountID int64
	price     float64
	volume    float64
	orderType OrderType
}

func NewOrderWithID(id int64, accountID int64, price float64, volume float64, orderType OrderType) *Order {
	order := NewOrder(accountID, price, volume, orderType)
	order.id = id

	return order
}

func NewOrder(accountID int64, price float64, volume float64, orderType OrderType) *Order {
	return &Order{-1, accountID, price, volume, orderType}
}

func (order Order) OrderType() OrderType {
	return order.orderType
}

func (order *Order) AdjustVolume(adjustedVolume float64) {
	order.volume = adjustedVolume
}

func (order *Order) SetID(id int64) {
	order.id = id
}

func (order Order) String() string {
	return fmt.Sprintf("Order: %s %d @ %f", orderTypeDisplayName(order.orderType), order.volume, order.price)
}

func (order Order) Type() OrderType {
	return order.orderType
}

func (order Order) Id() int64 {
	return order.id
}

func (order Order) Price() float64 {
	return order.price
}

func (order Order) Volume() float64 {
	return order.volume
}

func (order Order) AccountID() int64 {
	return order.id
}

func orderTypeDisplayName(orderType OrderType) string {
	if orderType == Ask {
		return "Ask"
	}

	return "Bid"
}
