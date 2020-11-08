package history

import (
	"github.com/cswank/mcli/internal/repo"
	"github.com/cswank/mcli/internal/schema"
)

type (
	fetcher interface {
		Save(r schema.Result) error
		Fetch(page, pageSize int, sortTerm repo.Sort) (*schema.Results, error)
	}

	SQLHistory struct {
		db fetcher
	}
)

func NewLocal(db fetcher) *SQLHistory {
	return &SQLHistory{db: db}
}

func (s *SQLHistory) Save(r schema.Result) error {
	return s.db.Save(r)
}

func (s *SQLHistory) Fetch(page, pageSize int, sortTerm repo.Sort) (*schema.Results, error) {
	return s.db.Fetch(page, pageSize, sortTerm)
}
