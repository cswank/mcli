package views

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"bitbucket.org/cswank/mcli/internal/history"
	"bitbucket.org/cswank/mcli/internal/source"
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
	playProgress chan progress
	source       source.Source
	pause        chan bool
	fastForward  chan bool
	vol          chan float64
	history      history.History
	queue        *queue
}

func newPlay(w, h int, s source.Source, pr chan<- progress) (*play, error) {
	hist, err := history.NewFileHistory()
	if err != nil {
		return nil, err
	}

	p := &play{
		width:        w,
		coords:       coords{x1: -1, y1: h - 2, x2: w, y2: h},
		source:       s,
		playProgress: make(chan progress),
		pause:        make(chan bool),
		fastForward:  make(chan bool),
		vol:          make(chan float64),
		history:      hist,
		queue:        newQueue(s, pr),
	}

	go p.loop()
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

}

func (p *play) removeFromQueue(i int) {

}

func (p *play) getQueue() []source.Result {
	return nil
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

func (p *play) play(r source.Result) {
	p.queue.add(r)
}

func (p *play) loop() {
	for {
		r := p.queue.next()
		if err := p.doPlay(r); err != nil {
			log.Fatal(err)
		}
	}
}

func (p *play) doPlay(result source.Result) error {
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

func (p *play) next(g *ui.Gui, v *ui.View) error {
	p.fastForward <- true
	return nil
}
