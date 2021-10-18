package exchange

import "time"

const LONG = true
const SHORT = false

type CancelType int

const (
	SELL CancelType = iota
	ALL
	BUY
)

//Exchange Interface is needed to provide data and execute commands to provide all futures
type Exchange interface {
	//SetOrder set an Order, reduceOnly optional, default is false
	SetOrder(side bool, ticker string, price float64, size float64, reduceOnly bool) (Order, error)

	OpenOrders(side bool, ticker string) ([]Order, error)
	//SetTriggerOrder set an TriggerOrder, reduceOnly optional, default is true
	SetTriggerOrder(side bool, ticker string, price float64, size float64, orderType string, reduceOnly bool) (TriggerOrder, error)
	//MarketPrice return the Market Price of the asked Ticker
	MarketPrice(ticker string) (float64, error)
	//Cancel All=0, Buy=1 Sell=-1 orders on given ticker. No ticker means all orders get cancelled. Return is the amount of orders that got cancelled
	Cancel(Side CancelType, Ticker string) error
	//CancelTrigger All=0, Buy=1 Sell=-1 orders on given ticker. No ticker means all orders get cancelled. Return is the amount of orders that got cancelled
	CancelTrigger(Side CancelType, Ticker string) error
	//Highest return the Highest Price of the ticker for the given duration
	Highest(ticker string, duration int64) (float64, error)
	//Lowest returns the Lowest Price of the ticker for the given duration
	Lowest(ticker string, duration int64) (float64, error)
	//FundingPayments returns the fundingpayments paid in the given period
	FundingPayments(ticker string, starttime int64, endtime int64) ([]FundingPayments, error)
	//FundingRates returns the Fundingrate of the ticker for the given period
	FundingRates(ticker string, starttime int64, endtime int64) ([]FundingRates, error)
	//FreeCollateral() return the free collateral in USD
	FreeCollateral() (float64, error)
	//OpenPositions returns all Open positions
	OpenPositions() (map[string]Position, error)
	//GetMarkets returns the ticker of all tradetable markets
	//GetMarkets() ([]string, error)
}

type Order struct {
	Ticker     string
	Side       string
	Size       float64
	Price      float64
	ReduceOnly bool
	Created    time.Time
}

type TriggerOrder struct {
	Ticker     string
	Side       bool
	Size       float64
	Price      float64
	ReduceOnly bool
	Created    time.Time
}

type Position struct {
	Side         string
	Future       string
	NotionalSize float64
	PositionSize float64
	UPNL         float64
	EntryPrice   float64
}

type FundingPayments struct {
	Ticker  string
	Payment float64
	Time    time.Time
}

type FundingRates struct {
	Ticker string
	Rate   float64
	Time   time.Time
}
