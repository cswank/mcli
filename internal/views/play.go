package views

import (
	"fmt"
	"strings"

	"bitbucket.org/cswank/mcli/internal/schema"
	ui "github.com/jroimartin/gocui"
)

type player struct {
	width  int
	coords coords
	client *client
}

func newPlay(w, h int, id string, c *client) *player {
	p := &player{
		width:  w,
		coords: coords{x1: -1, y1: h - 2, x2: w, y2: h},
		client: c,
	}

	c.PlayProgress(id, p.playProgress)
	return p
}

func (p *player) doPause() {
	p.client.Pause()
}

func (p *player) volume(v float64) float64 {
	return p.client.Volume(v)
}

func (p *player) addAlbumToQueue(album []schema.Result) {
	p.client.PlayAlbum(&schema.Results{Results: album})
}

func (p *player) removeFromQueue(i int) {
	p.client.RemoveFromQueue([]int{i})
}

func (p *player) getQueue() []schema.Result {
	return p.client.Queue().Results
}

func (p *player) playProgress(prog schema.Progress) {
	if prog.Total > 0 {
		g.Update(func(g *ui.Gui) error {
			v, _ := g.View("play")
			v.Clear()
			fmt.Fprint(v, fmt.Sprintf(strings.Repeat("|", p.width*prog.N/prog.Total)))
			return nil
		})
	}
}

func (p *player) clear() {
	v, _ := g.View("play")
	v.Clear()
}

func (p *player) play(r schema.Result) {
	p.client.Play(r)
}

func (p *player) next() {
	p.client.FastForward()
	p.clear()
}

func (p *player) rewind() {
	p.client.Rewind()
	p.clear()
}
