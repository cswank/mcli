package views

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"bitbucket.org/cswank/music/internal/source"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/speaker"
	ui "github.com/jroimartin/gocui"
)

type progress struct {
	n     int
	total int
}

type playlist struct {
	ids []string
}

// type progressRead struct {
// 	io.Reader
// 	t, l int
// }

// func newProxy(r io.Reader, l int) *progressRead {
// 	return &progressRead{r, 0, l}
// }

// func (r *progressRead) Read(p []byte) (n int, err error) {
// 	n, err = r.Reader.Read(p)
// 	r.t += n
// 	ui.Update(func() {
// 		progress.SetCurrent(r.t)
// 		progress.SetMax(r.l)
// 	})
// 	return
// }

// // Close the reader when it implements io.Closer
// func (r *progressRead) Close() (err error) {
// 	ui.Update(func() {
// 		progress.SetCurrent(0)
// 		progress.SetMax(1)
// 	})
// 	if closer, ok := r.Reader.(io.Closer); ok {
// 		return closer.Close()
// 	}
// 	return
// }

type play struct {
	coords   coords
	progress chan<- progress
	ch       chan playlist
	cancel   chan bool
	source   source.Source
}

func newPlay(w, h int, pr chan<- progress) *play {
	p := &play{
		coords:   coords{x1: -1, y1: h - 3, x2: w, y2: h - 1},
		progress: pr,
		ch:       make(chan playlist),
		cancel:   make(chan bool),
	}

	go p.play(p.ch, p.cancel)
	return p
}

func (p *play) play(ch <-chan playlist, cancel <-chan bool) error {
	for {
		select {
		case pl := <-ch:
			log.Println("playing", pl)
		case <-cancel:
		}
	}
}

func (p *play) render(g *ui.Gui, v *ui.View) {

}

func (p *play) doPlay() error {
	u := p.source.GetTrack("x")
	resp, err := http.Get(u)
	if err != nil {
		return err
	}

	in, err := ioutil.TempFile("", "")
	log.Println("flac", in.Name())
	if err != nil {
		return err
	}

	m, err := io.Copy(in, resp.Body)
	in.Close()
	log.Println("done writing", in.Name(), m, err)

	f, err := os.Open(in.Name())
	log.Println("open", err)
	if err != nil {
		return err
	}

	s, format, err := flac.Decode(f)
	if err != nil {
		log.Println("decode", err)
		return err
	}

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
		log.Println("speaker init", err)
		return err
	}

	done := make(chan struct{})
	log.Println("about to play", s, format)
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(done)
	})))
	<-done
	return s.Close()
}
