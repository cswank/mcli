package fetch

import (
	"fmt"
	"path/filepath"
	"sync"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/cswank/mcli/internal/schema"
)

type Local struct {
	pth  string
	db   *sql.DB
	lock sync.Mutex
}

func NewLocal(pth string, db *sql.DB) *Local {
	return &Local{
		pth: pth,
		db:  db,
	}
}

func (l Local) Name() string               { return "disk" }
func (l Local) Login(string, string) error { return nil }
func (l Local) Ping() bool                 { return true }
func (l Local) AlbumLink() string          { return "" }

func (l Local) FindArtist(term string, n int) (*schema.Results, error) {
	q := `SELECT id, name
FROM artists
WHERE name LIKE ?;`
	return l.doFind(q, fmt.Sprintf("%%%s%%", term), "artist search")
}

func (l Local) FindAlbum(term string, n int) (*schema.Results, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name
FROM albums AS al
JOIN artists AS ar ON ar.id = al.artist_id
WHERE al.name LIKE ?;`
	return l.doFind(q, fmt.Sprintf("%%%s%%", term), "album search")
}

func (l Local) FindTrack(term string, n int) (*schema.Results, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name, t.id, t.name
FROM tracks AS t
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
WHERE t.name LIKE ?;`
	return l.doFind(q, fmt.Sprintf("%%%s%%", term), "album")
}

func (l Local) doFind(q string, term interface{}, t string) (*schema.Results, error) {
	rows, err := l.db.Query(q, term)
	if err != nil {
		return nil, err
	}

	var out []schema.Result
	var maxTitle int

	for rows.Next() {
		var args []interface{}
		var res schema.Result
		var title string
		switch t {
		case "artist search":
			args = []interface{}{&res.Artist.ID, &res.Artist.Name}
			title = res.Artist.Name
		case "album search":
			args = []interface{}{&res.Artist.ID, &res.Artist.Name, &res.Album.ID, &res.Album.Title}
			title = res.Album.Title
		case "album":
			args = []interface{}{&res.Artist.ID, &res.Artist.Name, &res.Album.ID, &res.Album.Title, &res.Track.ID, &res.Track.Title}
			title = res.Track.Title
		}

		if err := rows.Scan(args...); err != nil {
			return nil, err
		}

		out = append(out, res)

		if len(title) > maxTitle {
			maxTitle = len(title)
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

func (l Local) GetAlbum(id int64) (*schema.Results, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name, t.id, t.name
FROM tracks AS t
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
WHERE al.id = ?;`
	return l.doFind(q, id, "album")
}

func (l Local) GetArtistAlbums(id int64, n int) (*schema.Results, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name
FROM albums AS al
JOIN artists AS ar ON ar.id = al.artist_id
WHERE ar.id = ?;`
	return l.doFind(q, id, "album search")
}

func (l Local) GetArtistTracks(id int64, n int) (*schema.Results, error) {
	q := `SELECT ar.id, ar.name, al.id, al.name, t.id, t.name
FROM tracks AS t
JOIN albums AS al ON al.id = t.album_id
JOIN artists AS ar ON ar.id = al.artist_id
WHERE ar.id = ?;`
	return l.doFind(q, id, "album")
}

func (l Local) GetPlaylists() (*schema.Results, error) {
	return &schema.Results{}, nil
}

func (l Local) GetPlaylist(int64, int) (*schema.Results, error) {
	return &schema.Results{}, nil
}

func (l Local) InitDB() error {
	q := `create table
	artists (
	  id integer not null primary key,
	  name text
	);`
	_, err := l.db.Exec(q)
	if err != nil {
		return fmt.Errorf("unable to create artists table: %s", err)
	}

	q = `create table
	albums (
	  id integer not null primary key,
	  artist_id integer,
	  name text
	);`
	_, err = l.db.Exec(q)
	if err != nil {
		return fmt.Errorf("unable to create albums table: %s", err)
	}

	q = `create table
	tracks (
	  id integer not null primary key,
	  album_id integer,
	  name text
	);`
	_, err = l.db.Exec(q)
	if err != nil {
		return fmt.Errorf("unable to create tracks table: %s", err)
	}

	g := filepath.Join(l.pth, "*", "*", "*.flac")
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
		tracks = append(tracks, track)
		art[album] = tracks
		m[artist] = art
	}

	artID := 1
	albID := 1
	trackID := 1

	for artist, albums := range m {
		_, err := l.db.Exec("insert into artists (id, name) values (?, ?)", artID, artist)
		if err != nil {
			return err
		}
		for album, tracks := range albums {
			_, err = l.db.Exec("insert into albums (id, name, artist_id) values (?, ?, ?)", albID, album, artID)
			if err != nil {
				return err
			}

			for _, track := range tracks {
				_, err = l.db.Exec("insert into tracks (id, name, album_id) values (?, ?, ?)", trackID, track, albID)
				if err != nil {
					return err
				}
				trackID++
			}
			albID++
		}
		artID++
	}

	return l.db.Close()
}
