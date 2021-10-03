package parser

import (
	"fmt"
	"strconv"

	"github.com/DawnKosmos/ftxwebapp/exchange"
	"github.com/DawnKosmos/ftxwebapp/lexer"
)

type Order struct {
	Side   string
	Ticker string
	A      Amount
	P      Price
}

func ParseOrder(Side string, tk []lexer.Token) (o *Order, err error) {
	o = &Order{}
	o.Side = Side
	var a Amount

	if len(tk) < 3 {
		return nil, nerr(empty, "Error Parse Order, False Input")
	}

	if tk[0].Type == lexer.VARIABLE {
		o.Ticker = tk[0].Content
	} else {
		return nil, nerr(empty, fmt.Sprintf("Error no Ticker is %s", tk[0].Content))
	}

	switch tk[1].Type {
	case lexer.FLOAT:
		a.Type = COIN
	case lexer.UFLOAT:
		a.Type = FIAT
	case lexer.PERCENT:
		a.Type = ACCOUNTSIZE
	case lexer.POSITION:
		a.Type = POSITIONSIZE
	default:
		return nil, nerr(empty, fmt.Sprintf("Error Parse Order, false Order Size of type"))
	}
	a.Ticker = o.Ticker
	a.Value, err = strconv.ParseFloat(tk[1].Content, 64)
	if err != nil {
		return nil, nerr(err, fmt.Sprintf("Parse Error Wrong Value should be a Float is %s", tk[1].Content))
	}

	o.A = a

	o.P, err = ParsePrice(tk[2:])
	return o, err

}

func (o *Order) Evaluate(w Communicator, f exchange.Exchange) error {
	size, err := o.A.Evaluate(f)
	if err != nil {
		return err
	}

	return o.P.Evaluate(f, w, o.Side, o.Ticker, size)
}
