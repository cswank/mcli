package repo

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/cswank/mcli/internal/schema"
)

type Sort string

const (
	Time  Sort = "time"
	Count Sort = "count"
)

var (
	empty struct{}
)

type SQLite struct {
	db  *sql.DB
	pth string
}

func NewSQL(cfg schema.Config) (*SQLite, error) {
	db, err := sql.Open("sqlite3", filepath.Join(cfg.Home, "mcli.db"))
	return &SQLite{
		db:  db,
		pth: cfg.Pth,
	}, err
}

func (s SQLite) FindArtist(term string, n int) ([]schema.Result, error) {
	q := `SELECT id, name
FROM artists
WHERE name LIKE ?;`
	return s.doFind(q, fmt.Sprintf("%%%s%%", term), artistArgs)
}

func (s SQLite) FindAlbum(term string, n int) ([]schema.Result, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name
FROM albums AS al
JOIN artists AS ar ON ar.id = al.artist_id
WHERE al.name LIKE ?;`
	return s.doFind(q, fmt.Sprintf("%%%s%%", term), albumArgs)
}

func (s SQLite) FindTrack(term string, n int) ([]schema.Result, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name, t.id, t.name
FROM tracks AS t
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
WHERE t.name LIKE ?;`
	return s.doFind(q, fmt.Sprintf("%%%s%%", term), trackArgs)
}

func (s SQLite) GetAlbum(id int64) ([]schema.Result, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name, t.id, t.name
FROM tracks AS t
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
WHERE al.id = ?;`
	return s.doFind(q, id, trackArgs)
}

func (s SQLite) GetArtistAlbums(id int64, n int) ([]schema.Result, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name
FROM albums AS al
JOIN artists AS ar ON ar.id = al.artist_id
WHERE ar.id = ?;`
	return s.doFind(q, id, artistArgs)
}

func (s SQLite) GetArtistTracks(id int64, n int) ([]schema.Result, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name, t.id, t.name
FROM tracks AS t
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
WHERE ar.id = ?;`
	return s.doFind(q, id, trackArgs)
}

func (s SQLite) GetPlaylists() ([]schema.Result, error) {
	return nil, nil
}

func (s SQLite) GetPlaylist(int64, int) ([]schema.Result, error) {
	return nil, nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) Save(res schema.Result) error {
	var count int64
	s.db.QueryRow("select count from history where id = ?", res.Track.ID).Scan(&count)
	var err error
	if count == 0 {
		_, err = s.db.Exec("insert into history (id, count, time) values (?, 1, ?)", res.Track.ID, time.Now().Format(time.RFC3339))
	} else {
		_, err = s.db.Exec("update history set count = ?, time = ? where id = ?", count+1, time.Now().Format(time.RFC3339), res.Track.ID)
	}

	if err != nil {
		return fmt.Errorf("unable to write to history: %s", err)
	}

	return nil
}

func (s *SQLite) Fetch(page, pageSize int, sortTerm Sort) (*schema.Results, error) {
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

func (s SQLite) doFind(q string, term interface{}, f func(schema.Result) []interface{}) ([]schema.Result, error) {
	rows, err := s.db.Query(q, term)
	if err != nil {
		return nil, err
	}

	var out []schema.Result
	for rows.Next() {
		var res schema.Result
		args := f(res)
		if err := rows.Scan(args...); err != nil {
			return nil, err
		}
		out = append(out, res)
	}

	return out, nil
}

func (s *SQLite) Track(id int64) (string, error) {
	q := `SELECT ar.name, al.name, t.name
FROM tracks AS t
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
WHERE t.id = ?;`

	var ar, al, t string
	err := s.db.QueryRow(q, id).Scan(&ar, &al, &t)
	pth := fmt.Sprintf("%s.flac", filepath.Join(s.pth, ar, al, t))
	return pth, err
}

func (s SQLite) Init() error {
	q := `CREATE TABLE IF NOT EXISTS history (id integer not null primary key, count integer, time text);`
	_, err := s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("unable to create history table: %s", err)
	}

	q = `CREATE TABLE IF NOT EXISTS
	artists (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  name TEXT NOT NULL
	);`
	_, err = s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("unable to create artists table: %s", err)
	}

	q = `CREATE TABLE IF NOT EXISTS
	albums (
	  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	  artist_id INTEGER NOT NULL,
	  name TEXT NOT NULL
	);`
	_, err = s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("unable to create albums table: %s", err)
	}

	q = `CREATE TABLE IF NOT EXISTS
	tracks (
	  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	  album_id INTEGER NOT NULL,
	  name TEXT NOT NULL
	);`
	_, err = s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("unable to create tracks table: %s", err)
	}

	return nil
}

func (s SQLite) InsertOrGetTrack(name string, albumID int64) (int64, error) {
	return s.insertOrGet("tracks", "insert into tracks (name, album_id) values (?, ?)", name, albumID)
}

func (s SQLite) InsertOrGetArtist(name string) (int64, error) {
	return s.insertOrGet("artists", "insert into artists (name) values (?)", name)
}

func (s SQLite) InsertOrGetAlbum(name string, artistID int64) (int64, error) {
	return s.insertOrGet("albums", "insert into albums (name, artist_id) values (?, ?)", name, artistID)
}

func (s SQLite) insertOrGet(table, q string, name string, args ...interface{}) (int64, error) {
	var id int64
	err := s.db.QueryRow(fmt.Sprintf("select id from %s where name = ?", table), name).Scan(&id)
	if err == nil {
		return id, nil
	}

	if err != nil && !strings.Contains(err.Error(), "no rows") {
		return 0, err
	}

	args = append([]interface{}{name}, args...)
	_, err = s.db.Exec(q, args...)
	if err != nil {
		return 0, err
	}

	return s.insertOrGet(table, q, name, args)
}

func artistArgs(res schema.Result) []interface{} {
	return []interface{}{&res.Artist.ID, &res.Artist.Name}
}

func albumArgs(res schema.Result) []interface{} {
	return []interface{}{&res.Artist.ID, &res.Artist.Name, &res.Album.ID, &res.Album.Title}
}

func trackArgs(res schema.Result) []interface{} {
	return []interface{}{&res.Artist.ID, &res.Artist.Name, &res.Album.ID, &res.Album.Title, &res.Track.ID, &res.Track.Title}
}
