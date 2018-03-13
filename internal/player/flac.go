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

func (p *Flac) NextSong(f func(Result)) {
	p.nextSong = f
}

func (p *Flac) Play(r Result) {
	p.queue.add(r)
}

func (p *Flac) History(page, pageSize int) (*Results, error) {
	return p.history.Fetch(page, pageSize)
}

func (p *Flac) PlayAlbum(album []Result) {
	for _, r := range album {
		p.Play(r)
	}
}

func (p *Flac) Pause() {
	if p.playing {
		p.pause <- true
	}
}

func (p *Flac) Volume(v float64) {
	if p.playing {
		p.vol <- v
	}
}

func (p *Flac) Queue() []Result {
	return p.queue.playlist()
}

func (p *Flac) RemoveFromQueue(i int) {
	p.queue.remove(i)
}

func (p *Flac) FastForward() {
	if p.playing {
		p.fastForward <- true
	}
}

func (p *Flac) downloadLoop() {
	for {
		prog := <-p.downloadProgress
		if p.downloadCB != nil {
			p.downloadCB(prog)
		}
	}
}

func (p *Flac) loop() {
	for {
		r := p.queue.next()
		if p.nextSong != nil {
			p.nextSong(r)
		}
		p.playing = true
		if err := p.doPlay(r); err != nil {
			log.Fatal(err)
		}
		p.playing = false
	}
}

func (p *Flac) doPlay(result Result) error {
	if err := p.history.Save(result); err != nil {
		return err
	}

	f, err := os.Open(result.Path)
	if err != nil {
		return err
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
			if p.playCB != nil {
				p.playCB(Progress{N: pos, Total: l})
			}
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

func (p *Flac) DownloadProgress(f func(Progress)) {
	p.downloadCB = f
}

func (p *Flac) PlayProgress(f func(Progress)) {
	p.playCB = f
}
