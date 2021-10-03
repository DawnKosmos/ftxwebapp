package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

type TokenType int

const (
	VARIABLE     TokenType = iota //x
	TICKER                        // btc-perp
	SIDE                          // buy, sell
	STOP                          // stop
	FLOAT                         // 100 => 100 $ of btc
	UFLOAT                        // u100  u = unitFloat => buying 100 btc
	PERCENT                       //
	DFLOAT                        // -200 d = differenceFloat =>
	ASSIGN                        // =
	FLAG                          // -l -le
	FUNC                          // func(a,b,c) creating function
	DURATION                      // 4h 1d 30
	LBRACKET                      // (
	RBRACKET                      // )
	MARKET                        // -market
	SOURCE                        // -high -low -open -close
	CANCEL                        // cancel
	CLOSE                         // close
	FUNDINGPAYS                   // fpay | fundingpayments
	POSITION                      // -position
	FUNDINGRATES                  //fundingrates
	ACCOUNT
)

type Token struct {
	Type    TokenType
	Content string
}

type lexerError struct {
	input      string
	err        error
	errmessage string
}

func nerr(input string, err error, errmessage string) *lexerError {
	return &lexerError{input, err, errmessage}
}

func (l *lexerError) Error() string {
	return fmt.Sprintf("Text: %s, Error: %v + %s", l.input, l.err, l.errmessage)
}

func Lexer(input string) (t []Token, err error) {
	in := strings.Split(input, " ")

	for _, s := range in {
		if len(s) == 0 {
			continue
		}
		last := len(s) - 1
		switch s {
		case "buy", "sell":
			t = append(t, Token{SIDE, s})
		case "stop":
			t = append(t, Token{STOP, s})
		case "=":
			t = append(t, Token{ASSIGN, "="})
		case "cancel":
			t = append(t, Token{CANCEL, "cancel"})
		case "fpays", "fundingpays":
			t = append(t, Token{FUNDINGPAYS, "fundingpays"})
		case "frates", "fundingrates":
			t = append(t, Token{FUNDINGRATES, "fundingrates"})
		case "account", "acc":
			t = append(t, Token{ACCOUNT, "account"})
		case "close":
			t = append(t, Token{CLOSE, s})
		default:
			if (s[last] == 'h' || s[last] == 'm' || s[last] == 'd') && len(s) > 1 {
				_, err := strconv.Atoi(s[:last])
				if err == nil {
					t = append(t, Token{DURATION, s})
					continue
				}
			}

			if len(s) > 6 {
				if s[:4] == "func" {
					t = append(t, Token{FUNC, "func"})
					t = append(t, lexFunc([]byte(s[4:]))...)
					continue
				}
			}

			if s[0] == '-' {
				_, err := strconv.ParseFloat(s[1:], 64)

				if err == nil {
					t = append(t, Token{DFLOAT, s[1:]})
				} else {
					ss := s[1:]
					switch ss {
					case "low", "high", "open", "close":
						t = append(t, Token{SOURCE, ss})
					case "position":
						t = append(t, Token{POSITION, "1.0"})
					case "market":
						t = append(t, Token{MARKET, "1.0"})
					default:
						t = append(t, Token{FLAG, ss})
					}
				}
				continue
			}

			if s[0] == 'u' && len(s) > 1 {
				_, err := strconv.Atoi(s[1:])
				if err == nil {
					t = append(t, Token{UFLOAT, s[1:]})
				} else {
					t = append(t, Token{VARIABLE, s})
				}
				continue
			}

			if s[last] == '%' {
				_, err := strconv.ParseFloat(s[:last], 64)
				if err != nil {
					return t, nerr(s, err, "A variable can't end with %")
				}
				t = append(t, Token{PERCENT, s[:len(s)-1]})
				continue
			}

			_, err := strconv.ParseFloat(s, 64)
			if err == nil {
				t = append(t, Token{FLOAT, s})
				continue
			}
			fmt.Println(string(s))
			t = append(t, lexVariable([]byte(s))...)
		}
	}

	return
}

func lexFunc(s []byte) []Token {
	var temp []byte
	var tk []Token
	for _, v := range s {
		switch v {
		case '(':
			tk = append(tk, Token{LBRACKET, ""})
		case ')':
			tk = append(tk, Token{VARIABLE, string(temp)}, Token{RBRACKET, ""})
			temp = []byte("")
		case ',':
			tk = append(tk, Token{VARIABLE, string(temp)})
			temp = []byte("")
		default:
			temp = append(temp, v)
		}
	}

	return tk
}

func lexVariable(s []byte) []Token {
	var temp []byte
	var tk []Token

	for _, v := range s {
		switch v {
		case '(':
			tk = append(tk, Token{VARIABLE, string(temp)}, Token{LBRACKET, ""})
			temp = []byte("")
		case ')':
			l, _ := Lexer(string(temp))
			tk = append(tk, l...)
			tk = append(tk, Token{RBRACKET, ""})
			temp = []byte("")
		case ',':
			temp = append(temp, ' ')
		default:
			temp = append(temp, v)
		}
	}

	if len(temp) != 0 {
		tk = append(tk, Token{VARIABLE, string(temp)})
	}
	return tk
}
