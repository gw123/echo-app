package echoapp

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
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

func (b *BaseController) AppErr(ctx echo.Context, appError AppError) error {
	response := Response{
		ErrorCode: appError.GetCode(),
		Msg:       appError.GetOuter(),
	}
	ctx.JSON(http.StatusOK, response)
	if appError.GetInner().Error() == appError.GetOuter() {
		return appError.GetInner()
	}
	return errors.Wrap(appError.GetInner(), appError.GetOuter())
}
