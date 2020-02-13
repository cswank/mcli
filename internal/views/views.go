package views

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/cswank/mcli/internal/colors"
	"bitbucket.org/cswank/mcli/internal/fetch"
	"bitbucket.org/cswank/mcli/internal/play"
	"bitbucket.org/cswank/mcli/internal/repo"
	ui "github.com/jroimartin/gocui"
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
	repo.History
}

// Start is what main calls to get the app rolling
func Start(p play.Player, f fetch.Fetcher, hist repo.History) error {
	col = colors.Get()
	cli := &client{Player: p, Fetcher: f, History: hist}
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

	cli.Done(s.id)
	g.Close()
	return nil
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
