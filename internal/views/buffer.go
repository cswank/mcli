package views

import ui "github.com/jroimartin/gocui"

type buffer struct {
	coords coords
}

func newBuffer(w, h int) *buffer {
	return &buffer{coords: coords{x1: -1, y1: h - 2, x2: w - 1, y2: h}}
}

func (b *buffer) render(g *ui.Gui, v *ui.View) error {
	return nil
}
