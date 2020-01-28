package play

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"bitbucket.org/cswank/mcli/internal/schema"
)

type Disk struct {
	pth string
}

func NewDisk() *Disk {
	return &Disk{
		pth: os.Getenv("MCLI_DISK_LOCATION"),
	}
}

func (d *Disk) Name() string               { return "disk" }
func (d *Disk) Login(string, string) error { return nil }
func (d *Disk) Ping() bool                 { return true }
func (d *Disk) AlbumLink() string          { return "" }

func (d *Disk) FindArtist(term string, n int) (*schema.Results, error) {
	glob := filepath.Join(d.pth, fmt.Sprintf("*%s*", term))
	return d.doFind(glob, "artist search")
}

func (d *Disk) FindAlbum(term string, n int) (*schema.Results, error) {
	glob := filepath.Join(d.pth, "*", fmt.Sprintf("*%s*", term))
	return d.doFind(glob, "album search")
}

func (d *Disk) FindTrack(term string, n int) (*schema.Results, error) {
	glob := filepath.Join(d.pth, "*", "*", fmt.Sprintf("*%s*.flac", term))
	return d.doFind(glob, "album")
}

func (d *Disk) doFind(glob, t string) (*schema.Results, error) {
	albums, err := filepath.Glob(glob)
	if err != nil {
		return nil, nil
	}

	var maxTitle int
	out := make([]schema.Result, len(albums))
	for i, s := range albums {
		r := d.resultFromPath(s)
		if len(r.Album.Title) > maxTitle {
			maxTitle = len(r.Album.Title)
		}
		out[i] = r
	}

	f := fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)

	return &schema.Results{
		Header:  fmt.Sprintf(f, "Title", "Artist"),
		Type:    t,
		Fmt:     f,
		Results: out,
	}, nil
}

func (d *Disk) GetAlbum(id string) (*schema.Results, error) {
	tracks, err := filepath.Glob(filepath.Join(d.pth, id, "*.flac"))
	if err != nil {
		return nil, nil
	}

	out := make([]schema.Result, len(tracks))
	var maxTitle int

	for i, tr := range tracks {
		res := d.resultFromPath(tr)
		if len(res.Track.Title) > maxTitle {
			maxTitle = len(res.Track.Title)
		}

		//dur, _ := tr.Duration.Int64() TODO: get duration from flac lib
		out[i] = res
	}

	f := fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)
	return &schema.Results{
		Album:   out[0].Album,
		Header:  fmt.Sprintf(f, "Title", "Length"),
		Type:    "album",
		Fmt:     f,
		Results: out,
	}, nil
}

func (d *Disk) GetTrack(id string) (string, error) {
	return id, nil
}

func (d *Disk) GetArtistAlbums(id string, n int) (*schema.Results, error) {
	log.Println("getartistalbums", id)
	glob := filepath.Join(d.pth, id, "*")
	return d.doFind(glob, "album search")
}

func (d *Disk) GetArtistTracks(id string, n int) (*schema.Results, error) {
	glob := filepath.Join(d.pth, id, "*", "*.flac")
	return d.doFind(glob, "album")
}

func (d *Disk) GetPlaylists() (*schema.Results, error) {
	return &schema.Results{}, nil
}

func (d *Disk) GetPlaylist(string, int) (*schema.Results, error) {
	return &schema.Results{}, nil
}

func (d *Disk) resultFromPath(pth string) schema.Result {
	pth = strings.Replace(pth, d.pth, "", -1)
	if strings.Index(pth, "/") == 0 {
		pth = pth[1:]
	}

	parts := strings.Split(pth, string(filepath.Separator))

	var album schema.Album
	var artist schema.Artist
	var track schema.Track

	artist = schema.Artist{
		ID:   parts[0],
		Name: parts[0],
	}

	if len(parts) >= 2 {
		album = schema.Album{
			ID:    filepath.Join(parts[0], parts[1]),
			Title: parts[1],
		}
	}

	if len(parts) >= 3 {
		track = schema.Track{
			ID:    filepath.Join(d.pth, parts[0], parts[1], parts[2]),
			Title: strings.Replace(parts[2], ".flac", "", -1),
			URI:   filepath.Join("tracks", parts[0], parts[1], parts[2]),
		}
	}

	return schema.Result{
		Artist: artist,
		Album:  album,
		Track:  track,
	}
}
