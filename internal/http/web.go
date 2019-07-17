package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"bitbucket.org/cswank/mcli/internal/player"
	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Start(cli player.Client, box *rice.Box) error {
	srv, err := newServer(cli, box)
	if err != nil {
		return err
	}

	log.Println("http listening on ", srv.server.Addr)
	return srv.server.ListenAndServe()
}

type pager interface {
	scripts([]string)
	stylesheets([]string)
	getTemplate() string
}

type handlerFunc func(http.ResponseWriter, *http.Request) error
type renderFunc func(http.ResponseWriter, *http.Request) (pager, error)

type server struct {
	cli       player.Client
	templates *templates
	server    *http.Server
	disk      string
}

func newServer(cli player.Client, box *rice.Box) (*server, error) {
	t, err := newTemplates(box)
	if err != nil {
		return nil, err
	}

	s := &server{
		cli:       cli,
		templates: t,
		disk:      os.Getenv("MCLI_DISK_LOCATION"),
	}

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", s.handle(s.doRender(s.history))).Methods("GET")
	r.HandleFunc("/history", s.handle(s.doRender(s.history))).Methods("GET")
	r.HandleFunc("/playlists", s.handle(s.doRender(s.playlists))).Methods("GET")
	r.HandleFunc("/playlists/{id}", s.handle(s.doRender(s.playlist))).Methods("GET")
	r.HandleFunc("/search", s.handle(s.doRender(s.search))).Methods("GET")
	r.HandleFunc("/search/{type}", s.handle(s.doRender(s.searchForm))).Methods("GET")
	r.HandleFunc("/albums/{id}", s.handle(s.doRender(s.album))).Methods("GET")
	r.HandleFunc("/artists/{id}", s.handle(s.doRender(s.albums))).Methods("GET")
	r.HandleFunc("/queue", s.handle(s.play)).Methods("POST")
	r.HandleFunc("/queue/album", s.handle(s.playAlbum)).Methods("POST")
	r.HandleFunc("/queue", s.handle(s.doRender(s.queue))).Methods("GET")
	r.HandleFunc("/queue/edit", s.handle(s.doRender(s.editQueue))).Methods("GET")
	r.HandleFunc("/queue/update", s.handle(s.updateQueue)).Methods("POST")
	r.HandleFunc("/volume", s.handle(s.volume)).Methods("POST")
	r.HandleFunc("/pause", s.handle(s.pause)).Methods("POST")
	r.HandleFunc("/rewind", s.handle(s.rewind)).Methods("POST")
	r.HandleFunc("/fastforward", s.handle(s.fastForward)).Methods("POST")
	r.HandleFunc("/ws/play-progress", s.handle(s.playProgress))
	r.HandleFunc("/ws/download-progress", s.handle(s.downloadProgress))
	r.HandleFunc("/tracks/{artist}/{album}/{track}", s.handle(s.getTrack))
	r.PathPrefix("/static/").Handler(static(box)).Methods("GET")

	s.server = &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	return s, err
}

func (s *server) handle(h handlerFunc) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		err := h(w, req)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *server) doRender(f renderFunc) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		pg, err := f(w, r)
		if err != nil || pg == nil {
			return err
		}
		t, scripts, stylesheets := s.templates.get(pg.getTemplate())
		pg.scripts(scripts)
		pg.stylesheets(stylesheets)
		return t.ExecuteTemplate(w, "base", pg)
	}
}

type link struct {
	Name     string
	Link     string
	Selected string
	Children []link
}

type page struct {
	ShowAlbum   bool
	Name        string
	Results     *player.Results
	Links       []link
	Scripts     []string
	Stylesheets []string
	template    string
}

func (p *page) scripts(s []string) {
	p.Scripts = s
}

func (p *page) stylesheets(s []string) {
	p.Stylesheets = s
}

func (p *page) getTemplate() string {
	return p.template
}

func (s *server) getList(r *player.Results) []link {
	out := make([]link, len(r.Results))
	for i, row := range r.Results {
		out[i] = link{Name: row.Artist.Name}
	}

	return out
}

func (s *server) queue(w http.ResponseWriter, req *http.Request) (pager, error) {
	return &page{
		Results:  s.cli.Queue(),
		template: "album.html",
	}, nil
}

func (s *server) editQueue(w http.ResponseWriter, req *http.Request) (pager, error) {
	return &page{
		Results:  s.cli.Queue(),
		template: "queue.html",
	}, nil
}

func (s *server) updateQueue(w http.ResponseWriter, req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	indices, err := s.getInts(req.Form["delete"])
	if err != nil {
		return err
	}

	s.cli.RemoveFromQueue(indices)
	http.Redirect(w, req, "/queue", http.StatusFound)
	return nil
}

func (s *server) getInts(vals []string) ([]int, error) {
	out := make([]int, len(vals))
	for i, val := range vals {
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		out[i] = int(n)
	}
	return out, nil
}

func (s *server) album(w http.ResponseWriter, req *http.Request) (pager, error) {
	vars := mux.Vars(req)
	id := vars["id"]
	a, err := s.cli.GetAlbum(id)
	return &page{
		ShowAlbum: true,
		Results:   a,
		template:  "album.html",
	}, err
}

//albums are those of an artist
func (s *server) albums(w http.ResponseWriter, req *http.Request) (pager, error) {
	vars := mux.Vars(req)
	id := vars["id"]
	a, err := s.cli.GetArtistAlbums(id, 100)
	return &page{
		Results:  a,
		template: "albums.html",
	}, err
}

func (s *server) play(w http.ResponseWriter, req *http.Request) error {
	var result player.Result
	if err := json.NewDecoder(req.Body).Decode(&result); err != nil {
		return err
	}

	s.cli.Play(result)
	return nil
}

func (s *server) playAlbum(w http.ResponseWriter, req *http.Request) error {
	var results player.Results
	if err := json.NewDecoder(req.Body).Decode(&results); err != nil {
		return err
	}

	s.cli.PlayAlbum(&results)
	return nil
}

func (s *server) pause(w http.ResponseWriter, req *http.Request) error {
	s.cli.Pause()
	return nil
}

type vol struct {
	Volume float64 `json:"volume"`
}

func (s *server) volume(w http.ResponseWriter, req *http.Request) error {
	var v vol
	if err := json.NewDecoder(req.Body).Decode(&v); err != nil {
		return err
	}

	s.cli.Volume(v.Volume)
	return nil
}

func (s *server) rewind(w http.ResponseWriter, req *http.Request) error {
	s.cli.Rewind()
	return nil
}

func (s *server) fastForward(w http.ResponseWriter, req *http.Request) error {
	s.cli.FastForward()
	return nil
}

var upgrader = websocket.Upgrader{} // use default options

type wsMessage struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func (s *server) playProgress(w http.ResponseWriter, req *http.Request) error {
	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return err
	}

	id := randString(10)
	s.cli.PlayProgress(id, func(p player.Progress) {
		d, _ := json.Marshal(wsMessage{Type: "play progress", Value: p})
		err := c.WriteMessage(websocket.TextMessage, d)
		if err != nil {
			return
		}
	})

	s.cli.DownloadProgress(id, func(p player.Progress) {
		d, _ := json.Marshal(wsMessage{Type: "download progress", Value: p})
		err := c.WriteMessage(websocket.TextMessage, d)
		if err != nil {
			return
		}
	})

	s.cli.NextSong(id, func(r player.Result) {
		d, _ := json.Marshal(wsMessage{Type: "next song", Value: r})
		err := c.WriteMessage(websocket.TextMessage, d)
		if err != nil {
			return
		}
	})

	defer c.Close()
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			s.cli.Done(id)
			return nil
		}
	}
}

func (s *server) getTrack(w http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	pth := filepath.Join(s.disk, vars["artist"], vars["album"], vars["track"])
	f, err := os.Open(pth)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}

func (s *server) downloadProgress(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func (s *server) history(w http.ResponseWriter, req *http.Request) (pager, error) {
	var sort player.Sort
	switch req.URL.Query().Get("sort_by") {
	case "recent":
		sort = player.Time
	case "count":
		sort = player.Count
	default:
		sort = player.Time
	}
	r, err := s.cli.History(0, 1000, sort)
	return &page{
		Results:  r,
		template: "history.html",
	}, err
}

func (s *server) playlists(w http.ResponseWriter, req *http.Request) (pager, error) {
	r, err := s.cli.GetPlaylists()
	return &page{
		Results:  r,
		template: "playlists.html",
	}, err
}

func (s *server) playlist(w http.ResponseWriter, req *http.Request) (pager, error) {
	v := mux.Vars(req)
	r, err := s.cli.GetPlaylist(v["id"], 1000)
	return &page{
		Results:   r,
		ShowAlbum: true,
		template:  "album.html",
	}, err
}

func (s *server) search(w http.ResponseWriter, req *http.Request) (pager, error) {
	q := req.URL.Query()
	term := q.Get("term")
	t := q.Get("type")
	if term == "" || t == "" {
		return &page{
			template: "search.html",
		}, nil
	}

	var res *player.Results
	var err error
	switch t {
	case "album":
		res, err = s.cli.FindAlbum(term, 100)
	case "artist":
		res, err = s.cli.FindArtist(term, 100)
	case "track":
		res, err = s.cli.FindTrack(term, 100)
	}

	return &page{
		Results:  res,
		template: fmt.Sprintf("%ss.html", t),
	}, err
}

type searchPage struct {
	page
	SearchType string
}

func (s *server) searchForm(w http.ResponseWriter, req *http.Request) (pager, error) {
	vars := mux.Vars(req)
	return &searchPage{
		page: page{
			template: "search-form.html",
		},
		SearchType: vars["type"],
	}, nil
}

func static(box *rice.Box) http.HandlerFunc {
	s := http.FileServer(box.HTTPBox())
	return func(w http.ResponseWriter, req *http.Request) {
		s.ServeHTTP(w, req)
	}
}
