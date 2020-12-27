package observer

import (
	"fmt"
	"testing"
)

type FakeObserver struct {
}

type FakeSubject struct {
	SubjectImp
}

func (f FakeObserver) Name() string {
	return "fake"
}

func (f FakeObserver) OnEvent(event interface{}) error {
	fmt.Println("FakeObserver OnEvent")
	return nil
}

func TestSubjectImp_Notify(t *testing.T) {
	sbj := FakeSubject{}
	observer := &FakeObserver{}
	sbj.AddObserver(observer)
	sbj.Notify("s123")
}
