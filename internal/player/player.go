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

type Player struct {
	queue        *queue
	source       Source
	playProgress chan Progress
	history      History
	playing      bool
	pause        chan bool
	vol          chan float64
	fastForward  chan bool
}

func NewPlayer(s Source, download chan Progress, play chan Progress) (*Player, error) {
	hist, err := NewFileHistory()
	if err != nil {
		return nil, err
	}

	p := &Player{
		playProgress: play,
		source:       s,
		history:      hist,
		queue:        newQueue(s, download),
		pause:        make(chan bool),
		fastForward:  make(chan bool),
		vol:          make(chan float64),
	}

	go p.loop()
	return p, nil
}

func (p *Player) Play(r Result) {
	p.queue.add(r)
}

func (p *Player) History(page, pageSize int) (*Results, error) {
	return p.history.Fetch(page, pageSize)
}

func (p *Player) PlayAlbum(album []Result) {
	for _, r := range album {
		p.Play(r)
	}
}

func (p *Player) Pause() {
	if p.playing {
		p.pause <- true
	}
}

func (p *Player) Volume(v float64) {
	if p.playing {
		p.vol <- v
	}
}

func (p *Player) Queue() []Result {
	return p.queue.playlist()
}

func (p *Player) RemoveFromQueue(i int) {
	p.queue.remove(i)
}

func (p *Player) FastForward() {
	if p.playing {
		p.fastForward <- true
	}
}

func (p *Player) loop() {
	for {
		r := p.queue.next()
		p.playing = true
		if err := p.doPlay(r); err != nil {
			log.Fatal(err)
		}
		p.playing = false
	}
}

func (p *Player) doPlay(result Result) error {
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

	//song := fmt.Sprintf("%s %s", result.Track.Title, time.Duration(result.Track.Duration)*time.Second)
	//msg := fmt.Sprintf(fmt.Sprintf("%%%ds", (p.width/2)+(len(song)/2)), song)
	//p.progress <- progress{msg: msg}

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
			p.playProgress <- Progress{N: pos, Total: l}
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
