package views

import (
	"fmt"

	ui "github.com/jroimartin/gocui"
)

const (
	volumeWidth  = 3
	volumeHeight = 15
)

type volume struct {
	step   float64
	coords coords
	vol    float64
}

func newVolume(w, h int, v float64) *volume {
	x1 := w - 2
	x2 := w + 1
	y1 := 0
	y2 := h

	return &volume{
		step:   7.0 / float64(h),
		coords: coords{x1: x1, y1: y1, x2: x2, y2: y2},
		vol:    v,
	}
}

func (v *volume) clear() error {
	vw, err := g.View("volume")
	if err != nil {
		return err
	}
	vw.Clear()
	return nil
}

func (v *volume) render(g *ui.Gui, vw *ui.View) error {
	vw.Editable = false
	vw.Frame = false
	vw.Clear()
	c := 0
	for i := 2.0; i >= -5.0; i -= v.step {
		vw.SetCursor(0, c)
		s := "\n"
		if i <= v.vol {
			s = "█\n"
		}
		fmt.Fprint(vw, s)
		c++
	}

	return nil
}
