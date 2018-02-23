package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"bitbucket.org/cswank/music/internal/views"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/speaker"
	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

var (
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
}

func cleanup() {
	if logfile != nil {
		logfile.Close()
	}
}

func main() {
	defer cleanup()
	err := views.Start()
	if err != nil {
		log.Fatal(err)
	}

}

func play() {
	f, _ := os.Open("Ahmad Jamal - Sometimes I Feel Like A Motherless Child.flac")
	s, format, _ := flac.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan struct{})
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(done)
	})))
	<-done
}
