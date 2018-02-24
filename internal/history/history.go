package history

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"bitbucket.org/cswank/music/internal/source"
)

type History interface {
	Save(source.Result) error
	Fetch(int, int) (*source.Results, error)
}

type FileHistory struct {
	pth string
}

func NewFileHistory() (*FileHistory, error) {
	pth := fmt.Sprintf("%s/.music/history.csv", os.Getenv("HOME"))
	e, err := exists(pth)
	if err != nil {
		return nil, err
	}

	if !e {
		f, err := os.Create(pth)
		if err != nil {
			return nil, err
		}
		f.Close()
	}

	return &FileHistory{pth: pth}, nil
}

func (f *FileHistory) Save(r source.Result) error {
	file, err := os.OpenFile(f.pth, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)

	if err := w.Write(r.ToCSV()); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (f *FileHistory) Fetch(page, pageSize int) (*source.Results, error) {
	file, err := os.Open(f.pth)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var res []source.Result
	var maxTitle, maxAlbum int
	seen := map[string]bool{}
	for _, row := range rows {
		if len(res) == pageSize {
			break
		}
		r := &source.Result{}
		if err := r.FromCSV(row); err != nil {
			return nil, err
		}
		if seen[r.Track.ID] {
			continue
		} else {
			seen[r.Track.ID] = true
		}
		if len(r.Track.Title) > maxTitle {
			maxTitle = len(r.Track.Title)
		}
		if len(r.Album.Title) > maxAlbum {
			maxAlbum = len(r.Album.Title)
		}
		res = append(res, *r)
	}

	format := fmt.Sprintf("%%-%ds%%-%ds%%s\n", maxTitle+4, maxAlbum+4)
	return &source.Results{
		Header: fmt.Sprintf(format, "Title", "Album", "Artist"),
		Type:   "album",
		Print: func(w io.Writer, r source.Result) error {
			_, err := fmt.Fprintf(w, format, r.Track.Title, r.Album.Title, r.Artist.Name)
			return err
		},
		Results: res,
	}, nil
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
