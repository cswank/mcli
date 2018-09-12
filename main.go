package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"bitbucket.org/cswank/mcli/internal/http"
	"bitbucket.org/cswank/mcli/internal/player"
	"bitbucket.org/cswank/mcli/internal/rpc"
	"bitbucket.org/cswank/mcli/internal/views"
	rice "github.com/GeertJohan/go.rice"
	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

var (
	srv     = kingpin.Flag("server", "start grpc and http servers").Short('s').Bool()
	cli     = kingpin.Flag("client", "start grpc client").Short('c').Bool()
	cache   = kingpin.Flag("cache", "cache songs to disk").Bool()
	addr    = kingpin.Flag("address", "address of grpc server").Short('a').Default("localhost:50051").String()
	logout  = kingpin.Flag("log", "log location (for debugging)").Short('l').String()
	logfile *os.File
)

func init() {
	if os.Getenv("MCLI_HOME") == "" {
		os.Setenv("MCLI_HOME", fmt.Sprintf("%s/.mcli", os.Getenv("HOME")))
	}

	kingpin.Parse()
	if *srv {
		return
	}

	if *logout != "" {
		f, err := os.Create(*logout)
		if err != nil {
			log.Fatal(err)
		}
		logfile = f
		log.SetOutput(f)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

func cleanup() {
	if logfile != nil {
		logfile.Close()
	}
}

func main() {
	defer cleanup()
	if *srv {
		doServe()
	} else {
		gui()
	}
}

func doServe() {
	cli, err := player.NewDisk(nil)
	if err != nil {
		log.Fatal("cli ", err)
	}

	go func() {
		if err := rpc.Start(cli); err != nil {
			log.Fatal("rpc ", err)
		}
	}()

	box := rice.MustFindBox("internal/http/html")
	if err := http.Start(cli, box); err != nil {
		log.Fatal("http ", err)
	}
}

func gui() {
	var p player.Player
	if *cli {
		c, err := rpc.NewClient(*addr)
		if err != nil {
			log.Fatal(err)
		}
		if err := views.Start(c); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := views.Start(p); err != nil {
			log.Fatal(err)
		}
	}
}
