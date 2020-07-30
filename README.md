# errorx
<a href="https://github.com/goapt/errorx/actions"><img src="https://github.com/goapt/errorx/workflows/test/badge.svg" alt="Build Status"></a>
<a href="https://codecov.io/gh/goapt/errorx"><img src="https://codecov.io/gh/goapt/errorx/branch/master/graph/badge.svg" alt="codecov"></a>
<a href="https://goreportcard.com/report/github.com/goapt/errorx"><img src="https://goreportcard.com/badge/github.com/goapt/errorx" alt="Go Report Card
"></a>
<a href="https://pkg.go.dev/github.com/goapt/errorx"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square" alt="GoDoc"></a>
<a href="https://opensource.org/licenses/mit-license.php" rel="nofollow"><img src="https://badges.frapsoft.com/os/mit/mit.svg?v=103"></a>

golang combine error package

```shell script
go get github.com/goapt/errorx  
```

## 使用
errorx包定义了全局的基本错误如下

* ErrSystem 系统错误，最顶层的错误定义，所有全局定义的错误必须包裹这个错误
* ErrDatabase 数据库错误
* ErrRedis Redis错误
* ErrNetwork 网络错误

其中分别提供了四种错误包裹的快捷方法

* errorx.System(err) 系统错误包裹
* errorx.Database(err) 数据库错误包裹
* errorx.Redis(err) Redis错误包裹
* errorx.Network(err) 网络错误包裹

四种错误包裹都提供额外的msg参数来重置错误信息如下

```go
errorx.Database(err, "订单不存在")
```

同样我们保留了`DbPrettyNoMoreRows`和`DbFilterNoMoreRows`两个方法，来辅助快捷的重置无数据错误信息，以及过滤无数据错误

```go
errorx.DbPrettyNoMoreRows(err, "订单不存在")
errorx.DbFilterNoMoreRows(err)
```

## New一个普通错误，带调用栈
```go
err := errorx.New("这是一个错误")
```

## 包裹一个普通错误，希望他拥有调用栈
```go
err := errors.New("普通的错误")
errorx.Wrap(err)
```

## 判定错误
有了包裹错误的方式，我们就可以使用官方的errors.Is来判定错误是否是系统错误，或者数据库错误以及官方包的错误等

```go
errors.Is(err,errorx.System)
errors.Is(err,errorx.Database)
errors.Is(err,sql.ErrNoRows)
```

## 打印调用堆栈
有时候我们不仅仅只想获取错误信息，还希望知道错误的调用堆栈，我们可以使用`%+v`来打印如下
```go
fmt.Sprintf("%+v",err)
```
输出如下
```
sql: no rows in result set
github.com/goapt/errorx.TestCombErr_Error.func1
	/Users/fifsky/wwwroot/go/github.com/goapt/errorx/error_test.go:57
github.com/goapt/errorx.TestCombErr_Error.func1
	/Users/fifsky/wwwroot/go/github.com/goapt/errorx/error_test.go:58
github.com/goapt/errorx.TestCombErr_Error.func1
	/Users/fifsky/wwwroot/go/github.com/goapt/errorx/error_test.go:59
github.com/goapt/errorx.TestCombErr_Error.func1
	/Users/fifsky/wwwroot/go/github.com/goapt/errorx/error_test.go:60
```

## 业务层错误
由于我们业务层统一规范结构如下
```json
{
    "code":"InvalidEmail",
    "msg": "无效的邮箱"
}
```

因此我们提供了辅助方法来定义这种结构的错误如下

```go
var ErrInvalidEmail = errorx.NewCode("InvalidEmail","无效的邮箱")
```

第二个参数，支持string,error，因此你可以轻松的包裹一个错误如下

```go
var err := errorx.NewCode("ConnectFail",errorx.Database(sql.ErrNoRows))
//下面都是成立的
errors.Is(err,errorx.ErrSystem)
errors.Is(err,errorx.ErrDatabase)
errors.Is(err,sql.ErrNoRows)
```
然后在业务层的response就可以轻松的断言出错误是否为 err.(errorx.CodeError) 