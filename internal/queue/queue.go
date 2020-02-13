package queue

import (
	"sync"

	"github.com/cswank/mcli/internal/schema"
)

type Queue struct {
	queue []schema.Result
	lock  sync.Mutex
	ready chan bool
}

func New() *Queue {
	return &Queue{
		ready: make(chan bool),
	}
}

func (q *Queue) Add(r schema.Result) {
	q.lock.Lock()
	q.queue = append(q.queue, r)
	q.lock.Unlock()
	go func() {
		q.ready <- true
	}()
}

func (q *Queue) Remove(i int) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.queue) == 0 || i >= len(q.queue) {
		return
	}
	q.queue = append(q.queue[:i], q.queue[i+1:]...)
}

func (q *Queue) Queue() []schema.Result {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.queue
}

func (q *Queue) Next() schema.Result {
	<-q.ready
	if len(q.queue) == 0 {
		return q.Next()
	}

	q.lock.Lock()
	r := q.queue[0]
	q.queue = q.queue[1:]
	q.lock.Unlock()
	return r
}
