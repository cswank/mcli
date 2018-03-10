package views

import (
	"fmt"
	"strings"
	"time"

	"bitbucket.org/cswank/mcli/internal/player"
	ui "github.com/jroimartin/gocui"
)

type buffer struct {
	width    int
	coords   coords
	progress chan player.Progress
	song     chan player.Result
}

func newBuffer(w, h int, ch chan player.Progress, song chan player.Result) *buffer {
	b := &buffer{
		width:    w - 1,
		coords:   coords{x1: -1, y1: h - 3, x2: w, y2: h - 1},
		progress: ch,
		song:     song,
	}

	go b.render()
	return b
}

func (b *buffer) render() {
	var v *ui.View
	var text string
	for {
		select {
		case <-time.After(time.Second):
			if v == nil {
				v, _ = g.View("buffer")
			}
			if text != "" && text != v.Buffer() {
				g.Update(func(g *ui.Gui) error {
					v.Clear()
					fmt.Fprint(v, text)
					return nil
				})
			}
		case r := <-b.song:
			text = b.center(fmt.Sprintf("%s %s", r.Track.Title, time.Duration(r.Track.Duration)*time.Second))
			g.Update(func(g *ui.Gui) error {
				v.Clear()
				fmt.Fprint(v, text)
				return nil
			})
		case p := <-b.progress:
			g.Update(func(g *ui.Gui) error {
				v.Clear()
				fmt.Fprint(v, fmt.Sprintf(strings.Repeat("|", b.width*p.N/p.Total)))
				return nil
			})
		}
	}
}

func (b *buffer) clear() {
	v, _ := g.View("buffer")
	v.Clear()
}

func (b *buffer) center(s string) string {
	return fmt.Sprintf(fmt.Sprintf("%%%ds", (b.width/2)+(len(s)/2)), s)
}
