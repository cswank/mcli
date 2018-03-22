package player

import (
	"fmt"
	"os"
	"time"

	"github.com/asdine/storm"
)

type StormEntry struct {
	ID     string `storm:"id"`
	Count  int    `storm:"index"`
	Time   string `storm:"index"`
	Result Result
}

type StormHistory struct {
	db *storm.DB
}

func NewStormHistory() (*StormHistory, error) {
	pth := fmt.Sprintf("%s/history.db", os.Getenv("MCLI_HOME"))
	db, err := storm.Open(pth)
	return &StormHistory{db: db}, err
}

func (b *StormHistory) Close() error {
	return b.db.Close()
}

func (b *StormHistory) Save(r Result) error {
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

func (b *StormHistory) Fetch(page, pageSize int) (*Results, error) {
	var entries []StormEntry
	err := b.db.Select().OrderBy("Time").Reverse().Limit(pageSize).Skip(page * pageSize).Find(&entries)
	if err != nil {
		return nil, err
	}

	out := make([]Result, len(entries))
	var maxTitle, maxAlbum int
	for i, e := range entries {
		if len(e.Result.Track.Title) > maxTitle {
			maxTitle = len(e.Result.Track.Title)
		}
		if len(e.Result.Album.Title) > maxAlbum {
			maxAlbum = len(e.Result.Album.Title)
		}
		out[i] = e.Result
	}

	f := fmt.Sprintf("%%-%ds%%-%ds%%s\n", maxTitle+4, maxAlbum+4)

	return &Results{
		Header:  fmt.Sprintf(f, "Title", "Album", "Artist"),
		Type:    "history",
		Fmt:     f,
		Results: out,
	}, nil
}
