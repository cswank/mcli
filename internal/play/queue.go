package play

import (
	"sync"

	"bitbucket.org/cswank/mcli/internal/schema"
)

type queue struct {
	queue []schema.Result
	lock  sync.Mutex
	ready chan bool
}

func newQueue() *queue {
	return &queue{
		ready: make(chan bool),
	}
}

func (q *queue) clear() {
	q.lock.Lock()
	q.queue = []schema.Result{}
	q.lock.Unlock()
}

func (q *queue) Add(r schema.Result) {
	q.lock.Lock()
	q.queue = append(q.queue, r)
	q.lock.Unlock()
	go func() {
		q.ready <- true
	}()
}

func (q *queue) Remove(i int) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.queue) == 0 || i >= len(q.queue) {
		return
	}
	q.queue = append(q.queue[:i], q.queue[i+1:]...)
}

func (q *queue) Playlist() []schema.Result {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.queue
}

func (q *queue) Next() schema.Result {
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
