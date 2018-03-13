package player

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/cswank/tidal"
)

type Tidal struct {
	client *tidal.Tidal
}

//NewTidal returns a Client composed of a Flac player and Tidal Fetcher
func NewTidal() (Client, error) {
	t, err := newTidal()
	if err != nil {
		return nil, err
	}
	return newFlac(t)
}

func newTidal() (*Tidal, error) {
	pth := getTidalPath()
	e, err := exists(pth)
	if err != nil {
		return nil, err
	}

	if !e {
		return &Tidal{client: &tidal.Tidal{}}, nil
	}

	f, err := os.Open(pth)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var t tidal.Tidal
	err = json.NewDecoder(f).Decode(&t)

	return &Tidal{client: &t}, err
}

func getTidalPath() string {
	return fmt.Sprintf("%s/tidal.json", os.Getenv("MCLI_HOME"))
}

func (t *Tidal) save(cli *tidal.Tidal) error {
	f, err := os.Create(getTidalPath())
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(cli); err != nil {
		return err
	}
	t.client = cli
	return nil
}

func (t *Tidal) Login(username, pw string) error {
	cli, err := tidal.New(username, pw)
	if err != nil {
		return err
	}

	if cli.SessionID == "" {
		return fmt.Errorf("couldn't log into tidal")
	}

	return t.save(cli)
}

func (t *Tidal) Ping() bool {
	ok, err := t.client.CheckSession()
	return ok && err == nil
}

func (t *Tidal) Name() string {
	return "tidal"
}

func (t *Tidal) AlbumLink() string {
	return "https://listen.tidal.com/album"
}

func (t *Tidal) GetPlaylists() (*Results, error) {
	l, err := t.client.GetUserPlaylists()
	if err != nil {
		return nil, err
	}
	out := make([]Result, len(l))
	for i, item := range l {
		out[i] = Result{
			Album: Album{
				ID:    item.UUID,
				Title: item.Title,
			},
		}
	}

	f := "%s\n"
	return &Results{
		Header: fmt.Sprintf(f, "Title"),
		Type:   "playlists",
		Print: func(w io.Writer, r Result) error {
			_, err := fmt.Fprintf(w, f, r.Album.Title)
			return err
		},
		Results: out,
	}, nil
}

func (t *Tidal) GetPlaylist(id string, limit int) (*Results, error) {
	tracks, err := t.client.GetPlaylistTracks(id, strconv.Itoa(limit))
	if err != nil {
		return nil, err
	}
	return t.getTracks(tracks, "playlist")
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

		if len(a.Title) > max {
			max = len(a.Title)
		}
		out[i] = Result{
			Artist: Artist{
				ID:   a.Artists[0].ID.String(),
				Name: a.Artists[0].Name,
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
	return t.getTracks(tracks, "album")
}

func (t *Tidal) getTracks(tracks []tidal.Track, tp string) (*Results, error) {
	out := make([]Result, len(tracks))
	var maxTitle int

	for i, tr := range tracks {
		if len(tr.Title) > maxTitle {
			maxTitle = len(tr.Title)
		}
		dur, _ := tr.Duration.Int64()
		out[i] = Result{
			Service: t.Name(),
			Artist: Artist{
				Name: tr.Artists[0].Name,
				ID:   tr.Artists[0].ID.String(),
			},
			Album: Album{
				ID:    tr.Album.ID.String(),
				Title: tr.Album.Title,
			},
			Track: Track{
				ID:       fmt.Sprintf("%s", tr.ID),
				Title:    tr.Title,
				Duration: int(dur),
			},
		}
	}

	f := fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)
	return &Results{
		Header: fmt.Sprintf(f, "Title", "Length"),
		Type:   tp,
		Print: func(w io.Writer, r Result) error {
			d := time.Duration(r.Track.Duration) * time.Second
			_, err := fmt.Fprintf(w, f, r.Track.Title, fmtDuration(d))
			return err
		},
		Results: out,
	}, nil
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%d:%02d", m, s)
}

func (t *Tidal) FindAlbum(term string, limit int) (*Results, error) {
	albums, err := t.client.SearchAlbums(term, fmt.Sprintf("%d", limit))
	if err != nil {
		return nil, err
	}
	out := make([]Result, len(albums))
	var maxTitle int
	for i, a := range albums {
		if len(a.Title) > maxTitle {
			maxTitle = len(a.Title)
		}
		out[i] = Result{
			Artist: Artist{
				ID:   a.Artists[0].ID.String(),
				Name: a.Artists[0].Name,
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
				Name: t.Artists[0].Name,
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
