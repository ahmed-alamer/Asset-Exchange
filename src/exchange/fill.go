package exchange

import (
	"time"
	"accounting"
)

type Fill struct {
	Price         float64
	Volume        float64
	ExecutionTime time.Time
	CreditAccount int64
	DebitAccount  int64
}

func (fill Fill) LedgerEntries() (credit *accounting.LedgerEntry, debit *accounting.LedgerEntry) {
	return accounting.NewCreditEntry(fill.CreditAccount, fill.Volume), accounting.NewDebitEntry(fill.DebitAccount, fill.Volume)
}
