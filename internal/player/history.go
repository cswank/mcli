package player

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type History interface {
	Save(Result) error
	Fetch(int, int) (*Results, error)
}

type FileHistory struct {
	pth        string
	archivePth string
	cache      []Result
	fmt        string
}

func NewFileHistory() (*FileHistory, error) {
	pth := fmt.Sprintf("%s/history.csv", os.Getenv("MCLI_HOME"))
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

	aPth := fmt.Sprintf("%s/history-archive", os.Getenv("MCLI_HOME"))

	e, err = exists(aPth)
	if err != nil {
		return nil, err
	}

	if !e {
		err := os.Mkdir(aPth, 0700)
		if err != nil {
			return nil, err
		}
	}

	return &FileHistory{pth: pth, archivePth: aPth}, nil
}

func (f *FileHistory) Save(r Result) error {
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
	f.cache = []Result{}
	return nil
}

func (f *FileHistory) Fetch(page, pageSize int) (*Results, error) {
	if len(f.cache) != 0 {
		return f.historyFromCache(page, pageSize)
	}

	file, err := os.Open(f.pth)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] > rows[j][0]
	})

	var res []Result
	var maxTitle, maxAlbum int
	seen := map[string]bool{}
	for _, row := range rows {
		r := &Result{}
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
	f.cache = res
	f.fmt = fmt.Sprintf("%%-%ds%%-%ds%%s\n", maxTitle+4, maxAlbum+4)

	if len(rows) > 1000 {
		if err := f.archive(rows[900:]); err != nil {
			return nil, err
		}
	}
	return f.historyFromCache(page, pageSize)

}

func (f *FileHistory) historyFromCache(page, pageSize int) (*Results, error) {
	start := page * pageSize
	if start > len(f.cache) {
		return nil, fmt.Errorf("no more pages")
	}

	end := start + pageSize
	if end > len(f.cache) {
		end = len(f.cache) - 1
	}

	return &Results{
		Header:  fmt.Sprintf(f.fmt, "Title", "Album", "Artist"),
		Type:    "history",
		Fmt:     f.fmt,
		Results: f.cache[start:end],
	}, nil
}

func (f *FileHistory) archive(rows [][]string) error {
	if err := os.Rename(f.pth, filepath.Join(f.archivePth, fmt.Sprintf("%s.csv", time.Now().Format(time.RFC3339)))); err != nil {
		return err
	}

	file, err := os.Create(f.pth)
	if err != nil {
		return err
	}
	w := csv.NewWriter(file)
	if err := w.WriteAll(rows); err != nil {
		return err
	}
	w.Flush()
	return file.Close()
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
