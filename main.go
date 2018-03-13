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
	if *srv || *cli {
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
	} else if *cli {
		doClient()
	} else {
		gui()
	}
}

func doServe() {
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func doClient() {
	if err := server.StartClient(*addr); err != nil {
		log.Fatal(err)
	}
}

func gui() {
	if err := views.Start(); err != nil {
		log.Fatal(err)
	}
}
