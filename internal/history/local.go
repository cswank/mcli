package history

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cswank/mcli/internal/schema"
)

type SQLHistory struct {
	db *sql.DB
}

func NewLocal(db *sql.DB) *SQLHistory {
	return &SQLHistory{db: db}
}

func (s *SQLHistory) Close() error {
	return s.db.Close()
}

func (s *SQLHistory) Save(r schema.Result) error {
	var count int64
	if err := s.db.QueryRow("select count from history where id = ?", r.Track.ID).Scan(&count); err != nil {
		return err
	}

	_, err := s.db.Exec("update history set count = ?, time = ? where id = ?", count+1, time.Now().Format(time.RFC3339), r.Track.ID)
	return err
}

func (s *SQLHistory) Fetch(page, pageSize int, sortTerm Sort) (*schema.Results, error) {
	offset := page * pageSize

	q := fmt.Sprintf(`SELECT ar.id, ar.name, al.id, al.name, t.id, t.name, h.count
FROM history AS h
JOIN tracks AS t ON t.id = h.id
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
ORDER BY %s DESC
LIMIT %d OFFSET %d;`, sortTerm, pageSize, offset)

	rows, err := s.db.Query(q)
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
