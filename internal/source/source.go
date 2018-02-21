package source

import "io"

type Result struct {
	Artist string
	Title  string
	URL    string
}

type Results struct {
	Type    string
	Header  string
	Results []Result
	Print   func(io.Writer, Result) error
}

type Source interface {
	FindArtist(string, int) (*Results, error)
	FindAlbum(string, int) (*Results, error)
	FindTrack(string, int) (*Results, error)
	GetAlbum(string) (*Results, error)
	GetTrack(string) string
}
