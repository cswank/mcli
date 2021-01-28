package views

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	ui "github.com/awesome-gocui/gocui"
)

const (
	seekWidth  = 30
	seekHeight = 2
)

type seek struct {
	coords coords
	doSeek func(int) error
}

func newSeek(w, h int, cb func(i int) error) *seek {
	maxX, maxY := g.Size()
	x1 := maxX/2 - seekWidth/2
	x2 := maxX/2 + seekWidth/2
	y1 := maxY/2 - seekHeight/2
	y2 := maxY/2 + seekHeight/2 + seekHeight%2

	return &seek{
		coords: coords{x1: x1, y1: y1, x2: x2, y2: y2},
		doSeek: cb,
	}
}

func (s *seek) render(g *ui.Gui, v *ui.View) error {
	v.Editable = true
	v.Frame = true
	v.Title = "seek"
	return nil
}

func (s *seek) exit(g *ui.Gui, v *ui.View) error {
	str := strings.TrimSpace(v.Buffer())
	var i int
	v.Clear()

	if len(str) == 0 {
		return nil
	}

	if strings.Contains(str, ":") {
		parts := strings.Split(str, ":")
		min := fmt.Sprintf("%sm", parts[0])
		sec := fmt.Sprintf("%ss", parts[1])
		d, err := time.ParseDuration(min + sec)
		if err != nil {
			return err
		}
		i = int(d / time.Second)
	} else {
		var err error
		i, err = strconv.Atoi(str)
		if err != nil {
			return err
		}
	}
	return s.doSeek(i)
}

func (s *seek) Edit(v *ui.View, key ui.Key, ch rune, mod ui.Modifier) {
	in := string(ch)
	if key == 0 && !strings.Contains("0123456789:", in) {
		return
	}

	buf := strings.TrimSpace(v.Buffer())
	if key == 127 && len(buf) > 0 {
		v.Clear()
		buf = buf[:len(buf)-1]
		v.Write([]byte(col.C1(buf)))
	} else {
		fmt.Fprint(v, col.C1(in))
		buf = v.Buffer()
	}

	v.SetCursor(len(buf), 0)
}
