package views

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"bitbucket.org/cswank/music/internal/history"
	"bitbucket.org/cswank/music/internal/source"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/speaker"
	ui "github.com/jroimartin/gocui"
)

type progress struct {
	n     int
	total int
}

type playlist struct {
	tracks []source.Result
}

type play struct {
	width        int
	coords       coords
	progress     chan<- progress
	playProgress chan progress
	ch           chan playlist
	source       source.Source
	pause        chan bool
	vol          chan float64
	history      history.History
	queue        *queue
}

func newPlay(w, h int, pr chan<- progress) (*play, error) {
	hist, err := history.NewFileHistory()
	if err != nil {
		return nil, err
	}
	p := &play{
		width:        w,
		queue:        newQueue(),
		coords:       coords{x1: -1, y1: h - 2, x2: w, y2: h},
		progress:     pr,
		playProgress: make(chan progress),
		ch:           make(chan playlist),
		pause:        make(chan bool),
		history:      hist,
		vol:          make(chan float64),
	}

	go p.play(p.ch)
	go p.playNext()
	go p.render()
	return p, nil
}

func (p *play) doPause() {
	p.pause <- true
}

func (p *play) volume(v float64) {
	p.vol <- v
}

func (p *play) play(ch <-chan playlist) {
	for {
		pl := <-ch
		p.queue.Put(pl.tracks[0])
	}
}

func (p *play) clear() {
	v, _ := g.View("play")
	v.Clear()
}

func (p *play) render() {
	var v *ui.View
	for {
		prog := <-p.playProgress
		g.Update(func(g *ui.Gui) error {
			if v == nil {
				v, _ = g.View("play")
			}
			v.Clear()
			fmt.Fprint(v, fmt.Sprintf(strings.Repeat("|", p.width*prog.n/prog.total)))
			return nil
		})
	}
}

func (p *play) playNext() {
	for {
		log.Println("wait for queue")
		result := p.queue.Pop()
		log.Println("got from queue", result)
		if err := p.doPlay(result); err != nil {
			log.Printf("couldn't play %v: %s", result, err)
		}
		log.Println("done playing", result)
	}
}

func (p *play) doPlay(result source.Result) error {
	if err := p.history.Save(result); err != nil {
		return err
	}

	in, f, err := p.getFile(result)
	if err != nil {
		return err
	}

	if f == nil {
		u, err := p.source.GetTrack(result.Track.ID)
		if err != nil {
			return err
		}
		resp, err := http.Get(u)
		if err != nil {
			return err
		}
		r := newProgressRead(resp.Body, int(resp.ContentLength), p.progress)

		_, err = io.Copy(in, r)
		if err != nil {
			return err
		}

		in.Close()
		r.Close()
		f, err = os.Open(in.Name())
		if err != nil {
			return err
		}
	}

	s, format, err := flac.Decode(f)
	if err != nil {
		return err
	}

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/2)); err != nil {
		return err
	}

	done := make(chan struct{})
	vol := &effects.Volume{
		Streamer: s,
		Base:     2,
	}

	ctrl := &beep.Ctrl{
		Streamer: vol,
	}
	speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
		close(done)
	})))

	start := int(time.Now().Unix())
	var paused bool
	var pauseTime time.Duration
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			n := int(time.Now().Unix())
			if !paused {
				p.playProgress <- progress{n: n - start, total: result.Track.Duration}
			} else {
				pauseTime += 100 * time.Millisecond
			}
		case <-done:
			return s.Close()
		case v := <-p.vol:
			speaker.Lock()
			vol.Volume += v
			speaker.Unlock()
		case <-p.pause:
			paused = !paused
			speaker.Lock()
			ctrl.Paused = paused
			speaker.Unlock()
		}
	}
}

func (p *play) getFile(result source.Result) (*os.File, *os.File, error) {
	dir := fmt.Sprintf("%s/.music/cache/%s/%s/%s", os.Getenv("HOME"), p.source.Name(), result.Artist.Name, result.Album.Title)
	e, err := exists(dir)
	if err != nil {
		return nil, nil, err
	}

	if !e {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return nil, nil, err
		}
	}

	pth := fmt.Sprintf("%s/.music/cache/%s/%s/%s/%s.flac", os.Getenv("HOME"), p.source.Name(), result.Artist.Name, result.Album.Title, result.Track.Title)
	e, err = exists(pth)
	if err != nil {
		return nil, nil, err
	}
	if e {
		f, err := os.Open(pth)
		return nil, f, err
	}

	f, err := os.Create(pth)
	return f, nil, err
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

// queue is a FIFO queue where Pop() operation is blocking if no items exists
type queue struct {
	lock       sync.Mutex
	notifyLock sync.Mutex
	monitor    *sync.Cond
	queue      []source.Result
}

func newQueue() *queue {
	bq := &queue{}
	bq.monitor = sync.NewCond(&bq.notifyLock)
	return bq
}

// Put any value to queue back. Returns false if queue closed
func (bq *queue) Put(value source.Result) bool {
	bq.lock.Lock()
	bq.queue = append(bq.queue, value)
	bq.lock.Unlock()

	bq.notifyLock.Lock()
	bq.monitor.Signal()
	bq.notifyLock.Unlock()
	return true
}

// Pop front value from queue. Returns nil and false if queue closed
func (bq *queue) Pop() source.Result {
	for {
		bq.notifyLock.Lock()
		bq.monitor.Wait()
		val := bq.getUnblock()
		bq.notifyLock.Unlock()
		return val
	}
}

func (bq *queue) getUnblock() source.Result {
	bq.lock.Lock()
	defer bq.lock.Unlock()
	elem := bq.queue[len(bq.queue)-1]
	bq.queue = bq.queue[0 : len(bq.queue)-1]
	return elem
}
