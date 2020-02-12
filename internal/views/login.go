package views

import (
	"fmt"
	"strings"

	ui "github.com/jroimartin/gocui"
)

const (
	promptWidth  = 30
	promptHeight = 2
)

type prompt struct {
	name   string
	title  string
	coords coords

	prompt func(string, string) error
}

func newPrompt(w, h int, cb func(string, string) error) *prompt {
	maxX, maxY := g.Size()
	x1 := maxX/2 - promptWidth/2
	x2 := maxX/2 + promptWidth/2
	y1 := maxY/2 - promptHeight/2
	y2 := maxY/2 + promptHeight/2 + promptHeight%2

	return &prompt{
		title:  "username",
		coords: coords{x1: x1, y1: y1, x2: x2, y2: y2},
		prompt: cb,
	}
}

func (p *prompt) show(g *ui.Gui, v *ui.View) error {
	v.Editable = true
	v.Frame = true
	v.Title = p.title
	return v.SetCursor(0, 0)
}

func (p *prompt) Edit(v *ui.View, key ui.Key, ch rune, mod ui.Modifier) {
	in := string(ch)
	buf := strings.TrimSpace(v.Buffer())
	if key == 127 && len(buf) > 0 {
		v.Clear()
		buf = buf[:len(buf)-1]
		v.Write([]byte(c1(buf)))
		v.SetCursor(len(buf), 0)
	} else {
		fmt.Fprint(v, c1(in))
		buf = v.Buffer()
		v.SetCursor(len(buf)-1, 0)
	}
}
