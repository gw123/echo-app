package echoapp

type CompanyBrief struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Addr  string `json:"addr"`
	Phone string `json:"phone"`
}

type Company struct {
	CompanyBrief
	WxMiniAppId    string `json:"wx_mini_app_id"`
	WxMinSecret    string `json:"wx_min_secret"`
	WxMiniAesKey   string `json:"wx_mini_aes_key"`
	WxPaymentAppId string `json:"wx_payment_app_id"`
	WxPaymentMchId string `json:"wx_payment_mch_id"`
	WxPaymentKey   string `json:"wx_payment_key"`
	XcxCover       string `json:"xcx_cover"`
	WechatCover    string `json:"wechat_cover"`
}

type Nav struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Href  string `json:"href"`
	Type  string `json:"type"`
	RefId string `json:"ref_id"`
}

func (*Nav) TableName() string {
	return "navs"
}

type CompanyService interface {
	GetCompanyById(comId int) (*Company, error)
	GetCompanyList(offsetId int, limit int) ([]*Company, error)
	GetCachedCompanyById(comId int) (*Company, error)
	UpdateCachedCompany(company *Company) (err error)
	GetQuickNav(comId int) ([]*Nav, error)
}
