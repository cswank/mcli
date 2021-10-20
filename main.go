package main

import (
	"os"
	"path/filepath"

	"github.com/cswank/mcli/internal/app"
	"github.com/cswank/mcli/internal/schema"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	homeDir = filepath.Join(os.Getenv("HOME"), ".config", "mcli")
)

var (
	_              = kingpin.New("mcli", "A command-line music player.")
	srv            = kingpin.Command("serve", "start the grpc server")
	_              = kingpin.Command("initdb", "initialize the database")
	_              = kingpin.Command("duration", "parse all songs in the database to get length")
	_              = kingpin.Command("ui", "mci UI").Default()
	addr           = kingpin.Flag("address", "address of grpc server").Short('a').Envar("MCLI_HOST").String()
	pth            = kingpin.Flag("music", "path to the flac files").Short('m').Envar("MCLI_MUSIC_LOCATION").String()
	home           = kingpin.Flag("home", "path to the directory where the database file lives").Default(homeDir).Envar("MCLI_HOME").String()
	remoteSpeakers = kingpin.Flag("remote", "play music on the server").Short('r').Default("false").Envar("MCLI_REMOTE_SPEAKERS").Bool()
	logout         = kingpin.Flag("log", "log location (for debugging)").Short('l').String()
	db             = kingpin.Flag("db", "which db?").Short('d').Default("sqlite").Enum("sqlite", "storm")

	speakers = srv.Flag("speakers", "play music on the server").Short('s').Default("true").Bool()
)

func main() {
	cmd := kingpin.Parse()

	cfg := schema.NewConfig(*addr, *pth, *home, *logout, *db, *remoteSpeakers, speakers)

	switch cmd {
	case "serve":
		app.Serve(cfg)
	case "initdb":
		app.InitDB(cfg)
	case "duration":
		app.Duration(cfg)
	default:
		app.UI(cfg)
	}
}
