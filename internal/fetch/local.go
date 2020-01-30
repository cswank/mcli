package fetch

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"bitbucket.org/cswank/mcli/internal/schema"
)

type Local struct {
	pth string
}

func NewLocal(pth string) *Local {
	return &Local{
		pth: pth,
	}
}

func (l Local) Name() string               { return "disk" }
func (l Local) Login(string, string) error { return nil }
func (l Local) Ping() bool                 { return true }
func (l Local) AlbumLink() string          { return "" }

func (l Local) FindArtist(term string, n int) (*schema.Results, error) {
	glob := filepath.Join(l.pth, fmt.Sprintf("*%s*", term))
	return l.doFind(glob, "artist search")
}

func (l Local) FindAlbum(term string, n int) (*schema.Results, error) {
	glob := filepath.Join(l.pth, "*", fmt.Sprintf("*%s*", term))
	return l.doFind(glob, "album search")
}

func (l Local) FindTrack(term string, n int) (*schema.Results, error) {
	glob := filepath.Join(l.pth, "*", "*", fmt.Sprintf("*%s*.flac", term))
	return l.doFind(glob, "album")
}

func (l Local) doFind(glob, t string) (*schema.Results, error) {
	albums, err := filepath.Glob(glob)
	if err != nil {
		return nil, nil
	}

	var maxTitle int
	out := make([]schema.Result, len(albums))
	for i, s := range albums {
		r := l.resultFromPath(s)
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

func (l Local) GetAlbum(id string) (*schema.Results, error) {
	tracks, err := filepath.Glob(filepath.Join(l.pth, id, "*.flac"))
	if err != nil {
		return nil, nil
	}

	out := make([]schema.Result, len(tracks))
	var maxTitle int

	for i, tr := range tracks {
		res := l.resultFromPath(tr)
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

func (l Local) GetTrack(id string) (string, error) {
	return id, nil
}

func (l Local) GetArtistAlbums(id string, n int) (*schema.Results, error) {
	log.Println("getartistalbums", id)
	glob := filepath.Join(l.pth, id, "*")
	return l.doFind(glob, "album search")
}

func (l Local) GetArtistTracks(id string, n int) (*schema.Results, error) {
	glob := filepath.Join(l.pth, id, "*", "*.flac")
	return l.doFind(glob, "album")
}

func (l Local) GetPlaylists() (*schema.Results, error) {
	return &schema.Results{}, nil
}

func (l Local) GetPlaylist(string, int) (*schema.Results, error) {
	return &schema.Results{}, nil
}

func (l Local) resultFromPath(pth string) schema.Result {
	pth = strings.Replace(pth, l.pth, "", -1)
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
		//uri := filepath.Join("tracks", parts[0], parts[1], parts[2])
		uri := fmt.Sprintf("file://%s", filepath.Join(l.pth, parts[0], parts[1], parts[2]))
		log.Println(uri, l.pth)
		track = schema.Track{
			ID:    filepath.Join(l.pth, parts[0], parts[1], parts[2]),
			Title: strings.Replace(parts[2], ".flac", "", -1),
			URI:   uri,
		}
	}

	return schema.Result{
		Artist: artist,
		Album:  album,
		Track:  track,
	}
}
