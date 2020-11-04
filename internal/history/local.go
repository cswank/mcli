package history

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	db, err := sql.Open("sqlite3", filepath.Join(dir, "database.sql"))
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `create table history (id integer not null primary key, count integer, time text);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

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

	for _, entry := range entries {
		fmt.Printf("%+v\n", entry)
		pth, ok := entry.ID.(string)
		if !ok {
			continue
		}

		parts := strings.Split(pth, "/")
		if len(parts) < 4 {
			continue
		}

		t := parts[len(parts)-1]
		q := `select id from tracks where name = ?;`
		var id int64
		if err := db.QueryRow(q, t).Scan(&id); err != nil {
			log.Printf("unable to find %s", t)
			continue
		}

		_, err := db.Exec("insert into history (id, count, time) values (?, ?, ?);", id, entry.Count, entry.Time)
		if err != nil {
			log.Printf("unable to write %s", err)
		}
	}

	db.Close()

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
