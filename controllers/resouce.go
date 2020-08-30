package controllers

import (
	"encoding/json"
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
	goodsSvc    echoapp.GoodsService
	echoapp.BaseController
}

func NewResourceController(resourceSvr echoapp.ResourceService, goodsSvr echoapp.GoodsService) *ResourceController {
	return &ResourceController{
		resourceSvc: resourceSvr,
		goodsSvc:    goodsSvr,
	}
}

func (rCtrl *ResourceController) SaveResource(ctx echo.Context) error {
	params := &echoapp.Resource{}
	if err := ctx.Bind(params); err != nil {
		return rCtrl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	err := rCtrl.resourceSvc.SaveResource(params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rCtrl.Fail(ctx, echoapp.CodeNotFound, echoapp.ErrNotFoundDb.Error(), err)
		} else {
			return rCtrl.Fail(ctx, echoapp.CodeInnerError, echoapp.ErrNotFoundEtcd.Error(), err)
		}
	}
	return rCtrl.Success(ctx, nil)
}
func (resourceCtrl *ResourceController) GetResourceById(c echo.Context) error {
	id := c.QueryParam("id")
	id_int64, _ := strconv.ParseInt(id, 10, 64)
	res, err := resourceCtrl.resourceSvc.GetResourceById(c, id_int64)
	if err != nil {
		return resourceCtrl.Fail(c, echoapp.CodeNotFound, "", errors.Wrap(err, "GetResourceById"))
	}
	return resourceCtrl.Success(c, res)
}

func (resourceCtrl *ResourceController) GetResourcesByTagId(c echo.Context) error {

	id := c.QueryParam("tagId")
	id_int64, _ := strconv.ParseInt(id, 10, 64)
	from := c.QueryParam("from")
	limit := c.QueryParam("limit")
	fromint, _ := strconv.Atoi(from)
	limitint, _ := strconv.Atoi(limit)
	res, err := resourceCtrl.resourceSvc.GetResourcesByTagId(c, id_int64, fromint, limitint)
	if err != nil {
		return resourceCtrl.Fail(c, echoapp.CodeNotFound, "resourceCtrl.resourceSvc", errors.Wrap(err, "GetResourceByTagId"))
	}
	return resourceCtrl.Success(c, res)
}

func (resourceCtrl *ResourceController) GetUserPaymentResources(c echo.Context) error {
	userId, _ := echoapp_util.GetCtxtUserId(c)
	from := c.QueryParam("from")
	limit := c.QueryParam("limit")
	fromint, _ := strconv.Atoi(from)
	limitint, _ := strconv.Atoi(limit)
	res, err := resourceCtrl.resourceSvc.GetUserPaymentResources(c, userId, fromint, limitint)
	if err != nil {
		return resourceCtrl.Fail(c, echoapp.CodeNotFound, "", errors.Wrap(err, "GetResourceByTagId"))
	}
	return resourceCtrl.Success(c, res)
}

func (rCtrl *ResourceController) GetSelfResources(c echo.Context) error {
	userId, _ := echoapp_util.GetCtxtUserId(c)
	from := c.QueryParam("from")
	limit := c.QueryParam("limit")
	fromint, _ := strconv.Atoi(from)
	limitint, _ := strconv.Atoi(limit)
	res, err := rCtrl.resourceSvc.GetSelfResources(c, userId, fromint, limitint)
	if err != nil {
		return rCtrl.Fail(c, echoapp.CodeNotFound, "rCtrl.resourceSvc.GetSelfResources", err)
	}
	return rCtrl.Success(c, res)
}

func (rCtrl *ResourceController) DownloadResource(c echo.Context) error {
	userId, err := echoapp_util.GetCtxtUserId(c)
	name := c.QueryParam("name")
	downloadPath := c.QueryParam("downloadPath")
	resource, err := rCtrl.resourceSvc.GetResourceByName(name)
	if err != nil {
		return rCtrl.Fail(c, echoapp.CodeNotFound, "resourceSvc.GetResourceByName", errors.Wrap(err, "rCtrl ->resourceSvc->GetResourceByName"))
	}
	if resource.UserId != userId {
		return rCtrl.Fail(c, echoapp.CodeNoAuth, echoapp.ErrNotAuth.Error(), nil)
	}
	filetype := echoapp_util.GetFileType(name)
	path := echoapp.ConfigOpts.ResourceOptions.BaseURL + "/" + filetype + "/" + name

	if err = echoapp_util.DownloadFile(path, downloadPath); err != nil {
		return rCtrl.Fail(c, echoapp.CodeCacheError, "", errors.Wrap(err, "rCtrl->echoapp_util.DownloadFile"))
	}
	res := map[string]interface{}{
		"FileNanme":    name,
		"DownloadPath": downloadPath,
		"FileInfo":     resource,
	}
	return rCtrl.Success(c, res)
}
func (rCtrl *ResourceController) GetResourceList(c echo.Context) error {
	from := c.QueryParam("from")
	limit := c.QueryParam("limit")
	fromint, _ := strconv.Atoi(from)
	limitint, _ := strconv.Atoi(limit)
	filelist, err := rCtrl.resourceSvc.GetResourceList(c, fromint, limitint)
	if err != nil {
		return rCtrl.Fail(c, echoapp.CodeNotFound, "", errors.Wrap(err, "ResourceController->GetResourceList"))
	}
	echoapp_util.ExtractEntry(c).Infof("from:%s,limit:%s", from, limit)
	return rCtrl.Success(c, filelist)
}

func (rCtrl *ResourceController) GetResourceByName(c echo.Context) error {
	path := c.QueryParam("name")
	res, err := rCtrl.resourceSvc.GetResourceByName(path)
	if err != nil {
		return rCtrl.Fail(c, echoapp.CodeNotFound, "", errors.Wrap(err, "ResourceController->GetResourceByName"))
	}
	return rCtrl.Success(c, res)
}

func (rCtrl *ResourceController) UploadResource(ctx echo.Context) error {
	newFile, err := rCtrl.resourceSvc.UploadFile(ctx, "file", echoapp.ConfigOpts.Asset.ResourceRoot, echoapp.ConfigOpts.ResourceOptions.UploadMaxFileSize)
	if err != nil {
		return rCtrl.Fail(ctx, echoapp.CodeNotAllow, "resourceSvc->UploadFile", err)
	}
	// userId, err := echoapp_util.GetCtxtUserId(ctx)
	// if err != nil {
	// 	return rCtrl.Fail(ctx, echoapp.CodeNotFound, "echoapp_util.GetCtxtUserId", err)
	// }

	md5path := newFile.Md5[:2] + "/" + newFile.Md5 + path.Ext(newFile.Name)
	if _, err := rCtrl.resourceSvc.GetResourceByMd5Path(ctx, md5path); err == nil {
		return rCtrl.Fail(ctx, echoapp.CodeNotFound, "file exits", errors.New("file exits"))
	}
	// if err := echoapp_util.Copy(echoapp.ConfigOpts.Asset.StorageRoot+"/"+newFile.Type+"/"+md5path, newFile.Path); err != nil {
	// 	return rCtrl.Fail(ctx, echoapp.CodeCacheError, " UploadResource echoapp_util.Copy", err)
	// }
	// putret, err := echoapp_util.UploadFileToQiniu(newFile.Path, "/"+newFile.Type+"/"+md5path)
	// if err != nil {
	// 	return rCtrl.Fail(ctx, echoapp.CodeCacheError, " UploadResource echoapp_util->UploadFileToQiniu", err)
	// }
	urlstrarr, err := echoapp_util.GetPPTCoverUrl(newFile.Url)
	if err != nil {
		return rCtrl.Fail(ctx, echoapp.CodeCacheError, "GetPPTCoverUrl", err)
	}
	var data []byte
	if len(urlstrarr) > 9 {
		data, _ = json.Marshal(urlstrarr[:9])
	} else {
		data, _ = json.Marshal(urlstrarr)
	}
	goods := &echoapp.Goods{
		GoodsBrief: echoapp.GoodsBrief{
			UserID:     uint(newFile.UserId),
			Name:       newFile.Name,
			RealPrice:  0,
			Price:      0,
			GoodsType:  newFile.Type,
			CoversStr:  string(data),
			SmallCover: urlstrarr[0],
		},
	}
	if err := rCtrl.goodsSvc.Save(goods); err != nil {
		return rCtrl.Fail(ctx, echoapp.CodeDBError, "", errors.Wrap(err, "rCtrl->goodsSvc.SaveGoods"))
	}

	oldGoods, err := rCtrl.goodsSvc.GetGoodsByName(path.Base(newFile.Name))
	if err != nil {
		return rCtrl.Fail(ctx, echoapp.CodeNotFound, "", errors.Wrap(err, "rCtrl->goodsSvc->GetGoodsByName"))
	}

	res_tag, err := rCtrl.goodsSvc.GetTagByName(path.Dir(newFile.Path))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			tag := &echoapp.GoodsTag{
				Name: path.Dir(newFile.Path),
			}
			if err := rCtrl.goodsSvc.SaveTag(tag); err != nil {
				return rCtrl.Fail(ctx, echoapp.CodeDBError, "", errors.Wrap(err, "rCtrl->resourceSvc->SaveTags"))
			}
		}
		return rCtrl.Fail(ctx, echoapp.CodeNotFound, "rCtrl->goodsSvc->GetGoodsByName", err)
	}

	resource := &echoapp.Resource{
		GoodsId:    int64(oldGoods.ID),
		UserId:     int64(newFile.UserId),
		Type:       newFile.Type,
		Name:       newFile.Name,
		Covers:     string(data),
		SmallCover: urlstrarr[0],
		TagId:      int64(res_tag.ID),
		Pages:      len(urlstrarr),
		Path:       newFile.Path,
		Status:     "user_upload",
	}
	if err := rCtrl.resourceSvc.SaveResource(resource); err != nil {
		return rCtrl.Fail(ctx, echoapp.CodeDBError, "rCtrl->resourceSvc->SaveResource", err)
	}

	result := map[string]interface{}{
		"goods":    goods,
		"resource": resource,
		"qiniures": newFile.Bucket,
	}
	return rCtrl.Success(ctx, result)
}

func (rCtrl *ResourceController) UploadImage(ctx echo.Context) error {
	newFile, err := rCtrl.resourceSvc.UploadFile(ctx, "image", echoapp.ConfigOpts.Asset.StorageRoot, echoapp.ConfigOpts.ResourceOptions.UploadMaxFileSize)
	if err != nil {
		return rCtrl.Fail(ctx, echoapp.CodeNotAllow, "resourceSvc->UploadFile", err)
	}
	return rCtrl.Success(ctx, map[string]string{"path": newFile.Url})
}
