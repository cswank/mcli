package download

import (
	"io"

	"github.com/cswank/mcli/internal/schema"
)

type Downloader interface {
	Download(id int64, w io.Writer, f func(pg schema.Progress))
}
