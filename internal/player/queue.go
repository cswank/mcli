package player

import (
	"sync"
)

type Progress struct {
	N     int
	Total int
}

type queue struct {
	queue []Result
	lock  sync.Mutex
	ready chan bool
}

func newQueue() *queue {
	return &queue{
		ready: make(chan bool),
	}
}

func (q *queue) Add(r Result) {
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

func (q *queue) Playlist() []Result {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.queue
}

func (q *queue) Next() Result {
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
