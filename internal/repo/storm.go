package repo

import (
	"fmt"
	"path/filepath"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	_ "github.com/mattn/go-sqlite3"

	"github.com/cswank/mcli/internal/schema"
)

type (
	Storm struct {
		db  *storm.DB
		pth string
	}

	artist struct {
		ID   int64  `storm:"unique,index,increment"`
		Name string `storm:"unique,index"`
	}

	album struct {
		ID       int64  `storm:"unique,index,increment"`
		Name     string `storm:"index"`
		ArtistID int64  `storm:"index"`
	}

	track struct {
		ID       int64  `storm:"unique,index,increment"`
		Name     string `storm:"index"`
		AlbumID  int64  `storm:"index"`
		ArtistID int64  `storm:"index"`
		Count    int64  `storm:"index"`
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

func (s Storm) FindArtist(term string, n int) ([]schema.Result, error) {
	var a []artist
	if err := s.db.Select(q.Re("Name", term)).Find(&a); err != nil {
		if err == storm.ErrNotFound {
			return nil, nil
		}

		return nil, err
	}

	out := make([]schema.Result, len(a))
	for i, art := range a {
		out[i] = schema.Result{Artist: schema.Artist{ID: art.ID, Name: art.Name}}
	}

	return out, nil
}

func (s Storm) FindAlbum(term string, n int) ([]schema.Result, error) {
	return nil, nil
}

func (s Storm) FindTrack(term string, n int) ([]schema.Result, error) {
	return nil, nil
}

func (s Storm) GetAlbum(id int64) ([]schema.Result, error) {
	var alb album
	if err := s.db.One("ID", id, &alb); err != nil {
		return nil, err
	}

	var art artist
	if err := s.db.One("ID", alb.ArtistID, &art); err != nil {
		return nil, err
	}

	var t []track
	if err := s.db.Select(q.Eq("AlbumID", id)).Find(&t); err != nil {
		if err == storm.ErrNotFound {
			return []schema.Result{}, nil
		}

		return nil, err
	}

	out := make([]schema.Result, len(t))
	for i, track := range t {
		out[i] = schema.Result{
			Track:  schema.Track{ID: track.ID, Title: track.Name},
			Album:  schema.Album{ID: alb.ID, Title: alb.Name},
			Artist: schema.Artist{ID: art.ID, Name: art.Name},
		}
	}

	return out, nil
}

func (s Storm) GetArtistAlbums(id int64, n int) ([]schema.Result, error) {
	var a []album
	if err := s.db.Select(q.Eq("ArtistID", id)).Find(&a); err != nil {
		if err == storm.ErrNotFound {
			return []schema.Result{}, nil
		}

		return nil, err
	}

	out := make([]schema.Result, len(a))
	for i, alb := range a {
		out[i] = schema.Result{Album: schema.Album{ID: alb.ID, Title: alb.Name}}
	}

	return out, nil
}

func (s Storm) GetArtistTracks(id int64, n int) ([]schema.Result, error) {
	return nil, nil
}

func (s Storm) GetPlaylists() ([]schema.Result, error) {
	return nil, nil
}

func (s Storm) GetPlaylist(int64, int) ([]schema.Result, error) {
	return nil, nil
}

func (s *Storm) Close() error {
	return nil
}

func (s *Storm) Save(res schema.Result) error {
	return nil
}

func (s *Storm) Fetch(page, pageSize int, sortTerm Sort) (*schema.Results, error) {
	var entries []track
	err := s.db.Select().OrderBy(string(sortTerm)).Reverse().Limit(pageSize).Skip(page * pageSize).Find(&entries)
	if err != nil {
		return nil, err
	}

	out := make([]schema.Result, len(entries))
	for i, t := range entries {
		out[i] = schema.Result{
			Track: schema.Track{ID: t.ID, Title: t.Name},
		}
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
	var t track
	var al album
	var ar artist

	if err := s.db.One("ID", id, &t); err != nil {
		return "", err
	}

	if err := s.db.One("ID", t.AlbumID, &al); err != nil {
		return "", err
	}

	if err := s.db.One("ID", al.ArtistID, &ar); err != nil {
		return "", err
	}

	return filepath.Join(s.pth, ar.Name, al.Name, fmt.Sprintf("%s.flac", t.Name)), nil
}

func (s Storm) Init() error {
	return nil
}

func (s Storm) InsertOrGetArtist(name string) (int64, error) {
	a := artist{Name: name}
	err := s.db.One("Name", name, &a)
	fmt.Println("not found artist?", err == storm.ErrNotFound)
	if err == storm.ErrNotFound {
		return a.ID, s.db.Save(&a)
	}

	return a.ID, err
}

func (s Storm) InsertOrGetAlbum(name string, artistID int64) (int64, error) {
	a := album{Name: name, ArtistID: artistID}
	err := s.db.One("Name", name, &a)
	if err == storm.ErrNotFound {
		return a.ID, s.db.Save(&a)
	}

	return a.ID, err
}

func (s Storm) InsertOrGetTrack(name string, albumID int64) (int64, error) {
	t := track{Name: name, AlbumID: albumID}
	err := s.db.One("Name", name, &t)
	if err == storm.ErrNotFound {
		var a album
		if err := s.db.One("ID", albumID, &a); err != nil {
			return 0, err
		}

		t.ArtistID = a.ArtistID
		return t.ID, s.db.Save(&t)
	}

	return t.ID, err
}
