package parser

import (
	"fmt"

	"github.com/DawnKosmos/ftxwebapp/lexer"
)

type VariableType int

const (
	FUNCTION VariableType = iota
	CONSTANT
)

type Variable struct {
	Type    VariableType
	Content interface{}
}

func ParseAssign(name string, tk []lexer.Token, w Communicator) (err error) {
	if len(tk) == 0 {
		return nerr(empty, "Nothing can't be assigned to a Variable")
	}

	switch tk[0].Type {
	case lexer.FUNC:
		r, err := ParseAssignFunc(tk[1:])
		if err != nil {
			return err
		}
		w.AddVariable(name, Variable{FUNCTION, r})
	default:
		w.AddVariable(name, Variable{CONSTANT, tk[1:]})
	}
	w.Write([]byte(fmt.Sprintf("Variable %s assigned succesfully", name)))
	return nil
}

func ParseAssignFunc(tk []lexer.Token) (f function, err error) {
	if len(tk) == 0 {
		return f, nerr(empty, "Empty Func can't be assigned to a variable")
	}
	if tk[0].Type != lexer.LBRACKET {
		return f, nerr(empty, fmt.Sprintf("%s INVALID Syntax, after a func there must be a '(' "))
	}

	m := make(map[string]int) //A map that track which variable is on which position of the tokenlist

	nl := tk[1:]
	var count int

L:
	for _, v := range nl {
		switch v.Type {
		case lexer.RBRACKET:
			break L
		case lexer.VARIABLE:
			m[v.Content] = count
		default:
			return f, nerr(empty, fmt.Sprintf("INVALID VARIABLE NAME %s: ", v.Content))
		}
		count++
	}

	f.ParameterPosition = make([]int, count)
	tk = tk[count+2:]

	f.NumbersOfParameters = count
	for i, v := range tk {
		if v.Type == lexer.VARIABLE {
			n, ok := m[v.Content]
			if ok {
				f.ParameterPosition[n] = i
				delete(m, v.Content)
			}
		}
		f.Token = append(f.Token, v)
	}

	if len(m) != 0 {
		return f, nerr(empty, fmt.Sprintf("Not all Variables got assigned: %+v", m))
	}

	return f, nil
}
