package main

import (
	"fmt"
	"time"
)

func example() {
	{
		a := make(chan int)
		close(a)
		<-a // 0, false
	}
	{
		a := make(chan int)
		close(a)
		a <- 10 // panic
	}
	{
		var a chan int // nil
		<-a            // вечная блокировка
	}
	{
		var a chan int // nil
		a <- 0         // вечная блокировка
	}
}

func Union(aCh, bCh chan int) chan int {
	cCh := make(chan int)
	go func() {
		for {
			select {
			case a, aOk := <-aCh:
				if aOk {
					cCh <- a
				} else {
					aCh = nil
				}
			case b, bOk := <-bCh:
				if bOk {
					cCh <- b
				} else {
					bCh = nil
				}
			}
			if aCh == nil && bCh == nil {
				close(cCh)
				return
			}
		}
	}()
	return cCh
}

func main() {
	a := make(chan int, 3)
	b := make(chan int, 3)
	a <- 1
	a <- 2
	a <- 3
	go func() {
		for i := 4; i < 14; i++ {
			time.Sleep(time.Second)
			b <- i
		}
		close(b)
	}()
	close(a)

	cCh := Union(a, b)

	for c := range cCh {
		fmt.Println(c)
	}
}
