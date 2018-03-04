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
	msg   string
	flash bool
}

type play struct {
	width        int
	coords       coords
	progress     chan<- progress
	playProgress chan progress
	ch           chan source.Result
	source       source.Source
	pause        chan bool
	fastForward  chan bool
	vol          chan float64
	history      history.History
	queue        chan source.Result
	done         chan bool
	playlist     []source.Result
	current      *source.Result
	lock         sync.Mutex
	sep          string
}

func newPlay(w, h int, pr chan<- progress) (*play, error) {
	hist, err := history.NewFileHistory()
	if err != nil {
		return nil, err
	}

	p := &play{
		sep:          string(filepath.Separator),
		width:        w,
		queue:        make(chan source.Result),
		coords:       coords{x1: -1, y1: h - 2, x2: w, y2: h},
		progress:     pr,
		playProgress: make(chan progress),
		ch:           make(chan source.Result),
		pause:        make(chan bool),
		fastForward:  make(chan bool),
		history:      hist,
		vol:          make(chan float64),
		done:         make(chan bool),
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

func (p *play) addAlbumToQueue(album []source.Result) {
	p.ch <- album[0]
	p.lock.Lock()
	p.playlist = album[1:]
	p.lock.Unlock()
}

func (p *play) removeFromQueue(i int) {
	if i > len(p.playlist)-1 {
		return
	}
	p.lock.Lock()
	p.playlist = append(p.playlist[:i], p.playlist[i+1:]...)
	p.lock.Unlock()
}

func (p *play) getQueue() []source.Result {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.playlist
}

func (p *play) play(ch <-chan source.Result) {
	var count int
	for {
		select {
		case r := <-ch:
			if count == 0 {
				p.queue <- r
			} else {
				p.lock.Lock()
				p.playlist = append(p.playlist, r)
				p.lock.Unlock()
			}
			count++
		case <-p.done:
			count--
			if len(p.playlist) > 0 {
				p.lock.Lock()
				r := p.playlist[0]
				p.playlist = p.playlist[1:]
				p.lock.Unlock()
				p.queue <- r
			}
		}
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
		result := <-p.queue
		if err := p.doPlay(result); err != nil {
			log.Printf("couldn't play %v: %s", result, err)
		}
		p.done <- true
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

	vol := &effects.Volume{
		Streamer: s,
		Base:     2,
	}

	ctrl := &beep.Ctrl{
		Streamer: vol,
	}
	speaker.Play(ctrl)

	song := fmt.Sprintf("%s %s", result.Track.Title, time.Duration(result.Track.Duration)*time.Second)
	msg := fmt.Sprintf(fmt.Sprintf("%%%ds", (p.width/2)+(len(song)/2)), song)
	p.progress <- progress{msg: msg}

	var done bool
	var paused bool
	l := s.Len()
	var i int
	for !done {
		select {
		case <-time.After(200 * time.Millisecond):
			pos := s.Position()
			done = pos >= l
			i++
			p.playProgress <- progress{n: pos, total: l}
		case v := <-p.vol:
			speaker.Lock()
			vol.Volume += v
			speaker.Unlock()
		case <-p.pause:
			paused = !paused
			speaker.Lock()
			ctrl.Paused = paused
			speaker.Unlock()
		case <-p.fastForward:
			done = true
		}
	}

	return s.Close()
}

func (p *play) clean(s string) string {
	return strings.Replace(s, p.sep, "", -1)
}

func (p *play) next(g *ui.Gui, v *ui.View) error {
	p.fastForward <- true
	return nil
}

func (p *play) getFile(result source.Result) (*os.File, *os.File, error) {
	dir := fmt.Sprintf("%s/.music/cache/%s/%s/%s", os.Getenv("HOME"), p.source.Name(), p.clean(result.Artist.Name), p.clean(result.Album.Title))
	e, err := exists(dir)
	if err != nil {
		return nil, nil, err
	}

	if !e {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return nil, nil, err
		}
	}

	pth := fmt.Sprintf("%s/.music/cache/%s/%s/%s/%s.flac", os.Getenv("HOME"), p.source.Name(), p.clean(result.Artist.Name), p.clean(result.Album.Title), p.clean(result.Track.Title))
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
