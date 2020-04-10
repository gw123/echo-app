package app

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/services"
	"github.com/pkg/errors"
)

var App *EchoApp

type EchoApp struct {
	areaSvc echoapp.AreaService
}

func init() {
	App = &EchoApp{}
	//做一些初始化工作
}

func GetAreaService() (echoapp.AreaService, error) {
	if App.areaSvc != nil {
		return App.areaSvc, nil
	}
	areaSvc, err := services.NewAreaService(echoapp.Config.Asset.AreaRoot)
	if err != nil {
		return nil, errors.Wrap(err, "GetAreaService")
	}
	App.areaSvc = areaSvc
	return areaSvc, nil
}

func MustGetAreaService() echoapp.AreaService {
	areaSvc, err := GetAreaService()
	if err != nil {
		panic(err)
	}
	return areaSvc
}
