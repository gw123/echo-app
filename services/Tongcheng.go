package services

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TongchengService struct {
	ConsumeNoticeUrl   string
	tongchengOptionMap map[string]echoapp.TongchengOption
	mu                 sync.Mutex
}

func NewTongchengService(options map[string]echoapp.TongchengOption) *TongchengService {
	return &TongchengService{
		tongchengOptionMap: options,
	}
}

func (mSvr TongchengService) CheckTicket(info *echoapp.CheckTicketJob) error {
	comInfo, ok := mSvr.tongchengOptionMap[info.ComId]
	if !ok {
		return errors.New("com not found")
	}

	encryptBody, err := mSvr.encryptBody([]byte(comInfo.Key), info)
	if err != nil {
		return errors.Wrap(err, "encryptBody")
	}

	t := echoapp.TongchengConsumeNoticeRequest{
		RequestHead: echoapp.TongchengRequestHead{
			UserId:    comInfo.UserId,
			Method:    "ConsumeNotice",
			Timestamp: time.Now().Unix(),
			Version:   "v1.0",
			Sign:      "",
		},
		RequestBody: info.TongchengRequestBody,
	}
	tt, _ := json.Marshal(t)
	log.Infof("tt : %s", string(tt))

	tongchengRequest := echoapp.TongchengRequest{
		RequestHead: echoapp.TongchengRequestHead{
			UserId:    comInfo.UserId,
			Method:    "ConsumeNotice",
			Timestamp: time.Now().Unix(),
			Version:   "v1.0",
			Sign:      "",
		},
		RawRequestBody: encryptBody,
	}
	signStr := mSvr.Sign(comInfo.Key, tongchengRequest)
	tongchengRequest.RequestHead.Sign = signStr
	log.Infof("tongchengRequest.SingStr: %s", signStr)
	log.Infof("tongchengRequest.Body: %s", tongchengRequest.RawRequestBody)

	base64ResponseData, err := mSvr.DoRequest(mSvr.ConsumeNoticeUrl, tongchengRequest)
	if err != nil {
		return errors.Wrap(err, "CheckTicket->DoResponse:"+string(base64ResponseData))
	}

	responseData := make([]byte, 0)
	if _, err = base64.StdEncoding.Decode(responseData, base64ResponseData); err != nil {
		return errors.Wrap(err, "CheckTicket->Decode:"+string(base64ResponseData))
	}

	tongchengResponse := &echoapp.TongchengResponse{}
	if err = json.Unmarshal(responseData, tongchengResponse); err != nil {
		return errors.Wrap(err, "CheckTicket->Decode:"+string(base64ResponseData))
	}

	if tongchengResponse.ResponseHead.ResCode == "1000" {
		return nil
	}
	//echoapp_util.DecryptDESECB(tongchengResponse.ResponseBody)
	return errors.Errorf("Code:%s,Msg:%s", tongchengResponse.ResponseHead.ResCode, tongchengResponse.ResponseHead.ResMsg)
}

func (mSvr TongchengService) Sign(key string, request echoapp.TongchengRequest) string {
	strBuf := strings.Builder{}
	strBuf.WriteString(request.RequestHead.UserId)
	strBuf.WriteString(request.RequestHead.Method)
	strBuf.WriteString(strconv.Itoa(int(request.RequestHead.Timestamp)))
	strBuf.WriteString(request.RequestHead.Version)
	strBuf.WriteString(request.RawRequestBody)
	strBuf.WriteString(key)
	hash := md5.New()
	hash.Write([]byte(strBuf.String()))
	strBuf.Reset()
	return hex.EncodeToString(hash.Sum(nil))
}

func (mSvr TongchengService) DoRequest(url string, params interface{}) ([]byte, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->Marshal")
	}
	base64Body := make([]byte, len(body)*2)
	base64.StdEncoding.Encode(base64Body, body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(base64Body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
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

func (mSvr TongchengService) encryptBody(key []byte, body interface{}) (string, error) {
	toEncryptStr, err := json.Marshal(body)
	if err != nil {
		return "", errors.Wrap(err, "encryptyBody")
	}
	data, err := echoapp_util.EntryptDesECB(toEncryptStr, key)
	if err != nil {
		return "", errors.Wrap(err, "encryptyBody")
	}
	return data, nil
}
