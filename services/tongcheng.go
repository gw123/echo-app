package services

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
)

type TongchengService struct {
	ConsumeNoticeUrl   string
	tongchengOptionMap map[string]echoapp.TongchengOption
	mu                 sync.Mutex
}

func NewTongchengService(config echoapp.TongchengConfig) *TongchengService {
	return &TongchengService{
		ConsumeNoticeUrl:   config.NotifyUrl,
		tongchengOptionMap: config.ClientMap,
	}
}

func (mSvr TongchengService) SyncPartnerCode(comID uint, info *echoapp.SyncPartnerCodeRequestBody) error {
	_, err := mSvr.TongchengRequest(strconv.Itoa(int(comID)), "RetCodeNotice", info)
	if err != nil {
		return err
	}
	return nil
}

func (mSvr TongchengService) CheckTicket(comID uint, info *echoapp.CheckTicketRequestBody) error {
	_, err := mSvr.TongchengRequest(strconv.Itoa(int(comID)), "ConsumeNotice", info)
	if err != nil {
		return err
	}
	return nil
}

func (mSvr TongchengService) TongchengRequest(token string, method string, requestBody interface{}) (*echoapp.TongchengResponse, error) {
	comInfo, ok := mSvr.tongchengOptionMap[token]
	if !ok {
		return nil, errors.Errorf("com: %d not found", token)
	}

	rawbody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal")
	}

	encryptBody, err := echoapp_util.EntryptDesECB([]byte(rawbody), []byte(comInfo.Key))
	if err != nil {
		return nil, errors.Wrap(err, "encryptbody")
	}
	tongchengRequest := echoapp.TongchengRequest{
		RequestHead: echoapp.TongchengRequestHead{
			UserId:    comInfo.UserId,
			Method:    method,
			Timestamp: time.Now().Unix(),
			Version:   "v1.0",
			Sign:      "",
		},
		EncryptRequestBody: encryptBody,
	}
	signStr := mSvr.Sign(comInfo.Key, tongchengRequest)
	tongchengRequest.RequestHead.Sign = signStr
	log.Infof("tongchengRequest.SingStr: %s", signStr)
	log.Infof("tongchengRequest.Body: %+v", tongchengRequest)

	base64ResponseData, err := mSvr.DoRequest(mSvr.ConsumeNoticeUrl, tongchengRequest)
	if err != nil {
		return nil, errors.Wrap(err, "CheckTicket->DoResponse:"+string(base64ResponseData))
	}

	var responseData []byte
	if responseData, err = base64.StdEncoding.DecodeString(string(base64ResponseData)); err != nil {
		return nil, errors.Wrap(err, "CheckTicket->Base64Decode:"+string(base64ResponseData))
	}
	echoapp_util.DefaultLogger().Info(string(responseData))
	tongchengResponse := &echoapp.TongchengResponse{}
	if err = json.Unmarshal(responseData, tongchengResponse); err != nil {
		return nil, errors.Wrap(err, "CheckTicket->JsonDecode:"+string(base64ResponseData))
	}

	if tongchengResponse.ResponseHead.ResCode == "1000" {
		return tongchengResponse, nil
	}
	//echoapp_util.DecryptDESECB(tongchengResponse.ResponseBody)
	return nil, errors.Errorf("Code:%s,Msg:%s", tongchengResponse.ResponseHead.ResCode, tongchengResponse.ResponseHead.ResMsg)
}

/***
9a86097b-b95d-4fd4-bbb9-a18aaafc84b1ConsumeNotice1588219228v1.0sBE4yQDodGqnKpe0BfeLzxdb6ntDQdaRlIEbgsS8OViJwcbTMydj8WVpT8Hgrd3Jq+lT4dz1ULPPWnew344FmLysGcYYLLRF5k1xLiNMaFsJv3ykoK1hao1OuBKZekWp4MTU1KBG
*/
func (mSvr TongchengService) Sign(key string, request echoapp.TongchengRequest) string {
	strBuf := strings.Builder{}
	strBuf.WriteString(request.RequestHead.UserId)
	strBuf.WriteString(request.RequestHead.Method)
	strBuf.WriteString(strconv.Itoa(int(request.RequestHead.Timestamp)))
	strBuf.WriteString(request.RequestHead.Version)
	strBuf.WriteString(request.EncryptRequestBody)
	strBuf.WriteString(key)
	log.Info("MD5 Content: " + strBuf.String())
	rawMd5 := md5.Sum([]byte(strBuf.String()))
	sign := fmt.Sprintf("%x", rawMd5)
	defer strBuf.Reset()
	return sign
}

func (mSvr TongchengService) DoRequest(url string, params interface{}) ([]byte, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->Marshal")
	}
	log.Info(string(body))
	str := base64.StdEncoding.EncodeToString(body)
	log.Info("请求base64 ： " + string(str))
	req, err := http.NewRequest("POST", url, strings.NewReader(str))
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->http.NewRequest")
	}
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
