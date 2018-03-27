package views

import (
	"fmt"

	"bitbucket.org/cswank/mcli/internal/player"
	ui "github.com/jroimartin/gocui"
)

type screen struct {
	view   string
	width  int
	height int

	header      *header
	body        *body
	play        *play
	buffer      *buffer
	search      *search
	history     *history
	historySort player.Sort
	login       *login
	help        *help

	keys []key

	client player.Client
	stack  stack
}

func newScreen(width, height int, p player.Player) (*screen, error) {
	cli, err := player.NewTidal(p)
	if err != nil {
		return nil, err
	}

	s := &screen{
		client:      cli,
		view:        "body",
		width:       width,
		height:      height,
		historySort: player.Time,
		play:        newPlay(width, height, cli),
		body:        newBody(width, height, cli.AlbumLink()),
		buffer:      newBuffer(width, height, cli),
		header:      newHeader(width, height),
		help:        newHelp(width, height),
	}

	s.login = newLogin(width, height, s.doLogin)
	s.search = newSearch(width, height, s.doSearch)
	s.history = newHistory(width, height, s.showHistory)
	s.keys = s.getKeys()
	return s, nil
}

func (s *screen) playAlbum(g *ui.Gui, v *ui.View) error {
	if s.body.results.Type != "album" || s.body.results == nil || len(s.body.results.Results) == 0 {
		return nil
	}

	s.play.addAlbumToQueue(s.body.results.Results)
	return nil
}

func (s *screen) enter(g *ui.Gui, v *ui.View) error {
	r := s.body.view[s.body.cursor]
	c := s.body.cursor
	switch s.body.results.Type {
	case "album search":
		s.body.cursor = 0
		results, err := s.client.GetAlbum(r.Album.ID)
		if err != nil {
			return err
		}
		results.Print = results.PrintAlbumTracks()
		s.body.newResults(results)
		s.header.header = results.Header
		s.stack.add(results, c)
	case "artist search":
		s.body.cursor = 0
		results, err := s.client.GetArtistAlbums(r.Artist.ID, s.height*5)
		if err != nil {
			return err
		}
		results.Print = results.PrintArtist()
		s.body.newResults(results)
		s.header.header = results.Header
		s.stack.add(results, c)
	case "artist albums":
		s.body.cursor = 0
		results, err := s.client.GetAlbum(r.Album.ID)
		if err != nil {
			return err
		}

		results.Print = results.PrintAlbum()
		s.body.newResults(results)
		s.header.header = results.Header
		s.stack.add(results, c)
	case "playlists":
		s.body.cursor = 0
		results, err := s.client.GetPlaylist(r.Album.ID, s.height)
		if err != nil {
			return err
		}
		results.Print = results.PrintAlbum()
		s.body.newResults(results)
		s.header.header = results.Header
		s.stack.add(results, c)
	case "album":
		s.play.play(r)
	case "history":
		s.play.play(r)
	case "playlist":
		s.play.play(r)
	}
	return nil
}

func (s *screen) playlists(g *ui.Gui, v *ui.View) error {
	results, err := s.client.GetPlaylists()
	if err != nil {
		return err
	}

	results.Print = results.PrintPlaylists()
	s.body.cursor = 0
	s.body.newResults(results)
	s.header.header = results.Header
	s.stack.clear()
	s.stack.add(results, s.body.cursor)
	return nil
}

func (s *screen) queue(g *ui.Gui, v *ui.View) error {
	s.body.cursor = 0
	return s.queueNoCursor(g, v)
}

func (s *screen) queueNoCursor(g *ui.Gui, v *ui.View) error {
	items := s.play.getQueue()
	if len(items) == 0 {
		return nil
	}

	var maxTitle, maxAlbum int
	for _, item := range items {
		if len(item.Track.Title) > maxTitle {
			maxTitle = len(item.Track.Title)
		}
		if len(item.Album.Title) > maxAlbum {
			maxAlbum = len(item.Album.Title)
		}
	}

	f := fmt.Sprintf("%%-%ds%%-%ds%%s\n", maxTitle+4, maxAlbum+4)
	results := &player.Results{
		Header: fmt.Sprintf(f, "Title", "Album", "Artist"),
		Type:   "queue",
		Print: func(r player.Result) string {
			return fmt.Sprintf(f, r.Track.Title, r.Album.Title, r.Artist.Name)
		},
		Results: items,
	}

	s.body.newResults(results)
	s.header.header = results.Header
	return nil
}

func (s *screen) removeFromQueue(g *ui.Gui, v *ui.View) error {
	if s.body.results.Type != "queue" {
		return nil
	}
	s.play.removeFromQueue(s.body.cursor)
	if err := s.queueNoCursor(g, v); err != nil {
		return err
	}

	if len(s.body.results.Results) > 0 && s.body.cursor >= len(s.body.results.Results) {
		s.body.cursor = len(s.body.results.Results) - 1
	}
	return nil
}

func (s *screen) escapeSearch(g *ui.Gui, v *ui.View) error {
	s.search.searchType = ""
	if s.view == "search-type" {
		s.view = "body"
	} else {
		s.view = "search-type"
		s.search.searchType = ""
	}
	return nil
}

func (s *screen) goToAlbum(g *ui.Gui, v *ui.View) error {
	r := s.body.view[s.body.cursor]
	c := s.body.cursor
	results, err := s.client.GetAlbum(r.Album.ID)
	if err != nil {
		return err
	}

	results.Print = results.PrintAlbumTracks()
	s.body.newResults(results)
	s.header.header = results.Header
	s.stack.add(results, c)
	s.body.cursor = 0
	return nil
}

func (s *screen) goToArtist(g *ui.Gui, v *ui.View) error {
	r := s.body.view[s.body.cursor]
	c := s.body.cursor
	results, err := s.client.GetArtistAlbums(r.Artist.ID, s.height)
	if err != nil {
		return err
	}

	results.Print = results.PrintArtist()
	s.body.newResults(results)
	s.header.header = results.Header
	s.stack.add(results, c)
	s.body.cursor = 0
	return nil
}

func (s *screen) goToArtistTracks(g *ui.Gui, v *ui.View) error {
	r := s.body.view[s.body.cursor]
	c := s.body.cursor
	results, err := s.client.GetArtistTracks(r.Artist.ID, s.height)
	if err != nil {
		return err
	}

	results.Print = results.PrintArtistTracks()
	s.body.newResults(results)
	s.header.header = results.Header
	s.stack.add(results, c)
	s.body.cursor = 0
	return nil
}

func (s *screen) volumeUp(g *ui.Gui, v *ui.View) error {
	s.play.volume(0.5)
	return nil
}

func (s *screen) volumeDown(g *ui.Gui, v *ui.View) error {
	s.play.volume(-0.5)
	return nil
}

func (s *screen) showHelp(g *ui.Gui, v *ui.View) error {
	s.view = "help"
	return nil
}

func (s *screen) hideHelp(g *ui.Gui, v *ui.View) error {
	s.view = "body"
	return s.help.hide(g, v)
}

func (s *screen) showHistory(sort player.Sort) error {
	g.DeleteView("history-type")
	s.historySort = sort
	res, err := s.client.History(0, s.height*10, sort)
	if err != nil {
		return err
	}

	s.view = "body"
	res.Print = res.PrintHistory()
	s.header.header = res.Header
	s.body.newResults(res)
	s.stack.add(res, s.body.cursor)
	s.body.cursor = 0
	return nil
}

func (s *screen) doLogin(username, password string) error {
	s.view = "search-type"
	err := s.client.Login(username, password)
	if err != nil {
		return err
	}

	return g.DeleteView("login")
}

func (s *screen) doSearch(searchType, term string) error {
	if searchType != "" && term == "" {
		s.view = "search"
		return g.DeleteView("search-type")
	}

	if term != "" {
		var results *player.Results
		var err error
		s.view = "body"
		switch searchType {
		case "album":
			results, err = s.client.FindAlbum(term, s.body.height*5)
			results.Print = results.PrintArtist()
		case "artist":
			results, err = s.client.FindArtist(term, s.body.height*5)
			results.Print = results.PrintArtists()
		case "track":
			results, err = s.client.FindTrack(term, s.body.height*5)
			results.Print = results.PrintTracks()
		}
		if err != nil {
			return err
		}

		s.stack.clear()
		s.body.cursor = 0
		s.body.newResults(results)
		s.header.header = results.Header
		s.stack.add(results, s.body.cursor)
		return g.DeleteView("search")
	}
	return nil
}

func (s *screen) pause(g *ui.Gui, v *ui.View) error {
	s.play.doPause()
	return nil
}

func (s *screen) escape(g *ui.Gui, v *ui.View) error {
	s.stack.pop()
	if s.stack.len() == 0 {
		s.body.clear()
		s.play.clear()
		s.buffer.clear()
		s.header.clear()
		s.view = "search-type"
		return nil
	}
	r, c := s.stack.top()
	s.body.newResults(r)
	s.body.cursor = c
	s.header.header = r.Header
	return nil
}

func (s *screen) quit(g *ui.Gui, v *ui.View) error {
	return ui.ErrQuit
}

func (s *screen) showSearchDialog(g *ui.Gui, v *ui.View) error {
	s.view = "search-type"
	return nil
}

func (s *screen) showHistoryDialog(g *ui.Gui, v *ui.View) error {
	s.view = "history-type"
	return nil
}

func (s *screen) getLayout(width, height int) func(*ui.Gui) error {
	if !s.client.Ping() {
		s.view = "login"
	} else {
		res, err := s.client.History(0, s.height*10, s.historySort)
		if err == nil && len(res.Results) > 0 {
			res.Print = res.PrintHistory()
			s.header.header = res.Header
			s.body.newResults(res)
			s.stack.add(res, 0)
			s.view = "body"
		} else {
			s.view = "search-type"
		}
	}

	return func(g *ui.Gui) error {
		g.Cursor = false
		if s.view == "login" {
			v, err := g.SetView("login", s.login.coords.x1, s.login.coords.y1, s.login.coords.x2, s.login.coords.y2)
			if err != nil && err != ui.ErrUnknownView {
				return err
			}

			ui.DefaultEditor = s.login
			v.Editable = true
			v.Frame = true
			v.Title = s.login.title
		} else if s.view == "help" {
			if err := s.help.show(g, s.keys); err != nil {
				return err
			}
		} else if s.view == "history-type" {
			v, err := g.SetView(s.view, s.history.coords.x1, s.history.coords.y1, s.history.coords.x2, s.history.coords.y2)
			if err != nil && err != ui.ErrUnknownView {
				return err
			}

			if err := s.history.render(g, v); err != nil {
				return err
			}
		} else if s.view == "search-type" || s.view == "search" {
			if s.view == "search" {
				g.Cursor = true
			}
			v, err := g.SetView(s.view, s.search.coords.x1, s.search.coords.y1, s.search.coords.x2, s.search.coords.y2)
			if err != nil && err != ui.ErrUnknownView {
				return err
			}

			if err := s.search.render(g, v); err != nil {
				return err
			}
		} else {
			g.DeleteView("search")
			g.DeleteView("search-type")
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
			v.Editable = false
			if err := s.body.render(g, v); err != nil {
				return err
			}

			v, err = g.SetView("play", s.play.coords.x1, s.play.coords.y1, s.play.coords.x2, s.play.coords.y2)
			if err != nil && err != ui.ErrUnknownView {
				return err
			}
			v.Frame = false
			v.Editable = true

			v, err = g.SetView("buffer", s.buffer.coords.x1, s.buffer.coords.y1, s.buffer.coords.x2, s.buffer.coords.y2)
			if err != nil && err != ui.ErrUnknownView {
				return err
			}
			v.Frame = false
			v.Editable = true
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

type stack struct {
	//current results
	topR *player.Results
	//current cursor
	topC    int
	stack   []player.Results
	cursors []int
}

func (s *stack) len() int {
	return len(s.stack)
}

func (s *stack) top() (*player.Results, int) {
	return s.topR, s.topC
}

func (s *stack) add(r *player.Results, c int) {
	s.topR = r
	s.topC = c
	s.stack = append(s.stack, *r)
	s.cursors = append(s.cursors, c)
}

func (s *stack) clear() {
	s.topR = nil
	s.topC = 0
	s.cursors = []int{}
	s.stack = []player.Results{}
}

func (s *stack) pop() {
	if len(s.stack) == 0 {
		return
	}

	i := len(s.stack) - 1
	s.stack = s.stack[0:i]
	c := s.cursors[len(s.cursors)-1]
	s.cursors = s.cursors[0:i]
	if len(s.stack) == 0 {
		s.topR = nil
		s.topC = c
	} else {
		s.topR = &s.stack[len(s.stack)-1]
		s.topC = c
	}
}
