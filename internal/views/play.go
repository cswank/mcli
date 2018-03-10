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
	song         chan player.Result
	client       player.Client
}

func newPlay(w, h int, c player.Client, ch chan player.Progress, song chan player.Result) *play {
	p := &play{
		width:        w,
		coords:       coords{x1: -1, y1: h - 2, x2: w, y2: h},
		client:       c,
		playProgress: ch,
		song:         song,
	}

	c.NextSong(p.nextSong)
	go p.render()
	return p
}

func (p *play) doPause() {
	p.client.Pause()
}

func (p *play) volume(v float64) {
	p.client.Volume(v)
}

func (p *play) addAlbumToQueue(album []player.Result) {
	p.client.PlayAlbum(album)
}

func (p *play) removeFromQueue(i int) {
	p.client.RemoveFromQueue(i)
}

func (p *play) getQueue() []player.Result {
	return p.client.Queue()
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

//nextSong gets called by player.Player whenever the
//next song begins playing.
func (p *play) nextSong(r player.Result) {
	p.song <- r
}

func (p *play) play(r player.Result) {
	p.client.Play(r)
}

func (p *play) next(g *ui.Gui, v *ui.View) error {
	p.client.FastForward()
	return nil
}
