package views

import (
	"fmt"
	"strings"

	ui "github.com/awesome-gocui/gocui"
)

type header struct {
	width  int
	coords coords
	header string
}

func newHeader(w, h int) *header {
	return &header{
		width:  w,
		coords: coords{x1: -1, y1: -1, x2: w, y2: 1},
	}
}

func (h *header) render(g *ui.Gui, v *ui.View) error {
	v.Clear()
	t := fmt.Sprintf("%%s%%%ds", h.width-len(h.header))
	fmt.Fprintf(v, col.C2(t), strings.TrimSuffix(h.header, "\n"), "type 'h' for help")
	return nil
}

func (h *header) clear() {
	v, _ := g.View("header")
	v.Clear()
}
