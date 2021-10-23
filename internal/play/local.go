package play

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cswank/mcli/internal/download"
	"github.com/cswank/mcli/internal/history"
	"github.com/cswank/mcli/internal/queue"
	"github.com/cswank/mcli/internal/schema"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"

	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/speaker"
)

type Local struct {
	queue         *queue.Queue
	history       history.History
	playing       bool
	sep           string
	pause         chan bool
	seek          chan int
	vol           chan float64
	volOut        chan float64
	volume        float64
	fastForward   chan bool
	rewind        chan bool
	playCB        func(schema.Progress)
	downloadCB    func(schema.Progress)
	nextSongCB    func(r schema.Result)
	onDeck        chan song
	onDeckResult  *schema.Result
	currentResult *schema.Result
	dl            download.Downloader
	pth           string
	sampleRate    int
}

type flacSettings struct {
	Volume float64 `json:"volume"`
}

func NewLocal(dir, home string, opts ...func(*Local)) (*Local, error) {
	e, err := exists(dir)
	if err != nil {
		return nil, err
	}

	pth := fmt.Sprintf("%s/flac.json", home)
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

	if err := speaker.Init(44100, 44100/2); err != nil {
		return nil, fmt.Errorf("unable to init speaker: %s", err)
	}

	l := &Local{
		sep:         string(filepath.Separator),
		queue:       queue.New(),
		pause:       make(chan bool),
		fastForward: make(chan bool),
		rewind:      make(chan bool),
		seek:        make(chan int),
		vol:         make(chan float64),
		volOut:      make(chan float64),
		onDeck:      make(chan song),
		volume:      s.Volume,
		pth:         pth,
		sampleRate:  44100,
	}

	for _, opt := range opts {
		opt(l)
	}

	go l.playLoop()
	go l.downloadLoop()
	return l, nil
}

func LocalDownload(dl download.Downloader) func(*Local) {
	return func(l *Local) {
		l.dl = dl
	}
}

func LocalHistory(h history.History) func(*Local) {
	return func(l *Local) {
		l.history = h
	}
}

func (l *Local) NextSong(id string, fn func(schema.Result)) {
	l.nextSongCB = fn
}

func (l *Local) callNextSong() {
	if l.currentResult != nil && l.nextSongCB != nil {
		l.nextSongCB(*l.currentResult)
	}
}

func (l *Local) Play(r schema.Result) {
	l.queue.Add(r)
}

func (l *Local) PlayAlbum(album *schema.Results) {
	for _, r := range album.Results {
		l.Play(r)
	}
}

func (l *Local) Pause() {
	if l.playing {
		l.pause <- true
	}
}

func (l *Local) Volume(v float64) float64 {
	var out float64
	if l.playing {
		l.vol <- v
		out = <-l.volOut
	} else {
		l.volume += v
		out = l.volume
	}

	return out
}

func (l *Local) Queue() *schema.Results {
	var r []schema.Result
	if l.onDeckResult != nil {
		r = []schema.Result{*l.onDeckResult}
	}
	return &schema.Results{
		Results: append(r, l.queue.Queue()...),
	}
}

func (l *Local) RemoveFromQueue(indices []int) {
	sort.Sort(sort.Reverse(sort.IntSlice(indices)))
	for _, i := range indices {
		if i == 0 {
			<-l.onDeck
		} else {
			l.queue.Remove(i - 1)
		}
	}
}

func (l *Local) Done(id string) {

}

func (l *Local) Close() {
	file, err := os.Create(l.pth)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(flacSettings{Volume: l.volume})
}

func (l *Local) FastForward() {
	if l.playing {
		l.fastForward <- true
	}
}

func (l *Local) Seek(i int) {
	if l.playing {
		l.seek <- i
	}
}

func (l *Local) Rewind() {
	if l.playing {
		l.rewind <- true
	}
}

func (l *Local) downloadLoop() {
	for {
		r := l.queue.Next()
		l.onDeckResult = &r
		song := l.download(&r)
		if song == nil {
			continue
		}
		l.onDeck <- *song
		l.onDeckResult = nil
	}
}

func (l *Local) playLoop() {
	for {
		s := <-l.onDeck
		l.playing = true
		if err := l.doPlay(s); err != nil {
			l.nextSongCB(schema.Result{Error: err.Error()})
			log.Printf("error playing song %+v: %s", s, err)
		}
		l.playing = false
	}
}

func (l *Local) doPlay(s song) error {
	music, format, err := flac.Decode(s.reader())
	if err != nil {
		return err
	}

	if int(format.SampleRate) != l.sampleRate {
		speaker.Close()
		if err := speaker.Init(format.SampleRate, int(format.SampleRate)/2); err == nil {
			return fmt.Errorf("unable to init speaker: %s", err)
		}
		l.sampleRate = int(format.SampleRate)
	}

	ln := music.Len()
	s.result.Track.Duration = ln / int(format.SampleRate)
	l.currentResult = &s.result
	l.callNextSong()

	if err := l.history.Save(s.result); err != nil {
		return fmt.Errorf("unable to save history for result %+v: %s", s.result, err)
	}

	vol := &effects.Volume{
		Streamer: music,
		Base:     2,
		Volume:   l.volume,
	}

	ctrl := &beep.Ctrl{
		Streamer: vol,
	}

	speaker.Play(ctrl)

	var done bool
	var paused bool
	var i int
	for !done {
		select {
		case <-time.After(500 * time.Millisecond):
			pos := music.Position()
			done = pos >= ln
			i++
			if l.playCB != nil {
				l.playCB(schema.Progress{N: pos, Total: ln})
			}
		case v := <-l.vol:
			speaker.Lock()
			if (l.volume < 2.0 && v > 0.0) || (l.volume > -5.0 && v < 0.0) {
				vol.Volume += v
				l.volume = vol.Volume
			}
			speaker.Unlock()
			l.volOut <- l.volume
		case <-l.pause:
			paused = !paused
			speaker.Lock()
			ctrl.Paused = paused
			speaker.Unlock()
		case <-l.fastForward:
			done = true
		case i := <-l.seek:
			speaker.Lock()
			t := music.Len() / int(format.SampleRate)
			pos := music.Len() * i / t
			music.Seek(pos)
			speaker.Unlock()
		case <-l.rewind:
			music.Close()
			s.r.Seek(0, 0)
			return l.doPlay(s)
		}
	}

	l.currentResult = nil
	return music.Close()
}

func (l *Local) DownloadProgress(id string, fn func(schema.Progress)) {
	l.downloadCB = fn
}

func (l *Local) PlayProgress(id string, fn func(schema.Progress)) {
	l.playCB = fn
}

func (l *Local) download(r *schema.Result) *song {
	out, err := l.doDownload(*r)
	if err != nil {
		log.Printf("error downloading song: %s", err)
		return nil
	}
	return out
}

func (l *Local) doDownload(rs schema.Result) (*song, error) {
	buf := bytes.Buffer{}
	l.dl.Download(rs.Track.ID, &buf, l.downloadCB)
	return &song{
		result: rs,
		r:      bytes.NewReader(buf.Bytes()),
	}, nil
}

func (l *Local) clean(s string) string {
	return strings.Replace(s, l.sep, "", -1)
}

type progressRead struct {
	io.Reader
	t, l, reads int
	cb          map[string]func(schema.Progress)
}

func newProgressRead(r io.Reader, l int, cb map[string]func(schema.Progress)) *progressRead {
	return &progressRead{Reader: r, t: 0, l: l, cb: cb}
}

func (r *progressRead) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.t += n
	r.reads++
	if r.cb != nil && r.reads%100 == 0 {
		for k, cb := range r.cb {
			if cb != nil {
				cb(schema.Progress{N: r.t, Total: r.l})
			} else {
				delete(r.cb, k)
			}
		}
	}
	return n, err
}

// Close the reader when it implements io.Closer
func (r *progressRead) Close() error {
	for k, cb := range r.cb {
		if cb != nil {
			cb(schema.Progress{N: 0, Total: r.t})
		} else {
			delete(r.cb, k)
		}
	}

	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

type songBuffer struct {
	closed bool
	buf    io.ReadSeeker
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
	result schema.Result
	r      io.ReadSeeker
}

func (s *song) reader() io.ReadSeeker {
	return &songBuffer{
		buf: s.r,
	}
}

func (s *song) reset() error {
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
