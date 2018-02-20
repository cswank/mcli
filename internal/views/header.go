package views

import ui "github.com/jroimartin/gocui"

type header struct {
	coords coords
}

func newHeader(w, h int) *header {
	return &header{coords: coords{x1: -1, y1: -1, x2: w - 1, y2: 1}}
}

func (h *header) render(g *ui.Gui, v *ui.View) error {
	return nil
}
