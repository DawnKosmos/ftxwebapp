package parser

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/DawnKosmos/ftxwebapp/exchange"
	"github.com/DawnKosmos/ftxwebapp/lexer"
)

type FundingPays struct {
	Ticker    []string
	Time      int64
	Position  bool
	Summarize bool
}

func ParseFundingPays(tl []lexer.Token) (*FundingPays, error) {
	var fund FundingPays
	var err error

	fund.Time = 36000
	for _, v := range tl {
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

func (e *FundingPays) Evaluate(w Communicator, f exchange.Exchange) (err error) {
	var ticker []string
	var fp []exchange.FundingPayments
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
		fp, err = f.FundingPayments("", tnow-e.Time, tnow)
		if err != nil {
			return nerr(err, "Error Evaluate Funding Pays:")
		}
	} else {
		for _, v := range ticker {
			temp, err := f.FundingPayments(v, tnow-e.Time, tnow)
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
			nm[v.Ticker] = v.Payment
		}

		for k, v := range nm {
			w.Write([]byte(fmt.Sprintf("%s: %.3f", k, v)))
		}
		return nil
	}

	for _, v := range fp {
		w.Write([]byte(fmt.Sprintf("%s %s: %.3f", v.Time.Format("Jan _2 15"), v.Ticker, v.Payment)))
	}

	return nil
}

func parseDuration(ss string) (int64, error) {
	n, err := strconv.Atoi(ss[:len(ss)-1])
	if err != nil {
		return 0, err
	}
	switch ss[len(ss)-1] {
	case 'h':
		n *= 3600
	case 'm':
		n *= 60
	case 'd':
		n *= 3600 * 24
	default:
		return 0, errors.New(ss + " I dont know how you fucked that up")
	}

	return int64(n), nil
}
