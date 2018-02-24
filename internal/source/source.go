package source

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

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
	Track  Track
	Artist Artist
	Album  Album
}

func (r *Result) ToCSV() []string {
	return []string{
		time.Now().Format(time.RFC3339),
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
	if len(row) < 8 {
		return fmt.Errorf("invalid history csv row: %v", row)
	}

	id := row[1]
	d, err := strconv.ParseInt(row[3], 10, 64)
	if err != nil {
		return err
	}

	*r = Result{
		Track: Track{
			ID:       id,
			Title:    row[2],
			Duration: int(d),
		},
		Album: Album{
			ID:    row[4],
			Title: row[5],
		},
		Artist: Artist{
			ID:   row[6],
			Name: row[7],
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

type Source interface {
	Name() string
	FindArtist(string, int) (*Results, error)
	FindAlbum(string, int) (*Results, error)
	FindTrack(string, int) (*Results, error)
	GetAlbum(string) (*Results, error)
	GetTrack(string) (string, error)
	GetArtistAlbums(string, int) (*Results, error)
}
