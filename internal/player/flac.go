package player

import (
	"log"
	"os"
	"time"

	"github.com/cswank/beep"
	"github.com/cswank/beep/effects"
	"github.com/cswank/beep/flac"
	"github.com/cswank/beep/speaker"
)

type Flac struct {
	Fetcher
	queue            *queue
	source           Fetcher
	history          History
	playing          bool
	pause            chan bool
	vol              chan float64
	fastForward      chan bool
	playCB           func(Progress)
	downloadProgress chan Progress
	downloadCB       func(Progress)
	nextSong         func(r Result)
}

func newFlac(f Fetcher) (*Flac, error) {
	hist, err := NewFileHistory()
	if err != nil {
		return nil, err
	}

	dlp := make(chan Progress)
	q, err := newQueue(f.Name(), f.GetTrack, dlp)
	if err != nil {
		return nil, err
	}

	p := &Flac{
		Fetcher:          f,
		history:          hist,
		queue:            q,
		pause:            make(chan bool),
		fastForward:      make(chan bool),
		vol:              make(chan float64),
		downloadProgress: dlp,
	}

	go p.loop()
	go p.downloadLoop()
	return p, nil
}

func (f *Flac) NextSong(fn func(Result)) {
	f.nextSong = fn
}

func (f *Flac) Play(r Result) {
	f.queue.add(r)
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
	return f.queue.playlist()
}

func (f *Flac) RemoveFromQueue(i int) {
	f.queue.remove(i)
}

func (f *Flac) FastForward() {
	if f.playing {
		f.fastForward <- true
	}
}

func (f *Flac) downloadLoop() {
	for {
		prog := <-f.downloadProgress
		if f.downloadCB != nil {
			f.downloadCB(prog)
		}
	}
}

func (f *Flac) loop() {
	for {
		r := f.queue.next()
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
