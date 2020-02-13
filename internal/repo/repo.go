package repo

import "github.com/cswank/mcli/internal/schema"

type History interface {
	Save(schema.Result) error
	Fetch(int, int, Sort) (*schema.Results, error)
}

type Sort string

const (
	Time  Sort = "Time"
	Count Sort = "Count"
)
