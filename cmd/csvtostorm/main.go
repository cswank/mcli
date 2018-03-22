package main

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/cswank/mcli/internal/player"
)

func init() {
	if os.Getenv("MCLI_HOME") == "" {
		os.Setenv("MCLI_HOME", fmt.Sprintf("%s/.mcli", os.Getenv("HOME")))
	}
}

func main() {
	fh, err := player.NewFileHistory()
	if err != nil {
		log.Fatal(err)
	}

	results, err := fh.Fetch(0, 10000)
	if err != nil {
		log.Fatal(err)
	}

	sh, err := player.NewStormHistory()
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range results.Results {
		if err := sh.Save(r); err != nil {
			log.Fatal(err)
		}
	}

	sh.Close()
}
