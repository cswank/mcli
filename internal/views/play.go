package views

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"bitbucket.org/cswank/music/internal/source"
	"github.com/faiface/beep"
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
	r.ch <- progress{n: r.t, total: r.t}
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

type play struct {
	coords   coords
	progress chan<- progress
	ch       chan playlist
	cancel   chan bool
	source   source.Source
	pause    chan bool
}

func newPlay(w, h int, pr chan<- progress) *play {
	p := &play{
		coords:   coords{x1: -1, y1: h - 3, x2: w, y2: h - 1},
		progress: pr,
		ch:       make(chan playlist),
		cancel:   make(chan bool),
		pause:    make(chan bool),
	}

	go p.play(p.ch, p.cancel)
	return p
}

func (p *play) doPause() {
	log.Println("pause")
	p.pause <- true
}

func (p *play) play(ch <-chan playlist, cancel <-chan bool) error {
	for {
		select {
		case pl := <-ch:
			p.doPlay(pl.tracks[0])
		case <-cancel:
		}
	}
}

func (p *play) clear() {
	v, _ := g.View("play")
	v.Clear()
}

func (p *play) render(g *ui.Gui, v *ui.View) {

}

func (p *play) doPlay(result source.Result) error {
	in, f, err := p.getFile(result)
	if err != nil {
		return err
	}

	if f == nil {
		u, err := p.source.GetTrack(result.ID)
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
		f, err = os.Open(in.Name())
		if err != nil {
			return err
		}
	}

	s, format, err := flac.Decode(f)
	if err != nil {
		return err
	}

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second)); err != nil {
		return err
	}

	done := make(chan struct{})
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(done)
	})))

	start := int(time.Now().Unix())
	var locked bool
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			n := int(time.Now().Unix())
			p.progress <- progress{n: n - start, total: result.Duration}
		case <-done:
			return s.Close()
		case <-p.pause:
			if locked {
				locked = false
				speaker.Unlock()
			} else {
				locked = true
				speaker.Lock()
			}
		}
	}
}

func (p *play) getFile(result source.Result) (*os.File, *os.File, error) {
	dir := fmt.Sprintf("%s/.music/cache/%s/%s/%s", os.Getenv("HOME"), p.source.Name(), result.Artist, result.Album)
	e, err := exists(dir)
	if err != nil {
		return nil, nil, err
	}

	if !e {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return nil, nil, err
		}
	}

	pth := fmt.Sprintf("%s/.music/cache/%s/%s/%s/%s", os.Getenv("HOME"), p.source.Name(), result.Artist, result.Album, result.ID)
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
