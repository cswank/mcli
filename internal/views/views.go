package views

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/cswank/mcli/internal/colors"
	"bitbucket.org/cswank/mcli/internal/player"
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

//Start is what main calls to get the app rolling
func Start(cli player.Player) error {
	dir := os.Getenv("MCLI_HOME")
	e, err := exists(dir)
	if err != nil {
		return err
	}

	if !e {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}

	g, err = ui.NewGui(ui.Output256)
	if err != nil {
		return fmt.Errorf("could not create gui: %s", err)
	}

	w, h := g.Size()
	s, err := newScreen(w, h, cli)
	if err != nil {
		return err
	}

	g.SetManagerFunc(s.getLayout(w, h))
	g.Cursor = true
	g.InputEsc = true

	if err := s.keybindings(g); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil {
		if err != ui.ErrQuit {
			log.SetOutput(os.Stderr)
			log.Println(err)
			return err
		}
	}

	g.Close()
	return nil
}

func getColors() (ui.Attribute, colors.Colorer, colors.Colorer, colors.Colorer) {
	bg = colors.GetBackground(os.Getenv("MCLI_COLOR0"))
	c1 := colors.Get(os.Getenv("MCLI_COLOR1"))
	if c1 == nil {
		c1 = colors.Get("white")
	}
	c2 := colors.Get(os.Getenv("MCLI_COLOR2"))
	if c2 == nil {
		c2 = colors.Get("green")
	}
	c3 := colors.Get(os.Getenv("MCLI_COLOR3"))
	if c3 == nil {
		c3 = colors.Get("yellow")
	}
	return bg, c1, c2, c3
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
