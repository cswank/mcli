package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"bitbucket.org/cswank/mcli/internal/fetch"
	"bitbucket.org/cswank/mcli/internal/play"
	"bitbucket.org/cswank/mcli/internal/server"
	"bitbucket.org/cswank/mcli/internal/views"
	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

var (
	app     = kingpin.New("mcli", "A command-line music player.")
	srv     = app.Command("serve", "start grpc and http servers")
	addr    = app.Flag("address", "address of grpc server").Short('a').Default(os.Getenv("MCLI_HOST")).String()
	logout  = app.Flag("log", "log location (for debugging)").Short('l').String()
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
		cleanup := doLog(*logout)
		defer cleanup()
		gui()
	}
}

func doServe() {
	p, err := play.NewLocal(*addr)
	if err != nil {
		log.Fatal("cli ", err)
	}

	f := fetch.NewLocal()

	if err := server.Start(p, f); err != nil {
		log.Fatal("rpc ", err)
	}
}

func gui() {
	f := fetch.NewLocal()
	p, err := play.NewLocal("")
	if err != nil {
		log.Fatal(err)
	}

	if err := views.Start(p, f); err != nil {
		log.Fatal(err)
	}

	// f, err := fetch.NewRemote()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var p play.Player
	// var f fetch.Fetcher

	// if !*srv {
	// 	c, err := rpc.NewClient(*addr, rpc.LocalPlay(!*remote, *addr))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	if err := views.Start(c); err != nil {
	// 		log.Fatal(err)
	// 	}
	// } else {
	// 	cli, err := player.NewDisk(p, *addr)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	if err := views.Start(cli); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
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
