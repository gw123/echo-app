package echoapp

import (
	"strings"

	"github.com/gw123/glog"

	"github.com/go-redis/redis/v7"
	"github.com/gw123/echo-app/libs/etcd"
	"github.com/gw123/gworker"
	"github.com/spf13/viper"
)

const (
	DefaultConfigFile = "config.yaml"
)

var ConfigOpts ConfigOptions
var Viper *viper.Viper

type ConfigOptions struct {
	Server            *Server
	Asset             Asset                         `yaml:"asset" mapstructure:"asset"`
	UserServer        *Server                       `yaml:"user_server" mapstructure:"user_server"`
	GoodsServer       *Server                       `yaml:"goods_server" mapstructure:"goods_server"`
	OrderServer       *Server                       `yaml:"order_server" mapstructure:"order_server"`
	CommentServer     *Server                       `yaml:"comment_server" mapstructure:"comment_server"`
	MessageServer     *Server                       `yaml:"message_server" mapstructure:"message_server"`
	TestpaperServer   *Server                       `yaml:"testpaper_server" mapstructure:"testpaper_server"`
	FileServer        *Server                       `yaml:"file_server" mapstructure:"file_server"`
	SiteServer        *Server                       `yaml:"site_server" mapstructure:"site_server"`
	ActivityServer    *Server                       `yaml:"activity_server" mapstructure:"activity_server"`
	ResourceOptions   *ResourceServerOption         `yaml:"resource" mapstructure:"resource"`
	SmsOptionTokenMap map[string]SmsOption          `yaml:"sms_tokens" mapstructure:"sms_tokens"`
	DBMap             map[string]DBOption           `yaml:"database" mapstructure:"database"`
	RedisMap          map[string]*redis.Options     `yaml:"cache" mapstructure:"cache"`
	MQMap             map[string]RabbitMqOption     `yaml:"rabbit_mq" mapstructure:"rabbit_mq"`
	TongchengConfig   TongchengConfig               `yaml:"tongcheng" mapstructure:"tongcheng"`
	ReportTicketMap   map[string]ReportTicketOption `yaml:"report_tickets" mapstructure:"report_tickets"`
	Jws               JwsHelperOpt                  `yaml:"jws" mapstructure:"jws"`
	RecommendOptions  *RecommendOptions             `yaml:"recommend_options" mapstructure:"recommend_options"`
	ApiVersion        string                        `yaml:"api_version" mapstructure:"api_version"`
	Wechat            *Wechat                       `yaml:"wechat" mapstructure:"wechat"`
	Es                *EsOptions                    `yaml:"es" mapstructure:"es"`
	Redis             *redis.Options
	Job               *gworker.Options `yaml:job  mapstructure:"job"`
}

type EsOptions struct {
	Scheme   string   `yaml:"scheme" mapstructure:"scheme"`
	Username string   `yaml:"username" mapstructure:"username"`
	Password string   `yaml:"password" mapstructure:"password"`
	Sniff    bool     `yaml:"sniff" mapstructure:"sniff"`
	URLs     []string `yaml:"urls" mapstructure:"urls"`
}

type RecommendOptions struct {
	AttributesWight []float64 `json:"attributes_wight" yaml:"attributes_wight" mapstructure:"attributes_wight"`
	ParamA          float64   `json:"param_a" yaml:"param_a" mapstructure:"param_a"`
	ParamB          float64   `json:"param_b" yaml:"param_b" mapstructure:"param_b"`
}

type ResourceServerOption struct {
	BucketName        string `json:"bucket_name" yaml:"bucket_name" mapstructure:"bucket_name"`
	CallbackUrl       string `json:"callback_url" yaml:"callback_url" mapstructure:"callback_url"`
	AccessKey         string `json:"access_key" yaml:"access_key" mapstructure:"access_key"`
	SecretKey         string `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key"`
	XytUrl            string `yaml:"xyt_url" mapstructure:"xyt_url"`
	UploadMaxFileSize int64  `yaml:"max_file_size" mapstructure:"max_file_size"`
	BaseURL           string `yaml:"base_url" mapstructure:"base_url"`
}

type Server struct {
	Addr    string   `yaml:"addr" mapstructure:"addr"`
	Mode    string   `yaml:"mode" mapstructure:"mode"`
	Origins []string `yaml:"origins" mapstructure:"origins"`
	AppMode string   `yaml:"app_mode" mapstructure:"app_mode"`
}

type JwsHelperOpt struct {
	Audience string `json:"audience"`
	Issuer   string `json:"issuer"`
	//单位秒
	Timeout int64 `json:"timeout"`
	//接受者需要提供公钥(不需要提供私钥)
	PublicKeyPath string `json:"public_key_path" yaml:"public_key_path"  mapstructure:"public_key_path"`
	//签发者需要知道私钥
	PrivateKeyPath string `json:"private_key_path" yaml:"private_key_path" mapstructure:"private_key_path"`
	//配置后使用hashIds　混淆UserId
	HashIdsSalt string `json:"hash_ids_salt" yaml:"hash_ids_salt" mapstructure:"hash_ids_salt"`
}

type Asset struct {
	PublicRoot      string   `yaml:"public_root" mapstructure:"public_root"`
	ResourceRoot    string   `yaml:"resource_root" mapstructure:"resource_root"`
	StorageRoot     string   `yaml:"storage_root" mapstructure:"storage_root"`
	AreaRoot        string   `yaml:"area_root" mapstructure:"area_root"`
	ViewRoot        string   `yaml:"view_root" mapstructure:"view_root"`
	Version         string   `yaml:"version" mapstructure:"version"`
	PublicHost      string   `yaml:"public_host" mapstructure:"public_host"`
	InnerIpPrefix   []string `yaml:"inner_ip_prefix" mapstructure:"inner_ip_prefix"`
	PublicHostInner string   `yaml:"public_host_inner" mapstructure:"public_host_inner"`
}

type SmsOption struct {
	AccessKey    string `yaml:"access_key" mapstructure:"access_key"`
	AccessSecret string `yaml:"access_secret" mapstructure:"access_secret"`
	SignName     string `yaml:"sign_name" mapstructure:"sign_name"`
	TemplateCode string `yaml:"template_code" mapstructure:"template_code"`
}

type RabbitMqOption struct {
	Url string `yaml:"url" mapstructure:"url"`
}

type TongchengConfig struct {
	NotifyUrl string                     `yaml:"notify_url" mapstructure:"notify_url"`
	ClientMap map[string]TongchengOption `yaml:"client_map" mapstructure:"client_map"`
}

type TongchengOption struct {
	Key    string `yaml:"key" mapstructure:"key"`
	UserId string `yaml:"usr_id" mapstructure:"user_id"`
}

type ReportTicketOption struct {
	ComId      int    `yaml:"com_id" mapstructure:"com_id"`
	AppKey     string `yaml:"app_key" mapstructure:"app_key"`
	BaseUrl    string `yaml:"base_url" mapstructure:"base_url"`
	ScenicCode string `yaml:"scenic_code" mapstructure:"scenic_code"`
	//
	LoginName string `yaml:"loginName" mapstructure:"loginName"`
	Pwd       string `yaml:"pwd" mapstructure:"pwd"`
}

type DBOption struct {
	Driver    string `yaml:"dirver" mapstructure:"driver"`
	DSN       string `yaml:"dsn" mapstructure:"dsn"`
	KeepAlive int    `yaml:"keey_alive" mapstructure:"keey_alive"`
	MaxOpens  int    `yaml:"max_opens" mapstructure:"max_opens"`
	MaxIdles  int    `yaml:"max_idles" mapstructure:"max_idles"`
}

type Wechat struct {
	AuthRedirectUrl string `yaml:"auth_redirect_url" mapstructure:"auth_redirect_url"`
	MessageHost     string `yaml:"message_host" mapstructure:"message_host"`
	JsHost          string `yaml:"js_host" mapstructure:"js_host"`
}

func LoadFromFile(cfgFile string) {
	if cfgFile == "" {
		cfgFile = DefaultConfigFile
	}

	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.ReadInConfig()
	if err := viper.Unmarshal(&ConfigOpts); err != nil {
		panic(err)
	}
}

func LoadFromEtcd(endpoints []string, namespace, username, password string) {
	etcdCli, err := etcd.NewEtcdConfig(etcd.EtcdOptions{
		Endpoints: endpoints,
		Namespace: namespace,
		Username:  username,
		Password:  password,
	})

	if err != nil {
		glog.DefaultLogger().Fatal(err)
		return
	}

	cfgData, err := etcdCli.Get("/config.prod.yaml")
	if err != nil {
		glog.DefaultLogger().Fatal(err)
		return
	}

	// 要先去设置文件类型,再去ReadConfig
	viper.SetConfigType("yaml") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
	viper.ReadConfig(strings.NewReader(cfgData))
	if err != nil {
		glog.DefaultLogger().Fatal(err)
		return
	}
	// unmarshal config
	err = viper.Unmarshal(&ConfigOpts)
	if err != nil {
		glog.DefaultLogger().Fatal(err)
		return
	}
}
