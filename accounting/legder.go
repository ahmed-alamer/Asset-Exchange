package accounting

type LedgerEntry struct {
	id        int64
	accountID int64
	credit    float64
	debit     float64
}

type Ledger struct {
	idGenerator int64
	entries     []*LedgerEntry
}

func NewLedger() *Ledger {
	return &Ledger{idGenerator: 1, entries: make([]*LedgerEntry, 0)}
}

func NewCreditEntry(accountID int64, volume float64) *LedgerEntry {
	return &LedgerEntry{
		accountID: accountID,
		credit:    volume,
		debit:     0,
	}

}

func NewDebitEntry(accountID int64, volume float64) *LedgerEntry {
	return &LedgerEntry{
		accountID: accountID,
		credit:    0,
		debit:     volume,
	}
}

func (ledger *Ledger) AddEntry(entry *LedgerEntry) {
	ledger.idGenerator++
	entry.id = ledger.idGenerator
	ledger.entries = append(ledger.entries, entry)
}

func (ledger Ledger) GetAccountBalance(accountID int64) float64 {
	var credit, debit float64 = 0, 0

	for _, entry := range ledger.entries {
		if entry.accountID == accountID {
			credit += entry.credit
			debit += entry.debit
		}
	}

	return credit - debit
}

func (ledger *Ledger) CreditAccount(accountID int64, amount float64) {
	ledger.AddEntry(NewCreditEntry(accountID, amount))
}

func (ledger *Ledger) DebitAccount(accountID int64, amount float64) {
	ledger.AddEntry(NewDebitEntry(accountID, amount))
}
