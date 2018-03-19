package exchange

import (
	"testing"
	"orders"
	"log"
)

func TestAskLessThanBidMatching(t *testing.T) {
	orderBook := NewLoggingExchange()

	orderBook.SubmitAskOrder(orders.NewOrder(1, 1.5, 10, orders.Limit))
	orderBook.SubmitBidOrder(orders.NewOrder(2, 1.5, 50, orders.Limit))

	fills := orderBook.Fills()

	if len(fills) != 0 {
		t.Errorf("Shouldn't have any fills")
	}
}

func TestAskLargerThanBidMatching(t *testing.T) {
	orderBook := NewLoggingExchange()

	orderBook.SubmitAskOrder(orders.NewOrder(1, 1.5, 50, orders.Limit))
	orderBook.SubmitBidOrder(orders.NewOrder(2, 1.5, 10, orders.Limit))

	fills := orderBook.Fills()

	log.Printf(" Fills: %s", fills)

	if len(fills) != 1 {
		t.Errorf("Incorrect fills!")
	}

	if orderBook.AskDepth().Size() != 1 {
		t.Errorf("Incorrec ask fills")
	}

	if orderBook.BidDepth().Size() != 0 {
		t.Errorf("Incorrect bid fills")
	}

}

func TestAskEqualsBidMatching(t *testing.T) {
	orderBook := NewLoggingExchange()

	orderBook.SubmitAskOrder(orders.NewOrder(1, 1.5, 10, orders.Limit))
	orderBook.SubmitBidOrder(orders.NewOrder(2, 1.5, 10, orders.Limit))

	fills := orderBook.Fills()

	if len(fills) != 1 {
		t.Errorf("Incorrect fills!")
	}
}

func TestCancelOrder(t *testing.T) {
	exchangeService := NewLoggingExchange()

	exchangeService.SubmitAskOrder(orders.NewOrder(1, 10, 10, orders.Limit))
	order2ID := exchangeService.SubmitAskOrder(orders.NewOrder(1, 10, 10, orders.Limit))

	cancelled, order := exchangeService.CancelAskOrder(order2ID)

	if !cancelled {
		t.Errorf("Didn't cancel order!")
	}

	if order2ID != order.Id() {
		t.Errorf("Didn't delete the correct order")
	}

	if exchangeService.AskDepth().Size() != 1 {
		t.Errorf("Wrong ask book size!")
	}
}
