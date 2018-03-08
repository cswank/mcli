package views

import (
	"fmt"
	"strings"

	"bitbucket.org/cswank/mcli/internal/player"
	ui "github.com/jroimartin/gocui"
)

type play struct {
	width        int
	coords       coords
	playProgress chan player.Progress
	source       player.Source
	player       *player.Player
}

func newPlay(w, h int, s player.Source, pr chan player.Progress) (*play, error) {
	p := &play{
		width:        w,
		coords:       coords{x1: -1, y1: h - 2, x2: w, y2: h},
		source:       s,
		playProgress: make(chan player.Progress),
	}

	pl, err := player.NewPlayer(s, pr, p.playProgress)
	if err != nil {
		return nil, err
	}

	p.player = pl
	go p.render()
	return p, nil
}

func (p *play) doPause() {
	p.player.Pause()
}

func (p *play) volume(v float64) {
	p.player.Volume(v)
}

func (p *play) addAlbumToQueue(album []player.Result) {
	p.player.PlayAlbum(album)
}

func (p *play) removeFromQueue(i int) {
	p.player.RemoveFromQueue(i)
}

func (p *play) getQueue() []player.Result {
	return p.player.Queue()
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
			fmt.Fprint(v, fmt.Sprintf(strings.Repeat("|", p.width*prog.N/prog.Total)))
			return nil
		})
	}
}

func (p *play) play(r player.Result) {
	p.player.Play(r)
}

func (p *play) next(g *ui.Gui, v *ui.View) error {
	p.player.FastForward()
	return nil
}
