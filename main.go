package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/cswank/mcli/internal/download"
	"github.com/cswank/mcli/internal/fetch"
	"github.com/cswank/mcli/internal/history"
	"github.com/cswank/mcli/internal/play"
	"github.com/cswank/mcli/internal/server"
	"github.com/cswank/mcli/internal/views"
	"google.golang.org/grpc"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app        = kingpin.New("mcli", "A command-line music player.")
	srv        = kingpin.Flag("serve", "start the grpc server").Default("false").Bool()
	addr       = kingpin.Flag("address", "address of grpc server").Short('a').Default(os.Getenv("MCLI_HOST")).String()
	pth        = kingpin.Flag("music", "path to the flac files").Short('m').Default(os.Getenv("MCLI_MUSIC_LOCATION")).String()
	home       = kingpin.Flag("home", "path to the directory where the database file lives").Default(os.Getenv("MCLI_HOME")).String()
	remotePlay = kingpin.Flag("remote", "play music on the server").Short('r').Default("false").Bool()
	logout     = kingpin.Flag("log", "log location (for debugging)").Short('l').String()

	logfile *os.File
)

func main() {
	kingpin.Parse()

	if os.Getenv("MCLI_HOME") == "" {
		os.Setenv("MCLI_HOME", fmt.Sprintf("%s/.mcli", os.Getenv("HOME")))
	}

	if *srv {
		startServer()
	} else {
		defer setupLog(*logout)()
		startUI()
	}
}

func startUI() {
	var f fetch.Fetcher
	var p play.Player
	var h history.History
	var close func()

	switch *addr {
	case "":
		p, f, h, close = local()
	default:
		p, f, h, close = remote()
	}

	defer close()

	if err := views.Start(p, f, h); err != nil {
		log.Println(err)
	}
}

func remote() (play.Player, fetch.Fetcher, history.History, func()) {
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(*addr, err)
	}

	h := history.NewRemote(conn)
	f := fetch.NewRemote(conn)

	var p play.Player
	if *remotePlay {
		p = play.NewRemote(conn)
	} else {
		p, err = play.NewLocal(*pth, play.LocalDownload(download.NewRemote(conn)), play.LocalHistory(h))
		if err != nil {
			log.Fatal(*addr, err)
		}
	}

	return p, f, h, func() { conn.Close() }
}

func startServer() {
	h, err := history.NewLocal(*home)
	if err != nil {
		log.Fatal(err)
	}

	dl := download.NewLocal(*pth)
	p, err := play.NewLocal(*pth, play.LocalDownload(dl), play.LocalHistory(h))
	if err != nil {
		log.Fatal("unable to create player ", err)
	}

	f := fetch.NewLocal(*pth)
	if err := server.Start(p, f, h); err != nil {
		log.Fatal("unable to start server ", err)
	}
}

func local() (play.Player, fetch.Fetcher, history.History, func()) {
	h, err := history.NewLocal(*home)
	if err != nil {
		log.Fatal(*addr, err)
	}

	dl := download.NewLocal(*pth)
	p, err := play.NewLocal(*pth, play.LocalDownload(dl), play.LocalHistory(h))
	if err != nil {
		log.Fatal(err)
	}

	f := fetch.NewLocal(*pth)

	return p, f, h, func() {}
}

func setupLog(logout string) func() {
	out := func() {}
	if logout != "" {
		f, err := os.Create(logout)
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
