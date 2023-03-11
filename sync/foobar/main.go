package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type FooBar struct {
	n     int
	fooCh chan struct{}
	barCh chan struct{}
}

func NewFooBar(n int) *FooBar {
	fb := &FooBar{
		n:     n,
		fooCh: make(chan struct{}, 1),
		barCh: make(chan struct{}, 1),
	}
	fb.fooCh <- struct{}{}
	return fb
}

func (fb *FooBar) Foo() {
	for i := 0; i < fb.n; i++ {
		<-fb.fooCh
		fmt.Printf("foo")
		fb.barCh <- struct{}{}
	}
}
func (fb *FooBar) Bar() {
	for i := 0; i < fb.n; i++ {
		<-fb.barCh
		fmt.Printf("bar\n")
		fb.fooCh <- struct{}{}
	}
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	fb := NewFooBar(5)

	go func() {
		time.Sleep(time.Duration(rand.Intn(2)) * 50 * time.Millisecond)
		fb.Bar()
		wg.Done()
	}()
	go func() {
		time.Sleep(time.Duration(rand.Intn(3)) * 50 * time.Millisecond)
		fb.Foo()
		wg.Done()
	}()

	wg.Wait()
}
