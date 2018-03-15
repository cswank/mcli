package player

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cswank/beep"
	"github.com/cswank/beep/effects"
	"github.com/cswank/beep/flac"
	"github.com/cswank/beep/speaker"
)

type Flac struct {
	Fetcher
	queue        *queue
	history      History
	playing      bool
	sep          string
	pause        chan bool
	vol          chan float64
	fastForward  chan bool
	playCB       func(Progress)
	downloadCB   func(Progress)
	nextSong     func(r Result)
	onDeck       chan Result
	onDeckResult *Result
}

func newFlac(f Fetcher) (*Flac, error) {
	hist, err := NewFileHistory()
	if err != nil {
		return nil, err
	}

	p := &Flac{
		Fetcher:     f,
		history:     hist,
		sep:         string(filepath.Separator),
		queue:       newQueue(),
		pause:       make(chan bool),
		fastForward: make(chan bool),
		vol:         make(chan float64),
		onDeck:      make(chan Result),
	}

	go p.playLoop()
	go p.downloadLoop()
	return p, nil
}

func (f *Flac) NextSong(fn func(Result)) {
	f.nextSong = fn
}

func (f *Flac) Play(r Result) {
	f.queue.Add(r)
}

func (f *Flac) History(page, pageSize int) (*Results, error) {
	return f.history.Fetch(page, pageSize)
}

func (f *Flac) PlayAlbum(album []Result) {
	for _, r := range album {
		f.Play(r)
	}
}

func (f *Flac) Pause() {
	if f.playing {
		f.pause <- true
	}
}

func (f *Flac) Volume(v float64) {
	if f.playing {
		f.vol <- v
	}
}

func (f *Flac) Queue() []Result {
	var out []Result
	if f.onDeckResult != nil {
		out = []Result{*f.onDeckResult}
	}
	q := f.queue.Playlist()
	return append(out, q...)
}

func (f *Flac) RemoveFromQueue(i int) {
	//TODO: allow removing of onDeckResult
	f.queue.Remove(i - 1)
}

func (f *Flac) FastForward() {
	if f.playing {
		f.fastForward <- true
	}
}

func (f *Flac) downloadLoop() {
	for {
		r := f.queue.Next()
		f.onDeckResult = &r
		f.download(&r)
		f.onDeck <- r
		f.onDeckResult = nil
	}
}

func (f *Flac) playLoop() {
	for {
		r := <-f.onDeck
		if f.nextSong != nil {
			f.nextSong(r)
		}
		f.playing = true
		if err := f.doPlay(r); err != nil {
			log.Fatal(err)
		}
		f.playing = false
	}
}

func (f *Flac) doPlay(result Result) error {
	if err := f.history.Save(result); err != nil {
		return err
	}

	file, err := os.Open(result.Path)
	if err != nil {
		return err
	}

	s, format, err := flac.Decode(file)
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
			if f.playCB != nil {
				f.playCB(Progress{N: pos, Total: l})
			}
		case v := <-f.vol:
			speaker.Lock()
			vol.Volume += v
			speaker.Unlock()
		case <-f.pause:
			paused = !paused
			speaker.Lock()
			ctrl.Paused = paused
			speaker.Unlock()
		case <-f.fastForward:
			done = true
		}
	}

	return s.Close()
}

func (f *Flac) DownloadProgress(fn func(Progress)) {
	f.downloadCB = fn
}

func (f *Flac) PlayProgress(fn func(Progress)) {
	f.playCB = fn
}

func (f *Flac) download(r *Result) {
	pth, e := f.checkCache(*r)
	log.Println("check cache", pth, e)
	r.Path = pth
	if !e {
		err := f.doDownload(*r)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (f *Flac) doDownload(r Result) error {
	file, err := os.Create(r.Path)
	if err != nil {
		return fmt.Errorf("could not create file for %+v: %s", r, err)
	}

	u, err := f.Fetcher.GetTrack(r.Track.ID)
	if err != nil {
		return fmt.Errorf("could not get track %+v: %s", r, err)
	}

	resp, err := http.Get(u)
	if err != nil {
		return fmt.Errorf("could not get stream %+v: %s", r, err)
	}

	pr := newProgressRead(resp.Body, int(resp.ContentLength), f.downloadCB)
	_, err = io.Copy(file, pr)
	if err != nil {
		return err
	}

	file.Close()
	pr.Close()
	return nil
}

func (f *Flac) ensureCache() error {
	dir := fmt.Sprintf("%s/cache/%s", os.Getenv("MCLI_HOME"), f.Fetcher.Name())
	e, _ := exists(dir)
	if e {
		return nil
	}
	return os.MkdirAll(dir, 0700)
}

func (f *Flac) checkCache(result Result) (string, bool) {
	dir := fmt.Sprintf("%s/cache/%s/%s/%s", os.Getenv("MCLI_HOME"), f.Fetcher.Name(), f.clean(result.Artist.Name), f.clean(result.Album.Title))
	e, _ := exists(dir)
	if !e {
		os.MkdirAll(dir, 0700)
	}

	pth := fmt.Sprintf("%s/cache/%s/%s/%s/%s.flac", os.Getenv("MCLI_HOME"), f.Fetcher.Name(), f.clean(result.Artist.Name), f.clean(result.Album.Title), f.clean(result.Track.Title))
	e, _ = exists(pth)
	return pth, e
}

func (f *Flac) clean(s string) string {
	return strings.Replace(s, f.sep, "", -1)
}

type progressRead struct {
	io.Reader
	t, l, reads int
	cb          func(Progress)
}

func newProgressRead(r io.Reader, l int, cb func(Progress)) *progressRead {
	return &progressRead{Reader: r, t: 0, l: l, cb: cb}
}

func (r *progressRead) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.t += n
	r.reads++
	if r.reads%100 == 0 {
		r.cb(Progress{N: r.t, Total: r.l})
	}
	return n, err
}

// Close the reader when it implements io.Closer
func (r *progressRead) Close() error {
	r.cb(Progress{N: 0, Total: r.t})
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
