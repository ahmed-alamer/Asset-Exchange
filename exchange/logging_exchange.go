package exchange

import (
	"Asset-Exchange/orders"
	"github.com/sirupsen/logrus"
	"os"
)

//TODO: Make this log text to the console and json to a file
type LoggingExchangeService struct {
	impl   CoreExchangeService
	logger *logrus.Logger
}

func (exchange LoggingExchangeService) CancelOrder(orderID int64, orderType orders.OrderType) (bool, *orders.Order) {
	exchange.logger.WithFields(logrus.Fields{
		"orderID": orderID,
		"type":    orderType,
	}).Info("Cancelling Order")

	return exchange.impl.CancelOrder(orderID, orderType)
}

func NewLoggingExchange() *LoggingExchangeService {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true}
	logger.SetLevel(logrus.DebugLevel)

	return &LoggingExchangeService{impl: NewExchange(), logger: logger}
}

func (exchange LoggingExchangeService) SubmitOrder(order *orders.Order) int64 {
	orderID := exchange.impl.SubmitOrder(order)
	exchange.logger.WithFields(logrus.Fields{
		"orderID": orderID,
		"volume":  order.Volume(),
		"price":   order.Price(),
	}).Info("Submitted new order")

	return orderID
}

func (exchange LoggingExchangeService) FillOrder(price float64, volume float64, creditAccount int64, debitAccount int64) {
	exchange.impl.FillOrder(price, volume, creditAccount, debitAccount)
	exchange.logger.WithFields(logrus.Fields{
		"price":  price,
		"vol":    volume,
		"credit": creditAccount,
		"debit":  debitAccount,
	}).Info("Executed Fill")
}

func (exchange LoggingExchangeService) AskDepth() *orders.OrderList {
	askBook := exchange.impl.AskDepth()
	exchange.logger.WithField("depth", askBook.Size()).Info("Retrieved Ask Book")

	return askBook
}

func (exchange LoggingExchangeService) BidDepth() *orders.OrderList {
	bidBook := exchange.impl.BidDepth()
	exchange.logger.WithField("depth", bidBook.Size()).Info("Retrieved Bid Book")

	return bidBook
}

func (exchange LoggingExchangeService) Fills() []*Fill {
	fills := exchange.impl.Fills()

	overallVolume, overallValue := 0.0, 0.0
	for _, fill := range fills {
		overallValue += fill.Price * float64(fill.Volume)
		overallVolume += fill.Volume
	}

	exchange.logger.WithFields(logrus.Fields{
		"total":  len(fills),
		"volume": overallVolume,
		"value":  overallValue,
	}).Info("Retrieved All Fills")

	return fills
}
