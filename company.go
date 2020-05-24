package echoapp

type Company struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Desc           string `json:"desc"`
	Addr           string `json:"addr"`
	WxMiniAppId    string `json:"wx_mini_app_id"`
	WxMinSecret    string `json:"wx_min_secret"`
	WxMiniAesKey   string `json:"wx_mini_aes_key"`
	WxPaymentAppId string `json:"wx_payment_app_id"`
	WxPaymentMchId string `json:"wx_payment_mch_id"`
	WxPaymentKey   string `json:"wx_payment_key"`
}

type CompanyService interface {
	GetCompanyById(comId int) (*Company, error)
	GetCompanyList(offsetId int, limit int) ([]*Company, error)
	GetCachedCompanyById(comId int) (*Company, error)
	UpdateCachedCompany(company *Company) (err error)
}
