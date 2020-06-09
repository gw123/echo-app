package echoapp_util

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"time"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	echoapp "github.com/gw123/echo-app"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func DoHttpRequest(url string, method string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->http.NewRequest")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->Do")
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->ReadAll")
	}
	if res.StatusCode == http.StatusOK {
		return data, errors.Wrapf(err, "StatusCode:%d", res.StatusCode)
	}
	return data, nil
}

func GetPPTCoverUrl(pptUrl string) ([]string, error) {
	clientMap := echoapp.ConfigOpts.PPTImages
	for _, options := range clientMap {
		//echoapp_util.DefaultLogger().Infof("访问%s,com_id:%d", key, options.ComId)
		url := options.BaseUrl + "onlinePreview" + "?url=" + pptUrl
		// url err ?
		data, err := DoHttpRequest(url, "GET")
		if err != nil {
			return nil, errors.Wrap(err, "doHttpRequest")
		}

		reg := regexp.MustCompile(`\<img .*?title=\"查看大图\" .*?data-src=\"(\S*)\" .*?\>`)
		if reg == nil {
			return nil, errors.Wrap(err, "regexp.MustCompile err")
		}
		res := reg.FindAllStringSubmatch(string(data), -1)
		fmt.Println(len(res[0]))
		if len(res) > 0 && len(res[0]) > 0 {
			urls := make([]string, 0)
			for _, text := range res {
				url := text[1]
				urls = append(urls, url)
			}
			return urls, nil
		}
	}
	return nil, nil
}

func Copy(dst, src string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	if _, err = io.Copy(destination, source); err != nil {
		return err
	}
	return nil
}

type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

func UploadFileToQiniu(localFile, key string) (*MyPutRet, error) {

	bucket := "testxytfile"
	//key := "github-x.png"
	putPolicy := storage.PutPolicy{
		Scope:      bucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
	mac := qbox.NewMac(echoapp.ConfigOpts.QiniuKeys.AccessKey, echoapp.ConfigOpts.QiniuKeys.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)

	ret := MyPutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "xyt ppt",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return nil, errors.Wrap(err, "formUploader.PutFile")
	}
	return &ret, nil
}
func GetFileType(filename string) string {
	s := path.Ext(filename)
	return s[1:]
}
func Md5SumFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	rawMd5 := md5.Sum(data)
	sign := fmt.Sprintf("%x", rawMd5)
	return sign, nil
}
func UploadFile(c echo.Context, formname, uploadpath string, maxfilesize int64) (map[string]string, error) {
	//r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize)
	file, err := c.FormFile(formname)
	if err != nil {
		return nil, err
	}
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	filetype := GetFileType(file.Filename)
	fullPath := uploadpath + "/" + filetype + "/" + file.Filename
	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()
	size, err := io.Copy(dst, src)
	if err != nil {
		return nil, err
	}
	if size > maxfilesize {
		return nil, errors.Wrapf(err, "File size over limit")
	}
	// userId, err := GetCtxtUserId(c)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "echoapp_util.GetCtxtUserId")
	// }
	res := map[string]string{
		//"userId":     strconv.FormatInt(userId, 10),
		"fileName":   file.Filename,
		"uploadPath": fullPath,
	}
	return res, nil
}
func DownloadFile(durl, localpath string) error {
	uri, err := url.ParseRequestURI(durl)
	if err != nil {
		return errors.Wrap(err, "ParseRequestURI")
	}
	filename := path.Base(uri.Path)

	req, err := http.NewRequest("GET", durl, nil)
	if err != nil {
		return errors.Wrap(err, "DoRequest->http.NewRequest")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "DoRequest->Do")
	}
	http.DefaultClient.Timeout = time.Second * 60 //超时设置

	file, err := os.Create(localpath + filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		return errors.New("Body is Null")
	}
	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}
