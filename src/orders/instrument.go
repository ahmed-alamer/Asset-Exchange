package orders

type Instrument struct {
	symbol string
	Name   string
}

func NewInstrument(symbol string, name string) *Instrument {
	return &Instrument{symbol, name}
}


