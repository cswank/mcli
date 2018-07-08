package views

import (
	"bitbucket.org/cswank/mcli/internal/player"
	ui "github.com/jroimartin/gocui"
)

//artistDialog is a dialog that lets you choos between showing the albums
//or songs of an artist
type artistDialog struct {
	name     string
	coords   coords
	selected player.Result

	searchType string
	callback   func(string, string) error
}

func newArtistDialog(w, h int, cb func(string, string) error) *artistDialog {
	maxX, maxY := g.Size()
	x1 := maxX/2 - searchWidth/2
	x2 := maxX/2 + searchWidth/2
	y1 := maxY/2 - searchHeight/2
	y2 := maxY/2 + searchHeight/2 + searchHeight%2

	return &artistDialog{
		coords:   coords{x1: x1, y1: y1, x2: x2, y2: y2},
		callback: cb,
	}
}

func (a *artistDialog) render(g *ui.Gui, v *ui.View) error {
	v.Editable = false
	v.Frame = true
	v.Title = "show"
	v.Clear()
	_, err := v.Write([]byte(c1("albu") + c2("m") + c1("s / trac") + c2("k") + c1("s")))
	return err
}

func (a *artistDialog) albums(g *ui.Gui, v *ui.View) error {
	return a.callback(a.selected.Artist.ID, "albums")
}

func (a *artistDialog) tracks(g *ui.Gui, v *ui.View) error {
	return a.callback(a.selected.Artist.ID, "tracks")
}
