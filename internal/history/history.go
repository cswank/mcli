package history

import (
	"github.com/cswank/mcli/internal/repo"
	"github.com/cswank/mcli/internal/schema"
)

type History interface {
	Save(schema.Result) error
	Fetch(int, int, repo.Sort) (*schema.Results, error)
}
