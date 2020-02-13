package download

import (
	"io"

	"github.com/cswank/mcli/internal/schema"
)

type Downloader interface {
	Download(id string, w io.Writer, f func(pg schema.Progress))
}
