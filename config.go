package echoapp

import (
	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
)

const (
	DefaultConfigFile = "config.yaml"
)

var ConfigOpts ConfigOptions
var Viper *viper.Viper

type ConfigOptions struct {
	Asset             Asset `yaml:"asset" mapstructure:"asset"`
	Server            *Server
	UserServer        *Server                       `yaml:"user_server" mapstructure:"user_server"`
	GoodsServer       *Server                       `yaml:"goods_server" mapstructure:"goods_server"`
	OrderServer       *Server                       `yaml:"order_server" mapstructure:"order_server"`
	CommentServer     *Server                       `yaml:"comment_server" mapstructure:"comment_server"`
	FileServer        *Server                       `yaml:"file_server" mapstructure:"file_server"`
	ResourceOptions   *ResourceServerOption         `yaml:"resource" mapstructure:"resource"`
	SmsOptionTokenMap map[string]SmsOption          `yaml:"sms_tokens" mapstructure:"sms_tokens"`
	DBMap             map[string]DBOption           `yaml:"database" mapstructure:"database"`
	RedisMap          map[string]*redis.Options     `yaml:"cache" mapstructure:"cache"`
	Redis             *redis.Options                `yaml:"redis" mapstructure:"redis"`
	MQMap             map[string]RabbitMqOption     `yaml:"rabbit_mq" mapstructure:"rabbit_mq"`
	TongchengConfig   TongchengConfig               `yaml:"tongcheng" mapstructure:"tongcheng"`
	ReportTicketMap   map[string]ReportTicketOption `yaml:"report_tickets" mapstructure:"report_tickets"`
	Jws               JwsHelperOpt                  `yaml:"jws" mapstructure:"jws"`
	RecommendOptions  *RecommendOptions             `yaml:"recommend_options" mapstructure:"recommend_options"`
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
	//接受者需要提供公钥(不需要提供公钥)
	PublicKeyPath string `json:"public_key_path" yaml:"public_key_path"  mapstructure:"public_key_path"`
	//签发者需要知道私钥
	PrivateKeyPath string `json:"private_key_path" yaml:"private_key_path" mapstructure:"private_key_path"`
	//配置后使用hashIds　混淆UserId
	HashIdsSalt string `json:"hash_ids_salt" yaml:"hash_ids_salt" mapstructure:"hash_ids_salt"`
}

type Asset struct {
	PublicRoot   string `yaml:"public_root" mapstructure:"public_root"`
	ResourceRoot string `yaml:"resource_root" mapstructure:"resource_root"`
	StorageRoot  string `yaml:"storage_root" mapstructure:"storage_root"`
	AreaRoot     string `yaml:"area_root" mapstructure:"area_root"`
	ViewRoot     string `yaml:"view_root" mapstructure:"view_root"`
	Version      string `yaml:"version" mapstructure:"version"`
	PublicHost   string `yaml:"public_host" mapstructure:"public_host"`
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

func InitConfig(cfgFile string) {
	if cfgFile == "" {
		cfgFile = DefaultConfigFile
	}

	Viper = viper.New()
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&ConfigOpts); err != nil {
		panic(err)
	}
}
