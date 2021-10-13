package starter

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/DawnKosmos/ftxwebapp/exchange"
	"github.com/DawnKosmos/ftxwebapp/lexer"
	"github.com/DawnKosmos/ftxwebapp/parser"
)

/*
type Communicator interface {
	io.Writer
	io.Reader
	AddVariable(string, Variable)
	GetVariable(string) (*Variable, error)
}
*/

type cle struct {
	reader *bufio.Reader
	vm     map[string]parser.Variable
}

func (e *cle) Write(p []byte) (int, error) {
	fmt.Println(">", string(p))
	return 0, nil
}

func (e *cle) Read(p []byte) (int, error) {
	pp, err := e.reader.ReadSlice('\n')
	copy(p, pp[:len(pp)-2])
	return len(pp) - 2, err
}

func (w *cle) AddVariable(s string, v parser.Variable) {
	w.vm[s] = v
}

func (w *cle) GetVariable(s string) (parser.Variable, error) {
	v, ok := w.vm[s]
	if !ok {
		return parser.Variable{}, errors.New("Error getting Variable")
	}
	return v, nil
}

func (w *cle) ErrorMessage(err error) {
	fmt.Println(err)
}

func CommandLine(f exchange.Exchange) parser.Communicator {
	cl := &cle{}
	cl.vm = make(map[string]parser.Variable)
	cl.reader = bufio.NewReader(os.Stdin)
	return cl
}

func Run(w parser.Communicator, f exchange.Exchange) error {
	for {
		var b []byte = make([]byte, 128)
		l, _ := w.Read(b)
		b = b[:l]
		if strings.Compare(string(b), "exit") == 0 {
			w.Write([]byte("byebye"))
			break
		}
		t, err := lexer.Lexer(string(b))
		if err != nil {
			w.ErrorMessage(err)
			continue
		}
		p, err := parser.Parse(t, w)
		if p == nil {
			if err != nil {
				w.ErrorMessage(err)
				continue
			} else {
				continue
			}
		}

		if err != nil {
			w.ErrorMessage(err)
			continue
		}

		err = p.Evaluate(w, f)
		if err != nil {
			w.ErrorMessage(err)
			continue
		}
	}
	return nil
}

//Execute is a lightweight version with limited functions
func Execute(w io.Writer, r io.Reader, f exchange.Exchange) error {
	ex := newExecute(w, r)

	buff := make([]byte, 1024)
	l, _ := ex.Read(buff)
	buff = buff[:l]

	commands := strings.Split(string(buff), "\n")
	for _, v := range commands {
		tk, err := lexer.Lexer(v)
		if err != nil {
			return err
		}

		p, err := parser.Parse(tk, ex)
		if err != nil {
			return err
		}

		err = p.Evaluate(ex, f)
		if err != nil {
			return err
		}
	}
	return nil
}

type execute struct {
	w io.Writer
	r io.Reader
}

func newExecute(w io.Writer, r io.Reader) *execute {
	return &execute{w, r}
}

func (e *execute) Write(p []byte) (int, error) {
	return e.Write(p)
}

func (e *execute) Read(p []byte) (int, error) {
	return e.Read(p)
}

func (e *execute) AddVariable(string, parser.Variable) {
	return
}
func (e *execute) GetVariable(string) (parser.Variable, error) {
	return parser.Variable{}, errors.New("You are not allowed to assign Variables here")
}

func (e *execute) ErrorMessage(error) {
	return
}
