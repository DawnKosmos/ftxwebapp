package parser

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DawnKosmos/ftxwebapp/exchange"
	"github.com/DawnKosmos/ftxwebapp/lexer"
)

type FundingRates struct {
	Ticker    []string
	Time      int64
	Position  bool
	Summarize bool
}

func ParseFundingRates(tk []lexer.Token) (Parser, error) {
	var fund FundingRates
	var err error

	fund.Time = 36000
	for _, v := range tk {
		switch v.Type {
		case lexer.FLAG:
			if v.Content == "sum" {
				fund.Summarize = true
			}
		case lexer.POSITION:
			fund.Position = true
		case lexer.VARIABLE:
			fund.Ticker = append(fund.Ticker, v.Content)
		case lexer.DURATION:
			fund.Time, err = parseDuration(v.Content)
			if err != nil {
				return &fund, err
			}
		case lexer.FLOAT:
			ff, err := strconv.ParseFloat(v.Content, 64)
			if err != nil {
				return nil, err
			}
			fund.Time = int64(ff) * 3600
		default:
			return nil, nerr(empty, fmt.Sprintf("Error Parsing FundingPayments %d %s not supported", v.Type, v.Content))
		}
	}

	return &fund, nil
}

func (e *FundingRates) Evaluate(w Communicator, f exchange.Exchange) (err error) {
	var ticker []string
	var fp []exchange.FundingRates
	tnow := time.Now().Unix()
	if e.Position {
		op, err := f.OpenPositions()
		if err != nil {
			return nerr(err, "Error Evaluate Funding Pays:")
		}
		for k, _ := range op {
			ticker = append(ticker, k)
		}
	}

	for _, v := range e.Ticker {
		ticker = append(ticker, v)
	}

	nm := make(map[string]float64)

	if len(e.Ticker) == 0 {
		fp, err = f.FundingRates("", tnow-e.Time, tnow)
		if err != nil {
			return nerr(err, "Error Evaluate Funding Pays:")
		}
	} else {
		for _, v := range ticker {
			temp, err := f.FundingRates(v, tnow-e.Time, tnow)
			if err != nil {
				return err
			}
			fp = append(fp, temp...)
			if len(fp) > 0 {
				nm[fp[0].Ticker] = 0
			}
		}
	}

	if e.Summarize {
		for _, v := range fp {
			nm[v.Ticker] = v.Rate
		}

		for k, v := range nm {
			w.Write([]byte(fmt.Sprintf("%s: %.3f", k, v*100)))
		}
		return nil
	}

	for _, v := range fp {
		w.Write([]byte(fmt.Sprintf("%s %s: %.3f", v.Time.Format("Jan _2 15"), v.Ticker, v.Rate)))
	}

	return nil
}
