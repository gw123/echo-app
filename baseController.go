package echoapp

import (
	"net/http"

	"github.com/labstack/echo"
)


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
