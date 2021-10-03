package parser

import (
	"github.com/DawnKosmos/ftxwebapp/exchange"
)

type AmountType int

const (
	COIN AmountType = iota
	FIAT
	ACCOUNTSIZE
	POSITIONSIZE
)

type Amount struct {
	Ticker string
	Type   AmountType
	Value  float64
}

func (a *Amount) Evaluate(f exchange.Exchange) (float64, error) {
	switch a.Type {
	case COIN:
		return a.Value, nil
	case FIAT:
		m, err := f.MarketPrice(a.Ticker)
		return a.Value / m, err
	case ACCOUNTSIZE:
		m, err := f.MarketPrice(a.Ticker)
		if err != nil {
			return a.Value, err
		}
		collateral, err := f.FreeCollateral()
		az := collateral * a.Value / 100
		return az / m, nil
	case POSITIONSIZE:
		pz, err := f.OpenPositions()
		return pz[a.Ticker].PositionSize, err
	}
	return 0, nil
}
