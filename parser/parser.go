package parser

import (
	"errors"
	"fmt"
	"io"

	"github.com/DawnKosmos/ftxwebapp/exchange"
	"github.com/DawnKosmos/ftxwebapp/lexer"
)

type Parser interface {
	Evaluate(w Communicator, f exchange.Exchange) error
}

/*
order : [SIDE][market][SIZE][SOURCE][PRICE]
stop : stop [SIDE][market][SIZE][SOURCE][PRICE] [-ro(reduce only)]
fundingpays : fpays [market] [duration] [-sum] [-position]
fundingrates : frates [market] [duration] [-sum] [-position] [-highest] [-lowest]
function : func([parameters]) [*]
assign : [variableName] = [function, variable, expression]
cancel : cancel [market] [SIDE] [STOP] [Position] [all]
account: return account information
change [name]: changes account
fastmode : fast [ticker]


side : buy | sell
market : btc-perp | btc-usd
size : 50% (50% free collateral) | u50(50 units of the coin) | 50(50 $)
source(optional, default is market) : | -low | -high] [duration]
price : [38000 | -400(below/above 400 points from source) | 4%(below/above 4% of source)]
laddered order: -l | -le  [5] price price
duration : [10 | 7h | 10d | 21d]

*/

type Communicator interface {
	io.Writer
	io.Reader
	AddVariable(string, Variable)
	GetVariable(string) (Variable, error)
}

type parseError struct {
	err error
	msg string
}

var empty = errors.New("")

func nerr(err error, msg string) *parseError {
	return &parseError{err, msg}
}

func (e *parseError) Error() string {
	return fmt.Sprintf("Message:%s + %v", e.msg, e.err)
}

func Parse(tk []lexer.Token, c Communicator) (Parser, error) {
	nk := tk

	if len(tk) == 0 {
		return nil, nerr(empty, "Error nothing got lexed")
	}

	if tk[0].Type == lexer.VARIABLE {
		v, err := c.GetVariable(tk[0].Content)
		if err != nil {
			if len(tk) == 1 {
				return nil, nerr(empty, fmt.Sprintf("ERROR %s is an unknown Variable", tk[0].Content))
			}
			if tk[1].Type == lexer.ASSIGN {
				if err := ParseAssign(tk[0].Content, tk[1:], c); err == nil {
					return nil, err
				}
			} else {
				return nil, nerr(empty, fmt.Sprintf("ERROR %s is an unknown Variable", tk[0].Content))
			}
		}
		nk, err = ParseVariable(v, tk[1:])
		if err != nil {
			return nil, err
		}
	}

	/*nk, err := parseVariables(nk, c)
	if err != nil {
		return nil, err
	}*/
	var o Parser
	var err error

	switch nk[0].Type {
	case lexer.SIDE:
		o, err = ParseOrder(nk[0].Content, nk[1:])
	case lexer.STOP:
	//	o, err = ParseStop(nk[0].Content, nk[1:])
	case lexer.CANCEL:
		o, err = ParseCancel(nk[1:])
	case lexer.FUNDINGPAYS:
		o, err = ParseFundingPays(nk[1:])
	case lexer.FUNDINGRATES:
		o, err = ParseFundingRates(nk[1:])
	case lexer.CLOSE:
	//	o, err = ParseClose(nk[1:])
	default:
		return o, nerr(empty, fmt.Sprintf("Invalid Type Error during Parsing %v", nk[0].Type))
	}

	if err != nil {
		return o, err
	}

	return o, nil
}

func parseVariables(tk []lexer.Token, w Communicator) ([]lexer.Token, error) {
	var nk []lexer.Token
	if len(tk) > 30 {
		return tk, nerr(empty, "Error, Loop detected")
	}
	for i, v := range tk {
		switch v.Type {
		case lexer.VARIABLE:
			val, err := w.GetVariable(v.Content)
			if err != nil {
				nk = append(nk, v)
			} else {
				temp, err := ParseVariable(val, tk[i:])
				if err != nil {
					return nk, err
				}
				temp2, err := parseVariables(temp, w)
				if err != nil {
					return nk, err
				}
				return append(nk, temp2...), nil

			}
		default:
			nk = append(nk, v)
		}
	}
	return nk, nil
}
