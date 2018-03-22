package views

import (
	"path"

	"bitbucket.org/cswank/mcli/internal/player"
	"github.com/atotto/clipboard"
	ui "github.com/jroimartin/gocui"
)

type body struct {
	albumURL string
	coords   coords
	height   int
	results  *player.Results
	view     []player.Result
	cursor   int
	page     int
}

func newBody(w, h int, albumLink string) *body {
	return &body{
		albumURL: albumLink,
		height:   h - 3,
		coords:   coords{x1: -1, y1: 0, x2: w, y2: h - 2},
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

	for _, r := range b.view {
		if err := b.results.Print(v, r); err != nil {
			return err
		}
	}
	return nil
}

func (b *body) newResults(r *player.Results) {
	b.page = 0
	b.results = r
	b.makeView()
}

func (b *body) albumLink(g *ui.Gui, v *ui.View) error {
	if b.results == nil {
		return nil
	}
	c, _ := v.Cursor()
	r := b.results.Results[c]
	l := path.Join(b.albumURL, r.Album.ID)
	return clipboard.WriteAll(l)
}

func (b *body) clear() {
	v, _ := g.View("body")
	v.Clear()
}

func (b *body) nextPage(g *ui.Gui, v *ui.View) error {
	if b.page >= len(b.results.Results)-b.height {
		return nil
	}

	b.page++
	b.makeView()
	return nil
}

func (b *body) prevPage(g *ui.Gui, v *ui.View) error {
	if b.page == 0 {
		return nil
	}

	b.page--
	b.makeView()
	return nil
}

func (b *body) makeView() {
	start := b.page * b.height
	end := start + b.height
	if end >= len(b.results.Results) {
		end = len(b.results.Results)
	}

	b.view = b.results.Results[start:end]
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
