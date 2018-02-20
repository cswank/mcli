package views

import (
	"fmt"

	ui "github.com/jroimartin/gocui"
)

const (
	searchWidth  = 50
	searchHeight = 4
)

type search struct {
	name   string
	coords coords

	searchType string
}

func newSearch(w, h int) *search {
	maxX, maxY := g.Size()
	x1 := maxX/2 - searchWidth/2
	x2 := maxX/2 + searchWidth/2
	y1 := maxY/2 - searchHeight/2
	y2 := maxY/2 + searchHeight/2 + searchHeight%2

	return &search{
		coords: coords{x1: x1, y1: y1, x2: x2, y2: y2},
	}
}

func (s *search) render(g *ui.Gui, v *ui.View) error {
	if s.searchType == "" {
		v.Editable = false
		v.Frame = true
		v.Title = "search"
		v, err := g.SetView("search", s.coords.x1, s.coords.y1, s.coords.x2, s.coords.y2)
		if err != ui.ErrUnknownView {
			return err
		}

		v.Write([]byte(fmt.Sprintf(c1("albu%s / artis%s / trac%s"), c2("m"), c1("t"), c1("k"))))
		_, err = g.SetCurrentView("search")
		return err
	}

	return nil
}
