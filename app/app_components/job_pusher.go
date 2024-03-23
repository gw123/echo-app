package app_components

import (
	"sync"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/gworker"
)

var jobPusher gworker.Producer
var jobPusherOnce sync.Once

func GetJobPusher() (gworker.Producer, error) {
	var err error
	jobPusherOnce.Do(func() {
		opt := echoapp.ConfigOpts.Job
		jobPusher, err = gworker.NewPorducerManager(opt)
	})

	return jobPusher, err
}
