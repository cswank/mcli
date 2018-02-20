package views

import ui "github.com/jroimartin/gocui"

type play struct {
	coords coords
}

func newPlay(w, h int) *play {
	return &play{coords: coords{x1: -1, y1: h - 3, x2: w, y2: h - 1}}
}

func (p *play) render(g *ui.Gui, v *ui.View) error {
	return nil
}
