package views

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/cswank/mcli/internal/schema"
	ui "github.com/jroimartin/gocui"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// buffer displays the progress bar that is shown when a song is read
// into the play queue
type buffer struct {
	width    int
	coords   coords
	progress chan schema.Progress
	result   chan schema.Result
	text     string
}

func newBuffer(w, h int, id string, cli *client) *buffer {
	b := &buffer{
		width:    w - 1,
		coords:   coords{x1: -1, y1: h - 3, x2: w, y2: h - 1},
		progress: make(chan schema.Progress),
		result:   make(chan schema.Result),
	}

	cli.NextSong(id, b.nextSong)
	cli.DownloadProgress(id, b.downloadProgress)

	go b.render()
	return b
}

func (b *buffer) downloadProgress(prog schema.Progress) {
	b.progress <- prog
}

func (b *buffer) nextSong(r schema.Result) {
	b.result <- r
}

func (b *buffer) render() {
	for {
		select {
		case <-time.After(time.Second):
			if g != nil && b.text != "" {
				v, _ := g.View("buffer")
				if b.text != v.Buffer() {
					g.Update(func(g *ui.Gui) error {
						v.Clear()
						fmt.Fprint(v, b.text)
						return nil
					})
				}
			}
		case r := <-b.result:
			if r.Error == "" {
				b.text = b.center(fmt.Sprintf("%s %s", r.Track.Title, time.Duration(r.Track.Duration)*time.Second))
			} else {
				b.text = b.center(col.Error(r.Error))
			}
			g.Update(func(g *ui.Gui) error {
				v, _ := g.View("buffer")
				v.Clear()
				fmt.Fprint(v, b.text)
				return nil
			})
		case p := <-b.progress:
			if p.Total != 0 {
				g.Update(func(g *ui.Gui) error {
					v, _ := g.View("buffer")
					v.Clear()
					if p.Total >= 0 {
						fmt.Fprint(v, fmt.Sprintf(strings.Repeat("|", b.width*p.N/p.Total)))
					}
					return nil
				})
			}
		}
	}
}

func (b *buffer) clear() {
	b.text = ""
	v, _ := g.View("buffer")
	v.Clear()
}

func (b *buffer) center(s string) string {
	return fmt.Sprintf(fmt.Sprintf("%%%ds", (b.width/2)+(len(s)/2)), s)
}
