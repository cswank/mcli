package fetch

import (
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/cswank/mcli/internal/repo"
	"github.com/cswank/mcli/internal/schema"
)

type Local struct {
	pth string
	db  *repo.SQLLite
}

func NewLocal(pth string, db *repo.SQLLite) *Local {
	db.Init()
	return &Local{
		pth: pth,
		db:  db,
	}
}

func (l Local) Name() string               { return "disk" }
func (l Local) Login(string, string) error { return nil }
func (l Local) Ping() bool                 { return true }
func (l Local) AlbumLink() string          { return "" }

func (l Local) FindArtist(term string, n int) (*schema.Results, error) {
	return l.db.FindArtist(term, n)
}

func (l Local) FindAlbum(term string, n int) (*schema.Results, error) {
	return l.db.FindAlbum(term, n)
}

func (l Local) FindTrack(term string, n int) (*schema.Results, error) {
	return l.db.FindTrack(term, n)
}

func (l Local) GetAlbum(id int64) (*schema.Results, error) {
	return l.db.GetAlbum(id)
}

func (l Local) GetArtistAlbums(id int64, n int) (*schema.Results, error) {
	return l.db.GetArtistAlbums(id, n)
}

func (l Local) GetArtistTracks(id int64, n int) (*schema.Results, error) {
	return l.db.GetArtistTracks(id, n)
}

func (l Local) GetPlaylists() (*schema.Results, error) {
	return &schema.Results{}, nil
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
