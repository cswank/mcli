package player

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Disk struct {
	Player
	pth string
}

func NewDisk(p Player, host string) (Client, error) {
	d := &Disk{
		Player: p,
		pth:    os.Getenv("MCLI_DISK_LOCATION"),
	}

	if p == nil {
		return NewFlac(d, false, host)
	}

	return d, nil
}

func (d *Disk) Name() string               { return "disk" }
func (d *Disk) Login(string, string) error { return nil }
func (d *Disk) Ping() bool                 { return true }
func (d *Disk) AlbumLink() string          { return "" }

func (d *Disk) FindArtist(term string, n int) (*Results, error) {
	glob := filepath.Join(d.pth, fmt.Sprintf("*%s*", term))
	return d.doFind(glob, "artist search")
}

func (d *Disk) FindAlbum(term string, n int) (*Results, error) {
	glob := filepath.Join(d.pth, "*", fmt.Sprintf("*%s*", term))
	return d.doFind(glob, "album search")
}

func (d *Disk) FindTrack(term string, n int) (*Results, error) {
	glob := filepath.Join(d.pth, "*", "*", fmt.Sprintf("*%s*.flac", term))
	return d.doFind(glob, "album")
}

func (d *Disk) doFind(glob, t string) (*Results, error) {
	albums, err := filepath.Glob(glob)
	if err != nil {
		return nil, nil
	}

	var maxTitle int
	out := make([]Result, len(albums))
	for i, s := range albums {
		r := d.resultFromPath(s)
		if len(r.Album.Title) > maxTitle {
			maxTitle = len(r.Album.Title)
		}
		out[i] = r
	}

	f := fmt.Sprintf("%%-%ds%%s\n", maxTitle+4)

	return &Results{
		Header:  fmt.Sprintf(f, "Title", "Artist"),
		Type:    t,
		Fmt:     f,
		Results: out,
	}, nil
}

func (d *Disk) GetAlbum(id string) (*Results, error) {
	tracks, err := filepath.Glob(filepath.Join(id, "*.flac"))
	if err != nil {
		return nil, nil
	}

	out := make([]Result, len(tracks))
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
	return &Results{
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

func (d *Disk) GetArtistAlbums(id string, n int) (*Results, error) {
	glob := filepath.Join(id, "*")
	return d.doFind(glob, "album search")
}

func (d *Disk) GetArtistTracks(id string, n int) (*Results, error) {
	glob := filepath.Join(id, "*", "*.flac")
	return d.doFind(glob, "album")
}

func (d *Disk) GetPlaylists() (*Results, error) {
	return &Results{}, nil
}

func (d *Disk) GetPlaylist(string, int) (*Results, error) {
	return &Results{}, nil
}

func (d *Disk) resultFromPath(pth string) Result {
	pth = filepath.Join(d.pth, pth)
	parts := filepath.SplitList(pth)
	baseParts := filepath.SplitList(d.pth)

	parts = parts[len(baseParts)-1:]

	var album Album
	var artist Artist
	var track Track

	artist = Artist{
		ID:   parts[0],
		Name: parts[0],
	}

	if len(parts) >= 2 {
		album = Album{
			ID:    filepath.Join(parts[0], parts[1]),
			Title: parts[1],
		}
	}

	if len(parts) >= 3 {
		track = Track{
			ID:    filepath.Join(d.pth, parts[0], parts[1], parts[2]),
			Title: strings.Replace(parts[2], ".flac", "", -1),
			URI:   filepath.Join("tracks", parts[0], parts[1], parts[2]),
		}
	}

	return Result{
		Artist: artist,
		Album:  album,
		Track:  track,
	}
}
