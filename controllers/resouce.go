package controllers

import (
	"path"
	"strconv"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type ResourceController struct {
	resourceSvc echoapp.ResourceService
	//userSvc     echoapp.UserService
	echoapp.BaseController
}

func NewResourceController(resourceSvr echoapp.ResourceService) *ResourceController {
	return &ResourceController{
		resourceSvc: resourceSvr,
	}
}

func (rCtrl *ResourceController) SaveResource(ctx echo.Context) error {
	params := &echoapp.Resource{}
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
	res, err := resourceCtrl.resourceSvc.GetResourceById(c, param.ID)
	if err != nil {
		return resourceCtrl.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "GetResourceById"))
	}
	return resourceCtrl.Success(c, res)
}

type Params struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
	From   int  `json:"from"`
	Limit  int  `json:"limit"`
	TagID  uint `json:"tag_id"`
}

func (resourceCtrl *ResourceController) GetResourcesByTagId(c echo.Context) error {

	params := &Params{}
	if err := c.Bind(params); err != nil {
		return resourceCtrl.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "Bind"))
	}
	res, err := resourceCtrl.resourceSvc.GetResourcesByTagId(c, params.TagID, params.From, params.Limit)
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
func (rCtrl *ResourceController) UploadResource(ctx echo.Context) error {
	path, err := rCtrl.resourceSvc.UploadFile(ctx, "pdffile", echoapp.ConfigOpts.Asset.WatchRoot)
	if err != nil {
		return rCtrl.Fail(ctx, echoapp.Error_ArgumentError, "SaveFile", err)
	}
	return rCtrl.Success(ctx, path)
}
func (rCtrl *ResourceController) GetResourceList(c echo.Context) error {
	from := c.QueryParam("from")
	limit := c.QueryParam("limit")
	fromint, _ := strconv.Atoi(from)
	limitint, _ := strconv.Atoi(limit)
	filelist, err := rCtrl.resourceSvc.GetResourceList(c, fromint, limitint)
	if err != nil {
		return rCtrl.Fail(c, echoapp.Error_ArgumentError, "GetPDFList", err)
	}
	echoapp_util.ExtractEntry(c).Infof("from:%s,limit:%s", from, limit)
	return rCtrl.Success(c, filelist)
}
func (rCtrl *ResourceController) GetResourceByPath(c echo.Context) error {
	path := c.QueryParam("path")

	res, err := rCtrl.resourceSvc.GetResourceByPath(path)
	if err != nil {
		return rCtrl.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "GetResourceByPath"))
	}
	return rCtrl.Success(c, res)
}

// 路径有些问题 数据库存储的是 ev.name  文件所在的项目目录(/home/gh/.../resource/tmp)， path （是当前目录 ./resource/tmp）
func (rCtrl *ResourceController) GetResourceByName(c echo.Context) error {
	name := c.QueryParam("name")
	t := path.Ext(name)
	path := echoapp.ConfigOpts.Asset.WatchRoot + "/" + t[1:] + "/" + name
	echoapp_util.ExtractEntry(c).Info("path:", path)
	res, err := rCtrl.resourceSvc.GetResourceByPath(path)
	if err != nil {
		return rCtrl.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "GetResourceByPath"))
	}
	return rCtrl.Success(c, res)
}
