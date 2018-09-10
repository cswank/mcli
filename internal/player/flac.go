package player

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cswank/beep"

	"github.com/cswank/beep/effects"
	"github.com/cswank/beep/flac"
	"github.com/cswank/beep/speaker"
)

type Flac struct {
	Fetcher
	queue         *queue
	history       History
	playing       bool
	sep           string
	pause         chan bool
	vol           chan float64
	volOut        chan float64
	volume        float64
	fastForward   chan bool
	rewind        chan bool
	playCB        map[string]func(Progress)
	downloadCB    map[string]func(Progress)
	nextSongCB    map[string]func(r Result)
	onDeck        chan song
	onDeckResult  *Result
	currentResult *Result
	playLock      sync.Mutex
	downloadLock  sync.Mutex
	nextSongLock  sync.Mutex
	cacheOnDisk   bool
}

type flacSettings struct {
	Volume float64 `json:"volume"`
}

func getFlacPath() string {
	return fmt.Sprintf("%s/flac.json", os.Getenv("MCLI_HOME"))
}

func NewFlac(f Fetcher, cache bool) (*Flac, error) {
	hist, err := NewStormHistory()
	if err != nil {
		return nil, err
	}

	pth := getFlacPath()
	e, err := exists(pth)
	if err != nil {
		return nil, err
	}

	var s flacSettings
	if e {
		f, err := os.Open(pth)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		err = json.NewDecoder(f).Decode(&s)
	} else {
		f, err := os.Create(pth)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		json.NewEncoder(f).Encode(s)
	}

	p := &Flac{
		Fetcher:     f,
		history:     hist,
		sep:         string(filepath.Separator),
		queue:       newQueue(),
		pause:       make(chan bool),
		fastForward: make(chan bool),
		rewind:      make(chan bool),
		vol:         make(chan float64),
		volOut:      make(chan float64),
		onDeck:      make(chan song),
		volume:      s.Volume,
		playCB:      make(map[string]func(Progress)),
		downloadCB:  make(map[string]func(Progress)),
		nextSongCB:  make(map[string]func(Result)),
		cacheOnDisk: cache,
	}

	go p.playLoop()
	go p.downloadLoop()
	return p, nil
}

var FlacCacheOnDisk = func(f *Flac) {
	f.cacheOnDisk = true
}

func (f *Flac) NextSong(id string, fn func(Result)) {
	f.nextSongCB[id] = fn
	if fn != nil && f.currentResult != nil {
		fn(*f.currentResult)
	}
}

func (f *Flac) callNextSong() {
	f.nextSongLock.Lock()
	for id, fn := range f.nextSongCB {
		if fn != nil && f.currentResult != nil {
			fn(*f.currentResult)
		} else if fn == nil {
			delete(f.nextSongCB, id)
		}
	}
	f.nextSongLock.Unlock()
}

func (f *Flac) Play(r Result) {
	f.queue.Add(r)
}

func (f *Flac) History(page, pageSize int, sort Sort) (*Results, error) {
	return f.history.Fetch(page, pageSize, sort)
}

func (f *Flac) PlayAlbum(album *Results) {
	for _, r := range album.Results {
		f.Play(r)
	}
}

func (f *Flac) Pause() {
	if f.playing {
		f.pause <- true
	}
}

func (f *Flac) Volume(v float64) float64 {
	var out float64
	if f.playing {
		f.vol <- v
		out = <-f.volOut
	} else {
		f.volume += v
		out = f.volume
	}

	return out
}

func (f *Flac) Queue() *Results {
	var r []Result
	if f.onDeckResult != nil {
		r = []Result{*f.onDeckResult}
	}
	return &Results{
		Results: append(r, f.queue.Playlist()...),
	}
}

func (f *Flac) RemoveFromQueue(indices []int) {
	sort.Sort(sort.Reverse(sort.IntSlice(indices)))
	for _, i := range indices {
		if i == 0 {
			<-f.onDeck
		} else {
			f.queue.Remove(i - 1)
		}
	}
}

func (f *Flac) Done(id string) {
	f.playLock.Lock()
	delete(f.playCB, id)
	for k, v := range f.playCB {
		if v == nil {
			delete(f.playCB, k)
		}
	}
	f.playLock.Unlock()
	f.downloadLock.Lock()
	delete(f.downloadCB, id)
	for k, v := range f.downloadCB {
		if v == nil {
			delete(f.downloadCB, k)
		}
	}
	f.downloadLock.Unlock()

	f.nextSongLock.Lock()
	delete(f.nextSongCB, id)
	for k, v := range f.nextSongCB {
		if v == nil {
			delete(f.nextSongCB, k)
		}
	}
	f.nextSongLock.Unlock()
}

func (f *Flac) Close() {
	file, err := os.Create(getFlacPath())
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(flacSettings{Volume: f.volume})
}

func (f *Flac) FastForward() {
	if f.playing {
		f.fastForward <- true
	}
}

func (f *Flac) Rewind() {
	if f.playing {
		f.rewind <- true
	}
}

func (f *Flac) downloadLoop() {
	for {
		r := f.queue.Next()
		f.onDeckResult = &r
		song := f.download(&r)
		f.onDeck <- *song
		f.onDeckResult = nil
	}
}

func (f *Flac) playLoop() {
	for {
		s := <-f.onDeck
		f.playing = true
		if err := f.doPlay(s); err != nil {
			log.Fatal(err)
		}
		f.playing = false
	}
}

func (f *Flac) doPlay(s song) error {
	f.currentResult = &s.result
	f.callNextSong()

	if err := f.history.Save(s.result); err != nil {
		return err
	}

	music, format, err := flac.Decode(s.reader())
	if err != nil {
		return err
	}

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/2)); err != nil {
		return err
	}

	vol := &effects.Volume{
		Streamer: music,
		Base:     2,
		Volume:   f.volume,
	}

	ctrl := &beep.Ctrl{
		Streamer: vol,
	}
	speaker.Play(ctrl)

	var done bool
	var paused bool
	l := music.Len()
	var i int
	for !done {
		select {
		case <-time.After(500 * time.Millisecond):
			pos := music.Position()
			done = pos >= l
			i++
			f.playLock.Lock()
			for k, cb := range f.playCB {
				if cb != nil {
					cb(Progress{N: pos, Total: l})
				} else {
					delete(f.playCB, k)
				}
			}
			f.playLock.Unlock()
		case v := <-f.vol:
			speaker.Lock()
			if (f.volume < 2.0 && v > 0.0) || (f.volume > -5.0 && v < 0.0) {
				vol.Volume += v
				f.volume = vol.Volume
			}
			speaker.Unlock()
			f.volOut <- f.volume
		case <-f.pause:
			paused = !paused
			speaker.Lock()
			ctrl.Paused = paused
			speaker.Unlock()
		case <-f.fastForward:
			done = true
		case <-f.rewind:
			music.Close()
			s.reset()
			return f.doPlay(s)
		}
	}

	f.currentResult = nil
	return music.Close()
}

func (f *Flac) DownloadProgress(id string, fn func(Progress)) {
	f.downloadLock.Lock()
	f.downloadCB[id] = fn
	f.downloadLock.Unlock()
}

func (f *Flac) PlayProgress(id string, fn func(Progress)) {
	f.playLock.Lock()
	f.playCB[id] = fn
	f.playLock.Unlock()
}

func (f *Flac) download(r *Result) *song {
	var out *song
	pth, e := f.checkCache(*r)
	r.Path = pth
	if !e || !f.cacheOnDisk {
		s, err := f.doDownload(*r)
		if err != nil {
			log.Fatal(err)
		}
		out = s
	}

	return out
}

func (f *Flac) doDownload(r Result) (*song, error) {
	out := &song{
		result: r,
	}
	u, err := f.Fetcher.GetTrack(r.Track.ID)
	if err != nil {
		return out, fmt.Errorf("could not get track %+v: %s", r, err)
	}

	resp, err := f.getTrack(u)
	if err != nil {
		return out, fmt.Errorf("could not get stream %+v: %s", r, err)
	}

	defer resp.Body.Close()

	if f.cacheOnDisk {
		file, err := ioutil.TempFile(fmt.Sprintf("%s/tmp", os.Getenv("MCLI_HOME")), "")
		if err != nil {
			return out, fmt.Errorf("could not create file for %+v: %s", r, err)
		}
		defer file.Close()
		if err != nil {
			os.Rename(file.Name(), r.Path)
		}
		out.file = file
	} else {
		out.buf = &bytes.Buffer{}
	}

	pr := newProgressRead(resp.Body, int(resp.ContentLength), f.downloadCB)
	w := out.writer()
	_, err = io.Copy(w, pr)
	if err != nil {
		return out, err
	}
	pr.Close()

	return out, nil
}

func (f *Flac) getTrack(uri string) (*http.Response, error) {
	if strings.Index(uri, "file://") == 0 {
		f, err := os.Open(uri[7:])
		return &http.Response{
			Body: f,
		}, err
	}

	return http.Get(uri)
}

func (f *Flac) ensureCache() error {
	dir := fmt.Sprintf("%s/cache/%s", os.Getenv("MCLI_HOME"), f.Fetcher.Name())
	e, _ := exists(dir)
	if e {
		return nil
	}
	return os.MkdirAll(dir, 0700)
}

func (f *Flac) checkCache(result Result) (string, bool) {
	if !f.cacheOnDisk {
		return "", false
	}
	dir := fmt.Sprintf("%s/cache/%s/%s/%s", os.Getenv("MCLI_HOME"), f.Fetcher.Name(), f.clean(result.Artist.Name), f.clean(result.Album.Title))
	e, _ := exists(dir)
	if !e {
		os.MkdirAll(dir, 0700)
	}

	dir = fmt.Sprintf("%s/tmp", os.Getenv("MCLI_HOME"))
	e, _ = exists(dir)
	if !e {
		os.MkdirAll(dir, 0700)
	}

	pth := fmt.Sprintf("%s/cache/%s/%s/%s/%s.flac", os.Getenv("MCLI_HOME"), f.Fetcher.Name(), f.clean(result.Artist.Name), f.clean(result.Album.Title), f.clean(result.Track.Title))
	e, _ = exists(pth)
	return pth, e
}

func (f *Flac) clean(s string) string {
	return strings.Replace(s, f.sep, "", -1)
}

type progressRead struct {
	io.Reader
	t, l, reads int
	cb          map[string]func(Progress)
}

func newProgressRead(r io.Reader, l int, cb map[string]func(Progress)) *progressRead {
	return &progressRead{Reader: r, t: 0, l: l, cb: cb}
}

func (r *progressRead) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.t += n
	r.reads++
	if r.cb != nil && r.reads%100 == 0 {
		for k, cb := range r.cb {
			if cb != nil {
				cb(Progress{N: r.t, Total: r.l})
			} else {
				delete(r.cb, k)
			}
		}
	}
	return n, err
}

// Close the reader when it implements io.Closer
func (r *progressRead) Close() error {
	for _, cb := range r.cb {
		cb(Progress{N: 0, Total: r.t})
	}

	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

type songBuffer struct {
	closed bool
	buf    *bytes.Reader
}

func (s *songBuffer) Read(p []byte) (n int, err error) {
	if s.closed {
		s.closed = false
		return 0, io.EOF
	}
	return s.buf.Read(p)
}

func (s *songBuffer) Seek(offset int64, whence int) (int64, error) {
	return s.buf.Seek(offset, whence)
}

func (s *songBuffer) Close() error {
	s.closed = true
	return nil
}

type song struct {
	result Result
	file   *os.File
	buf    *bytes.Buffer
}

func (s *song) reader() io.ReadCloser {
	if s.file != nil {
		return s.file
	}

	return &songBuffer{
		buf: bytes.NewReader(s.buf.Bytes()),
	}
}

func (s *song) writer() io.Writer {
	if s.file != nil {
		return s.file
	}
	return s.buf
}

func (s *song) reset() error {
	if s.file != nil {
		_, err := s.file.Seek(0, 0)
		return err
	}
	return nil
}
