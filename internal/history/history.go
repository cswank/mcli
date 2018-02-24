package history

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

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

	if err := w.Write([]string{time.Now().Format(time.RFC3339), r.Title, r.Album, r.Artist, r.ID, strconv.Itoa(r.Duration)}); err != nil {
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
	start := page * pageSize
	if start >= len(rows) {
		return nil, nil
	}
	end := start + pageSize
	if end >= len(rows) {
		end = len(rows)
	}

	rows = rows[start:end]
	res := make([]source.Result, len(rows))
	var maxTitle, maxAlbum int
	for i, row := range rows {
		if len(row) < 6 {
			continue
		}
		if len(row[1]) > maxTitle {
			maxTitle = len(row[1])
		}
		if len(row[2]) > maxAlbum {
			maxAlbum = len(row[2])
		}

		d, err := strconv.ParseInt(row[5], 10, 64)
		if err != nil {
			return nil, err
		}
		res[i] = source.Result{
			ID:       row[4],
			Title:    row[1],
			Album:    row[2],
			Artist:   row[3],
			Duration: int(d),
		}
	}

	format := fmt.Sprintf("%%-%ds%%-%ds%%s\n", maxTitle+4, maxAlbum)
	return &source.Results{
		Header: fmt.Sprintf(format, "Title", "Album", "Artist"),
		Type:   "album",
		Print: func(w io.Writer, r source.Result) error {
			_, err := fmt.Fprintf(w, format, r.Title, r.Album, r.Artist)
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
