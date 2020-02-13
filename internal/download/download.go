package download

import (
	"io"

	"bitbucket.org/cswank/mcli/internal/schema"
)

type Downloader interface {
	Download(id string, w io.Writer, f func(pg schema.Progress))
}
