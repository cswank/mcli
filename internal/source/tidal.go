package source

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/the5heepdev/tidal"
)

type Tidal struct {
	client *tidal.Tidal
}

func getTidalPath() string {
	return fmt.Sprintf("%s/.music/tidal.json", os.Getenv("HOME"))
}

func saveTidal(t *tidal.Tidal) error {
	f, err := os.Create(getTidalPath())
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(t)
}

func getTidal() (*tidal.Tidal, error) {
	f, err := os.Open(getTidalPath())
	if err != nil {
		return nil, err
	}

	defer f.Close()
	var t tidal.Tidal
	return &t, json.NewDecoder(f).Decode(&t)
}

func GetTidal() (*Tidal, error) {
	t, err := getTidal()
	return &Tidal{client: t}, err
}

func NewTidal(username, pw string) (*Tidal, error) {
	var t *tidal.Tidal
	t, err := getTidal()
	if err != nil {
		t = tidal.New(username, pw)
	}

	if t.SessionID == "" {
		return nil, fmt.Errorf("couldn't log into tidal")
	}

	if err := saveTidal(t); err != nil {
		return nil, err
	}

	return &Tidal{client: t}, nil
}

func (t *Tidal) FindArtist(term string, limit int) (*Results, error) {
	return nil, nil
}

func (t *Tidal) GetTrack(id string) string {
	return t.client.GetStreamURL(id, "LOSSLESS")
}

func (t *Tidal) GetAlbum(id string) (*Results, error) {
	tracks := t.client.GetAlbumTracks(id)
	out := make([]Result, len(tracks))
	var maxTitle int

	for i, t := range tracks {
		artists := make([]string, len(t.Artists))
		for i, a := range t.Artists {
			artists[i] = a.Name
		}
		as := strings.Join(artists, ", ")
		if len(t.Title) > maxTitle {
			maxTitle = len(t.Title)
		}
		out[i] = Result{
			Artist: as,
			Title:  t.Title,
			ID:     fmt.Sprintf("%s", t.ID),
		}
	}

	f := fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)
	return &Results{
		Header: fmt.Sprintf(f, "Title", "Artist"),
		Type:   "album",
		Print: func(w io.Writer, r Result) error {
			_, err := fmt.Fprintf(w, f, r.Title, r.Artist)
			return err
		},
		Results: out,
	}, nil
}

func (t *Tidal) FindAlbum(term string, limit int) (*Results, error) {
	albums := t.client.SearchAlbums(term, fmt.Sprintf("%d", limit))
	out := make([]Result, len(albums))
	var maxTitle int
	for i, a := range albums {
		artists := make([]string, len(a.Artists))
		for i, a := range a.Artists {
			artists[i] = a.Name
		}
		as := strings.Join(artists, ", ")
		if len(a.Title) > maxTitle {
			maxTitle = len(a.Title)
		}
		out[i] = Result{
			Artist: as,
			Title:  a.Title,
			ID:     fmt.Sprintf("%s", a.ID),
		}
	}

	f := fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)
	return &Results{
		Header: fmt.Sprintf(f, "Title", "Artist"),
		Type:   "album search",
		Print: func(w io.Writer, r Result) error {
			_, err := fmt.Fprintf(w, f, r.Title, r.Artist)
			return err
		},
		Results: out,
	}, nil
}

func (t *Tidal) FindTrack(term string, limit int) (*Results, error) {
	return nil, nil
}
