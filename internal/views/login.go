package views

import (
	ui "github.com/jroimartin/gocui"
)

const (
	searchWidth  = 30
	searchHeight = 2
)

type login struct {
	name   string
	coords coords

	login string
}

func newLogin(w, h int) *login {
	maxX, maxY := g.Size()
	x1 := maxX/2 - searchWidth/2
	x2 := maxX/2 + searchWidth/2
	y1 := maxY/2 - searchHeight/2
	y2 := maxY/2 + searchHeight/2 + searchHeight%2

	return &login{
		coords: coords{x1: x1, y1: y1, x2: x2, y2: y2},
	}
}

func (s *login) render(g *ui.Gui, v *ui.View) error {
	v.Editable = false
	v.Frame = true
	v.Title = "sign in"
	return nil
}
