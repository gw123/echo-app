package echoapp

import (
	"github.com/labstack/echo"
	"net/http"
)

const (
	Err_NotFound   = 400
	Err_NoAuth     = 401
	Err_DBError    = 402
	Err_CacheError = 403
	Err_Argument   = 404
	Err_NotAllow   = 405
	Err_EtcdError  = 406
	Err_InnerError = 501
)

type Response struct {
	ErrorCode  int         `json:"code"`
	Msg        string      `json:"msg"`
	InnerError string      `json:"inner_error"`
	Data       interface{} `json:"data"`
}

type BaseController struct {
}

func (b *BaseController) Success(ctx echo.Context, data interface{}) error {
	response := Response{
		ErrorCode: 0,
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
