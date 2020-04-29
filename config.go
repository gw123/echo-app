package echoapp

import (
	"fmt"

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
	ReportTicketMap   map[string]ReportTicketOption `yaml:"report_ticket_map" mapstructure:"report_ticket_map"`
}

type Server struct {
	Addr      string   `yaml:"addr" mapstructure:"addr"`
	Origins   []string `yaml:"origins" mapstructure:"origins"`
	AppMode   string   `yaml:"app_mode" mapstructure:"app_mode"`
	JwtPubkey string   `yaml:"jwt_pubkey" mapstructure:"jwt_pubkey"`
}

type Asset struct {
	PublicRoot string `yaml:"public_root" mapstructure:"public_root"`
	AreaRoot   string `yaml:"area_root" mapstructure:"area_root"`
	ViewRoot   string `yaml:"view_root" mapstructure:"view_root"`
	Version    string `yaml:"version" mapstructure:"version"`
	PublicHost string `yaml:"public_host" mapstructure:"public_host"`
	WatchRoot  string `yanml:"watch_root"  mapstructure:"watch_root"`
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
	Addr     string `yaml:"addr" mapstructure:"addr"`
	ComId    int    `yaml:"com_id" mapstructure:"com_id"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
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
	fmt.Println("ConfigOpts.Asset:", ConfigOpts.Asset)
	fmt.Println("ConfigOpts.ReportTicketMa:", ConfigOpts.ReportTicketMap)
}
