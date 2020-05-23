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
		return rCtrl.Fail(ctx, echoapp.Err_Argument, err.Error(), err)
	}
	err := rCtrl.resourceSvc.SaveResource(params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return rCtrl.Fail(ctx, echoapp.Err_NotFound, "用户不存在", err)
		} else {
			return rCtrl.Fail(ctx, echoapp.Err_DBError, "系统异常", err)
		}
	}
	return rCtrl.Success(ctx, nil)
}
func (resourceCtrl *ResourceController) GetResourceById(c echo.Context) error {
	id := c.QueryParam("id")
	id_int64, _ := strconv.ParseInt(id, 10, 64)
	res, err := resourceCtrl.resourceSvc.GetResourceById(c, id_int64)
	if err != nil {
		return resourceCtrl.Fail(c, echoapp.Err_Argument, "", errors.Wrap(err, "GetResourceById"))
	}
	return resourceCtrl.Success(c, res)
}

// type Params struct {
// 	Id    uint `json:"id"`
// 	From  int  `json:"from"`
// 	Limit int  `json:"limit"`
// 	TagId uint `json:"tag_id"`
// }

func (resourceCtrl *ResourceController) GetResourcesByTagId(c echo.Context) error {

	id := c.QueryParam("tagId")
	id_int64, _ := strconv.ParseInt(id, 10, 64)
	from := c.QueryParam("from")
	limit := c.QueryParam("limit")
	fromint, _ := strconv.Atoi(from)
	limitint, _ := strconv.Atoi(limit)
	res, err := resourceCtrl.resourceSvc.GetResourcesByTagId(c, id_int64, fromint, limitint)
	if err != nil {
		return resourceCtrl.Fail(c, echoapp.Err_DBError, "resourceCtrl.resourceSvc", errors.Wrap(err, "GetResourceByTagId"))
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
		return resourceCtrl.Fail(c, echoapp.Err_NotFound, "", errors.Wrap(err, "GetResourceByTagId"))
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
		return rCtrl.Fail(c, echoapp.Err_Argument, "rCtrl.resourceSvc.GetSelfResources", err)
	}
	return rCtrl.Success(c, res)
}
func (rCtrl *ResourceController) UploadResource(ctx echo.Context) error {
	//formval := ctx.FormValue("PPT")
	fileOption, err := rCtrl.resourceSvc.UploadFile(ctx, "file", echoapp.ConfigOpts.Asset.ResourceRoot, echoapp.ConfigOpts.Asset.UploadMaxFileSize)
	if err != nil {
		return rCtrl.Fail(ctx, echoapp.Err_NotAllow, "resourceSvc->UploadFile", err)
	}
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return rCtrl.Fail(ctx, echoapp.Err_CacheError, "echoapp_util.GetCtxtUserId", err)
	}
	md5fileStr32, err := echoapp_util.Md5SumFile(fileOption["uploadpath"])
	if err != nil {
		return rCtrl.Fail(ctx, echoapp.Err_Argument, "rCtrl->echoapp_util->Md5SumFile", err)
	}
	md5path := md5fileStr32[:2] + md5fileStr32 + path.Ext(fileOption["filename"])
	if _, err := rCtrl.resourceSvc.GetResourceByMd5Path(ctx, md5path); err == nil {
		return rCtrl.Fail(ctx, echoapp.Err_DBError, "file exits", errors.New("file exits"))
	}
	filetype := echoapp_util.GetFileType(fileOption["filename"])
	if err := echoapp_util.Copy(echoapp.ConfigOpts.Asset.StorageRoot+"/"+filetype+"/"+md5path, fileOption["uploadpath"]); err != nil {
		return rCtrl.Fail(ctx, echoapp.Err_EtcdError, " UploadResource echoapp_util.Copy", err)
	}
	putret, err := echoapp_util.UploadFileToQiniu(fileOption["uploadpath"], "/"+filetype+"/"+md5path)
	if err != nil {

		return rCtrl.Fail(ctx, echoapp.Err_EtcdError, " UploadResource echoapp_util->UploadFileToQiniu", err)
	}
	urlstrarr, err := echoapp_util.GetPPTCoverUrl(echoapp.ConfigOpts.Asset.MyURL + "/" + filetype + "/" + fileOption["filename"])
	if err != nil {
		return rCtrl.Fail(ctx, echoapp.Err_Argument, "GetPPTCoverUrl", err)
	}
	var data []byte
	if len(urlstrarr) > 9 {
		data, _ = json.Marshal(urlstrarr[:9])
	} else {
		data, _ = json.Marshal(urlstrarr)
	}
	goods := &echoapp.Goods{
		UserId:     userId,
		Name:       fileOption["filename"],
		Price:      0.30,
		GoodType:   filetype,
		RealPrice:  0.50,
		Covers:     string(data),
		SmallCover: urlstrarr[0],
		TagStr:     path.Dir(fileOption["uploadpath"]),
		Pages:      len(urlstrarr),
	}
	if err := rCtrl.goodsSvc.SaveGoods(goods); err != nil {
		return rCtrl.Fail(ctx, echoapp.Err_DBError, "rCtrl->goodsSvc.SaveGoods", err)
	}
	res, err := rCtrl.goodsSvc.GetGoodsByName(path.Base(fileOption["filename"]))
	if err != nil {
		return rCtrl.Fail(ctx, echoapp.Err_DBError, "rCtrl->goodsSvc->GetGoodsByName", err)
	}
	tag := &echoapp.Tags{
		Name: path.Dir(fileOption["uploadpath"]),
	}
	if err := rCtrl.goodsSvc.SaveTags(tag); err != nil {
		return rCtrl.Fail(ctx, echoapp.Err_DBError, "rCtrl->resourceSvc->SaveTags", err)
	}
	res_tag, err := rCtrl.goodsSvc.GetTagsByName(path.Dir(fileOption["uploadpath"]))
	if err != nil {
		return rCtrl.Fail(ctx, echoapp.Err_DBError, "rCtrl->goodsSvc->GetGoodsByName", err)
	}

	resource := &echoapp.Resource{
		GoodsId:    res.ID,
		UserId:     userId,
		Type:       filetype,
		Name:       fileOption["filename"],
		Covers:     string(data),
		SmallCover: urlstrarr[0],
		TagId:      res_tag.ID,
		Pages:      len(urlstrarr),
		Path:       md5path,
		Status:     "user_upload",
	}
	if err := rCtrl.resourceSvc.SaveResource(resource); err != nil {
		return rCtrl.Fail(ctx, echoapp.Err_DBError, "rCtrl->resourceSvc->SaveResource", err)
	}

	result := map[string]interface{}{
		"goods":    goods,
		"resource": resource,
		"tag":      tag,
		"qiniures": putret,
	}
	return rCtrl.Success(ctx, result)
}
func (rCtrl *ResourceController) DownloadResource(c echo.Context) error {
	userId, err := echoapp_util.GetCtxtUserId(c)
	name := c.QueryParam("name")
	downloadPath := c.QueryParam("downloadPath")
	resource, err := rCtrl.resourceSvc.GetResourceByName(name)
	if err != nil {
		return rCtrl.Fail(c, echoapp.Err_DBError, "resourceSvc.GetResourceByName", err)
	}
	if resource.UserId != userId {
		return rCtrl.Fail(c, echoapp.Err_NoAuth, "User Do Not Have Authority", err)
	}
	filetype := echoapp_util.GetFileType(name)
	path := echoapp.ConfigOpts.Asset.MyURL + "/" + filetype + "/" + name

	if err = echoapp_util.DownloadFile(path, downloadPath); err != nil {
		return rCtrl.Fail(c, echoapp.Err_Argument, "resourceSvc.DownloadFile", err)
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
		return rCtrl.Fail(c, echoapp.Err_Argument, "GetResourceList", err)
	}
	echoapp_util.ExtractEntry(c).Infof("from:%s,limit:%s", from, limit)
	return rCtrl.Success(c, filelist)
}
func (rCtrl *ResourceController) GetResourceByName(c echo.Context) error {
	path := c.QueryParam("name")
	res, err := rCtrl.resourceSvc.GetResourceByName(path)
	if err != nil {
		return rCtrl.Fail(c, echoapp.Err_Argument, "", errors.Wrap(err, "ResourceController->GetResourceByName"))
	}
	return rCtrl.Success(c, res)
}
