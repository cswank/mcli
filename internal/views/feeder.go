package views

import (
	"bitbucket.org/cswank/music/internal/source"
)

//feeder feeds the screen the data that it craves
type feeder interface {
	header() string
	rows() ([]string, error)
	enter(row int) feeder
}

type root struct {
	cli    source.Source
	width  int
	height int

	//mode is artist, album, or track
	mode string
}

func newRoot(cli source.Source, width, height int) *root {
	return &root{
		cli:    cli,
		width:  width,
		height: height,
	}
}

func (r *root) header() string {
	return ""
}

func (r *root) rows() ([]string, error) {
	return nil, nil
}

func (r *root) enter() feeder {
	return nil
}
