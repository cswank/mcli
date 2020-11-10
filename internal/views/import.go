package views

import (
	ui "github.com/awesome-gocui/gocui"
)

const (
	importWidth  = 30
	importHeight = 2
)

type albumImport struct {
	name   string
	coords coords
	body   []byte
}

func newImport(w, h int) *albumImport {
	return &albumImport{
		name: "import",
	}
}

func getImportCoords(g *ui.Gui) coords {
	maxX, maxY := g.Size()
	x1 := maxX/2 - importWidth/2
	x2 := maxX/2 + importWidth/2
	y1 := maxY/2 - importHeight/2
	y2 := maxY/2 + importHeight/2 + importHeight%2
	return coords{x1: x1, y1: y1, x2: x2, y2: y2}
}

func (i *albumImport) render(g *ui.Gui, v *ui.View) error {
	v.Editable = false
	v.Frame = true
	v.Title = "import"
	return nil
}
