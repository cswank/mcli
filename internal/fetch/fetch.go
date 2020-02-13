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
	GetAlbum(string) (*schema.Results, error)
	GetArtistAlbums(string, int) (*schema.Results, error)
	GetArtistTracks(string, int) (*schema.Results, error)
	GetPlaylists() (*schema.Results, error)
	GetPlaylist(string, int) (*schema.Results, error)
}
