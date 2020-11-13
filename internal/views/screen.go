package views

import (
	"fmt"
	"time"

	ui "github.com/awesome-gocui/gocui"
	"github.com/cswank/mcli/internal/repo"
	"github.com/cswank/mcli/internal/schema"
)

type screen struct {
	view   string
	width  int
	height int
	id     string

	header       *header
	body         *body
	play         *player
	buffer       *buffer
	search       *search
	artistDialog *artistDialog
	history      *history
	volume       *volume
	historySort  repo.Sort
	help         *help

	keys []key

	client *client
	stack  stack

	volumeEvent chan bool
}

func newScreen(width, height int, cli *client) (*screen, error) {
	id := randString(10)
	s := &screen{
		id:          id,
		client:      cli,
		view:        "body",
		width:       width,
		height:      height,
		historySort: repo.Time,
		play:        newPlay(width, height, id, cli),
		body:        newBody(width, height, cli.AlbumLink()),
		buffer:      newBuffer(width, height, id, cli),
		header:      newHeader(width, height),
		help:        newHelp(width, height),
		volume:      newVolume(width, height, cli.Volume(0.0)),
		volumeEvent: make(chan bool),
	}

	go s.clearVolume()
	s.search = newSearch(width, height, s.doSearch)
	s.artistDialog = newArtistDialog(width, height, s.doShowArtist)
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
		results, err := s.client.GetArtistAlbums(r.Artist.ID, 0, s.height)
		if err != nil {
			return err
		}
		results.Print = results.PrintArtist()
		results.Page = func(p int) (*schema.Results, error) {
			return s.client.GetArtistAlbums(r.Artist.ID, p, s.body.height)
		}

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
		results.Page = func(p int) (*schema.Results, error) {
			return results, err
		}
		s.body.newResults(results)
		s.header.header = results.Header
		s.stack.add(results, c)
	case "playlists":
		s.body.cursor = 0
		results, err := s.client.GetPlaylist(r.Album.ID, 0, s.body.height)
		if err != nil {
			return err
		}
		results.Print = results.PrintAlbum()
		results.Page = func(p int) (*schema.Results, error) {
			return s.client.GetPlaylist(r.Album.ID, p, s.body.height)
		}
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

func (s *screen) importMusic(g *ui.Gui, v *ui.View) error {
	s.play.emptyQueue()
	s.play.client.Import(func(p schema.Progress) {
		s.play.playProgress(p)
	})
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
	results := &schema.Results{
		Header: fmt.Sprintf(f, "Title", "Album", "Artist"),
		Type:   "queue",
		Print: func(r schema.Result) string {
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
	if r.Album.ID == 0 && r.Artist.ID != 0 {
		return s.doShowArtist(r.Artist.ID, "albums")
	}
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
	s.view = "artist-dialog"
	s.artistDialog.selected = r
	return nil
}

func (s *screen) goToArtistTracks(g *ui.Gui, v *ui.View) error {
	r := s.body.view[s.body.cursor]
	c := s.body.cursor
	results, err := s.client.GetArtistTracks(r.Artist.ID, 0, s.height)
	if err != nil {
		return err
	}

	results.Print = results.PrintArtistTracks()
	results.Page = func(p int) (*schema.Results, error) {
		return s.client.GetArtistTracks(r.Artist.ID, p, s.body.height)
	}
	s.body.newResults(results)
	s.header.header = results.Header
	s.stack.add(results, c)
	s.body.cursor = 0
	return nil
}

func (s *screen) volumeUp(g *ui.Gui, v *ui.View) error {
	s.volume.vol = s.play.volume(0.5)
	s.view = "volume"
	s.volumeEvent <- true
	return nil
}

func (s *screen) volumeDown(g *ui.Gui, v *ui.View) error {
	s.volume.vol = s.play.volume(-0.5)
	s.view = "volume"
	s.volumeEvent <- true
	return nil
}

func (s *screen) clearVolume() {
	after := time.After(1000000 * time.Second)
	for {
		select {
		case <-after:
			s.view = "body"
			g.Update(func(g *ui.Gui) error {
				return s.volume.clear()
			})
			after = time.After(1000000 * time.Second)
		case <-s.volumeEvent:
			after = time.After(3 * time.Second)
		}
	}
}

func (s *screen) next(g *ui.Gui, v *ui.View) error {
	s.play.next()
	s.buffer.clear()
	return nil
}

func (s *screen) rewind(g *ui.Gui, v *ui.View) error {
	s.play.rewind()
	s.buffer.clear()
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

func (s *screen) showManual(g *ui.Gui, v *ui.View) error {
	s.view = "manual"
	return nil
}

func (s *screen) hideManual(g *ui.Gui, v *ui.View) error {
	s.view = "body"
	return s.help.hide(g, v)
}

func (s *screen) showHistory(sort repo.Sort) error {
	g.DeleteView("history-type")
	s.historySort = sort
	res, err := s.client.Fetch(0, s.body.height, sort)
	if err != nil {
		return err
	}

	s.view = "body"
	res.Print = res.PrintHistory()
	res.Page = func(p int) (*schema.Results, error) {
		return s.client.Fetch(p, s.body.height, sort)
	}
	s.header.header = res.Header
	s.body.newResults(res)
	s.stack.add(res, s.body.cursor)
	s.body.cursor = 0
	return nil
}

func (s *screen) doShowArtist(id int64, term string) error {
	s.view = "body"
	c := s.body.cursor
	if term == "albums" {
		s.body.cursor = 0
		results, err := s.client.GetArtistAlbums(id, 0, s.body.height)
		if err != nil {
			return err
		}
		results.Print = results.PrintArtist()
		results.Page = func(p int) (*schema.Results, error) {
			return s.client.GetArtistAlbums(id, p, s.body.height)
		}
		s.body.newResults(results)
		s.header.header = results.Header
		s.stack.add(results, c)
		return nil
	}

	//must be tracks
	s.body.cursor = 0
	results, err := s.client.GetArtistTracks(id, 0, s.height)
	if err != nil {
		return err
	}

	results.Print = results.PrintTracks()
	results.Page = func(p int) (*schema.Results, error) {
		return s.client.GetArtistTracks(id, p, s.body.height)
	}
	s.body.newResults(results)
	s.header.header = results.Header
	s.stack.add(results, c)
	return nil
}

func (s *screen) doSearch(searchType, term string) error {
	if searchType != "" && term == "" {
		s.view = "search"
		return g.DeleteView("search-type")
	}

	if term != "" {
		var results *schema.Results
		var err error
		s.view = "body"
		switch searchType {
		case "album":
			results, err = s.client.FindAlbum(term, 0, s.body.height)
			results.Print = results.PrintArtist()
			results.Page = func(p int) (*schema.Results, error) {
				return s.client.FindAlbum(term, p, s.body.height)
			}
		case "artist":
			results, err = s.client.FindArtist(term, 0, s.body.height)
			results.Print = results.PrintArtists()
			results.Page = func(p int) (*schema.Results, error) {
				return s.client.FindArtist(term, p, s.body.height)
			}
		case "track":
			results, err = s.client.FindTrack(term, 0, s.body.height)
			results.Print = results.PrintTracks()
			results.Page = func(p int) (*schema.Results, error) {
				return s.client.FindTrack(term, p, s.body.height)
			}
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
		// show error message
	} else {
		res, err := s.client.Fetch(0, s.body.height, s.historySort)
		if err == nil && len(res.Results) > 0 {
			res.Print = res.PrintHistory()
			res.Page = func(p int) (*schema.Results, error) {
				return s.client.Fetch(p, s.body.height, s.historySort)
			}
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
		switch s.view {
		case "help":
			if err := s.help.show(g, s.keys); err != nil {
				return err
			}
		case "volume":
			v, err := g.SetView(s.view, s.volume.coords.x1, s.volume.coords.y1, s.volume.coords.x2, s.volume.coords.y2, 0)
			if err != nil && !ui.IsUnknownView(err) {
				return err
			}

			if err := s.volume.render(g, v); err != nil {
				return err
			}
		case "history-type":
			v, err := g.SetView(s.view, s.history.coords.x1, s.history.coords.y1, s.history.coords.x2, s.history.coords.y2, 0)
			if err != nil && !ui.IsUnknownView(err) {
				return err
			}

			if err := s.history.render(g, v); err != nil {
				return err
			}
		case "search-type", "search":
			if s.view == "search" {
				g.Cursor = true
			}
			v, err := g.SetView(s.view, s.search.coords.x1, s.search.coords.y1, s.search.coords.x2, s.search.coords.y2, 0)
			if err != nil && !ui.IsUnknownView(err) {
				return err
			}

			if err := s.search.render(g, v); err != nil {
				return err
			}
		case "artist-dialog":
			v, err := g.SetView(s.view, s.artistDialog.coords.x1, s.artistDialog.coords.y1, s.artistDialog.coords.x2, s.artistDialog.coords.y2, 0)
			if err != nil && !ui.IsUnknownView(err) {
				return err
			}

			if err := s.artistDialog.render(g, v); err != nil {
				return err
			}
		default:
			g.DeleteView("search")
			g.DeleteView("search-type")
			g.DeleteView("artist-dialog")
			v, err := g.SetView("header", s.header.coords.x1, s.header.coords.y1, s.header.coords.x2, s.header.coords.y2, 0)
			if err != nil && !ui.IsUnknownView(err) {
				return err
			}

			v.Frame = false
			if err := s.header.render(g, v); err != nil {
				return err
			}

			v, err = g.SetView("body", s.body.coords.x1, s.body.coords.y1, s.body.coords.x2, s.body.coords.y2, 0)
			if err != nil && !ui.IsUnknownView(err) {
				return err
			}

			v.Frame = false
			v.Editable = false
			if err := s.body.render(g, v); err != nil {
				return err
			}

			v, err = g.SetView("play", s.play.coords.x1, s.play.coords.y1, s.play.coords.x2, s.play.coords.y2, 0)
			if err != nil && !ui.IsUnknownView(err) {
				return err
			}
			v.Frame = false
			v.Editable = true

			v, err = g.SetView("buffer", s.buffer.coords.x1, s.buffer.coords.y1, s.buffer.coords.x2, s.buffer.coords.y2, 0)
			if err != nil && !ui.IsUnknownView(err) {
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
	topR *schema.Results
	//current cursor
	topC    int
	stack   []schema.Results
	cursors []int
}

func (s *stack) len() int {
	return len(s.stack)
}

func (s *stack) top() (*schema.Results, int) {
	return s.topR, s.topC
}

func (s *stack) add(r *schema.Results, c int) {
	s.topR = r
	s.topC = c
	s.stack = append(s.stack, *r)
	s.cursors = append(s.cursors, c)
}

func (s *stack) clear() {
	s.topR = nil
	s.topC = 0
	s.cursors = []int{}
	s.stack = []schema.Results{}
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
