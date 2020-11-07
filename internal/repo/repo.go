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

type SQLLite struct {
	db  *sql.DB
	pth string
}

func New(cfg schema.Config) (*SQLLite, error) {
	db, err := sql.Open("sqlite3", filepath.Join(cfg.Home, "mcli.db"))
	return &SQLLite{
		db:  db,
		pth: cfg.Pth,
	}, err
}

func (s SQLLite) FindArtist(term string, n int) (*schema.Results, error) {
	q := `SELECT id, name
FROM artists
WHERE name LIKE ?;`
	return s.doFind(q, fmt.Sprintf("%%%s%%", term), "artist search")
}

func (s SQLLite) FindAlbum(term string, n int) (*schema.Results, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name
FROM albums AS al
JOIN artists AS ar ON ar.id = al.artist_id
WHERE al.name LIKE ?;`
	return s.doFind(q, fmt.Sprintf("%%%s%%", term), "album search")
}

func (s SQLLite) FindTrack(term string, n int) (*schema.Results, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name, t.id, t.name
FROM tracks AS t
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
WHERE t.name LIKE ?;`
	return s.doFind(q, fmt.Sprintf("%%%s%%", term), "album")
}

func (s SQLLite) GetAlbum(id int64) (*schema.Results, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name, t.id, t.name
FROM tracks AS t
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
WHERE al.id = ?;`
	return s.doFind(q, id, "album")
}

func (s SQLLite) GetArtistAlbums(id int64, n int) (*schema.Results, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name
FROM albums AS al
JOIN artists AS ar ON ar.id = al.artist_id
WHERE ar.id = ?;`
	return s.doFind(q, id, "album search")
}

func (s SQLLite) GetArtistTracks(id int64, n int) (*schema.Results, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name, t.id, t.name
FROM tracks AS t
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
WHERE ar.id = ?;`
	return s.doFind(q, id, "album")
}

func (s SQLLite) GetPlaylists() (*schema.Results, error) {
	return &schema.Results{}, nil
}

func (s SQLLite) GetPlaylist(int64, int) (*schema.Results, error) {
	return &schema.Results{}, nil
}

func (s *SQLLite) Close() error {
	return s.db.Close()
}

func (s *SQLLite) Save(res schema.Result) error {
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

func (s *SQLLite) Fetch(page, pageSize int, sortTerm Sort) (*schema.Results, error) {
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

func (s SQLLite) doFind(q string, term interface{}, t string) (*schema.Results, error) {
	rows, err := s.db.Query(q, term)
	if err != nil {
		return nil, err
	}

	var out []schema.Result
	var maxTitle int

	for rows.Next() {
		var res schema.Result
		title, args := s.args(&res, t)
		if err := rows.Scan(args...); err != nil {
			return nil, err
		}

		out = append(out, res)
		if len(*title) > maxTitle {
			maxTitle = len(*title)
		}
	}

	f := fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)

	return &schema.Results{
		Header:  fmt.Sprintf(f, "Title", "Artist"),
		Type:    t,
		Fmt:     f,
		Results: out,
	}, nil
}

func (s SQLLite) args(res *schema.Result, t string) (*string, []interface{}) {
	switch t {
	case "artist search":
		return &res.Artist.Name, []interface{}{&res.Artist.ID, &res.Artist.Name}
	case "album search":
		return &res.Album.Title, []interface{}{&res.Artist.ID, &res.Artist.Name, &res.Album.ID, &res.Album.Title}
	default:
		return &res.Track.Title, []interface{}{&res.Artist.ID, &res.Artist.Name, &res.Album.ID, &res.Album.Title, &res.Track.ID, &res.Track.Title}
	}
}

func (s *SQLLite) Track(id int64) (string, error) {
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

func (s SQLLite) Init() error {
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

	g := filepath.Join(s.pth, "*", "*", "*.flac")
	tracks, err := filepath.Glob(g)
	if err != nil {
		return err
	}

	m := map[string]map[string][]string{}

	for _, pth := range tracks {
		rest, track := filepath.Split(pth)
		album := filepath.Base(filepath.Dir(rest))
		artist := filepath.Base(filepath.Dir(rest[:len(rest)-len(album)]))
		fmt.Printf("track: %s, album: %s, artist: %s\n", track, album, artist)

		art, ok := m[artist]
		if !ok {
			art = map[string][]string{}
		}

		tracks := art[album]
		tracks = append(tracks, strings.TrimSuffix(track, ".flac"))
		art[album] = tracks
		m[artist] = art
	}

	for artist, albums := range m {
		artID, err := s.insertOrGet("artists", "insert into artists (name) values (?)", artist)
		if err != nil {
			return err
		}
		for album, tracks := range albums {
			albID, err := s.insertOrGet("albums", "insert into albums (name, artist_id) values (?, ?)", album, artID)
			if err != nil {
				return err
			}

			for _, track := range tracks {
				_, err = s.insertOrGet("tracks", "insert into tracks (name, album_id) values (?, ?)", track, albID)
				if err != nil {
					return err
				}
			}
		}
	}

	return s.db.Close()
}

func (s SQLLite) insertOrGet(table, q string, name string, args ...interface{}) (int, error) {
	var id int
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
