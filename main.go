package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"bitbucket.org/cswank/mcli/internal/player"
	"bitbucket.org/cswank/mcli/internal/server"
	"bitbucket.org/cswank/mcli/internal/views"
	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

var (
	srv     = kingpin.Flag("server", "start grpc server").Short('s').Bool()
	cli     = kingpin.Flag("client", "start grpc client").Short('c').Bool()
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
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func gui() {
	var p player.Player
	if *cli {
		c, err := server.NewClient(*addr)
		if err != nil {
			log.Fatal(err)
		}
		p = c
	}
	if err := views.Start(p); err != nil {
		log.Fatal(err)
	}

	if p != nil {
		c := p.(*server.Client)
		c.Done()
	}
}
