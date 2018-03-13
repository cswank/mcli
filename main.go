package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"bitbucket.org/cswank/mcli/internal/server"
	"bitbucket.org/cswank/mcli/internal/views"
	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

var (
	srv     = kingpin.Flag("server", "start grpc server").Short('s').Bool()
	logout  = kingpin.Flag("log", "log location (for debugging)").Short('l').String()
	logfile *os.File
)

func init() {
	kingpin.Parse()
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
	if os.Getenv("MCLI_HOME") == "" {
		os.Setenv("MCLI_HOME", fmt.Sprintf("%s/.mcli", os.Getenv("HOME")))
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
	if err := views.Start(); err != nil {
		log.Fatal(err)
	}
}
