package orders

import (
	"log"
	"testing"
)

func TestOrderList(t *testing.T) {
	orderList := NewOrderList()
	orderList.Add(NewOrderWithID(1, 1, 10, 1, Ask))
	orderList.Add(NewOrderWithID(2, 2, 10, 1, Ask))
	orderList.Add(NewOrderWithID(3, 3, 11, 1, Bid))

	for _, order := range orderList.Array() {
		log.Println(order)

		if order.Id() == 2 {
			orderList.Remove(order)
		}
	}

	expected := []*Order{
		NewOrderWithID(1, 1, 10, 1, Ask),
		NewOrderWithID(3, 2, 11, 1, Bid),
	}

	if orderList.Size() != len(expected) {
		t.Errorf("Incorrect list size: %d", orderList.Size())
	}

	for i := 0; i < len(expected); i++ {
		expectedOrder := expected[i]
		actualOrder := orderList.Array()[i]

		if expectedOrder.Id() != actualOrder.Id() {
			t.Errorf("Incorrect oreder!")
		}
	}
}

func TestOrdersStack(t *testing.T) {
	orderList := NewOrderList()
	orderList.Insert(NewOrderWithID(1, 1, 10, 1, Ask))
	orderList.Insert(NewOrderWithID(2, 2, 9, 1, Ask))
	orderList.Insert(NewOrderWithID(3, 3, 11, 1, Bid))

	log.Println(orderList.Size())
	for !orderList.IsEmpty() {
		order := orderList.Pop()
		log.Println(order)
	}
}
