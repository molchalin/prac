package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Job interface {
	Run() error
}

// 1) FIFO
// 2) Задачи выполняются последовательно
// 3) Push асинхронный
type Queue interface {
	Start() error
	Push(j Job) error
	Stop(ctx context.Context) error
}

type InMemoryQueue struct {
	mu     sync.Mutex
	jobs   chan Job
	done   chan struct{}
	closed bool
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{
		jobs: make(chan Job),
		done: make(chan struct{}),
	}
}

func (i *InMemoryQueue) Start() {
	jobDone := make(chan struct{})
	var jobs []Job
	var inWork bool
	jobsCh := i.jobs
	for len(jobs) > 0 || inWork || jobsCh != nil {
		select {
		case job, ok := <-jobsCh:
			if ok {
				jobs = append(jobs, job)
			} else {
				jobsCh = nil
			}
		case <-jobDone:
			inWork = false
		}

		if len(jobs) == 0 || inWork {
			continue
		}

		job := jobs[0]
		jobs = jobs[1:]
		inWork = true
		go func() {
			err := job.Run()
			if err != nil {
				log.Printf("job err: %v", err)
			}
			jobDone <- struct{}{}
		}()
	}
	close(i.done)
}

func (i *InMemoryQueue) Push(j Job) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.closed {
		return errors.New("closed")
	}
	i.jobs <- j
	return nil
}

func (i *InMemoryQueue) Stop(ctx context.Context) error {
	i.mu.Lock()
	if i.closed {
		return errors.New("closed")
	}
	i.closed = true
	close(i.jobs)
	i.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-i.done:
		return nil
	}
}

type funcJob func() error

func (f funcJob) Run() error {
	return f()
}

func main() {
	q := NewInMemoryQueue()
	go q.Start()
	// err := q.Push(funcJob(func() error {
	// 	time.Sleep(3 * time.Second)
	// 	fmt.Println("done")
	// 	return nil
	// }))
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	err := q.Stop(ctx)
	fmt.Println(err)
}
