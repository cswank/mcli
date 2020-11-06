package fetch

import (
	_ "github.com/mattn/go-sqlite3"

	"github.com/cswank/mcli/internal/repo"
	"github.com/cswank/mcli/internal/schema"
)

type Local struct {
	pth string
	db  *repo.SQLLite
}

func NewLocal(pth string, db *repo.SQLLite) *Local {
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
