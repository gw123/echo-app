package controllers

import (
	"fmt"
	echoapp "github.com/gw123/echo-app"
	"github.com/labstack/echo"
)

type ResourceController struct {
	resourceSvr echoapp.ResourceService
	echoapp.BaseController
}

func NewResourceController(resourceSvr echoapp.ResourceService) *ResourceController {
	return &ResourceController{
		resourceSvr: resourceSvr,
	}
}

func (rCtl *ResourceController) GetUploadToken(ctx echo.Context) error {
	token, err := rCtl.resourceSvr.GetUploadToken(0)
	if err != nil {
		return rCtl.Fail(ctx, echoapp.CodeArgument, "", err)
	}
	return rCtl.Success(ctx, token)
}

type UploadCallbackParams struct {
	Key    string `json:"key"`
	Hash   string `json:"hash"`
	Size   int    `json:"fsize"`
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

func (rCtl *ResourceController) UploadCallback(ctx echo.Context) error {
	params := UploadCallbackParams{}
	err := ctx.Bind(&params)
	if err != nil {
		return err
	}
	fmt.Println(params)
	return nil
}
