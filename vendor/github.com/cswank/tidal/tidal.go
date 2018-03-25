package tidal

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const baseurl = "https://api.tidalhifi.com/v1/"
const clientVersion = "1.9.1" // ayy that's the golang version too!
const token = "kgsOOmYk3zShYrNP"

var cookieJar, _ = cookiejar.New(nil)
var c = &http.Client{
	Jar: cookieJar, // I stole the cookie from the cookie jar
}

func (t *Tidal) get(dest string, query *url.Values, s interface{}) error {
	req, err := http.NewRequest("GET", baseurl+dest, nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-Tidal-SessionID", t.SessionID)
	query.Add("countryCode", t.CountryCode)
	req.URL.RawQuery = query.Encode()
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&s)
}

func (t *Tidal) CheckSession() (bool, error) {
	var out struct {
		ValidUntil string `json:"validUntil"`
		Status     string `json:"status"`
	}
	err := t.get(fmt.Sprintf("users/%s/subscription", t.UserID), &url.Values{}, &out)
	return out.Status == "ACTIVE" && out.ValidUntil > time.Now().Format(time.RFC3339Nano), err
}

// GetStreamURL func
func (t *Tidal) GetStreamURL(id, q string) (string, error) {
	var s struct {
		URL string `json:"url"`
	}
	err := t.get("tracks/"+id+"/streamUrl", &url.Values{
		"soundQuality": {q},
	}, &s)
	return s.URL, err
}

// GetAlbumTracks func
func (t *Tidal) GetAlbumTracks(id string) ([]Track, error) {
	var s struct {
		Items []Track `json:"items"`
	}
	return s.Items, t.get("albums/"+id+"/tracks", &url.Values{}, &s)
}

// GetPlaylistTracks func
func (t *Tidal) GetPlaylistTracks(id string, l string) ([]Track, error) {
	var s struct {
		Items []Track `json:"items"`
	}
	return s.Items, t.get("playlists/"+id+"/tracks", &url.Values{"limit": {l}}, &s)
}

// SearchTracks func
func (t *Tidal) SearchTracks(d, l string) ([]Track, error) {
	var s Search
	return s.Tracks.Items, t.get("search", &url.Values{
		"query": {d},
		"types": {"TRACKS"},
		"limit": {l},
	}, &s)
}

// SearchAlbums func
func (t *Tidal) SearchAlbums(d, l string) ([]Album, error) {
	var s Search
	return s.Albums.Items, t.get("search", &url.Values{
		"query": {d},
		"types": {"ALBUMS"},
		"limit": {l},
	}, &s)
}

// SearchArtists func
func (t *Tidal) SearchArtists(d, l string) ([]Artist, error) {
	var s Search
	return s.Artists.Items, t.get("search", &url.Values{
		"query": {d},
		"types": {"ARTISTS"},
		"limit": {l},
	}, &s)
}

// SearchArtists func
func (t *Tidal) GetArtistAlbums(artist, l string) ([]Album, error) {
	var s Search
	return s.Items, t.get(fmt.Sprintf("artists/%s/albums", artist), &url.Values{
		"limit": {l},
	}, &s)
}

// SearchArtistTracks func
func (t *Tidal) GetArtistTracks(artist, l string) ([]Track, error) {
	var s Search
	return s.Tracks.Items, t.get(fmt.Sprintf("artists/%s/toptracks", artist), &url.Values{
		"limit":  {l},
		"filter": {"ALL"},
	}, &s.Tracks)
}

func (t *Tidal) GetUserPlaylists() ([]Album, error) {
	var s Search
	return s.Items, t.get(fmt.Sprintf("users/%s/playlists", t.UserID), &url.Values{}, &s)
}

// helper function to generate a uuid
func uuid() string {
	b := make([]byte, 16)
	rand.Read(b[:])
	b[8] = (b[8] | 0x40) & 0x7F
	b[6] = (b[6] & 0xF) | (4 << 4)
	return fmt.Sprintf("%x", b)
}

// New func
func New(user, pass string) (*Tidal, error) {
	query := url.Values{
		"username":        {user},
		"password":        {pass},
		"token":           {token},
		"clientUniqueKey": {uuid()},
		"clientVersion":   {clientVersion},
	}
	res, err := http.PostForm(baseurl+"login/username", query)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected error code from tidal: %d", res.StatusCode)
	}

	defer res.Body.Close()
	var t Tidal
	return &t, json.NewDecoder(res.Body).Decode(&t)
}
