package views

import (
	"fmt"
	"os"

	"github.com/cswank/music/internal/colors"
	ui "github.com/jroimartin/gocui"
)

var (
	g *ui.Gui

	bg         ui.Attribute
	c1, c2, c3 colors.Colorer
)

func init() {
	bg, c1, c2, c3 = getColors()
}

type coords struct {
	x1 int
	x2 int
	y1 int
	y2 int
}

type gui struct {
}

func newGUI() error {
	var err error
	g, err = ui.NewGui(ui.Output256)
	if err != nil {
		return fmt.Errorf("could not create gui: %s", err)
	}

	s, err := newScreen()
	if err != nil {
		return err
	}

	w, h := g.Size()
	g.SetManagerFunc(s.getLayout(w, h))
}

func Start() error {
	return nil
}

func getColors() (ui.Attribute, colors.Colorer, colors.Colorer, colors.Colorer) {
	bg = colors.GetBackground(os.Getenv("MUSIC_COLOR0"))
	c1 := colors.Get(os.Getenv("MUSIC_COLOR1"))
	if c1 == nil {
		c1 = colors.Get("white")
	}
	c2 := colors.Get(os.Getenv("MUSIC_COLOR2"))
	if c2 == nil {
		c2 = colors.Get("green")
	}
	c3 := colors.Get(os.Getenv("MUSIC_COLOR3"))
	if c3 == nil {
		c3 = colors.Get("yellow")
	}
	return bg, c1, c2, c3
}
