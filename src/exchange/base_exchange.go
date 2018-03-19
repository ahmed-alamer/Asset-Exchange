package exchange

import (
	"time"
	"log"
	"math"
	"accounting"
	"orders"
)

type BaseExchangeService struct {
	askBook          *orders.OrderList
	bidBook          *orders.OrderList
	fills            []*Fill
	ledger           *accounting.Ledger
	orderIDGenerator int64
}

func (exchange *BaseExchangeService) SubmitAskOrder(order *orders.Order) int64 {
	order.SetID(exchange.generateOrderID())
	exchange.askBook.Add(order)

	exchange.execute()

	return order.Id()
}

func (exchange *BaseExchangeService) SubmitBidOrder(order *orders.Order) int64 {
	order.SetID(exchange.generateOrderID())

	exchange.bidBook.Add(order)

	exchange.execute()

	return order.Id()
}

func (exchange *BaseExchangeService) generateOrderID() int64 {
	exchange.orderIDGenerator++

	return exchange.orderIDGenerator
}

func (exchange *BaseExchangeService) CancelAskOrder(orderID int64) (bool, *orders.Order) {
	return exchange.askBook.RemoveByOrderId(orderID)
}

func (exchange *BaseExchangeService) CancelBidOrder(orderID int64) (bool, *orders.Order) {
	return exchange.bidBook.RemoveByOrderId(orderID)
}

func (exchange *BaseExchangeService) execute() {
	if exchange.askBook.IsEmpty() || exchange.bidBook.IsEmpty() {
		return
	}

	askOrder, bidOrder := exchange.askBook.Peek(), exchange.bidBook.Peek()

	if bidOrder.Price() != askOrder.Price() {
		return
	}

	remainingVolume := bidOrder.Volume() - askOrder.Volume()

	if remainingVolume > 0 {
		// bid larger than ask
		return
	}

	if remainingVolume < 0 {
		// ask larger than bid

		log.Printf("Processing Bid Match: %s", bidOrder)
		exchange.bidBook.Pop() // filled

		askOrder.AdjustVolume(math.Abs(remainingVolume))

	} else {
		// all cleared
		log.Println("Processing full match")
		exchange.bidBook.Pop() // filled
		exchange.askBook.Pop() // filled
	}

	exchange.fillOrder(askOrder.Price(), math.Abs(remainingVolume), askOrder.AccountID(), bidOrder.AccountID())
}

func (exchange *BaseExchangeService) fillOrder(price float64, volume float64, creditAccount int64, debitAccount int64) {
	fill := &Fill{
		Price:         price,
		Volume:        volume,
		ExecutionTime: time.Now(),
		CreditAccount: creditAccount,
		DebitAccount:  debitAccount,
	}

	exchange.fills = append(exchange.fills, fill)

	credit, debit := fill.LedgerEntries()
	exchange.ledger.AddEntry(credit)
	exchange.ledger.AddEntry(debit)

	log.Printf("WTF: %s", exchange.fills)
}

func NewExchange() *BaseExchangeService {
	return &BaseExchangeService{
		askBook:          orders.NewOrderList(),
		bidBook:          orders.NewOrderList(),
		fills:            make([]*Fill, 0),
		ledger:           accounting.NewLedger(),
		orderIDGenerator: 0,
	}
}

func (exchange *BaseExchangeService) Match() *orders.OrderList {
	fills := orders.NewOrderList()

	for _, bidOrder := range exchange.bidBook.Array() {
		for _, askOrder := range exchange.askBook.Array() {

			if bidOrder.Price() != askOrder.Price() {
				continue
			}

			remainingVolume := bidOrder.Volume() - askOrder.Volume()

			if remainingVolume > 0 {
				// bid larger than ask
				continue
			}

			if remainingVolume < 0 {
				// ask larger than bid
				log.Println("Filling Bid")
				exchange.bidBook.Remove(bidOrder) // filled

				askOrder.AdjustVolume(math.Abs(remainingVolume))

				fills.Add(bidOrder)

			} else {
				log.Println("Filling Ask & Bid")
				// all cleared
				exchange.bidBook.Remove(bidOrder) // filled
				exchange.askBook.Remove(askOrder) // filled

				fills.Add(askOrder)
				fills.Add(bidOrder)

			}

			exchange.fillOrder(bidOrder.Price(), math.Abs(bidOrder.Volume()), askOrder.Id(), bidOrder.Id())
		}

	}

	return fills
}

func (exchange BaseExchangeService) AskDepth() *orders.OrderList {
	return exchange.askBook
}

func (exchange BaseExchangeService) BidDepth() *orders.OrderList {
	return exchange.bidBook
}

func (exchange BaseExchangeService) Fills() []*Fill {
	return exchange.fills
}
