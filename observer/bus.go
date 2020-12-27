package observer

import (
	"sync"
)

type Bus struct {
	Observers []Observer
	Subject   []Subject
}

type Observer interface {
	// 观察者名字
	Name() string
	//
	OnEvent(event interface{}) error
}

type Subject interface {
	// 添加观察者
	AddObserver(observer Observer)
	// 删除观察者
	RemoveObserver(observer Observer)
	// 通知所有的观察者
	Notify(interface{})
}

// 订阅者中主题的实现
type SubjectImp struct {
	observer sync.Map
}

func NewSubjectImp() *SubjectImp {
	return &SubjectImp{
		observer: sync.Map{},
	}
}

func (s *SubjectImp) AddObserver(observer Observer) {
	s.observer.Store(observer.Name(), observer)
}

func (s *SubjectImp) RemoveObserver(observer Observer) {
	s.observer.Delete(observer.Name())
}

func (s *SubjectImp) Notify(event interface{}) {
	s.observer.Range(func(key, value interface{}) bool {
		observer, ok := value.(Observer)
		if ok {
			if err := observer.OnEvent(event); err != nil {

			}
		}
		return true
	})
}
