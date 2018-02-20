package views

import ui "github.com/jroimartin/gocui"

type body struct {
	coords coords
}

func newBody(w, h int) *body {
	return &body{
		coords: coords{x1: -1, y1: 0, x2: w - 1, y2: h - 2},
	}
}

func (b *body) render(g *ui.Gui, v *ui.View) error {
	return nil
}
