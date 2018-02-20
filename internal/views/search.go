package views

import (
	ui "github.com/jroimartin/gocui"
)

const (
	searchWidth  = 30
	searchHeight = 2
)

type searchType struct {
	name   string
	coords coords

	searchType string
}

func newSearchType(w, h int) *searchType {
	maxX, maxY := g.Size()
	x1 := maxX/2 - searchWidth/2
	x2 := maxX/2 + searchWidth/2
	y1 := maxY/2 - searchHeight/2
	y2 := maxY/2 + searchHeight/2 + searchHeight%2

	return &searchType{
		coords: coords{x1: x1, y1: y1, x2: x2, y2: y2},
	}
}

func (s *searchType) render(g *ui.Gui, v *ui.View) error {
	if s.searchType == "" {
		v.Editable = false
		v.Frame = true
		v.Title = "search"
		// v, err := g.SetView("search", s.coords.x1, s.coords.y1, s.coords.x2, s.coords.y2)
		// if err != ui.ErrUnknownView {
		// 	return err
		// }

		v.Clear()
		v.Write([]byte(c1("albu") + c2("m") + c1(" / artis") + c2("t") + c1(" / trac") + c2("k")))
		_, err := g.SetCurrentView("search")
		return err
	}

	return nil
}
