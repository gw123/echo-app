package echoapp

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	DefaultConfigFile = "config.yaml"
)

var Config ConfigOptions
var Viper *viper.Viper

type ConfigOptions struct {
	Asset  Asset `yaml:"asset" mapstructure:"asset"`
	Server *Server
	Redis  *CacheOptions
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
	fmt.Println("Config.Asset:", Config.Asset)
}
