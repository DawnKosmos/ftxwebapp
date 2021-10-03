package parser

import (
	"fmt"

	"github.com/DawnKosmos/ftxwebapp/lexer"
)

type function struct {
	NumbersOfParameters int
	ParameterPosition   []int
	Token               []lexer.Token
}

func (f *function) Parse(tk []lexer.Token) ([]lexer.Token, error) {
	if len(tk) < 3 {
		return tk, nerr(empty, "Function Syntax Error, Brackets Missing")
	}

	if tk[0].Type != lexer.LBRACKET {
		return tk, nerr(empty, "Function Syntax Error, no bracket")
	}
	var nk []lexer.Token
	for _, v := range tk[1:] {
		switch v.Type {
		case lexer.RBRACKET:
			break
		default:
			nk = append(nk, v)
		}
	}

	if len(nk) != f.NumbersOfParameters {
		return tk, nerr(empty, fmt.Sprintf("Function Error, wrong number of parameters. Want %d got %d", f.NumbersOfParameters, len(nk)))
	}

	for i, t := range nk {
		n := f.ParameterPosition[i]
		f.Token[n] = t
	}

	return f.Token, nil
}
