package views

import (
	"fmt"
	"strings"

	ui "github.com/jroimartin/gocui"
)

const (
	searchWidth  = 30
	searchHeight = 2
)

type search struct {
	name   string
	coords coords

	searchType string
	doSearch   func(string, string) error
}

func newSearch(w, h int, cb func(string, string) error) *search {
	maxX, maxY := g.Size()
	x1 := maxX/2 - searchWidth/2
	x2 := maxX/2 + searchWidth/2
	y1 := maxY/2 - searchHeight/2
	y2 := maxY/2 + searchHeight/2 + searchHeight%2

	return &search{
		coords:   coords{x1: x1, y1: y1, x2: x2, y2: y2},
		doSearch: cb,
	}
}

func (s *search) album(g *ui.Gui, v *ui.View) error {
	s.searchType = "album"
	v.Clear()
	return s.doSearch(s.searchType, "")
}

func (s *search) artist(g *ui.Gui, v *ui.View) error {
	s.searchType = "artist"
	v.Clear()
	return s.doSearch(s.searchType, "")
}

func (s *search) track(g *ui.Gui, v *ui.View) error {
	s.searchType = "track"
	v.Clear()
	return s.doSearch(s.searchType, "")
}

func (s *search) escape(g *ui.Gui, v *ui.View) error {
	s.searchType = ""
	return nil
}

func (s *search) exit(g *ui.Gui, v *ui.View) error {
	st := s.searchType
	s.searchType = ""
	t := strings.TrimSpace(v.Buffer())
	v.Clear()
	return s.doSearch(st, t)
}

func (s *search) render(g *ui.Gui, v *ui.View) error {
	if s.searchType == "" {
		v.Editable = false
		v.Frame = true
		v.Title = "search"
		v.Clear()
		_, err := v.Write([]byte(c1("albu") + c2("m") + c1(" / artis") + c2("t") + c1(" / trac") + c2("k")))
		return err
	}

	v.Editable = true
	v.Frame = true
	v.Title = fmt.Sprintf("search %s", s.searchType)
	return nil
}

func (s *search) Edit(v *ui.View, key ui.Key, ch rune, mod ui.Modifier) {
	in := string(ch)
	buf := strings.TrimSpace(v.Buffer())
	if key == 127 && len(buf) > 0 {
		v.Clear()
		buf = buf[:len(buf)-1]
		v.Write([]byte(c1(buf)))
		v.SetCursor(len(buf), 0)
	} else {
		fmt.Fprint(v, c1(in))
		buf = v.Buffer()
		v.SetCursor(len(buf)-1, 0)
	}
}
