package upload

import (
	"io"

	"bitbucket.org/cswank/mcli/internal/schema"
)

type Local struct {
	// pth is the location of the flac music files
	pth string
}

func NewLocal(pth string) *Local {
	return &Local{
		pth: pth,
	}
}

func (l Local) Upload(u schema.Upload, rd io.Reader, f func(pg schema.Progress)) error {
	return nil
}
