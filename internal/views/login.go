package views

import (
	"fmt"
	"strings"

	ui "github.com/jroimartin/gocui"
)

const (
	loginWidth  = 30
	loginHeight = 2
)

type login struct {
	name   string
	title  string
	coords coords

	username string
	password string
	login    func(string, string) error
}

func newLogin(w, h int, cb func(string, string) error) *login {
	maxX, maxY := g.Size()
	x1 := maxX/2 - loginWidth/2
	x2 := maxX/2 + loginWidth/2
	y1 := maxY/2 - loginHeight/2
	y2 := maxY/2 + loginHeight/2 + loginHeight%2

	return &login{
		title:  "username",
		coords: coords{x1: x1, y1: y1, x2: x2, y2: y2},
		login:  cb,
	}
}

func (l *login) show(g *ui.Gui, v *ui.View) error {
	v.Editable = true
	v.Frame = true
	v.Title = "username"
	return v.SetCursor(0, 0)
}

func (l *login) next(g *ui.Gui, v *ui.View) error {
	buf := strings.TrimSpace(v.Buffer())
	v.Clear()
	if l.username == "" {
		l.username = buf
		l.title = "password"
		if err := v.SetCursor(0, 0); err != nil {
			return err
		}
		return g.DeleteView("login")
	} else if l.password == "" {
		l.password = buf
		return l.login(l.username, l.password)
	}
	return nil
}

func (l *login) Edit(v *ui.View, key ui.Key, ch rune, mod ui.Modifier) {
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
