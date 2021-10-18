package parser

import (
	"errors"
	"fmt"
	"io"

	"github.com/DawnKosmos/ftxwebapp/exchange"
	"github.com/DawnKosmos/ftxwebapp/lexer"
)

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

//Parser is a Struct that is able to Evaluate itself and communicate the Result to a Communicator
type Parser interface {
	Evaluate(w Communicator, f exchange.Exchange) error
}

//Communicator is an Interface that Implements the Reader and Writer Interface for communication. It also handles the functions and values we assign a variable
type Communicator interface {
	io.Writer                             //Parser Writes the Evaluated Results
	io.Reader                             //Lexer Reads input from the communicator that get send to the Parser
	AddVariable(string, Variable)         //Add an Variable
	GetVariable(string) (Variable, error) //Return a Variable
	ErrorMessage(error)                   //To handle errors
}

//Parse returns a Parser which then gets Evaluated and returns
func Parse(tk []lexer.Token, c Communicator) (Parser, error) {
	nk := tk

	if len(tk) == 0 {
		return nil, nerr(empty, "Error nothing got lexed")
	}

	if tk[0].Type == lexer.VARIABLE {
		v, err := c.GetVariable(tk[0].Content)
		if err != nil { // If Varriable is unknown it will check if there is an assign
			if len(tk) == 1 { //
				return nil, nerr(err, fmt.Sprintf("ERROR %s is an unknown Variable", tk[0].Content))
			}
			if tk[1].Type == lexer.ASSIGN {
				if err := ParseAssign(tk[0].Content, tk[1:], c); err == nil { //Assigns the Variable. Returns a (nil, nil) if succesfull
					return nil, nil
				}
			} else {
				return nil, nerr(empty, fmt.Sprintf("ERROR %s is an unknown Variable", tk[0].Content))
			}
		}
		nk, err = ParseVariable(v, tk[1:]) //We parse the Variable, which then returns a new []Token, which gets parsed
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
	case lexer.SIDE: // buy, sell
		o, err = ParseOrder(nk[0].Content, nk[1:])
	case lexer.STOP: //stop
		o, err = ParseStop(nk[1:])
	case lexer.CANCEL: //cancel
		o, err = ParseCancel(nk[1:])
	case lexer.FUNDINGPAYS: //fpays
		o, err = ParseFundingPays(nk[1:])
	case lexer.FUNDINGRATES: //frates
		o, err = ParseFundingRates(nk[1:])
	case lexer.CLOSE: //fclose
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
