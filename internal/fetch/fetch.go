package fetch

import "github.com/cswank/mcli/internal/schema"

type Fetcher interface {
	Name() string
	Login(string, string) error
	Ping() bool
	AlbumLink() string
	FindArtist(string, int) (*schema.Results, error)
	FindAlbum(string, int) (*schema.Results, error)
	FindTrack(string, int) (*schema.Results, error)
	GetAlbum(int64) (*schema.Results, error)
	GetArtistAlbums(int64, int) (*schema.Results, error)
	GetArtistTracks(int64, int) (*schema.Results, error)
	GetPlaylists() (*schema.Results, error)
	GetPlaylist(int64, int) (*schema.Results, error)
	Import(fn func(schema.Progress)) error
}
