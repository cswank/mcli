package repo

import (
	"fmt"
	"os"
	"time"

	"bitbucket.org/cswank/mcli/internal/schema"
	"github.com/asdine/storm"
)

type StormEntry struct {
	ID     string `storm:"id"`
	Count  int    `storm:"index"`
	Time   string `storm:"index"`
	Result schema.Result
}

type StormHistory struct {
	db *storm.DB
}

func NewStorm() (*StormHistory, error) {
	pth := fmt.Sprintf("%s/history.db", os.Getenv("MCLI_HOME"))
	db, err := storm.Open(pth)
	return &StormHistory{db: db}, err
}

func (b *StormHistory) Close() error {
	return b.db.Close()
}

func (b *StormHistory) Save(r schema.Result) error {
	var entry StormEntry
	err := b.db.One("ID", r.Track.ID, &entry)
	if err == storm.ErrNotFound {
		return b.db.Save(&StormEntry{ID: r.Track.ID, Count: 1, Time: time.Now().Format(time.RFC3339), Result: r})
	}

	if err != nil {
		return err
	}

	return b.db.Update(&StormEntry{ID: r.Track.ID, Count: entry.Count + 1, Time: time.Now().Format(time.RFC3339), Result: r})
}

func (b *StormHistory) Fetch(page, pageSize int, sortTerm Sort) (*schema.Results, error) {
	var entries []StormEntry
	err := b.db.Select().OrderBy(string(sortTerm)).Reverse().Limit(pageSize).Skip(page * pageSize).Find(&entries)
	if err != nil {
		return nil, err
	}

	out := make([]schema.Result, len(entries))
	for i, e := range entries {

		e.Result.PlayCount = e.Count
		out[i] = e.Result
	}

	return &schema.Results{
		Type:    "history",
		Results: out,
	}, nil
}