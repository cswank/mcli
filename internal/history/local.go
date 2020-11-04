package history

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/cswank/mcli/internal/schema"
)

type StormHistory struct {
	db *sql.DB
}

func NewLocal(db *sql.DB) *StormHistory {
	return &StormHistory{db: db}
}

func (b *StormHistory) Close() error {
	return b.db.Close()
}

func (b *StormHistory) Save(r schema.Result) error {
	var count int64
	if err := b.db.QueryRow("select count from history where id = ?", r.Track.ID).Scan(&count); err != nil {
		return err
	}

	_, err := b.db.Exec("update history set count = ?, time = ? where id = ?", count+1, time.Now().Format(time.RFC3339), r.Track.ID)
	return err
}

func (b *StormHistory) Fetch(page, pageSize int, sortTerm Sort) (*schema.Results, error) {
	log.Println("history", page, pageSize, sortTerm)
	offset := page * pageSize

	q := fmt.Sprintf(`SELECT ar.id, ar.name, al.id, al.name, t.id, t.name, h.count
FROM history AS h
JOIN tracks AS t ON t.id = h.id
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
ORDER BY %s DESC
LIMIT %d OFFSET %d;`, sortTerm, pageSize, offset)

	rows, err := b.db.Query(q)
	if err != nil {
		return nil, err
	}

	var out []schema.Result

	for rows.Next() {
		var res schema.Result
		if err := rows.Scan(&res.Artist.ID, &res.Artist.Name, &res.Album.ID, &res.Album.Title, &res.Track.ID, &res.Track.Title, &res.PlayCount); err != nil {
			return nil, err
		}
		out = append(out, res)
	}

	return &schema.Results{
		Type:    "history",
		Results: out,
	}, nil
}
