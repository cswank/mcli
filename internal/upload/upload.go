package upload

import (
	"io"

	"bitbucket.org/cswank/mcli/internal/schema"
)

type Uploader interface {
	Upload(u schema.Upload, rd io.Reader, f func(pg schema.Progress)) error
}
