package parser

import (
	"fmt"

	"github.com/DawnKosmos/ftxwebapp/exchange"
	"github.com/DawnKosmos/ftxwebapp/lexer"
)

type Cancel struct {
	Side         exchange.CancelType
	Ticker       []string
	triggerOrder bool
}

func ParseCancel(tk []lexer.Token) (p Parser, err error) {

	var cancle Cancel
	cancle.Side = exchange.ALL
	cancle.Ticker = make([]string, 0)

	if len(tk) == 0 {
		return &cancle, nil
	}

	for _, v := range tk {
		switch v.Type {
		case lexer.SIDE:
			if v.Content == "buy" {
				cancle.Side = exchange.BUY
			} else {
				cancle.Side = exchange.SELL
			}
		case lexer.FLAG:
			switch v.Content {
			case "stop":
				cancle.triggerOrder = true
			default:
				return nil, nerr(empty, fmt.Sprintf("Error Parsing Cancel, Invalid flag %s", v.Content))
			}
		case lexer.VARIABLE:
			cancle.Ticker = append(cancle.Ticker, v.Content)
		default:
			return nil, nerr(empty, fmt.Sprintf("Error Parsing Cancel, Invalid Type %d %s", v.Type, v.Content))
		}
	}
	return &cancle, nil
}

func (c *Cancel) Evaluate(w Communicator, f exchange.Exchange) error {
	if len(c.Ticker) == 0 {
		return f.Cancel(c.Side, "")
	}
	for _, v := range c.Ticker {
		err := f.Cancel(c.Side, v)
		if err != nil {
			return err
		}
	}
	return nil
}
