package views

import (
	"fmt"
	"strings"

	ui "github.com/jroimartin/gocui"
)

type buffer struct {
	width    int
	coords   coords
	progress chan progress
}

func newBuffer(w, h int) *buffer {
	b := &buffer{
		width:    w - 1,
		coords:   coords{x1: -1, y1: h - 3, x2: w, y2: h - 1},
		progress: make(chan progress),
	}

	go b.render(b.progress)
	return b
}

func (b *buffer) render(ch <-chan progress) {
	var v *ui.View
	for {
		p := <-ch
		g.Update(func(g *ui.Gui) error {
			if v == nil {
				v, _ = g.View("buffer")
			}
			v.Clear()
			fmt.Fprint(v, fmt.Sprintf(strings.Repeat("|", b.width*p.n/p.total)))
			return nil
		})
	}
}

func (b *buffer) clear() {
	v, _ := g.View("buffer")
	v.Clear()
}
