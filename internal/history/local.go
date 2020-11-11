package history

import (
	"github.com/cswank/mcli/internal/repo"
	"github.com/cswank/mcli/internal/schema"
)

type (
	Historian interface {
		Save(r schema.Result) error
		History(page, pageSize int, sortTerm repo.Sort) ([]schema.Result, error)
	}

	LocalHistory struct {
		db Historian
	}
)

func NewLocal(db Historian) *LocalHistory {
	return &LocalHistory{db: db}
}

func (l *LocalHistory) Save(r schema.Result) error {
	return l.db.Save(r)
}

func (l *LocalHistory) Fetch(page, pageSize int, sortTerm repo.Sort) (*schema.Results, error) {
	r, err := l.db.History(page, pageSize, sortTerm)
	return &schema.Results{
		Type:    "history",
		Results: r,
	}, err
}
