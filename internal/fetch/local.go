package fetch

import (
	"fmt"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/cswank/mcli/internal/schema"
)

type (
	fetcher interface {
		FindArtist(term string, n int) ([]schema.Result, error)
		FindAlbum(term string, n int) ([]schema.Result, error)
		FindTrack(term string, n int) ([]schema.Result, error)
		GetAlbum(id int64) ([]schema.Result, error)
		GetArtistAlbums(id int64, n int) ([]schema.Result, error)
		GetArtistTracks(id int64, n int) ([]schema.Result, error)
		GetPlaylists() ([]schema.Result, error)
		GetPlaylist(int64, int) ([]schema.Result, error)
		InsertOrGetArtist(name string) (int64, error)
		InsertOrGetAlbum(name string, artistID int64) (int64, error)
		InsertOrGetTrack(name string, albumID int64) (int64, error)
		Init() error
	}

	Local struct {
		pth string
		db  fetcher
	}
)

func NewLocal(pth string, db fetcher) (*Local, error) {
	return &Local{pth: pth, db: db}, db.Init()
}

func (l Local) Name() string               { return "disk" }
func (l Local) Login(string, string) error { return nil }
func (l Local) Ping() bool                 { return true }
func (l Local) AlbumLink() string          { return "" }

func (l Local) FindArtist(term string, n int) (*schema.Results, error) {
	r, err := l.db.FindArtist(term, n)
	return l.doFind(r, "artist search", err, albumTitle)
}

func (l Local) FindAlbum(term string, n int) (*schema.Results, error) {
	r, err := l.db.FindAlbum(term, n)
	return l.doFind(r, "album search", err, albumTitle)
}

func (l Local) FindTrack(term string, n int) (*schema.Results, error) {
	r, err := l.db.FindTrack(term, n)
	return l.doFind(r, "album", err, trackTitle)
}

func (l Local) GetAlbum(id int64) (*schema.Results, error) {
	r, err := l.db.GetAlbum(id)
	return l.doFind(r, "album", err, trackTitle)
}

func (l Local) GetArtistAlbums(id int64, n int) (*schema.Results, error) {
	r, err := l.db.GetArtistAlbums(id, n)
	return l.doFind(r, "album search", err, albumTitle)
}

func (l Local) GetArtistTracks(id int64, n int) (*schema.Results, error) {
	r, err := l.db.GetArtistTracks(id, n)
	return l.doFind(r, "album", err, trackTitle)
}

func (l Local) GetPlaylists() (*schema.Results, error) {
	return &schema.Results{}, nil
}

func (l Local) doFind(res []schema.Result, t string, err error, f func(schema.Result) string) (*schema.Results, error) {
	var maxTitle int
	for _, r := range res {
		if len(f(r)) > maxTitle {
			maxTitle = len(f(r))
		}
	}

	tpl := fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)

	return &schema.Results{
		Header:  fmt.Sprintf(tpl, "Title", "Artist"),
		Type:    t,
		Fmt:     tpl,
		Results: res,
	}, err
}

func (l Local) GetPlaylist(int64, int) (*schema.Results, error) {
	return &schema.Results{}, nil
}

func (l Local) Import(fn func(schema.Progress)) error {
	g := filepath.Join(l.pth, "*", "*", "*.flac")
	tracks, err := filepath.Glob(g)
	if err != nil {
		return err
	}

	m := map[string]map[string][]string{}

	tot := len(tracks)
	for i, pth := range tracks {
		rest, track := filepath.Split(pth)
		album := filepath.Base(filepath.Dir(rest))
		artist := filepath.Base(filepath.Dir(rest[:len(rest)-len(album)]))
		art, ok := m[artist]
		if !ok {
			art = map[string][]string{}
		}

		tracks := art[album]
		tracks = append(tracks, strings.TrimSuffix(track, ".flac"))
		art[album] = tracks
		m[artist] = art
		if i%20 == 0 {
			fn(schema.Progress{N: i, Total: tot})
		}
	}

	var i int
	for artist, albums := range m {
		artID, err := l.db.InsertOrGetArtist(artist)
		if err != nil {
			return err
		}
		for album, tracks := range albums {
			albumID, err := l.db.InsertOrGetAlbum(album, artID)
			if err != nil {
				return err
			}

			for _, track := range tracks {
				_, err := l.db.InsertOrGetTrack(track, albumID)
				if err != nil {
					return err
				}
				if i%20 == 0 {
					fn(schema.Progress{N: i, Total: tot})
				}
				i++
			}
		}
	}

	return nil
}

func albumTitle(r schema.Result) string {
	return r.Album.Title
}

func artistName(r schema.Result) string {
	return r.Artist.Name
}

func trackTitle(r schema.Result) string {
	return r.Track.Title
}
