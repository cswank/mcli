package app

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/cswank/mcli/internal/download"
	"github.com/cswank/mcli/internal/fetch"
	"github.com/cswank/mcli/internal/history"
	"github.com/cswank/mcli/internal/play"
	"github.com/cswank/mcli/internal/repo"
	"github.com/cswank/mcli/internal/schema"
	"github.com/cswank/mcli/internal/server"
	"github.com/cswank/mcli/internal/views"
	"google.golang.org/grpc"
)

var (
	logfile *os.File
)

func UI(cfg schema.Config) {
	var f fetch.Fetcher
	var p play.Player
	var h history.History
	var close func()

	switch cfg.Addr {
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

func remote(cfg schema.Config) (play.Player, fetch.Fetcher, history.History, func()) {
	conn, err := grpc.Dial(cfg.Addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(cfg.Addr, err)
	}

	h := history.NewRemote(conn)
	f := fetch.NewRemote(conn)

	var p play.Player
	if cfg.RemotePlay {
		p = play.NewRemote(conn)
	} else {
		p, err = play.NewLocal(cfg.Pth, cfg.Home, play.LocalDownload(download.NewRemote(conn)), play.LocalHistory(h))
		if err != nil {
			log.Fatal(cfg.Addr, err)
		}
	}

	return p, f, h, func() { conn.Close() }
}

func Serve(cfg schema.Config) {
	//db, err := repo.NewSQL(cfg)
	db, err := repo.NewStorm(cfg)
	if err != nil {
		log.Fatal(err)
	}

	h := history.NewLocal(db)
	if err != nil {
		log.Fatal(err)
	}

	f, err := fetch.NewLocal(cfg.Pth, db)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Start(nil, f, h, db, cfg.Pth); err != nil {
		log.Fatal("unable to start server ", err)
	}
}

func local(cfg schema.Config) (play.Player, fetch.Fetcher, history.History, func()) {
	db, err := repo.NewSQL(cfg)
	if err != nil {
		log.Fatal(err)
	}

	h := history.NewLocal(db)
	dl := download.NewLocal(cfg.Pth, db)
	p, err := play.NewLocal(cfg.Pth, cfg.Home, play.LocalDownload(dl), play.LocalHistory(h))
	if err != nil {
		log.Fatal(err)
	}

	f, err := fetch.NewLocal(cfg.Pth, db)
	if err != nil {
		log.Fatal(err)
	}

	return p, f, h, func() {}
}

func SetupLog(cfg schema.Config) func() {
	out := func() {}
	if cfg.Log != "" {
		f, err := os.Create(cfg.Log)
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

func InitDB(cfg schema.Config) {
	db, err := repo.NewSQL(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Init()
	if err != nil {
		log.Fatal(err)
	}
}
