package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"bitbucket.org/cswank/mcli/internal/download"
	"bitbucket.org/cswank/mcli/internal/fetch"
	"bitbucket.org/cswank/mcli/internal/play"
	"bitbucket.org/cswank/mcli/internal/repo"
	"bitbucket.org/cswank/mcli/internal/server"
	"bitbucket.org/cswank/mcli/internal/views"
	"google.golang.org/grpc"
	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

var (
	app    = kingpin.New("mcli", "A command-line music player.")
	srv    = kingpin.Flag("serve", "start the grpc server").Default("false").Bool()
	addr   = kingpin.Flag("address", "address of grpc server").Short('a').Default(os.Getenv("MCLI_HOST")).String()
	pth    = kingpin.Flag("music", "path to the flac files").Short('m').Default(os.Getenv("MCLI_MUSIC_LOCATION")).String()
	home   = kingpin.Flag("home", "path to the directory where the database file lives").Default(os.Getenv("MCLI_HOME")).String()
	remote = kingpin.Flag("remote", "play music on the server").Short('r').Default("false").Bool()
	logout = kingpin.Flag("log", "log location (for debugging)").Short('l').String()

	logfile *os.File
)

func main() {
	kingpin.Parse()

	if os.Getenv("MCLI_HOME") == "" {
		os.Setenv("MCLI_HOME", fmt.Sprintf("%s/.mcli", os.Getenv("HOME")))
	}

	if *srv {
		doServe()
	} else {
		defer doLog(*logout)()
		gui()
	}
}

func doServe() {
	h, err := repo.NewLocal(*home)
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

func gui() {
	var f fetch.Fetcher
	var p play.Player
	var h repo.History
	var err error

	switch *addr {
	case "":
		h, err = repo.NewLocal(*home)
		if err != nil {
			log.Fatal(*addr, err)
		}

		dl := download.NewLocal(*pth)
		p, err = play.NewLocal(*pth, play.LocalDownload(dl), play.LocalHistory(h))
		if err != nil {
			log.Fatal(err)
		}

		f = fetch.NewLocal(*pth)
	default:
		conn, err := grpc.Dial(*addr, grpc.WithInsecure())
		if err != nil {
			log.Fatal(*addr, err)
		}

		h = repo.NewRemote(conn)
		f = fetch.NewRemote(conn)
		if *remote {
			p = play.NewRemote(conn)
		} else {
			p, err = play.NewLocal(*pth, play.LocalDownload(download.NewRemote(conn)), play.LocalHistory(h))
		}
		defer conn.Close()
	}

	if err := views.Start(p, f, h); err != nil {
		log.Fatal(err)
	}
}

func doLog(logout string) func() {
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
