package source

import "io"

type Result struct {
	Artist   string
	Album    string
	Title    string
	ID       string
	Duration int
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
