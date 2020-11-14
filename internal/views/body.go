package views

import (
	"fmt"

	ui "github.com/awesome-gocui/gocui"
	"github.com/cswank/mcli/internal/schema"
)

type body struct {
	albumURL string
	coords   coords
	height   int
	results  *schema.Results
	view     []schema.Result
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

	for i, r := range b.view {
		s := b.results.Print(r)
		if i == b.cursor {
			s = col.C2(s)
		} else {
			s = col.C1(s)
		}
		if _, err := fmt.Fprint(v, s); err != nil {
			return err
		}
	}
	return nil
}

func (b *body) newResults(r *schema.Results) {
	b.page = 0
	b.results = r
	b.makeView()
}

func (b *body) clear() {
	v, _ := g.View("body")
	v.Clear()
}

func (b *body) nextPage(g *ui.Gui, v *ui.View) error {
	if b.results.Page == nil {
		return nil
	}

	b.cursor = 0
	b.page++
	p := b.results.Print
	pg := b.results.Page
	results, err := b.results.Page(b.page)
	results.Print = p
	results.Page = pg
	b.results = results
	b.makeView()
	return err
}

func (b *body) prevPage(g *ui.Gui, v *ui.View) error {
	if b.page == 0 || b.results.Page == nil {
		return nil
	}

	b.cursor = 0
	b.page--
	p := b.results.Print
	pg := b.results.Page
	results, err := b.results.Page(b.page)
	results.Print = p
	results.Page = pg
	b.results = results
	b.makeView()
	return err
}

func (b *body) makeView() {
	end := b.height
	if end >= len(b.results.Results) {
		end = len(b.results.Results)
	}

	b.view = b.results.Results[0:end]
}

func (b *body) next(g *ui.Gui, v *ui.View) error {
	if b.cursor == b.height-1 || (b.results != nil && b.cursor == len(b.view)-1) {
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
