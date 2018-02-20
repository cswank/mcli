package main

import (
	"log"
	"os"
	"time"

	"github.com/cswank/music/internal/views"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/speaker"
)

func main() {
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
