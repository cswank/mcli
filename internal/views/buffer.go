package views

import (
	"fmt"
	"strings"
	"time"

	"bitbucket.org/cswank/mcli/internal/player"
	ui "github.com/jroimartin/gocui"
)

type buffer struct {
	width    int
	coords   coords
	progress chan player.Progress
	song     chan player.Result
	text     string
}

func newBuffer(w, h int, cli player.Player) *buffer {
	b := &buffer{
		width:    w - 1,
		coords:   coords{x1: -1, y1: h - 3, x2: w, y2: h - 1},
		progress: make(chan player.Progress),
		song:     make(chan player.Result),
	}

	cli.NextSong(b.nextSong)
	cli.DownloadProgress(b.downloadProgress)

	go b.render()
	return b
}

func (b *buffer) downloadProgress(prog player.Progress) {
	b.progress <- prog
}

func (b *buffer) nextSong(r player.Result) {
	b.song <- r
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
		case r := <-b.song:
			b.text = b.center(fmt.Sprintf("%s %s", r.Track.Title, time.Duration(r.Track.Duration)*time.Second))
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
					fmt.Fprint(v, fmt.Sprintf(strings.Repeat("|", b.width*p.N/p.Total)))
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
