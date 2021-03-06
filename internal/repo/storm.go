package repo

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

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
		Duration int
	}
)

func NewStorm(cfg schema.Config) (*Storm, error) {
	db, err := storm.Open(filepath.Join(cfg.Home, "storm.db"))
	return &Storm{
		db:  db,
		pth: cfg.Pth,
	}, err
}

func (s Storm) FindArtist(term string, p, ps int) ([]schema.Result, error) {
	var a []artist
	if err := s.db.Select(q.Re("Name", term)).Skip(p * ps).Limit(ps).Find(&a); err != nil {
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

func (s Storm) FindAlbum(term string, p, ps int) ([]schema.Result, error) {
	var a []album
	if err := s.db.Select(q.Re("Name", term)).Skip(p * ps).Limit(ps).Find(&a); err != nil {
		if err == storm.ErrNotFound {
			return nil, nil
		}

		return nil, err
	}

	out := make([]schema.Result, len(a))
	for i, alb := range a {
		out[i] = schema.Result{Album: schema.Album{ID: alb.ID, Title: alb.Name}}
	}

	return out, nil
}

func (s Storm) FindTrack(term string, p, ps int) ([]schema.Result, error) {
	var t []track
	if err := s.db.Select(q.Re("Name", term)).Skip(p * ps).Limit(ps).Find(&t); err != nil {
		if err == storm.ErrNotFound {
			return nil, nil
		}

		return nil, err
	}

	out := make([]schema.Result, len(t))
	for i, tr := range t {
		out[i] = schema.Result{
			Track:  schema.Track{ID: tr.ID, Title: tr.Name, Duration: tr.Duration},
			Album:  schema.Album{ID: tr.AlbumID},
			Artist: schema.Artist{ID: tr.ArtistID},
		}
	}

	return out, nil
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
			Track:  schema.Track{ID: track.ID, Title: track.Name, Duration: track.Duration},
			Album:  schema.Album{ID: alb.ID, Title: alb.Name},
			Artist: schema.Artist{ID: art.ID, Name: art.Name},
		}
	}

	return out, nil
}

func (s Storm) GetArtistAlbums(id int64, p, ps int) ([]schema.Result, error) {
	var ar artist
	if err := s.db.One("ID", id, &ar); err != nil {
		return nil, err
	}

	var a []album
	if err := s.db.Select(q.Eq("ArtistID", id)).Skip(p * ps).Limit(ps).Find(&a); err != nil {
		if err == storm.ErrNotFound {
			return []schema.Result{}, nil
		}

		return nil, err
	}

	out := make([]schema.Result, len(a))
	for i, alb := range a {
		out[i] = schema.Result{
			Album:  schema.Album{ID: alb.ID, Title: alb.Name},
			Artist: schema.Artist{ID: ar.ID, Name: ar.Name},
		}
	}

	return out, nil
}

func (s Storm) GetArtistTracks(id int64, p, ps int) ([]schema.Result, error) {
	var ar artist
	if err := s.db.One("ID", id, &ar); err != nil {
		return nil, err
	}

	var a []album
	if err := s.db.Select(q.Eq("ArtistID", id)).Find(&a); err != nil {
		if err == storm.ErrNotFound {
			return []schema.Result{}, nil
		}

		return nil, err
	}

	m := make(map[int64]string)
	ors := make([]q.Matcher, len(a))
	for _, alb := range a {
		m[alb.ID] = alb.Name
		ors = append(ors, q.Eq("AlbumID", alb.ID))
	}

	var t []track
	if err := s.db.Select(ors...).Find(&t); err != nil {
		return nil, err
	}

	out := make([]schema.Result, len(t))
	for i, tr := range t {
		out[i] = schema.Result{
			Album:  schema.Album{ID: tr.AlbumID, Title: m[tr.AlbumID]},
			Artist: schema.Artist{ID: ar.ID, Name: ar.Name},
			Track:  schema.Track{ID: tr.ID, Title: tr.Name, Duration: tr.Duration},
		}
	}

	return out, nil
}

func (s Storm) GetPlaylists() ([]schema.Result, error) {
	return nil, nil
}

func (s Storm) GetPlaylist(int64, int, int) ([]schema.Result, error) {
	return nil, nil
}

func (s *Storm) Close() error {
	return s.db.Close()
}

func (s *Storm) Save(res schema.Result) error {
	var t track
	if err := s.db.One("ID", res.Track.ID, &t); err != nil {
		return err
	}

	t.Count++
	t.Time = time.Now().Format(time.RFC3339)
	t.Duration = res.Track.Duration
	return s.db.Save(&t)
}

func (s *Storm) History(page, pageSize int, sortTerm Sort) ([]schema.Result, error) {
	var entries []track
	err := s.db.Select().OrderBy(strings.Title(string(sortTerm))).Reverse().Limit(pageSize).Skip(page * pageSize).Find(&entries)
	if err != nil {
		return nil, err
	}

	out := make([]schema.Result, len(entries))
	for i, t := range entries {
		var a album
		if err := s.db.One("ID", t.AlbumID, &a); err != nil {
			return nil, err
		}

		var ar artist
		if err := s.db.One("ID", a.ArtistID, &ar); err != nil {
			return nil, err
		}

		out[i] = schema.Result{
			PlayCount: int(t.Count),
			Track:     schema.Track{ID: t.ID, Title: t.Name, Duration: t.Duration},
			Album:     schema.Album{ID: a.ID, Title: a.Name},
			Artist:    schema.Artist{ID: ar.ID, Name: ar.Name},
		}
	}

	return out, nil
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

func (s *Storm) AllTracks() ([]int64, error) {
	var t []track
	if err := s.db.All(&t); err != nil {
		return nil, err
	}
	out := make([]int64, len(t))
	for i, track := range t {
		out[i] = track.ID
	}

	return out, nil
}

func (s Storm) SaveDuration(id int64, duration int) error {
	return s.db.UpdateField(&track{ID: id}, "Duration", duration)
}

func (s Storm) Init() error {
	return nil
}

func (s Storm) InsertOrGetArtist(name string) (int64, error) {
	a := artist{Name: name}
	err := s.db.One("Name", name, &a)
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
