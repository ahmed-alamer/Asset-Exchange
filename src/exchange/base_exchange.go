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

func (exchange *BaseExchangeService) generateOrderID() int64 {
	exchange.orderIDGenerator++

	return exchange.orderIDGenerator
}

func (exchange *BaseExchangeService) SubmitOrder(order *orders.Order) int64 {
	order.SetID(exchange.generateOrderID())
	if order.OrderType() == orders.Ask {
		exchange.askBook.Add(order)
	} else {
		exchange.bidBook.Add(order)
	}

	exchange.Execute()

	return order.Id()
}

func (exchange *BaseExchangeService) CancelOrder(orderID int64, orderType orders.OrderType) (bool, *orders.Order) {
	var targetBook *orders.OrderList
	if orderType == orders.Ask {
		targetBook = exchange.askBook
	} else {
		targetBook = exchange.bidBook
	}

	return targetBook.RemoveByOrderId(orderID)
}

func (exchange *BaseExchangeService) Execute() {
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

	exchange.FillOrder(askOrder.Price(), math.Abs(remainingVolume), askOrder.AccountID(), bidOrder.AccountID())
}

func (exchange *BaseExchangeService) FillOrder(price float64, volume float64, creditAccount int64, debitAccount int64) {
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

			exchange.FillOrder(bidOrder.Price(), math.Abs(bidOrder.Volume()), askOrder.Id(), bidOrder.Id())
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
