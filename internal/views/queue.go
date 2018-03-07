package views

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"bitbucket.org/cswank/mcli/internal/source"
)

type queue struct {
	source     source.Source
	sourceName string
	in         []source.Result
	queue      []source.Result
	out        chan source.Result
	buf        chan<- progress
	lock       sync.Mutex
	outLock    sync.Mutex
	sep        string
}

func newQueue(s source.Source, buf chan<- progress) *queue {
	q := &queue{
		source:     s,
		buf:        buf,
		sourceName: s.Name(),
		sep:        string(filepath.Separator),
		out:        make(chan source.Result),
	}

	go q.download()
	return q
}

func (q *queue) add(r source.Result) {
	q.lock.Lock()
	q.in = append(q.in, r)
	q.queue = append(q.queue, r)
	q.lock.Unlock()
}

func (q *queue) remove(i int) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.queue) == 0 || i >= len(q.queue) {
		return
	}
	q.queue = append(q.queue[:i], q.queue[i+1:]...)
}

func (q *queue) playlist() []source.Result {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.queue
}

func (q *queue) next() source.Result {
	r := <-q.out
	q.lock.Lock()
	if len(q.queue) > 0 {
		q.queue = q.queue[1:]
	}
	q.lock.Unlock()
	return r
}

func (q *queue) download() {
	for {
		time.Sleep(time.Second)
		q.lock.Lock()
		if len(q.in) == 0 {
			q.lock.Unlock()
			continue
		}

		r := q.in[len(q.in)-1]
		q.in = q.in[:len(q.in)-1]
		q.lock.Unlock()
		pth, e := q.checkCache(r)
		r.Path = pth
		if !e {
			err := q.doDownload(r)
			if err != nil {
				log.Println(err)
				continue
			}
		}
		q.out <- r
	}
}

func (q *queue) doDownload(r source.Result) error {
	f, err := os.Create(r.Path)
	if err != nil {
		return fmt.Errorf("could not create file for %+v: %s", r, err)
	}

	u, err := q.source.GetTrack(r.Track.ID)
	if err != nil {
		return fmt.Errorf("could not get track %+v: %s", r, err)
	}

	resp, err := http.Get(u)
	if err != nil {
		return fmt.Errorf("could not get stream %+v: %s", r, err)
	}

	pr := newProgressRead(resp.Body, int(resp.ContentLength), q.buf)
	_, err = io.Copy(f, pr)
	if err != nil {
		return err
	}

	f.Close()
	pr.Close()
	return nil
}

func (q *queue) checkCache(result source.Result) (string, bool) {
	dir := fmt.Sprintf("%s/cache/%s/%s/%s", os.Getenv("MCLI_HOME"), q.sourceName, q.clean(result.Artist.Name), q.clean(result.Album.Title))
	e, _ := exists(dir)
	if !e {
		os.MkdirAll(dir, 0700)
	}

	pth := fmt.Sprintf("%s/cache/%s/%s/%s/%s.flac", os.Getenv("MCLI_HOME"), q.sourceName, q.clean(result.Artist.Name), q.clean(result.Album.Title), q.clean(result.Track.Title))
	e, _ = exists(pth)
	return pth, e
}

func (q *queue) clean(s string) string {
	return strings.Replace(s, q.sep, "", -1)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

type progressRead struct {
	io.Reader
	t, l, reads int
	ch          chan<- progress
}

func newProgressRead(r io.Reader, l int, ch chan<- progress) *progressRead {
	return &progressRead{Reader: r, t: 0, l: l, ch: ch}
}

func (r *progressRead) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.t += n
	r.reads++
	if r.reads%100 == 0 {
		r.ch <- progress{n: r.t, total: r.l}
	}
	return n, err
}

// Close the reader when it implements io.Closer
func (r *progressRead) Close() error {
	r.ch <- progress{n: 0, total: r.t}
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
