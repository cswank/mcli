package schema

import (
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"
)

type Config struct {
	Addr       string
	Pth        string
	Home       string
	Log        string
	DB         string
	RemotePlay bool
	Speakers   bool
}

func NewConfig(addr, pth, home, log, db string, remotePlay bool, speakers *bool) Config {
	s := speakers != nil && *speakers
	return Config{
		Addr:       addr,
		Pth:        pth,
		Home:       home,
		Log:        log,
		DB:         db,
		RemotePlay: remotePlay,
		Speakers:   s,
	}
}

type Results struct {
	Album   Album                         `json:"album"`
	Type    string                        `json:"type"`
	Header  string                        `json:"header"`
	Results []Result                      `json:"results"`
	Error   string                        `json:"error"`
	Fmt     string                        `json:"fmt"`
	Print   func(Result) string           `json:"-" template:"-"`
	Page    func(i int) (*Results, error) `json:"-" template:"-"`
}

type Progress struct {
	N       int    `json:"n"`
	Total   int    `json:"total"`
	Payload []byte `json:"payload"`
}

type Track struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

type Artist struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Album struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type Result struct {
	Service   string `json:"service"`
	Path      string `json:"path"`
	PlayCount int    `json:"play_count"`
	Track     Track  `json:"track"`
	Artist    Artist `json:"artist"`
	Album     Album  `json:"album"`
	Playlist  Album  `json:"playlist"`
	Error     string `json:"message"`
}

func (r *Results) PrintPlaylists() func(res Result) string {
	r.Fmt = "%s\n"
	r.Header = fmt.Sprintf(r.Fmt, "Title")
	return func(res Result) string {
		return fmt.Sprintf(r.Fmt, res.Album.Title)
	}
}

func (r *Results) PrintAlbum() func(res Result) string {
	var maxTitle int
	for _, res := range r.Results {
		c := utf8.RuneCountInString(res.Track.Title)
		if c > maxTitle {
			maxTitle = c
		}
	}

	r.Fmt = fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)
	r.Header = fmt.Sprintf(r.Fmt, "Title", "Length")

	return func(res Result) string {
		d := time.Duration(res.Track.Duration) * time.Second
		return fmt.Sprintf(r.Fmt, res.Track.Title, d)
	}
}

func (r *Results) PrintArtist() func(res Result) string {
	var maxTitle int

	for _, res := range r.Results {
		c := utf8.RuneCountInString(res.Album.Title)
		if c > maxTitle && c < 50 {
			maxTitle = c
		}
	}

	r.Fmt = fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)
	r.Header = fmt.Sprintf(r.Fmt, "Title", "Artist")
	return func(res Result) string {
		title := res.Album.Title
		end := utf8.RuneCountInString(title)
		if end > 50 {
			end = 48
		}
		return fmt.Sprintf(r.Fmt, title[:end], res.Artist.Name)
	}
}

func (r *Results) PrintAlbumTracks() func(res Result) string {
	var maxTitle int
	for _, res := range r.Results {
		c := utf8.RuneCountInString(res.Track.Title)
		if c > maxTitle {
			maxTitle = c
		}
	}

	r.Fmt = fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)
	r.Header = fmt.Sprintf(r.Fmt, "Title", "Length")
	return func(res Result) string {
		d := time.Duration(res.Track.Duration) * time.Second
		return fmt.Sprintf(r.Fmt, res.Track.Title, FmtDuration(d))
	}
}

func (r *Results) PrintArtists() func(res Result) string {
	var max int
	for _, res := range r.Results {
		c := utf8.RuneCountInString(res.Artist.Name)
		if c > max {
			max = c
		}
	}

	r.Fmt = "%s\n"
	r.Header = fmt.Sprintf(r.Fmt, "Artist")
	return func(res Result) string {
		return fmt.Sprintf(r.Fmt, res.Artist.Name)
	}
}

func (r *Results) PrintTracks() func(res Result) string {
	var maxTitle int
	var maxAlbum int

	for _, res := range r.Results {
		c := utf8.RuneCountInString(res.Track.Title)
		if c > maxTitle {
			maxTitle = c
		}
		c = utf8.RuneCountInString(res.Album.Title)
		if c > maxAlbum {
			maxAlbum = c
		}
	}

	r.Fmt = fmt.Sprintf("%%-%ds%%-%ds%%s\n", maxTitle+4, maxAlbum)
	r.Header = fmt.Sprintf(r.Fmt, "Title", "Album", "Artist")
	return func(res Result) string {
		return fmt.Sprintf(r.Fmt, res.Track.Title, res.Album.Title, res.Artist.Name)
	}
}

func (r *Results) PrintArtistTracks() func(res Result) string {
	var maxTitle int

	for _, res := range r.Results {
		c := utf8.RuneCountInString(res.Track.Title)
		if c > maxTitle {
			maxTitle = c
		}
	}

	r.Fmt = fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)
	r.Header = fmt.Sprintf(r.Fmt, "Title", "Album")
	return func(res Result) string {
		return fmt.Sprintf(r.Fmt, res.Track.Title, res.Album.Title)
	}
}

func (r *Results) PrintHistory() func(res Result) string {
	col := 40
	r.Fmt = fmt.Sprintf("%%-%ds%%-%ds%%-%ds%%s\n", col+4, col+4, col+4)
	r.Header = fmt.Sprintf(r.Fmt, "Title", "Album", "Artist", "Played")

	return func(res Result) string {
		return fmt.Sprintf(r.Fmt, truncate(res.Track.Title, col), truncate(res.Album.Title, col), truncate(res.Artist.Name, col), strconv.Itoa(res.PlayCount))
	}
}

func truncate(s string, l int) string {
	c := utf8.RuneCountInString(s)
	if c < l {
		l = c
	}
	return s[:l]
}

func FmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%d:%02d", m, s)
}
