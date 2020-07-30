package errorx

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCode(t *testing.T) {
	t.Run("wrap database", func(t *testing.T) {
		err := NewCode("InvalidDatabase", ErrDatabase)
		assert.True(t, errors.Is(err, ErrSystem))
		assert.True(t, errors.Is(err, ErrDatabase))
		assert.EqualError(t, err, "数据库异常:系统异常")
		err3 := NewCode("InvalidDatabase", err)

		fmt.Printf("%+v\n", err3)

		err2 := NewCode("InvalidDatabase", Database(sql.ErrNoRows))
		assert.True(t, errors.Is(err2, ErrSystem))
		assert.True(t, errors.Is(err2, ErrDatabase))
		assert.True(t, errors.Is(err2, sql.ErrNoRows))
		assert.EqualError(t, err2, "sql: no rows in result set:数据库异常:系统异常")
	})

	t.Run("coustom", func(t *testing.T) {
		err := NewCode("InvalidEmail", "无效的邮箱")
		err2 := NewCode(err.Code, "邮箱不能为空")
		assert.EqualError(t, err, "无效的邮箱")
		assert.True(t, errors.Is(err, err2))
	})

	t.Run("struct error", func(t *testing.T) {
		p := &struct {
			Id int
		}{
			Id: 1,
		}

		err := NewCode("InvalidEmail", p)
		assert.EqualError(t, err, "[InvalidEmail]&{1}")
	})

	t.Run("unkonw error", func(t *testing.T) {
		err := NewCode("InvalidEmail", "")
		assert.EqualError(t, err, "[InvalidEmail]unknow error")
	})
}
