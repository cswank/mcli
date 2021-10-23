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
	"github.com/faiface/beep/flac"
	"google.golang.org/grpc"
	"gopkg.in/cheggaaa/pb.v1"
)

var (
	logfile *os.File
)

type (
	repository interface {
		history.Historian
		fetch.Repository
		server.Tracker
		AllTracks() ([]int64, error)
		SaveDuration(id int64, dur int) error
	}
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
	defer SetupLog(cfg)()

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
	r, err := db(cfg)
	if err != nil {
		log.Fatal(err)
	}

	h := history.NewLocal(r)
	if err != nil {
		log.Fatal(err)
	}

	f, err := fetch.NewLocal(cfg.Pth, r)
	if err != nil {
		log.Fatal(err)
	}

	var p play.Player
	if cfg.Speakers {
		p, err = play.NewLocal(cfg.Pth, cfg.Home, play.LocalDownload(download.NewLocal(cfg.Pth, r)), play.LocalHistory(h))
		if err != nil {
			log.Fatal(cfg.Addr, err)
		}
	}

	if err := server.Start(p, f, h, r, cfg.Pth); err != nil {
		log.Fatal("unable to start server ", err)
	}
}

func local(cfg schema.Config) (play.Player, fetch.Fetcher, history.History, func()) {
	r, err := db(cfg)
	if err != nil {
		log.Fatal(err)
	}

	h := history.NewLocal(r)
	dl := download.NewLocal(cfg.Pth, r)
	p, err := play.NewLocal(cfg.Pth, cfg.Home, play.LocalDownload(dl), play.LocalHistory(h))
	if err != nil {
		log.Fatal(err)
	}

	f, err := fetch.NewLocal(cfg.Pth, r)
	if err != nil {
		log.Fatal(err)
	}

	return p, f, h, func() {}
}

func InitDB(cfg schema.Config) {
	db, err := db(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Init()
	if err != nil {
		log.Fatal(err)
	}
}

func Duration(cfg schema.Config) {
	db, err := db(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ids, err := db.AllTracks()
	if err != nil {
		log.Fatal(err)
	}

	bar := pb.StartNew(len(ids))

	for _, id := range ids {
		pth, err := db.Track(id)
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Open(pth)
		if err != nil {
			log.Printf("unable to open %s: %s", pth, err)
			continue
		}

		music, format, err := flac.Decode(f)
		if err != nil {
			log.Printf("unable to parse %s: %s", pth, err)
			continue
		}

		ln := music.Len()
		d := ln / int(format.SampleRate)

		if err := db.SaveDuration(id, d); err != nil {
			log.Fatalf("unable to save duration for %s: %s", pth, err)
		}

		music.Close()
		bar.Increment()
	}
	bar.Finish()
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

func db(cfg schema.Config) (repository, error) {
	switch cfg.DB {
	case "storm":
		return repo.NewStorm(cfg)
	default:
		return repo.NewSQL(cfg)
	}
}
