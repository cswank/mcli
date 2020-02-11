package download

import (
	"io"
	"os"

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

func (l Local) Download(id string, w io.Writer, f func(pg schema.Progress)) error {
	file, err := os.Open(id)
	if err != nil {
		return err
	}

	fi, err := file.Stat()
	if err != nil {
		return err
	}

	tot := int(fi.Size())

	buf := make([]byte, 10000)
	for {
		n, err := file.Read(buf)
		if err != nil {
			return isEOF(err)
		} else {
			_, err := w.Write(buf[:n])
			if err != nil {
				return err
			}
			f(schema.Progress{N: int(n), Total: tot})
		}
	}
}

func isEOF(err error) error {
	if err == io.EOF {
		return nil
	}
	return err
}
