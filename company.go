package echoapp

type CompanyBrief struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Addr      string `json:"addr"`
	AddrBrief Addr   `json:"address" gorm:"-"`
	Phone     string `json:"phone"`
}

type Addr struct {
	CityId     int64 `json:"city_id"`
	DistrictId int64 `json:"district_id"`
	ProvinceId int64 `json:"province_id"`
}

type Company struct {
	CompanyBrief
	WxMiniAppId  string `json:"wx_mini_app_id"`
	WxMinSecret  string `json:"wx_min_secret"`
	WxMiniAesKey string `json:"wx_mini_aes_key"`

	WxOfficialAppId  string `json:"wx_official_app_id"`
	WxOfficialSecret string `json:"wx_official_secret"`

	WxPaymentAppId   string `json:"wx_payment_app_id"`
	WxPaymentMchId   string `json:"wx_payment_mch_id"`
	WxPaymentKey     string `json:"wx_payment_key"`
	XcxCover         string `json:"xcx_cover"`
	WechatCover      string `json:"wechat_cover"`
	WxToken          string `json:"wx_token"`
	WxOfficialAesKey string `json:"wx_official_aes_key"`

	SmsChannels map[string]*SmsChannel `json:"sms_channels"`
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
	GetCompanyById(comId uint) (*Company, error)
	GetCompanyList(offsetId uint, limit int) ([]*Company, error)
	GetCachedCompanyById(comId uint) (*Company, error)
	UpdateCachedCompany(company *Company) (err error)
	GetQuickNav(comId uint) ([]*Nav, error)
}
