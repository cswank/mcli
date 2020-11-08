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

	LocalHistory struct {
		db fetcher
	}
)

func NewLocal(db fetcher) *LocalHistory {
	return &LocalHistory{db: db}
}

func (l *LocalHistory) Save(r schema.Result) error {
	return l.db.Save(r)
}

func (l *LocalHistory) Fetch(page, pageSize int, sortTerm repo.Sort) (*schema.Results, error) {
	return l.db.Fetch(page, pageSize, sortTerm)
}
