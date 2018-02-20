package views

import ui "github.com/jroimartin/gocui"

type volume struct {
	coords coords
}

func newVolume(w, h int) *volume {
	return &volume{coords: coords{x1: w - 2, y1: -1, x2: w, y2: h}}
}

func (v *volume) render(g *ui.Gui, vw *ui.View) error {
	return nil
}
