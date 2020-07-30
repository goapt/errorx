package errorx

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
)

var (
	ErrGeneral  = errors.New("错误")
	ErrSystem   = errors.New("系统异常")
	ErrDatabase = fmt.Errorf("数据库异常:%w", ErrSystem)
	ErrRedis    = fmt.Errorf("redis异常:%w", ErrSystem)
	ErrNetwork  = fmt.Errorf("网络异常:%w", ErrSystem)
)

func New(msg string) error {
	return Wrap(errors.New(msg))
}

// 包裹普通错误，让其拥有调用栈信息
func Wrap(err error, msg ...string) error {
	return combine(ErrGeneral, err, msg...)
}

// 组合系统错误，使得这个错误可以使用errors.Is判别出这两个错误
func System(err error, msg ...string) error {
	return combine(ErrSystem, err, msg...)
}

// 组合数据库错误
func Database(err error, msg ...string) error {
	return combine(ErrDatabase, err, msg...)
}

// DbPrettyNoMoreRows 优化查无记录的错误
// 调用方无需额外判断是否为空记录，只需要指定遇到空记录时，替换的错误提示
func DbPrettyNoMoreRows(err error, msg string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return Database(err, msg)
	}
	return Database(err)
}

// DbFilterNoMoreRows 过滤sql.ErrNoRows错误
func DbFilterNoMoreRows(err error) error {
	if err == nil || errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return Database(err)
}

// 组合一个Redis错误
func Redis(err error, msg ...string) error {
	return combine(ErrRedis, err, msg...)
}

// 组合一个HTTP错误
func Network(err error, msg ...string) error {
	return combine(ErrNetwork, err, msg...)
}

// 组合一个错误，使得errors.Is可以判断其中的任意一个错误
// 如果传入msg则会覆盖第二个参数err的错误信息
func combine(perr, err error, msg ...string) error {
	m := ""
	if len(msg) > 0 {
		m = msg[0]
	}

	return &combErr{
		perr:  perr,
		err:   err,
		msg:   m,
		stack: callers(4),
	}
}

type combErr struct {
	perr error
	err  error
	msg  string
	*stack
}

func (e *combErr) Error() string {
	msg := e.msg
	if msg == "" {
		var err *combErr
		if errors.As(e.err, &err) {
			msg = e.err.Error()
		} else {
			msg = fmt.Sprintf("%s:%s", e.err.Error(), e.perr.Error())
		}
	}

	return msg
}

func (e *combErr) Unwrap() error {
	return e.err
}

func (e *combErr) Is(err error) bool {
	return errors.Is(e.perr, err) || errors.Is(e.err, err)
}

func (e *combErr) Format(s fmt.State, verb rune) {
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
