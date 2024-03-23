package app_components

import (
	"sync"

	"github.com/gw123/echo-app/external/sms_tpls"
	"github.com/gw123/glog"
)

var smsTplApi sms_tpls.SmsTplAPi
var smsTplApiOnce sync.Once

func GetSmsTplApi() (sms_tpls.SmsTplAPi, error) {
	var err2 error
	smsTplApiOnce.Do(func() {
		JobPusher, err := GetJobPusher()
		if err != nil {
			glog.DefaultLogger().Errorf("GetSmsTplApi : %s", err.Error())
			err2 = err
			return
		}

		smsSvr, err := GetSmsService()
		if err != nil {
			err2 = err
			glog.DefaultLogger().Errorf("GetSmsTplApi : %s", err.Error())
			return
		}
		smsTplApi = sms_tpls.NewSmsTplApi(JobPusher, smsSvr)
	})

	return smsTplApi, err2
}
