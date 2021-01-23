package echoapp

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/gw123/echo-app/libs/etcd"

	"github.com/spf13/viper"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/spf13/viper/remote"
)

func GetEtcd(t *testing.T) *etcd.EtcdConfig {
	etcdCli, err := etcd.NewEtcdConfig(etcd.EtcdOptions{
		Endpoints: []string{"http://sh2:22379"},
		Namespace: "/xyt",
		Username:  "root",
		Password:  "123456",
	})

	if err != nil {
		t.Fatal(err)
		return nil
	}
	return etcdCli
}

func Test_Config(t *testing.T) {
	var ConfigOpts ConfigOptions
	etcdCli := GetEtcd(t)
	cf, err := etcdCli.Get("/config.prod.yaml")
	if err != nil {
		t.Error(err)
		return
	}
	// 要先去设置文件类型,再去ReadConfig
	viper.SetConfigType("yaml") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
	viper.ReadConfig(bytes.NewBuffer([]byte(cf)))
	if err != nil {
		t.Error(err)
		return
	}

	// unmarshal config
	err = viper.Unmarshal(&ConfigOpts)
	if err != nil {
		t.Error(err)
		return
	}
	spew.Dump(ConfigOpts)
}

func TestIntoStruct(t *testing.T) {
	eletype := reflect.ValueOf(ConfigOpts)
	t.Log(eletype.String())

	t.Log(eletype.Field(0).Type().String())
	t.Log(eletype.Type().Field(0).Tag.Get("json"))
}
