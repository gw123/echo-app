package echoapp

const (
	TimeFormat   = "2006-01-02 15:04:05"
	DateFormat   = "2006-01-02"
	HostURL      = "host_url"
	RequestRoot  = "static_root"
	AssetVersion = "asset_version"

	//缓存key
	//redis 相关的key
	RedisUserKey          = "User:%d"
	RedisUserXCXOpenidKey = "UserXCXOpenid:%d"
	RedisSmsLoginCodeKey  = "SmsLoginCode"
	RedisCompanyKey       = "Company:%d"
	RedisGoodsKey         = "Goods:%d"
	//Goods:comId:position
	RedisBannerListKey = "Goods:%d:%s"

	RedisActivityKey = "Activity:%d"
	RedisArticleKey  = "Article:%d"

	RedisUserDefaultAddrKey = "UserDefaultAddr:%d"
	RedisUserCollectionKey  = "UserCollection:%d"

)
