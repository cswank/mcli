package views

import ui "github.com/jroimartin/gocui"

type screen struct {
	view string

	header *header
	body   *body
	volume *volume
	play   *play
	buffer *buffer
	search *search
}

func newScreen() (*screen, error) {
	return nil, nil
}

func (s *screen) getLayout(width, height int) func(*ui.Gui) error {
	//ui.DefaultEditor = s.search

	return func(g *ui.Gui) error {
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

		_, err = g.SetCurrentView(s.view)
		return err
	}
}
