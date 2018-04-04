package views

import (
	"bitbucket.org/cswank/mcli/internal/player"
	ui "github.com/jroimartin/gocui"
)

const (
	historyWidth  = 30
	historyHeight = 2
)

type history struct {
	coords    coords
	doHistory func(player.Sort) error
}

func newHistory(w, h int, cb func(player.Sort) error) *history {
	maxX, maxY := g.Size()
	x1 := maxX/2 - historyWidth/2
	x2 := maxX/2 + historyWidth/2
	y1 := maxY/2 - historyHeight/2
	y2 := maxY/2 + historyHeight/2 + historyHeight%2

	return &history{
		coords:    coords{x1: x1, y1: y1, x2: x2, y2: y2},
		doHistory: cb,
	}
}

func (h *history) recent(g *ui.Gui, v *ui.View) error {
	v.Clear()
	return h.doHistory(player.Time)
}

func (h *history) played(g *ui.Gui, v *ui.View) error {
	v.Clear()
	return h.doHistory(player.Count)
}

func (s *history) render(g *ui.Gui, v *ui.View) error {
	v.Editable = false
	v.Frame = true
	v.Title = "sort by"
	v.Clear()
	_, err := v.Write([]byte(c1("most ") + c2("r") + c1("ecent / most ") + c2("p") + c1("layed")))
	return err
}
