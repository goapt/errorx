package errorx

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	tests := []struct {
		name  string
		field error
		err   error
		want  bool
	}{
		{"system error", System(errors.New("system error")), ErrSystem, true},
		{"database error", Database(errors.New("new error")), ErrDatabase, true},
		{"system error", Database(errors.New("new error")), ErrSystem, true},
		{"database error2", Database(Database(errors.New("new error"))), ErrDatabase, true},
		{"sql no rows", Database(sql.ErrNoRows), sql.ErrNoRows, true},
		{"not match error", Database(errors.New("new error")), errors.New("not match error"), false},
		{"redis error", Redis(errors.New("redis error")), ErrRedis, true},
		{"system error", Redis(errors.New("redis error")), ErrSystem, true},
		{"network error", Network(errors.New("network error")), ErrNetwork, true},
		{"system error", Network(errors.New("network error")), ErrSystem, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Is(tt.field, tt.err); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}

	assert.EqualError(t, Database(sql.ErrNoRows), "sql: no rows in result set:数据库异常:系统异常")
	assert.EqualError(t, Database(sql.ErrNoRows, "订单不存在"), "订单不存在")
	assert.EqualError(t, DbPrettyNoMoreRows(sql.ErrNoRows, "订单不存在"), "订单不存在")
	assert.EqualError(t, DbPrettyNoMoreRows(errors.New("not sql no rows"), "订单不存在"), "not sql no rows:数据库异常:系统异常")
	assert.Nil(t, DbPrettyNoMoreRows(nil, "订单不存在"))
	assert.Nil(t, DbFilterNoMoreRows(sql.ErrNoRows))
	assert.EqualError(t, DbFilterNoMoreRows(errors.New("not sql no rows")), "not sql no rows:数据库异常:系统异常")

	assert.EqualError(t, Redis(errors.New("redis error")), "redis error:redis异常:系统异常")
	assert.EqualError(t, Redis(errors.New("redis error"), "缓存已过期"), "缓存已过期")

	assert.EqualError(t, Network(net.ErrWriteToConnected), "use of WriteTo with pre-connected connection:网络异常:系统异常")
	assert.EqualError(t, Network(net.ErrWriteToConnected, "网络连接错误"), "网络连接错误")
}

func TestCombErr_Error(t *testing.T) {
	t.Run("mutil wrap error", func(t *testing.T) {
		err := sql.ErrNoRows
		err1 := Database(err)
		err11 := Redis(err1)
		err2 := Database(err11)
		err22 := Redis(err2)

		st := fmt.Sprintf("%+v", err22)
		fmt.Println(st)
		assert.Contains(t, st, "TestCombErr_Error")
		assert.Contains(t, st, "error_test.go:57")
		assert.Contains(t, st, "error_test.go:58")
		assert.Contains(t, st, "error_test.go:59")
		assert.Contains(t, st, "error_test.go:60")
		assert.EqualError(t, err22, "sql: no rows in result set:数据库异常:系统异常")
		assert.True(t, errors.Is(err22, ErrRedis))
	})

	t.Run("repeat wrap error", func(t *testing.T) {
		err := sql.ErrNoRows
		err1 := Database(err)
		err2 := Database(err1)
		assert.EqualError(t, err2, "sql: no rows in result set:数据库异常:系统异常")
		assert.True(t, errors.Is(err2, err))
	})
}

func TestWrap(t *testing.T) {
	err := errors.New("自定义错误")
	err2 := Wrap(err)
	err3 := Wrap(err2)
	fmt.Printf("%+v\n", err3)
	assert.EqualError(t, err3, "自定义错误:错误")
}

func TestNew(t *testing.T) {
	err := New("自定义错误")
	assert.EqualError(t, err, "自定义错误:错误")
}
