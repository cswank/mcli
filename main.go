package main

import (
	"os"
	"path/filepath"

	"github.com/cswank/mcli/internal/app"
	"github.com/cswank/mcli/internal/schema"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	homeDir = filepath.Join(os.Getenv("HOME"), ".mcli")
)

var (
	_          = kingpin.New("mcli", "A command-line music player.")
	srv        = kingpin.Flag("serve", "start the grpc server").Default("false").Bool()
	initDB     = kingpin.Flag("initdb", "initialize the database").Default("false").Bool()
	duration   = kingpin.Flag("duration", "parse all songs in the database to get length").Default("false").Bool()
	addr       = kingpin.Flag("address", "address of grpc server").Short('a').Envar("MCLI_HOST").String()
	pth        = kingpin.Flag("music", "path to the flac files").Short('m').Envar("MCLI_MUSIC_LOCATION").String()
	home       = kingpin.Flag("home", "path to the directory where the database file lives").Default(homeDir).Envar("MCLI_HOME").String()
	remotePlay = kingpin.Flag("remote", "play music on the server").Short('r').Default("false").Bool()
	logout     = kingpin.Flag("log", "log location (for debugging)").Short('l').String()
	db         = kingpin.Flag("db", "which db?").Short('d').Default("sqlite").Enum("sqlite", "storm")
)

func main() {
	kingpin.Parse()

	cfg := schema.NewConfig(*addr, *pth, *home, *logout, *db, *remotePlay)

	if *srv {
		app.Serve(cfg)
	} else if *initDB {
		app.InitDB(cfg)
	} else if *duration {
		app.Duration(cfg)
	} else {
		app.UI(cfg)
	}
}
