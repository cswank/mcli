package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"bitbucket.org/cswank/mcli/internal/download"
	"bitbucket.org/cswank/mcli/internal/fetch"
	"bitbucket.org/cswank/mcli/internal/play"
	"bitbucket.org/cswank/mcli/internal/server"
	"bitbucket.org/cswank/mcli/internal/views"
	"google.golang.org/grpc"
	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

var (
	app    = kingpin.New("mcli", "A command-line music player.")
	srv    = app.Command("serve", "start the grpc server")
	addr   = app.Flag("address", "address of grpc server").Short('a').Default(os.Getenv("MCLI_HOST")).String()
	pth    = app.Flag("path", "path to the flac files").Short('p').Default(os.Getenv("MCLI_MUSIC_LOCATION")).String()
	remote = app.Flag("remote", "play music server").Short('r').Default("false").Bool()
	logout = app.Flag("log", "log location (for debugging)").Short('l').String()

	logfile *os.File
)

func main() {
	if os.Getenv("MCLI_HOME") == "" {
		os.Setenv("MCLI_HOME", fmt.Sprintf("%s/.mcli", os.Getenv("HOME")))
	}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case "serve":
		doServe()
	default:
		defer doLog(*logout)()
		gui()
	}
}

func doServe() {
	dl := download.NewLocal(*pth)
	p, err := play.NewLocal(play.LocalDownload(dl))
	if err != nil {
		log.Fatal("cli ", err)
	}

	f := fetch.NewLocal(*pth)

	if err := server.Start(p, f); err != nil {
		log.Fatal("rpc ", err)
	}
}

func gui() {
	var f fetch.Fetcher
	var p play.Player
	var err error

	switch {
	case *addr != "":
		conn, err := grpc.Dial(*addr, grpc.WithInsecure())
		if err != nil {
			log.Fatal(*addr, err)
		}

		f = fetch.NewRemote(conn)
		if !*remote {
			p = play.NewRemote(conn)
		} else {
			p, err = play.NewLocal(play.LocalDownload(download.NewRemote(conn)))
		}
	case *addr == "":
		dl := download.NewLocal(*pth)
		p, err = play.NewLocal(play.LocalDownload(dl))
		if err != nil {
			log.Fatal(err)
		}

		f = fetch.NewLocal(*pth)
	}

	if err := views.Start(p, f); err != nil {
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
