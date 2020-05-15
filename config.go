package echoapp

import (
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
	Redis             *CacheOptions
	SmsOptionTokenMap map[string]SmsOption          `yaml:"sms_tokens" mapstructure:"sms_tokens"`
	DBMap             map[string]DBOption           `yaml:"database" mapstructure:"database"`
	MQMap             map[string]RabbitMqOption     `yaml:"rabbit_mq" mapstructure:"rabbit_mq"`
	TongchengConfig   TongchengConfig               `yaml:"tongcheng" mapstructure:"tongcheng"`
	ReportTicketMap   map[string]ReportTicketOption `yaml:"report_tickets" mapstructure:"report_tickets"`
	PPTImages         map[string]PPTImagesOption    `yaml:"get_ppt_images" mapstructure:"get_ppt_images"`

	QiniuKeys QiniuKeyOption `yaml:"qiniu_key" mapstructure:"qiniu_key"`
}

type Server struct {
	Addr      string   `yaml:"addr" mapstructure:"addr"`
	Origins   []string `yaml:"origins" mapstructure:"origins"`
	AppMode   string   `yaml:"app_mode" mapstructure:"app_mode"`
	JwtPubkey string   `yaml:"jwt_pubkey" mapstructure:"jwt_pubkey"`
}

type Asset struct {
	PublicRoot        string `yaml:"public_root" mapstructure:"public_root"`
	AreaRoot          string `yaml:"area_root" mapstructure:"area_root"`
	ViewRoot          string `yaml:"view_root" mapstructure:"view_root"`
	Version           string `yaml:"version" mapstructure:"version"`
	PublicHost        string `yaml:"public_host" mapstructure:"public_host"`
	WatchRoot         string `yaml:"watch_root" mapstructure:"watch_root"`
	TmpRoot           string `yaml:"tmp_root" mapstructure:"tmp_root"`
	UploadMaxFileSize int64  `yaml:"max_file_size" mapstructure:"max_file_size"`
	MyURL             string `yaml:"my_url" mapstructure:"my_url"`
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
}

type PPTImagesOption struct {
	ComId int `yaml:"com_id" mapstructure:"com_id"`

	BaseUrl string `yaml:"base_url" mapstructure:"base_url"`
}
type DBOption struct {
	Driver    string `yaml:"dirver" mapstructure:"driver"`
	DSN       string `yaml:"dsn" mapstructure:"dsn"`
	KeepAlive int    `yaml:"keep_alive" mapstructure:"keep_alive"`
	MaxOpens  int    `yaml:"max_opens" mapstructure:"max_opens"`
	MaxIdles  int    `yaml:"max_idles" mapstructure:"max_idles"`
}
type QiniuKeyOption struct {
	AccessKey string `yaml:"access_key" mapstructure:"access_key"`
	SecretKey string `yaml:"secret_key" mapstructure:"secret_key"`
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
