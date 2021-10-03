package parser

import (
	"fmt"

	"github.com/DawnKosmos/ftxwebapp/lexer"
)

func ParseVariable(v Variable, tk []lexer.Token) ([]lexer.Token, error) {
	switch v.Type {
	case FUNCTION:
		nk, err := ParseFunc(v.Content, tk)
		if err != nil {
			return tk, err
		}
		return nk, nil
	case CONSTANT:
		nk, ok := v.Content.([]lexer.Token)
		if !ok {
			return tk, nerr(empty, "Error Parse Variable, Variable not existing")
		}
		return append(nk, tk[1:]...), nil
	}

	return tk, nerr(empty, "Error while Parsing a Variable Something went wrong")
}

func ParseFunc(v interface{}, tk []lexer.Token) ([]lexer.Token, error) {
	fun, ok := v.(function)
	if !ok {
		return tk, nerr(empty, fmt.Sprintf("Unexpected ERROR with %v ", fun))
	}
	e, err := fun.Parse(tk)
	return e, err

}
