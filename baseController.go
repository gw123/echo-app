package echoapp

import (
	"errors"
	"github.com/labstack/echo"
	"net/http"
)

const (
	CodeNotFound   = 400
	CodeNoAuth     = 401
	CodeDBError    = 402
	CodeCacheError = 403
	CodeArgument   = 404
	CodeNotAllow   = 405
	CodeEtcdError  = 406
	CodeInnerError = 501
)

var ErrNotFoundCache = errors.New("not found cache item")
var ErrNotFoundDb = errors.New("not found db item")
var ErrDb = errors.New("db exec err")
var ErrNotFoundEtcd = errors.New("not found etcd item")
var ErrArgument = errors.New("argument error")
var ErrNotLogin = errors.New("not login")
var ErrNotAuth = errors.New("not auth")
var ErrNotAllow = errors.New("not allow")

type Response struct {
	ErrorCode  int         `json:"code"`
	Msg        string      `json:"msg"`
	InnerError string      `json:"-"`
	Data       interface{} `json:"data"`
}

type BaseController struct {
}

func (b *BaseController) Success(ctx echo.Context, data interface{}) error {
	response := Response{
		ErrorCode: 200,
		Msg:       "success",
		Data:      data,
	}
	return ctx.JSON(http.StatusOK, response)
}

func (b *BaseController) Fail(ctx echo.Context, errcode int, msg string, innerErr error) error {
	response := Response{
		ErrorCode: errcode,
		Msg:       msg,
	}
	ctx.JSON(http.StatusOK, response)

	return innerErr
}
