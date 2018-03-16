package player

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Client interface {
	Player
	Fetcher
}

type Player interface {
	Play(Result)
	PlayAlbum([]Result)
	Volume(float64)
	Pause()
	FastForward()
	Rewind()
	Queue() []Result
	RemoveFromQueue(int)
	NextSong(func(Result))
	PlayProgress(func(Progress))
	DownloadProgress(func(Progress))
	History(int, int) (*Results, error)
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
	Service  string
	Track    Track
	Artist   Artist
	Album    Album
	Playlist Album
	Path     string
}

func (r *Result) ToCSV() []string {
	return []string{
		time.Now().Format(time.RFC3339),
		r.Service,
		r.Track.ID,
		r.Track.Title,
		strconv.Itoa(r.Track.Duration),
		r.Album.ID,
		r.Album.Title,
		r.Artist.ID,
		r.Artist.Name,
	}
}

func (r *Result) FromCSV(row []string) error {
	if len(row) < 9 {
		return fmt.Errorf("invalid history csv row: %v", row)
	}

	id := row[2]
	d, err := strconv.ParseInt(row[4], 10, 64)
	if err != nil {
		return err
	}

	*r = Result{
		Service: row[1],
		Track: Track{
			ID:       id,
			Title:    row[3],
			Duration: int(d),
		},
		Album: Album{
			ID:    row[5],
			Title: row[6],
		},
		Artist: Artist{
			ID:   row[7],
			Name: row[8],
		},
	}
	return nil
}

type Results struct {
	Type    string
	Header  string
	Results []Result
	Print   func(io.Writer, Result) error
}
