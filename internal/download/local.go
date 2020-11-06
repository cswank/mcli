package download

import (
	"io"
	"log"
	"os"

	"github.com/cswank/mcli/internal/repo"
	"github.com/cswank/mcli/internal/schema"
)

type Local struct {
	// pth is the location of the flac music files
	pth string
	db  *repo.SQLLite
}

func NewLocal(pth string, db *repo.SQLLite) *Local {
	return &Local{
		pth: pth,
		db:  db,
	}
}

func (l Local) Download(id int64, w io.Writer, f func(pg schema.Progress)) {
	track, err := l.track(id)
	if err != nil {
		log.Println(err)
		return
	}

	file, err := os.Open(track)
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

func (l Local) track(id int64) (string, error) {
	return l.db.Track(id)
}
