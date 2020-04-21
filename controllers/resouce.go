package controllers

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type ResourceController struct {
	resourceSvc echoapp.ResourceService
	userSvc     echoapp.UserService
	echoapp.BaseController
}

func NewResourceController() *ResourceController {
	return &ResourceController{
		resourceSvc: app.MustGetResService(),
	}
}

func (rCtrl *ResourceController) SaveResource(ctx echo.Context) error {
	params := echoapp.Resource{}
	if err := ctx.Bind(params); err != nil {
		return rCtrl.Fail(ctx, echoapp.Error_ArgumentError, err.Error(), err)
	}
	err := rCtrl.resourceSvc.SaveResource(params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rCtrl.Fail(ctx, echoapp.Error_NotFound, "用户不存在", err)
		} else {
			return rCtrl.Fail(ctx, echoapp.Error_DBError, "系统异常", err)
		}
	}
	return rCtrl.Success(ctx, nil)
}
func (resourceCtrl *ResourceController) GetResourceById(c echo.Context) error {
	param := &Params{}
	if err := c.Bind(param); err != nil {
		return resourceCtrl.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "Bind"))
	}
	res, err := resourceCtrl.resourceSvc.GetResourceById(c, param.UserID)
	if err != nil {
		return resourceCtrl.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "GetResourceById"))
	}
	return resourceCtrl.Success(c, res)
}

type Params struct {
	UserID uint `json:"user_id"`
	From   int  `json:"from"`
	Limit  int  `json:"limit"`
	TagID  uint `json:"tag_id"`
}

func (resourceCtrl *ResourceController) GetResourcesByTagID(c echo.Context) error {
	/*tagid := c.QueryParams("tagid")
	from := c.QueryParams("from")
	limit := c.QueryParams("limit")
	*/
	params := &Params{}
	if err := c.Bind(params); err != nil {
		return resourceCtrl.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "Bind"))
	}
	res, err := resourceCtrl.resourceSvc.GetResourcesByTagID(c, params.TagID, params.From, params.Limit)
	if err != nil {
		return resourceCtrl.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "GetResourceByTagId"))
	}
	return resourceCtrl.Success(c, res)
}

func (resourceCtrl *ResourceController) GetUserPaymentResources(c echo.Context) error {
	params := &Params{}
	if err := c.Bind(params); err != nil {
		return resourceCtrl.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "Bind"))
	}
	res, err := resourceCtrl.resourceSvc.GetUserPaymentResources(c, params.UserID, params.From, params.Limit)
	if err != nil {
		return resourceCtrl.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "GetResourceByTagId"))
	}
	return resourceCtrl.Success(c, res)
}

func (resourceCtrl *ResourceController) GetSelfResources(c echo.Context) error {
	return nil
}
