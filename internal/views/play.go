package views

import (
	"fmt"
	"strings"

	"bitbucket.org/cswank/mcli/internal/player"
	ui "github.com/jroimartin/gocui"
)

type play struct {
	width  int
	coords coords
	ch     chan player.Progress
	client player.Client
}

func newPlay(w, h int, c player.Client) *play {
	p := &play{
		width:  w,
		coords: coords{x1: -1, y1: h - 2, x2: w, y2: h},
		client: c,
		ch:     make(chan player.Progress),
	}

	c.PlayProgress(p.playProgress)
	return p
}

func (p *play) doPause() {
	p.client.Pause()
}

func (p *play) volume(v float64) float64 {
	return p.client.Volume(v)
}

func (p *play) addAlbumToQueue(album []player.Result) {
	p.client.PlayAlbum(&player.Results{Results: album})
}

func (p *play) removeFromQueue(i int) {
	p.client.RemoveFromQueue(i)
}

func (p *play) getQueue() []player.Result {
	return p.client.Queue().Results
}

func (p *play) playProgress(prog player.Progress) {
	if prog.Total > 0 {
		g.Update(func(g *ui.Gui) error {
			v, _ := g.View("play")
			v.Clear()
			fmt.Fprint(v, fmt.Sprintf(strings.Repeat("|", p.width*prog.N/prog.Total)))
			return nil
		})
	}
}

func (p *play) clear() {
	v, _ := g.View("play")
	v.Clear()
}

func (p *play) play(r player.Result) {
	p.client.Play(r)
}

func (p *play) next(g *ui.Gui, v *ui.View) error {
	p.client.FastForward()
	return nil
}

func (p *play) rewind(g *ui.Gui, v *ui.View) error {
	p.client.Rewind()
	return nil
}
