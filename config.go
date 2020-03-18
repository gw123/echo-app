package echoapp

import (
	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

const (
	DefaultConfigFile = "config.yaml"
)

var Config ConfigOptions
var Viper *viper.Viper

type ConfigOptions struct {
	Asset    Asset `yaml:"asset" mapstructure:"asset"`
	Env      *Environments
	Database *DatabaseOptions
	Redis    *CacheOptions
	Origins  []string
}

type Environments struct {
	Addr      string `yaml:"addr" mapstructure:"addr"`
	AppMode   string `yaml:"app_mode" mapstructure:"app_mode"`
	JwtPubkey string `yaml:"jwt_pubkey" mapstructure:"jwt_pubkey"`
}

type Asset struct {
	PublicRoot string `yaml:"public_root" mapstructure:"public_root"`
	ViewRoot   string `yaml:"view_root" mapstructure:"view_root"`
	Version    string `yaml:"version" mapstructure:"version"`
	PublicHost string `yaml:"public_host" mapstructure:"public_host"`
}

// staticFileURL 拼接静态文件路径
func StaticFileURL(tp, uri string) string {
	host := Config.Asset.PublicHost
	version := Config.Asset.Version

	if Config.Env.AppMode == "development" || Config.Env.AppMode == "dev" {
		if tp == "" {
			return host + "/" + uri + "?version=" + version
		} else {
			return host + "/" + tp + "/" + uri + "?version=" + version
		}
	}
	return host + "/mobile/" + tp + "/" + uri
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

	if err := viper.Unmarshal(&Config); err != nil {
		panic(err)
	}
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}


