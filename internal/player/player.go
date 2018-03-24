package player

import (
	"fmt"
	"io"
	"os"
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

func (r *Results) PrintPlaylists(w io.Writer, res Result) error {
	_, err := fmt.Fprintf(w, r.Fmt, res.Album.Title)
	return err
}

func (r *Results) PrintAlbum(w io.Writer, res Result) error {
	d := time.Duration(res.Track.Duration) * time.Second
	_, err := fmt.Fprintf(w, r.Fmt, res.Track.Title, d)
	return err
}

func (r *Results) PrintArtist(w io.Writer, res Result) error {
	_, err := fmt.Fprintf(w, r.Fmt, res.Album.Title, res.Artist.Name)
	return err
}

func (r *Results) PrintAlbumTracks(w io.Writer, res Result) error {
	d := time.Duration(res.Track.Duration) * time.Second
	_, err := fmt.Fprintf(w, r.Fmt, res.Track.Title, fmtDuration(d))
	return err
}

func (r *Results) PrintArtists(w io.Writer, res Result) error {
	_, err := fmt.Fprintf(w, r.Fmt, res.Artist.Name)
	return err
}

func (r *Results) PrintTracks(w io.Writer, res Result) error {
	_, err := fmt.Fprintf(w, r.Fmt, res.Track.Title, res.Album.Title, res.Artist.Name)
	return err
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
