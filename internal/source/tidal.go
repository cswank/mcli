package source

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cswank/tidal"
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
		t, err = tidal.New(username, pw)
		if err != nil {
			return nil, err
		}
	}

	if t.SessionID == "" {
		return nil, fmt.Errorf("couldn't log into tidal")
	}

	if err := saveTidal(t); err != nil {
		return nil, err
	}

	return &Tidal{client: t}, nil
}

func (t *Tidal) Name() string {
	return "tidal"
}

func (t *Tidal) FindArtist(term string, limit int) (*Results, error) {
	artists, err := t.client.SearchArtists(term, fmt.Sprintf("%d", limit))
	if err != nil {
		return nil, err
	}
	out := make([]Result, len(artists))
	var max int
	for i, a := range artists {
		if len(a.Name) > max {
			max = len(a.Name)
		}
		out[i] = Result{
			Artist: Artist{
				ID:   fmt.Sprintf("%s", a.ID),
				Name: a.Name,
			},
		}
	}

	f := "%s\n"
	return &Results{
		Header: fmt.Sprintf(f, "Artist"),
		Type:   "artist search",
		Print: func(w io.Writer, r Result) error {
			_, err := fmt.Fprintf(w, f, r.Artist.Name)
			return err
		},
		Results: out,
	}, nil
}

func (t *Tidal) GetArtistAlbums(id string, limit int) (*Results, error) {
	albums, err := t.client.GetArtistAlbums(id, fmt.Sprintf("%d", limit))
	if err != nil {
		return nil, err
	}
	out := make([]Result, len(albums))
	var max int
	for i, a := range albums {
		if len(a.Title) > max {
			max = len(a.Title)
		}
		artists := make([]string, len(a.Artists))
		for i, a := range a.Artists {
			artists[i] = a.Name
		}
		as := strings.Join(artists, ", ")
		if len(a.Title) > max {
			max = len(a.Title)
		}
		out[i] = Result{
			Artist: Artist{
				ID:   a.Artists[0].ID.String(),
				Name: as,
			},
			Album: Album{
				ID:    fmt.Sprintf("%s", a.ID),
				Title: a.Title,
			},
		}
	}

	f := fmt.Sprintf("%%-%ds%%s\n", max+4)
	return &Results{
		Header: fmt.Sprintf(f, "Album", "Artist"),
		Type:   "album search",
		Print: func(w io.Writer, r Result) error {
			_, err := fmt.Fprintf(w, f, r.Album.Title, r.Artist.Name)
			return err
		},
		Results: out,
	}, nil
}

func (t *Tidal) GetTrack(id string) (string, error) {
	return t.client.GetStreamURL(id, "LOSSLESS")
}

func (t *Tidal) GetAlbum(id string) (*Results, error) {
	tracks, err := t.client.GetAlbumTracks(id)
	if err != nil {
		return nil, err
	}
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
		dur, _ := t.Duration.Int64()
		out[i] = Result{
			Artist: Artist{
				Name: as,
				ID:   t.Artists[0].ID.String(),
			},
			Album: Album{
				ID:    t.Album.ID.String(),
				Title: t.Album.Title,
			},
			Track: Track{
				ID:       fmt.Sprintf("%s", t.ID),
				Title:    t.Title,
				Duration: int(dur),
			},
		}
	}

	f := fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)
	return &Results{
		Header: fmt.Sprintf(f, "Title", "Artist"),
		Type:   "album",
		Print: func(w io.Writer, r Result) error {
			_, err := fmt.Fprintf(w, f, r.Track.Title, r.Artist.Name)
			return err
		},
		Results: out,
	}, nil
}

func (t *Tidal) FindAlbum(term string, limit int) (*Results, error) {
	albums, err := t.client.SearchAlbums(term, fmt.Sprintf("%d", limit))
	if err != nil {
		return nil, err
	}
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
			Artist: Artist{
				ID:   a.Artists[0].ID.String(),
				Name: as,
			},
			Album: Album{
				ID:    a.ID.String(),
				Title: a.Title,
			},
		}
	}

	f := fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)
	return &Results{
		Header: fmt.Sprintf(f, "Title", "Artist"),
		Type:   "album search",
		Print: func(w io.Writer, r Result) error {
			_, err := fmt.Fprintf(w, f, r.Album.Title, r.Artist.Name)
			return err
		},
		Results: out,
	}, nil
}

func (t *Tidal) FindTrack(term string, limit int) (*Results, error) {
	tracks, err := t.client.SearchTracks(term, fmt.Sprintf("%d", limit))
	if err != nil {
		return nil, err
	}
	out := make([]Result, len(tracks))
	var maxTitle int
	var maxAlbum int
	for i, t := range tracks {
		artists := make([]string, len(t.Artists))
		for i, a := range t.Artists {
			artists[i] = a.Name
		}
		as := strings.Join(artists, ", ")
		if len(t.Title) > maxTitle {
			maxTitle = len(t.Title)
		}
		if len(t.Album.Title) > maxAlbum {
			maxAlbum = len(t.Album.Title)
		}
		dur, _ := t.Duration.Int64()
		out[i] = Result{
			Artist: Artist{
				ID:   t.Artists[0].ID.String(),
				Name: as,
			},
			Album: Album{
				ID:    t.Album.ID.String(),
				Title: t.Album.Title,
			},
			Track: Track{
				Title:    t.Title,
				ID:       fmt.Sprintf("%s", t.ID),
				Duration: int(dur),
			},
		}
	}

	f := fmt.Sprintf("%%-%ds%%-%ds%%s\n", maxTitle+4, maxAlbum)
	return &Results{
		Header: fmt.Sprintf(f, "Title", "Album", "Artist"),
		Type:   "album",
		Print: func(w io.Writer, r Result) error {
			_, err := fmt.Fprintf(w, f, r.Track.Title, r.Album.Title, r.Artist.Name)
			return err
		},
		Results: out,
	}, nil
}
