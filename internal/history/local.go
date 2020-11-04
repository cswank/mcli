package history

import (
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/asdine/storm"
	"github.com/cswank/mcli/internal/schema"
)

type StormEntry struct {
	ID     interface{} `storm:"id"`
	Count  int         `storm:"index"`
	Time   string      `storm:"index"`
	Result interface{}
}

type StormHistory struct {
	db *storm.DB
}

func Migrate(dir string) error {
	// db, err := sql.Open("sqlite3", filepath.Join(dir, "database.sql"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	pth := fmt.Sprintf("%s/history.db", dir)
	fmt.Println(pth)
	st, err := storm.Open(pth)
	if err != nil {
		return err
	}

	var entries []StormEntry
	err = st.All(&entries)
	if err != nil {
		return err
	}

	fmt.Println("got entries", entries)

	for _, entry := range entries {
		fmt.Printf("%+v\n", entry)
		pth, ok := entry.ID.(string)
		if !ok {
			continue
		}

		fmt.Println(pth, entry.Count)
	}

	return nil
}

func NewLocal(dir string) (*StormHistory, error) {
	e, err := exists(dir)
	if err != nil {
		return nil, err
	}

	if !e {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return nil, err
		}
	}
	pth := fmt.Sprintf("%s/history.db", dir)
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
	// for i, e := range entries {
	// 	e.Result.PlayCount = e.Count
	// 	out[i] = e.Result
	// }

	return &schema.Results{
		Type:    "history",
		Results: out,
	}, nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
