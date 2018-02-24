package views

import (
	"bitbucket.org/cswank/music/internal/source"
	ui "github.com/jroimartin/gocui"
)

type body struct {
	coords  coords
	height  int
	results *source.Results
	cursor  int
}

func newBody(w, h int) *body {
	return &body{
		height: h - 3,
		coords: coords{x1: -1, y1: 0, x2: w, y2: h - 2},
	}
}

func (b *body) render(g *ui.Gui, v *ui.View) error {
	v.Clear()
	if err := v.SetCursor(0, b.cursor); err != nil {
		return nil
	}
	if b.results == nil {
		return nil
	}

	for _, r := range b.results.Results {
		if err := b.results.Print(v, r); err != nil {
			return err
		}
	}
	return nil
}

func (b *body) clear() {
	v, _ := g.View("body")
	v.Clear()
}

func (b *body) next(g *ui.Gui, v *ui.View) error {
	if b.cursor == b.height-1 || (b.results != nil && b.cursor == len(b.results.Results)-1) {
		return nil
	}
	b.cursor++
	return nil
}

func (b *body) prev(g *ui.Gui, v *ui.View) error {
	if b.cursor == 0 {
		return nil
	}
	b.cursor--
	return nil
}
