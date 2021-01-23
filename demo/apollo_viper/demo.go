package main

import (
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"

	echoapp "github.com/gw123/echo-app"

	"github.com/shima-park/agollo"

	remote "github.com/shima-park/agollo/viper-remote"
	"github.com/spf13/viper"
)

type Config struct {
	AppSalt string         `mapstructure:"appsalt"`
	DB      DatabaseConfig `mapstructure:"database"`
}

type DatabaseConfig struct {
	Driver  string        `mapstructure:"driver"`
	Host    string        `mapstructure:"host"`
	Port    int           `mapstructure:"port"`
	Timeout time.Duration `mapstructure:"timeout"`
	// ...
}

func main() {
	//agollo.Init("sh2:18080", "echoapp-dev")
	remote.SetAppID("echoapp-dev")
	remote.SetAgolloOptions(
		agollo.WithLogger(agollo.NewLogger(agollo.LoggerWriter(os.Stdout))),
		agollo.AccessKey("a150a9e914644e608e78123fe7a7e960"))

	v := viper.New()
	v.SetConfigType("yaml") // 根据namespace实际格式设置对应type
	err := v.AddRemoteProvider("apollo", "sh2:18080", "echoapp")
	if err != nil {
		fmt.Println(err)
		return
	} // error handle...
	fmt.Println("app.AllSettings:", v.AllSettings())

	fmt.Println(viper.SupportedRemoteProviders)

	err = v.ReadRemoteConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	// error handle...

	// 直接反序列化到结构体中
	var conf echoapp.ConfigOptions
	err = v.Unmarshal(&conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	spew.Dump(conf)

	// 各种基础类型配置项读取
	fmt.Println("Host:", v.GetString("db.host"))
	fmt.Println("Port:", v.GetInt("db.port"))
	fmt.Println("Timeout:", v.GetDuration("db.timeout"))

	// 获取所有key，所有配置
	fmt.Println("AllKeys", v.AllKeys(), "AllSettings", v.AllSettings())
}
