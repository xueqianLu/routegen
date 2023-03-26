package tool

import (
	"errors"
	"sync"
)

var (
	ErrTaskPoolIsFull = errors.New("task pool is full")
)

type TaskHandle func(interface{})

type Tasks struct {
	tasknum  uint
	handler  TaskHandle
	taskpool chan interface{}
	wg       sync.WaitGroup
}

func NewTasks(routine uint, handle TaskHandle) *Tasks {
	return &Tasks{
		tasknum:  routine,
		handler:  handle,
		taskpool: make(chan interface{}, 1000000000),
	}
}

func (t *Tasks) AddTask(task interface{}) error {
	select {
	case t.taskpool <- task:
		return nil
	default:
		return ErrTaskPoolIsFull
	}
}

func (t *Tasks) Stop() {
	close(t.taskpool)
}

func (t *Tasks) Run() {
	for i := uint(0); i < t.tasknum; i++ {
		t.wg.Add(1)
		go func() {
			defer t.wg.Done()
			for {
				select {
				case task, ok := <-t.taskpool:
					if !ok {
						return
					}

					t.handler(task)
				}
			}
		}()
	}
}

func (t *Tasks) Done() {
	t.wg.Wait()
}
