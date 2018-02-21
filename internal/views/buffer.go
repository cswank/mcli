package views

import "log"

type buffer struct {
	coords   coords
	progress chan progress
}

func newBuffer(w, h int) *buffer {
	b := &buffer{
		coords:   coords{x1: -1, y1: h - 2, x2: w - 1, y2: h},
		progress: make(chan progress),
	}

	go b.render(b.progress)
	return b
}

func (b *buffer) render(ch <-chan progress) {
	for {
		p := <-ch
		log.Println(p)
	}
}
