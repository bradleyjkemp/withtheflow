package taskqueue

import (
	"context"
	"sync"
)

const bufferSize = 10

type IdQueue interface {
	AddTask(id int64)
	GetTask() int64
	GetTaskChan() chan struct{}
}

type idQueue struct {
	addTaskChan chan struct{}
	getTaskChan chan struct{}
	stack       []int64
	mutex       sync.Mutex
}

func CreateIdQueue(ctx context.Context) IdQueue {
	queue := &idQueue{
		addTaskChan: make(chan struct{}, bufferSize),
		getTaskChan: make(chan struct{}, bufferSize),
	}
	queue.setupInfiniteChannel(ctx)
	return queue
}

func (q *idQueue) setupInfiniteChannel(ctx context.Context) {
	var slots int
	in := q.addTaskChan
	out := q.getTaskChan

	go func() {
		for {
			if slots == 0 {
				out = nil
			} else {
				out = q.getTaskChan
			}

			select {
			case <-ctx.Done():
				return
			case <-in:
				slots++
			case out <- struct{}{}:
				slots--
			}
		}
	}()
}

func (q *idQueue) AddTask(id int64) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.addTaskChan <- struct{}{}
	q.stack = append(q.stack, id)
}

func (q *idQueue) GetTask() int64 {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	stackSize := len(q.stack)
	id := q.stack[stackSize-1]
	q.stack = q.stack[:stackSize-1]

	return id
}

func (q *idQueue) GetTaskChan() chan struct{} {
	return q.getTaskChan
}
