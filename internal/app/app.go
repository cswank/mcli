package app

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/cswank/mcli/internal/download"
	"github.com/cswank/mcli/internal/fetch"
	"github.com/cswank/mcli/internal/history"
	"github.com/cswank/mcli/internal/play"
	"github.com/cswank/mcli/internal/server"
	"github.com/cswank/mcli/internal/views"
	"google.golang.org/grpc"
)

var (
	logfile *os.File
)

type Config struct {
	addr       string
	pth        string
	home       string
	log        string
	remotePlay bool
}

func NewConfig(addr, pth, home, log string, remotePlay bool) Config {
	return Config{
		addr:       addr,
		pth:        pth,
		home:       home,
		log:        log,
		remotePlay: remotePlay,
	}
}

func UI(cfg Config) {
	var f fetch.Fetcher
	var p play.Player
	var h history.History
	var close func()

	switch cfg.addr {
	case "":
		p, f, h, close = local(cfg)
	default:
		p, f, h, close = remote(cfg)
	}

	defer close()

	if err := views.Start(p, f, h); err != nil {
		log.Println(err)
	}
}

func remote(cfg Config) (play.Player, fetch.Fetcher, history.History, func()) {
	conn, err := grpc.Dial(cfg.addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(cfg.addr, err)
	}

	h := history.NewRemote(conn)
	f := fetch.NewRemote(conn)

	var p play.Player
	if cfg.remotePlay {
		p = play.NewRemote(conn)
	} else {
		p, err = play.NewLocal(cfg.pth, cfg.home, play.LocalDownload(download.NewRemote(conn)), play.LocalHistory(h))
		if err != nil {
			log.Fatal(cfg.addr, err)
		}
	}

	return p, f, h, func() { conn.Close() }
}

func Serve(cfg Config) {
	db, err := sql.Open("sqlite3", filepath.Join(cfg.home, "database.sql"))
	if err != nil {
		log.Fatal(err)
	}

	h := history.NewLocal(db)
	if err != nil {
		log.Fatal(err)
	}

	f := fetch.NewLocal(cfg.pth, db)
	if err := server.Start(nil, f, h, db, cfg.pth); err != nil {
		log.Fatal("unable to start server ", err)
	}
}

func local(cfg Config) (play.Player, fetch.Fetcher, history.History, func()) {
	db, err := sql.Open("sqlite3", filepath.Join(cfg.home, "database.sql"))
	if err != nil {
		log.Fatal(err)
	}

	h := history.NewLocal(db)
	dl := download.NewLocal(cfg.pth, db)
	p, err := play.NewLocal(cfg.pth, cfg.home, play.LocalDownload(dl), play.LocalHistory(h))
	if err != nil {
		log.Fatal(err)
	}

	f := fetch.NewLocal(cfg.pth, db)

	return p, f, h, func() {}
}

func SetupLog(cfg Config) func() {
	out := func() {}
	if cfg.log != "" {
		f, err := os.Create(cfg.log)
		if err != nil {
			log.Fatal(err)
		}
		logfile = f
		log.SetOutput(f)
		out = func() {
			f.Close()
		}
	} else {
		log.SetOutput(ioutil.Discard)
	}

	return out
}

func InitDB(cfg Config) {
	db, err := sql.Open("sqlite3", filepath.Join(cfg.pth, "database.sql"))
	if err != nil {
		log.Fatal(err)
	}

	f := fetch.NewLocal(cfg.pth, db)
	err = f.InitDB()
	if err != nil {
		log.Fatal(err)
	}
}
