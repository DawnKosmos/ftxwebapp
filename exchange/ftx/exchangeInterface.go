package ftx

import (
	"net/http"

	"github.com/DawnKosmos/ftxcmd/ftx"
	"github.com/DawnKosmos/ftxwebapp/exchange"
)

type FTX struct {
	f *ftx.Client
}

func NewClient(cl *http.Client, key, secret, sub string) *FTX {
	f := ftx.NewClient(cl, key, secret, sub)
	return &FTX{f}
}

func (f *FTX) SetOrder(side bool, ticker string, price float64, size float64, reduceOnly bool) (exchange.Order, error) {
	ss := "sell"
	if side {
		ss = "buy"
	}
	e, err := f.f.SetOrder(ticker, ss, price, size, "limit", reduceOnly)
	eo := exchange.Order{
		Ticker:     e.Result.Ticker,
		Side:       ss,
		Size:       e.Result.Size,
		Price:      e.Result.Price,
		ReduceOnly: e.Result.ReduceOnly,
		Created:    e.Result.CreatedAt,
	}
	return eo, err
}

func (f *FTX) SetTriggerOrder(side bool, ticker string, price float64, size float64, orderType string, reduceOnly bool) (exchange.TriggerOrder, error) {
	ss := "sell"
	if side {
		ss = "buy"
	}
	e, err := f.f.SetTriggerOrder(ticker, ss, price, size, orderType, reduceOnly)
	eo := exchange.TriggerOrder{
		Ticker:     ticker,
		Side:       side,
		Size:       e.Result.Size,
		Price:      e.Result.TriggerPrice,
		ReduceOnly: reduceOnly,
		Created:    e.Result.CreatedAt,
	}
	return eo, err
}

func (f *FTX) MarketPrice(ticker string) (float64, error) {
	e, err := f.f.GetMarket(ticker)
	eo := (e.Ask + e.Bid + e.Last) / 3
	return eo, err
}

func (f *FTX) Cancel(Side exchange.CancelType, Ticker string) error {
	switch Side {
	case exchange.SELL:
		return f.f.CancelOrders(Ticker, "sell", false)
	case exchange.ALL:
		return f.f.CancelOrders(Ticker, "", false)
	case exchange.BUY:
		return f.f.CancelOrders(Ticker, "buy", false)
	}
	return nil
}

func (f *FTX) CancelTrigger(Side exchange.CancelType, Ticker string) error {
	switch Side {
	case exchange.SELL:
		return f.f.CancelOrders(Ticker, "sell", true)
	case exchange.ALL:
		return f.f.CancelOrders(Ticker, "", true)
	case exchange.BUY:
		return f.f.CancelOrders(Ticker, "buy", true)
	}
	return nil
}

func (f *FTX) Highest(ticker string, duration int64) (float64, error) {
	return f.f.GetPriceSource("high", ticker, duration)
}

func (f *FTX) Lowest(ticker string, duration int64) (float64, error) {
	return f.f.GetPriceSource("low", ticker, duration)
}

func (f *FTX) FundingPayments(ticker string, starttime int64, endtime int64) ([]exchange.FundingPayments, error) {
	e, err := f.f.GetFundingPayments(ticker, starttime, endtime)
	var eo []exchange.FundingPayments
	for _, v := range e {
		temp := exchange.FundingPayments{
			Ticker:  v.Future,
			Payment: v.Payment,
			Time:    v.Time,
		}
		eo = append(eo, temp)
	}
	return eo, err
}

func (f *FTX) FundingRates(ticker string, starttime int64, endtime int64) ([]exchange.FundingRates, error) {
	e, err := f.f.GetFundingRates(ticker, starttime, endtime)
	var eo []exchange.FundingRates
	for _, v := range e {
		temp := exchange.FundingRates{
			Ticker: ticker,
			Rate:   v.Rate,
			Time:   v.Time,
		}
		eo = append(eo, temp)
	}
	return eo, err
}

func (f *FTX) FreeCollateral() (float64, error) {
	a, err := f.f.GetAccount()

	return a.FreeCollateral, err
}

func (f *FTX) OpenPositions() (map[string]exchange.Position, error) {
	e, err := f.f.GetPosition()

	var eo map[string]exchange.Position = make(map[string]exchange.Position)
	for _, v := range e {
		temp := exchange.Position{
			Side:         v.Side,
			Future:       v.Future,
			NotionalSize: v.NotionalSize,
			PositionSize: v.PositionSize,
			UPNL:         v.UPNL,
			EntryPrice:   v.EntryPrice,
		}
		eo[v.Future] = temp
	}
	return eo, err
}

func (f *FTX) OpenOrders(side bool, ticker string) ([]exchange.Order, error) {
	ss := "sell"
	if side {
		ss = "buy"
	}
	o, err := f.f.GetOpenOrders(ticker)
	var out []exchange.Order
	if err != nil {
		return out, err
	}
	for _, v := range o {
		if v.Side == ss {
			temp := exchange.Order{
				Ticker:     v.Ticker,
				Side:       v.Side,
				Size:       v.Size,
				Price:      v.Price,
				ReduceOnly: v.ReduceOnly,
				Created:    v.CreatedAt,
			}
			out = append(out, temp)
		}
	}
	return out, nil
}

/*
type Exchange interface {

	//Cancel All=0, Buy=1 Sell=-1 orders on given ticker. No ticker means all orders get cancelled. Return is the amount of orders that got cancelled
	Cancel(Side int, Ticker ...string) (int, error)
	//CancelTrigger All=0, Buy=1 Sell=-1 orders on given ticker. No ticker means all orders get cancelled. Return is the amount of orders that got cancelled
	CancelTrigger(Side int, Ticker ...string) (int, error)
	//Returns the Position of the asked ticker
	Position(ticker string) (Position, error)
	//GetMarkets returns the ticker of all tradetable markets
	GetMarkets() ([]string, error)
}*/
