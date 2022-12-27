package repo

import (
	"database/sql"
	"fmt"
	"log"
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

func (s SQLite) FindArtist(term string, p, ps int) ([]schema.Result, error) {
	q := `SELECT id, name
FROM artists
WHERE name LIKE ?
LIMIT ? OFFSET ?;`
	return s.doFind(q, fmt.Sprintf("%%%s%%", term), artistArgs, p, ps)
}

func (s SQLite) FindAlbum(term string, p, ps int) ([]schema.Result, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name
FROM albums AS al
JOIN artists AS ar ON ar.id = al.artist_id
WHERE al.name LIKE ?;`
	return s.doFind(q, fmt.Sprintf("%%%s%%", term), albumArgs, p, ps)
}

func (s SQLite) FindTrack(term string, p, ps int) ([]schema.Result, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name, t.id, t.name, t.duration
FROM tracks AS t
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
WHERE t.name LIKE ?;`
	return s.doFind(q, fmt.Sprintf("%%%s%%", term), trackArgs, p, ps)
}

func (s SQLite) GetAlbum(id int64) ([]schema.Result, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name, t.id, t.name, t.duration
FROM tracks AS t
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
WHERE al.id = ?;`
	return s.doFind(q, id, trackArgs, -1, -1)
}

func (s SQLite) GetArtistAlbums(id int64, p, ps int) ([]schema.Result, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name
FROM albums AS al
JOIN artists AS ar ON ar.id = al.artist_id
WHERE ar.id = ?;`
	return s.doFind(q, id, albumArgs, p, ps)
}

func (s SQLite) GetArtistTracks(id int64, p, ps int) ([]schema.Result, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name, t.id, t.name, t.duration
FROM tracks AS t
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
WHERE ar.id = ?;`
	return s.doFind(q, id, trackArgs, p, ps)
}

func (s SQLite) GetPlaylists() ([]schema.Result, error) {
	return nil, nil
}

func (s SQLite) GetPlaylist(int64, int, int) ([]schema.Result, error) {
	return nil, nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) Save(res schema.Result) error {
	var count, duration int64
	err := s.db.QueryRow("select count, duration from history join tracks on tracks.id = history.id where tracks.id = ?", res.Track.ID).Scan(&count, &duration)
	if err != nil && !strings.Contains(err.Error(), "no rows in result") {
		return fmt.Errorf("unable to query history: %s", err)
	}

	if count == 0 {
		_, err = s.db.Exec("insert into history (id, count, time) values (?, 1, ?)", res.Track.ID, time.Now().Format(time.RFC3339))
	} else {
		_, err = s.db.Exec("update history set count = ?, time = ? where id = ?", count+1, time.Now().Format(time.RFC3339), res.Track.ID)
	}

	if err != nil {
		return fmt.Errorf("unable to write to history for track %d: %s", res.Track.ID, err)
	}

	if duration == 0 {
		_, err = s.db.Exec("update tracks set duration = ? where id = ?", res.Track.Duration, res.Track.ID)
	}

	return err
}

func (s *SQLite) History(page, pageSize int, sortTerm Sort) ([]schema.Result, error) {
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

	return out, nil
}

func (s SQLite) doFind(q string, term interface{}, f func(*schema.Result) []interface{}, p, ps int) ([]schema.Result, error) {
	a := []interface{}{term}
	if p >= 0 {
		o := p * ps
		a = append(a, ps, o)
	}
	rows, err := s.db.Query(q, a...)
	if err != nil {
		return nil, err
	}

	var out []schema.Result
	for rows.Next() {
		var res schema.Result
		args := f(&res)
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
	log.Printf("Track %s", pth)
	return pth, err
}

func (s SQLite) DeleteArtist(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM history h
JOIN albums a ON a.id = t.album_id
JOIN tracks t ON t.album_id = a.id
WHERE a.artist_id = ?;`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE FROM tracks
JOIN albums a ON a.id = t.album_id
WHERE a.artist_id = ?;`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE FROM albums WHERE artist_id = ?;`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = s.db.Exec(`DELETE FROM artists WHERE id = ?;`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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
	  name TEXT NOT NULL,
      duration INTEGER NOT NULL DEFAULT 0
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

func (s SQLite) AllTracks() ([]int64, error) {
	rows, err := s.db.Query(`SELECT id FROM tracks`)
	if err != nil {
		return nil, err
	}

	var out []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, nil
}

func (s SQLite) SaveDuration(id int64, duration int) error {
	_, err := s.db.Exec(`UPDATE tracks SET duration = ? WHERE id = ?`, duration, id)
	return err
}

func artistArgs(res *schema.Result) []interface{} {
	return []interface{}{&res.Artist.ID, &res.Artist.Name}
}

func albumArgs(res *schema.Result) []interface{} {
	return []interface{}{&res.Artist.ID, &res.Artist.Name, &res.Album.ID, &res.Album.Title}
}

func trackArgs(res *schema.Result) []interface{} {
	return []interface{}{&res.Artist.ID, &res.Artist.Name, &res.Album.ID, &res.Album.Title, &res.Track.ID, &res.Track.Title, &res.Track.Duration}
}
