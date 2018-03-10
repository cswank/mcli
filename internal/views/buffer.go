package views

import (
	"fmt"
	"strings"

	"bitbucket.org/cswank/mcli/internal/player"
	ui "github.com/jroimartin/gocui"
)

type buffer struct {
	width    int
	coords   coords
	progress chan player.Progress
}

func newBuffer(w, h int, ch chan player.Progress) *buffer {
	b := &buffer{
		width:    w - 1,
		coords:   coords{x1: -1, y1: h - 3, x2: w, y2: h - 1},
		progress: ch,
	}

	go b.render(b.progress)
	return b
}

func (b *buffer) render(ch <-chan player.Progress) {
	var v *ui.View
	for {
		p := <-ch
		g.Update(func(g *ui.Gui) error {
			if v == nil {
				v, _ = g.View("buffer")
			}
			v.Clear()
			fmt.Fprint(v, fmt.Sprintf(strings.Repeat("|", b.width*p.N/p.Total)))
			return nil
		})
	}
}

func (b *buffer) clear() {
	v, _ := g.View("buffer")
	v.Clear()
}
