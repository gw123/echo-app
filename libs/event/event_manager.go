package event

import (
	"errors"

	"github.com/gw123/gworker"
)

type EventHandle func(job gworker.Job) error

type EventManager struct {
	eventListeners map[string][]EventHandle
}

func NewEventManager() *EventManager {
	return &EventManager{eventListeners: map[string][]EventHandle{}}
}

func (e *EventManager) Register(eventName string, handle EventHandle) {
	if _, ok := e.eventListeners[eventName]; !ok {
		e.eventListeners[eventName] = make([]EventHandle, 0)
	}
	e.eventListeners[eventName] = append(e.eventListeners[eventName], handle)
}

func (e *EventManager) PushSyncJob(job gworker.Job) (errs []error) {
	if _, ok := e.eventListeners[job.GetName()]; ok {
		for _, handel := range e.eventListeners[job.GetName()] {
			if err := handel(job); err != nil {
				errs = append(errs, err)
			}
		}
	} else {
		errs = append(errs, errors.New("未注册eventHandle "+job.GetName()))
	}
	return errs
}
