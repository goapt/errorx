package errorx

import (
	"fmt"
	"io"
)

type CodeError struct {
	err error
	*stack
	Code string
	Msg  string
}

func (e CodeError) Error() string {
	return e.Msg
}

func (e *CodeError) Is(err error) bool {
	if er, ok := err.(*CodeError); ok {
		return er.Code == e.Code
	}
	return false
}

func (e *CodeError) Unwrap() error {
	return e.err
}

func (e *CodeError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", e.Unwrap())
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}

func NewCode(code string, msg interface{}) *CodeError {
	ex := &CodeError{
		Code:  code,
		stack: callers(3),
	}

	m := ""
	switch e := msg.(type) {
	case error:
		m = e.Error()
		ex.err = e
	case string:
		m = e
	default:
		m = fmt.Sprintf("[%s]%v", code, e)
	}

	if m == "" {
		m = fmt.Sprintf("[%s]unknow error", code)
	}

	ex.Msg = m

	return ex
}
