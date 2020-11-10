package views

import (
	"fmt"
	"log"
	"os"

	ui "github.com/awesome-gocui/gocui"
	"github.com/cswank/mcli/internal/colors"
	"github.com/cswank/mcli/internal/fetch"
	hist "github.com/cswank/mcli/internal/history"
	"github.com/cswank/mcli/internal/play"
)

var (
	g   *ui.Gui
	col colors.Colors
)

type coords struct {
	x1 int
	x2 int
	y1 int
	y2 int
}

type client struct {
	play.Player
	fetch.Fetcher
	hist.History
}

// Start is what main calls to get the app rolling
func Start(p play.Player, f fetch.Fetcher, hist hist.History) error {
	col = colors.New()
	cli := &client{Player: p, Fetcher: f, History: hist}

	var err error
	g, err = ui.NewGui(ui.Output256, false)
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

	cli.Done(s.id)
	g.Close()
	return nil
}
