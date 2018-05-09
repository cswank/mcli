package http

import (
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"bitbucket.org/cswank/mcli/internal/player"

	"github.com/GeertJohan/go.rice"
)

var (
	funcs = template.FuncMap{
		"duration": func(d int) string {
			return player.FmtDuration(time.Duration(d) * time.Second)
		},
	}
)

type tmpl struct {
	template    *template.Template
	files       []string
	scripts     []string
	stylesheets []string
	funcs       template.FuncMap
	bare        bool
}

type templates struct {
	templates map[string]tmpl
}

func newTemplates(box *rice.Box) (*templates, error) {
	data := map[string]string{}
	html := getHTML(box)
	for _, pth := range html {
		s, err := box.String(pth)
		if err != nil {
			return nil, err
		}
		data[pth] = s
	}

	tmpls := map[string]tmpl{
		"history.html":     {},
		"playlists.html":   {},
		"album.html":       {funcs: funcs, files: []string{"row.html"}},
		"queue.html":       {},
		"artists.html":     {},
		"tracks.html":      {},
		"albums.html":      {},
		"search.html":      {},
		"search-form.html": {},
	}

	base := []string{"head.html", "base.html", "navbar.html", "footer.html", "base.js"}

	for key, val := range tmpls {
		t := template.New(key)
		if val.funcs != nil {
			t = t.Funcs(val.funcs)
		}
		var err error
		files := append([]string{key}, val.files...)
		files = append(files, base...)
		for _, f := range files {
			t, err = t.Parse(data[f])
			if err != nil {
				log.Fatal(err)
			}
		}
		val.template = t
		tmpls[key] = val
	}

	return &templates{
		templates: tmpls,
	}, nil
}

func getHTML(box *rice.Box) []string {
	var html []string
	box.Walk("/", func(pth string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(pth, ".html") || strings.HasSuffix(pth, ".js") {
			if box.IsEmbedded() {
				pth = pth[1:] //workaround until https://github.com/GeertJohan/go.rice/issues/71 is fixed (which is probably never)
			}
			html = append(html, pth)
		}
		return nil
	})
	return html
}

func (t *templates) get(k string) (*template.Template, []string, []string) {
	tpl := t.templates[k]
	return tpl.template, tpl.scripts, tpl.stylesheets
}
