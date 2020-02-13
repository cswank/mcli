package download

import (
	"io"
	"log"
	"os"

	"github.com/cswank/mcli/internal/schema"
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

func (l Local) Download(id string, w io.Writer, f func(pg schema.Progress)) {
	file, err := os.Open(id)
	if err != nil {
		log.Println(err)
		return
	}

	fi, err := file.Stat()
	if err != nil {
		log.Println(err)
		return
	}

	tot := int(fi.Size())

	buf := make([]byte, 10000)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			return
		} else if err != nil {
			log.Println(err)
			return
		} else {
			_, err := w.Write(buf[:n])
			if err != nil {
				log.Println(err)
			}
			f(schema.Progress{N: int(n), Total: tot})
		}
	}
}
