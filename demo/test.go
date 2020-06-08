package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Pool struct {
	queue chan int
	wg    *sync.WaitGroup
}

func NewPoll(size int) *Pool {
	if size <= 0 {
		size = 1
	}
	return &Pool{
		queue: make(chan int, size),
		wg:    &sync.WaitGroup{},
	}
}
func (p *Pool) Add(delta int) {
	for i := 0; i < delta; i++ {
		p.queue <- i
	}
	for i := 0; i > delta; i-- {
		<-p.queue
	}
	p.wg.Add(delta)
}
func (p *Pool) Done() {
	<-p.queue
	p.wg.Done()
}
func (p *Pool) Wait() {
	p.wg.Wait()
}
func main() {
	pool := NewPoll(5)
	fmt.Println("the NumGoroutine begin is:", runtime.NumGoroutine())
	for i := 1; i <= 20; i++ {
		pool.Add(1)
		go func(i int) {
			time.Sleep(1 * time.Second)
			fmt.Println("the numGoroutien continue is:", runtime.NumGoroutine(), i)
			pool.Done()
		}(i)
		// if i%5 == 0 {
		// 	time.Sleep(time.Second)
		// }

	}
	pool.Wait()
	fmt.Println("the NumGoroutine done is :", runtime.NumGoroutine())
}
