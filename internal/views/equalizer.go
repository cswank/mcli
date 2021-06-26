package views

import (
	"bytes"
	"log"

	ui "github.com/awesome-gocui/gocui"
)

var (
	equalizerWidth  = 100
	equalizerHeight = 26
)

type equalizer struct {
	name   string
	coords coords
	body   []byte
}

func newEqualizer(w, h int) *equalizer {
	return &equalizer{
		name: "equalizer",
	}
}

func getEqualizerCoords(g *ui.Gui) coords {
	maxX, maxY := g.Size()
	x1 := maxX/2 - equalizerWidth/2
	x2 := maxX/2 + equalizerWidth/2
	y1 := maxY/2 - equalizerHeight/2
	y2 := maxY/2 + equalizerHeight/2 + equalizerHeight%2
	return coords{x1: x1, y1: y1, x2: x2, y2: y2}
}

func (h *equalizer) show(g *ui.Gui, keys []key) error {
	coords := getEqualizerCoords(g)
	v, err := g.SetView("equalizer", coords.x1, coords.y1, coords.x2, coords.y2, 0)
	if !ui.IsUnknownView(err) {
		return err
	}

	v.Editable = false
	if h.body == nil {
		h.body = h.getBody(keys)
	}
	v.Title = h.name
	v.Write([]byte(h.body))
	_, err = g.SetCurrentView("equalizer")
	return err
}

func (h *equalizer) hide(g *ui.Gui, v *ui.View) error {
	v.Clear()
	return g.DeleteView(h.name)
}

func (h *equalizer) next(g *ui.Gui, v *ui.View) error {
	log.Println("next band")
	return nil
}

func (h *equalizer) getBody(keys []key) []byte {
	out := &bytes.Buffer{}
	return []byte(out.String())
}
