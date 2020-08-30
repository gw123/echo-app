package echoapp_util

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type ResponseToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}
type AccessTokenOption struct {
	GrantType string `json:"grant_type"`
	AppId     string `json:"appid"`
	Secret    string `json:"secret"`
}

func GetAccessToken() (*ResponseToken, error) {
	DefaultJsonLogger().Infof("获取微信AccessToken")
	getTokenRequestParam := &AccessTokenOption{
		GrantType: "client_credential",
		AppId:     "wx97a1f596e4b87a83",
		Secret:    "ea389cb3939a61965a6cd5aba79e09e0",
	}
	baseurl := "https://api.weixin.qq.com/cgi-bin/token"
	url := baseurl + fmt.Sprintf("?grant_type=%s&appid=%s&secret=%s",
		getTokenRequestParam.GrantType,
		getTokenRequestParam.AppId,
		getTokenRequestParam.Secret)

	res, err := DoHttpRequest(url, "GET", nil)
	if err != nil {
		return nil, errors.Wrap(err, "echoapp_util.DoHttpRequest")
	}
	responseToken := &ResponseToken{}
	if err = json.Unmarshal(res, responseToken); err != nil {
		return nil, errors.Wrap(err, "doReportHttpRequest")
	}
	DefaultLogger().Info(string(res))
	return responseToken, nil
}

type ResCallbackIp struct {
	IpList  []string `json:"ip_list"`
	Errcode int      `json:"errcode"`
	Errmsg  string   `json:"errmsg"`
}

//baseurl:="https://api.weixin.qq.com/cgi-bin"
func GetCallbackIp(baseurl string) (*ResCallbackIp, error) {
	accessTokenstruct, err := GetAccessToken()

	if err != nil {
		return nil, errors.Wrap(err, "GetAccessToken()")
	}
	var ResIpList = &ResCallbackIp{}
	if accessTokenstruct.AccessToken != "" {
		url := baseurl + "/getcallbackip?access_token=" + accessTokenstruct.AccessToken
		res, err := DoHttpRequest(url, "GET", nil)
		if err != nil {
			return nil, errors.Wrap(err, "DoHttpRequest")
		}
		if err = json.Unmarshal(res, ResIpList); err != nil {
			return nil, errors.Wrap(err, "Unmarshal")
		}
		DefaultLogger().Info(string(res))
		return ResIpList, nil
	} else {
		return nil, errors.New("AccessToken is nil")
	}
}
func GetApiDomainIp(baseurl string) (*ResCallbackIp, error) {
	accessTokenstruct, err := GetAccessToken()

	if err != nil {
		return nil, errors.Wrap(err, "GetAccessToken()")
	}
	var ResIpList = &ResCallbackIp{}
	if accessTokenstruct.AccessToken != "" {
		url := baseurl + "/get_api_domain_ip?access_token=" + accessTokenstruct.AccessToken
		res, err := DoHttpRequest(url, "GET", nil)
		if err != nil {
			return nil, errors.Wrap(err, "DoHttpRequest")
		}
		if err = json.Unmarshal(res, ResIpList); err != nil {
			return nil, errors.Wrap(err, "Unmarshal")
		}
		DefaultLogger().Info(string(res))
		return ResIpList, nil
	} else {
		return nil, errors.New("AccessToken is nil")
	}
}

type TemplateApiSetIndustryParams struct {
	IndustryId1 string `json:"industry_id1"`
	IndustryId2 string `json:"industry_id2"`
}

//POST https://api.weixin.qq.com/cgi-bin/template/api_set_industry?access_token=ACCESS_TOKEN
func PostTemplateApiSetIndustry(baseurl string, industryParams *TemplateApiSetIndustryParams) error {
	accessTokenstruct, err := GetAccessToken()
	if err != nil {
		return errors.Wrap(err, "GetAccessToken()")
	}
	if accessTokenstruct.AccessToken != "" {
		url := baseurl + "/template/api_set_industry?access_token=" + accessTokenstruct.AccessToken
		temp, err := json.Marshal(industryParams)
		if err != nil {
			return errors.Wrap(err, ".Marshal")
		}
		res, err := DoHttpRequest(url, "POST", temp)
		if err != nil {
			return errors.Wrap(err, "DoHttpRequest")
		}
		DefaultLogger().Info(string(res))
		return nil
	} else {
		return errors.New("AccessToken is nil")
	}
}

type ResponseIndustryParams struct {
	// 	PrimaryIndustray   map[string]string `json:"primary_industry"`
	// 	SecondaryIndustray map[string]string `json:"secondary_industry"`
	PrimaryIndustray   *Class `json:"primary_industry"`
	SecondaryIndustray *Class `json:"secondary_industry"`
}

type Class struct {
	FirstClass  string `json:"first_class"`
	SecondClass string `json:"second_class"`
}

//GET https://api.weixin.qq.com/cgi-bin/template/get_industry?access_token=ACCESS_TOKEN

func GetTemplateGetIndustry(baseurl string) (*TemplateApiSetIndustryParams, error) {
	accessTokenstruct, err := GetAccessToken()
	if err != nil {
		return nil, errors.Wrap(err, "GetAccessToken()")
	}
	var ResIpList = &TemplateApiSetIndustryParams{}
	if accessTokenstruct.AccessToken != "" {
		url := baseurl + "/template/get_industry?access_token=" + accessTokenstruct.AccessToken

		res, err := DoHttpRequest(url, "GET", nil)
		if err != nil {
			return nil, errors.Wrap(err, "DoHttpRequest")
		}
		DefaultLogger().Info(string(res))
		if err = json.Unmarshal(res, ResIpList); err != nil {
			return nil, errors.Wrap(err, "Unmarshal")
		}
		return ResIpList, nil
	} else {
		return nil, errors.New("AccessToken is nil")
	}
}

type TemplateID struct {
	TemplateIdShort string
}
type ResponseTemplateID struct {
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	TemplateId string `json:"template_id"`
}

//POST https://api.weixin.qq.com/cgi-bin/template/api_add_template?access_token=ACCESS_TOKEN
func PostApiAddTemplate(baseurl string, industryParams *TemplateID) (*ResponseTemplateID, error) {
	accessTokenstruct, err := GetAccessToken()
	if err != nil {
		return nil, errors.Wrap(err, "GetAccessToken()")
	}
	var response = &ResponseTemplateID{}
	if accessTokenstruct.AccessToken != "" {
		url := baseurl + "/template/api_add_template?access_token=" + accessTokenstruct.AccessToken
		temp, err := json.Marshal(industryParams)
		if err != nil {
			return nil, errors.Wrap(err, ".Marshal")
		}
		res, err := DoHttpRequest(url, "POST", temp)
		if err != nil {
			return nil, errors.Wrap(err, "DoHttpRequest")
		}
		DefaultLogger().Info(string(res))
		if err = json.Unmarshal(res, response); err != nil {
			return nil, errors.Wrap(err, "Unmarshal")
		}
		return response, nil
	} else {
		return nil, errors.New("AccessToken is nil")
	}
}

type ResponseTemplateList struct {
	TemplateList []*TemplateOptions
}
type TemplateOptions struct {
	TemplateId      string `json:"template_id"`
	Title           string `json:"title"`
	PrimaryIndustry string `json:"primary_industry"`
	DeputyIndustry  string `json:"deputy_industry"`
	Content         string `json:"content" `
	Example         string `son:"example" `
}

//GET https://api.weixin.qq.com/cgi-bin/template/get_all_private_template?access_token=ACCESS_TOKEN
func GetTemplateList(baseurl string) (*ResponseTemplateList, error) {
	accessTokenstruct, err := GetAccessToken()
	if err != nil {
		return nil, errors.Wrap(err, "GetAccessToken()")
	}
	var ResIpList = &ResponseTemplateList{}
	if accessTokenstruct.AccessToken != "" {
		url := baseurl + "/template/get_all_private_template?access_token=" + accessTokenstruct.AccessToken

		res, err := DoHttpRequest(url, "GET", nil)
		if err != nil {
			return nil, errors.Wrap(err, "DoHttpRequest")
		}
		DefaultLogger().Info(string(res))
		if err = json.Unmarshal(res, ResIpList); err != nil {
			return nil, errors.Wrap(err, "Unmarshal")
		}
		return ResIpList, nil
	} else {
		return nil, errors.New("AccessToken is nil")
	}
}

type DelTemplateOp struct {
	TemplateId string `json:"template_id"`
}
type Response struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

//http请求方式：POST https://api.weixin.qq.com/cgi-bin/template/del_private_template?access_token=ACCESS_TOKEN
func PostDelPrivateTemplate(baseurl string, industryParams *DelTemplateOp) (*Response, error) {
	accessTokenstruct, err := GetAccessToken()
	if err != nil {
		return nil, errors.Wrap(err, "GetAccessToken()")
	}
	var response = &Response{}
	if accessTokenstruct.AccessToken != "" {
		url := baseurl + "/template/del_private_template?access_token=" + accessTokenstruct.AccessToken
		temp, err := json.Marshal(industryParams)
		if err != nil {
			return nil, errors.Wrap(err, ".Marshal")
		}
		res, err := DoHttpRequest(url, "POST", temp)
		if err != nil {
			return nil, errors.Wrap(err, "DoHttpRequest")
		}
		DefaultLogger().Info(string(res))
		if err = json.Unmarshal(res, response); err != nil {
			return nil, errors.Wrap(err, "Unmarshal")
		}
		return response, nil
	} else {
		return nil, errors.New("AccessToken is nil")
	}
}

type SendTemplateOptions struct {
	Touser      string       `json:"touser"`
	TemplateId  string       `json:"template_id"`
	Url         string       `json:"url"`
	Miniprogram *Miniprogram `json:"miniprogram"`
	Data        map[string]*DataOp
}
type DataOp struct {
	Value string `json:"value"`
	Color string `json:"color"`
}
type Miniprogram struct {
	AppId    string `json:"app_id"`
	Pagepath string `json:"pagepath"`
}
type ResponseSend struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	MsgId   int64  `json:"msg_id"`
}

//POST https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=ACCESS_TOKEN
func PostSendPrivateTemplate(baseurl string, industryParams *SendTemplateOptions) (*ResponseSend, error) {
	accessTokenstruct, err := GetAccessToken()
	if err != nil {
		return nil, errors.Wrap(err, "GetAccessToken()")
	}
	var response = &ResponseSend{}
	if accessTokenstruct.AccessToken != "" {
		url := baseurl + "/message/template/send?access_token=" + accessTokenstruct.AccessToken
		temp, err := json.Marshal(industryParams)
		if err != nil {
			return nil, errors.Wrap(err, ".Marshal")
		}
		res, err := DoHttpRequest(url, "POST", temp)
		if err != nil {
			return nil, errors.Wrap(err, "DoHttpRequest")
		}
		DefaultLogger().Info(string(res))
		if err = json.Unmarshal(res, response); err != nil {
			return nil, errors.Wrap(err, "Unmarshal")
		}
		return response, nil
	} else {
		return nil, errors.New("AccessToken is nil")
	}
}
