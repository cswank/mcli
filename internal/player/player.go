package player

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type Client interface {
	Player
	Fetcher
}

type Player interface {
	Play(Result)
	PlayAlbum(*Results)
	Volume(float64)
	Pause()
	FastForward()
	Rewind()
	Queue() *Results
	RemoveFromQueue(int)
	NextSong(func(Result))
	PlayProgress(func(Progress))
	DownloadProgress(func(Progress))
	History(int, int, Sort) (*Results, error)
	Done()
}

type Fetcher interface {
	Name() string
	Login(string, string) error
	Ping() bool
	AlbumLink() string
	FindArtist(string, int) (*Results, error)
	FindAlbum(string, int) (*Results, error)
	FindTrack(string, int) (*Results, error)
	GetAlbum(string) (*Results, error)
	GetTrack(string) (string, error)
	GetArtistAlbums(string, int) (*Results, error)
	GetPlaylists() (*Results, error)
	GetPlaylist(string, int) (*Results, error)
}

type Progress struct {
	N     int
	Total int
}

type Track struct {
	ID       string
	Title    string
	Duration int
}

type Artist struct {
	ID   string
	Name string
}

type Album struct {
	ID    string
	Title string
}

type Result struct {
	Service   string
	Path      string
	PlayCount int
	Track     Track
	Artist    Artist
	Album     Album
	Playlist  Album
}

type Results struct {
	Type    string
	Header  string
	Results []Result
	Fmt     string
	Print   func(io.Writer, Result) error
}

func (r *Results) PrintPlaylists() func(w io.Writer, res Result) error {
	r.Fmt = "%s\n"
	r.Header = fmt.Sprintf(r.Fmt, "Title")
	return func(w io.Writer, res Result) error {

		_, err := fmt.Fprintf(w, r.Fmt, res.Album.Title)
		return err
	}
}

func (r *Results) PrintAlbum() func(w io.Writer, res Result) error {
	var maxTitle int
	for _, res := range r.Results {
		if len(res.Track.Title) > maxTitle {
			maxTitle = len(res.Track.Title)
		}
	}

	r.Fmt = fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)
	r.Header = fmt.Sprintf(r.Fmt, "Title", "Length")

	return func(w io.Writer, res Result) error {
		d := time.Duration(res.Track.Duration) * time.Second
		_, err := fmt.Fprintf(w, r.Fmt, res.Track.Title, d)
		return err
	}
}

func (r *Results) PrintArtist() func(w io.Writer, res Result) error {
	var maxTitle int

	for _, res := range r.Results {
		if len(res.Album.Title) > maxTitle {
			maxTitle = len(res.Album.Title)
		}
	}

	r.Fmt = fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)
	r.Header = fmt.Sprintf(r.Fmt, "Title", "Artist")
	return func(w io.Writer, res Result) error {
		_, err := fmt.Fprintf(w, r.Fmt, res.Album.Title, res.Artist.Name)
		return err
	}
}

func (r *Results) PrintAlbumTracks() func(w io.Writer, res Result) error {
	var maxTitle int
	for _, res := range r.Results {
		if len(res.Track.Title) > maxTitle {
			maxTitle = len(res.Track.Title)
		}
	}

	r.Fmt = fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)
	r.Header = fmt.Sprintf(r.Fmt, "Title", "Length")
	return func(w io.Writer, res Result) error {
		d := time.Duration(res.Track.Duration) * time.Second
		_, err := fmt.Fprintf(w, r.Fmt, res.Track.Title, fmtDuration(d))
		return err
	}
}

func (r *Results) PrintArtists() func(w io.Writer, res Result) error {
	var max int
	for _, res := range r.Results {
		if len(res.Artist.Name) > max {
			max = len(res.Artist.Name)
		}
	}

	r.Fmt = "%s\n"
	r.Header = fmt.Sprintf(r.Fmt, "Artist")
	return func(w io.Writer, res Result) error {
		_, err := fmt.Fprintf(w, r.Fmt, res.Artist.Name)
		return err
	}
}

func (r *Results) PrintTracks() func(w io.Writer, res Result) error {
	var maxTitle int
	var maxAlbum int

	for _, res := range r.Results {
		if len(res.Track.Title) > maxTitle {
			maxTitle = len(res.Track.Title)
		}
		if len(res.Album.Title) > maxAlbum {
			maxAlbum = len(res.Album.Title)
		}
	}

	r.Fmt = fmt.Sprintf("%%-%ds%%-%ds%%s\n", maxTitle+4, maxAlbum)
	r.Header = fmt.Sprintf(r.Fmt, "Title", "Album", "Artist")
	return func(w io.Writer, res Result) error {
		_, err := fmt.Fprintf(w, r.Fmt, res.Track.Title, res.Album.Title, res.Artist.Name)
		return err
	}
}

func (r *Results) PrintHistory() func(w io.Writer, res Result) error {
	col := 40
	r.Fmt = fmt.Sprintf("%%-%ds%%-%ds%%-%ds%%s\n", col+4, col+4, col+4)
	r.Header = fmt.Sprintf(r.Fmt, "Title", "Album", "Artist", "Played")

	return func(w io.Writer, res Result) error {
		_, err := fmt.Fprintf(w, r.Fmt, truncate(res.Track.Title, col), truncate(res.Album.Title, col), truncate(res.Artist.Name, col), strconv.Itoa(res.PlayCount))
		return err
	}
}

func truncate(s string, l int) string {
	if len(s) < l {
		l = len(s)
	}
	return s[:l]
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%d:%02d", m, s)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
