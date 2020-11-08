package repo

import (
	"path/filepath"

	"github.com/asdine/storm"
	_ "github.com/mattn/go-sqlite3"

	"github.com/cswank/mcli/internal/schema"
)

type (
	Storm struct {
		db  *storm.DB
		pth string
	}

	artist struct {
		ID   int    `storm:"unique,index"`
		Name string `storm:"unique,index"`
	}

	album struct {
		ID       int    `storm:"unique,index"`
		Name     string `storm:"index"`
		ArtistID int    `storm:"index"`
	}

	track struct {
		ID       int    `storm:"unique,index"`
		Name     string `storm:"index"`
		AlbumID  int    `storm:"index"`
		ArtistID int    `storm:"index"`
		Count    int    `storm:"index"`
		Time     string `storm:"index"`
	}
)

func NewStorm(cfg schema.Config) (*Storm, error) {
	db, err := storm.Open(filepath.Join(cfg.Home, "storm.db"))
	return &Storm{
		db:  db,
		pth: cfg.Pth,
	}, err
}

func (s Storm) FindArtist(term string, n int) (*schema.Results, error) {
	return nil, nil
}

func (s Storm) FindAlbum(term string, n int) (*schema.Results, error) {
	return nil, nil
}

func (s Storm) FindTrack(term string, n int) (*schema.Results, error) {
	return nil, nil
}

func (s Storm) GetAlbum(id int64) (*schema.Results, error) {
	return nil, nil
}

func (s Storm) GetArtistAlbums(id int64, n int) (*schema.Results, error) {
	return nil, nil
}

func (s Storm) GetArtistTracks(id int64, n int) (*schema.Results, error) {
	return nil, nil
}

func (s Storm) GetPlaylists() (*schema.Results, error) {
	return &schema.Results{}, nil
}

func (s Storm) GetPlaylist(int64, int) (*schema.Results, error) {
	return &schema.Results{}, nil
}

func (s *Storm) Close() error {
	return nil
}

func (s *Storm) Save(res schema.Result) error {
	return nil
}

func (s *Storm) Fetch(page, pageSize int, sortTerm Sort) (*schema.Results, error) {
	var entries []Track
	err := s.db.Select().OrderBy(string(sortTerm)).Reverse().Limit(pageSize).Skip(page * pageSize).Find(&entries)
	if err != nil {
		return nil, err
	}

	out := make([]schema.Result, len(entries))
	for i, e := range entries {
		e.Result.PlayCount = e.Count
		out[i] = e.Result
	}

	return &schema.Results{
		Type:    "history",
		Results: out,
	}, nil
}

func (s Storm) doFind(q string, term interface{}, t string) (*schema.Results, error) {

	return &schema.Results{
		Type:    t,
		Fmt:     "",
		Results: nil,
	}, nil
}

func (s *Storm) Track(id int64) (string, error) {
	return "", nil
}

func (s Storm) Init() error {
	return nil
}

func (s Storm) InsertOrGetArtist(name string) (int, error) {
	var a artist
	return a.ID, s.insertOrGet(name, &a)
}

func (s Storm) InsertOrGetAlbum(name string, artistID int) (int, error) {
	a := album{ArtistID: artistID}
	return a.ID, s.insertOrGet(name, &a)
}

func (s Storm) InsertOrGetTrack(name string, albumID int) (int, error) {
	t := track{AlbumID: albumID}
	return t.ID, s.insertOrGet(name, &t)
}

func (s Storm) insertOrGet(name string, item interface{}) error {
	err := s.db.One("Name", name, item)
	if err == storm.ErrNotFound {
		if err := s.db.Save(&item); err != nil {
			return err
		}
		return s.insertOrGet(name, item)
	}

	return err
}
