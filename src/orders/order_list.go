package orders

import (
	"strings"
)

const LastItemIndex = -1

type OrderList struct {
	backingSlice []*Order
	size         int
}

func NewOrderList() *OrderList {
	return &OrderList{}
}

func (list *OrderList) Add(order *Order) {
	if list.backingSlice != nil {
		list.backingSlice = append(list.backingSlice, order)
		list.size += 1
	} else {
		list.init(order)
	}

}

func (list *OrderList) Insert(order *Order) {
	if list.backingSlice == nil {
		list.init(order)
		return
	}

	if index := list.getInsertionIndex(order); index == LastItemIndex {
		list.backingSlice = append(list.backingSlice, order)
	} else {
		list.backingSlice = append(list.backingSlice, nil)

		copy(list.backingSlice[index+1:], list.backingSlice[index:]) // Very Erlang! idk! fuck!
		list.backingSlice[index] = order
	}

	list.size++
}

func (list *OrderList) Pop() *Order {
	var top *Order = nil
	if list.size > 1 {
		top, list.backingSlice = list.backingSlice[0], list.backingSlice[1:]
	} else {
		top, list.backingSlice = list.backingSlice[0], list.backingSlice[:0]
	}

	list.size--

	return top
}

func (list *OrderList) Peek() Order {
	return *list.backingSlice[0]
}

func (list *OrderList) getInsertionIndex(order *Order) int {
	for index, existingOrder := range list.backingSlice {
		if order.price > existingOrder.price {
			return index
		}
	}

	return LastItemIndex
}

func (list *OrderList) Remove(order *Order) {
	list.RemoveByOrderId(order.id)
}

func (list *OrderList) RemoveByOrderId(orderId int64) (removed bool, order *Order) {
	backingListPointer := list.backingSlice[:0]
	for _, existingOrder := range list.backingSlice {
		if orderId == existingOrder.id {
			list.size--
			order = existingOrder
			removed = true

			continue
		}

		backingListPointer = append(backingListPointer, existingOrder)
	}

	list.backingSlice = backingListPointer

	return removed, order
}

func (list *OrderList) Array() []*Order {
	return list.backingSlice
}

func (list *OrderList) Size() int {
	return list.size
}

func (list *OrderList) IsEmpty() bool {
	return list.Size() == 0
}

func (list OrderList) String() string {
	display := make([]string, len(list.backingSlice))

	for i, order := range list.backingSlice {
		display[i] = order.String()
	}

	return strings.Join(display, "\t")
}

func (list *OrderList) init(order *Order) {
	list.backingSlice = make([]*Order, 1)
	list.backingSlice[0] = order
	list.size = 1
}
