package views

import (
	"log"

	"bitbucket.org/cswank/music/internal/source"
	ui "github.com/jroimartin/gocui"
)

type screen struct {
	view   string
	width  int
	height int

	header *header
	body   *body
	volume *volume
	play   *play
	buffer *buffer
	search *search
	login  *login

	keys []key

	source source.Source
}

func newScreen(width, height int) (*screen, error) {
	s := &screen{
		view:   "body",
		width:  width,
		height: height,
		header: newHeader(width, height),
		body:   newBody(width, height),
		play:   newPlay(width, height),
		buffer: newBuffer(width, height),
		volume: newVolume(width, height),
	}

	l := newLogin(width, height, s.doLogin)
	s.search = newSearch(width, height, s.doSearch)
	s.login = l
	s.keys = s.getKeys()
	return s, nil
}

func (s *screen) doLogin(username, passoword string) error {
	s.view = "search-type"
	return g.DeleteView("login")
}

func (s *screen) doSearch() error {
	if s.search.searchType != "" && s.search.searchTerm == "" {
		s.view = "search"
		return g.DeleteView("search-type")
	}

	if s.search.searchTerm != "" {
		s.view = "body"
		log.Printf("searching %s for %s\n", s.search.searchType, s.search.searchTerm)
		return g.DeleteView("search")
	}
	return nil
}

func (s *screen) quit(g *ui.Gui, v *ui.View) error {
	return ui.ErrQuit
}

func (s *screen) showSearch(g *ui.Gui, v *ui.View) error {
	s.view = "search-type"
	return nil
}

func (s *screen) getLayout(width, height int) func(*ui.Gui) error {
	s.view = "login"
	return func(g *ui.Gui) error {
		if s.view == "login" {
			v, err := g.SetView("login", s.login.coords.x1, s.login.coords.y1, s.login.coords.x2, s.login.coords.y2)
			if err != nil && err != ui.ErrUnknownView {
				return err
			}

			g.Cursor = true
			g.InputEsc = true
			ui.DefaultEditor = s.login

			v.Editable = true
			v.Frame = true
			v.Title = s.login.title
		} else if s.view == "search-type" || s.view == "search" {
			v, err := g.SetView(s.view, s.search.coords.x1, s.search.coords.y1, s.search.coords.x2, s.search.coords.y2)
			if err != nil && err != ui.ErrUnknownView {
				return err
			}

			if err := s.search.render(g, v); err != nil {
				return err
			}
		} else {
			v, err := g.SetView("header", s.header.coords.x1, s.header.coords.y1, s.header.coords.x2, s.header.coords.y2)
			if err != nil && err != ui.ErrUnknownView {
				return err
			}

			v.Frame = false
			if err := s.header.render(g, v); err != nil {
				return err
			}

			v, err = g.SetView("body", s.body.coords.x1, s.body.coords.y1, s.body.coords.x2, s.body.coords.y2)
			if err != nil && err != ui.ErrUnknownView {
				return err
			}

			v.Frame = false
			if err := s.body.render(g, v); err != nil {
				return err
			}

			v, err = g.SetView("play", s.play.coords.x1, s.play.coords.y1, s.play.coords.x2, s.play.coords.y2)
			if err != nil && err != ui.ErrUnknownView {
				return err
			}
			v.Frame = false
			v.Editable = true

			if err := s.play.render(g, v); err != nil {
				return err
			}

			v, err = g.SetView("buffer", s.buffer.coords.x1, s.buffer.coords.y1, s.buffer.coords.x2, s.buffer.coords.y2)
			if err != nil && err != ui.ErrUnknownView {
				return err
			}
			v.Frame = false
			v.Editable = true

			if err := s.buffer.render(g, v); err != nil {
				return err
			}

			v, err = g.SetView("volume", s.volume.coords.x1, s.volume.coords.y1, s.volume.coords.x2, s.volume.coords.y2)
			if err != nil && err != ui.ErrUnknownView {
				return err
			}
			v.Frame = false
			v.Editable = true

			if err := s.volume.render(g, v); err != nil {
				return err
			}

		}

		_, err := g.SetCurrentView(s.view)
		return err
	}
}

func (s *screen) keybindings(g *ui.Gui) error {
	for _, k := range s.keys {
		for _, view := range k.views {
			for _, kb := range k.keys {
				if err := g.SetKeybinding(view, kb, ui.ModNone, k.keybinding); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
