package play

import (
	"github.com/cswank/mcli/internal/schema"
)

type Sort string

const (
	Time  Sort = "time"
	Count Sort = "count"
)

type Player interface {
	Play(schema.Result)
	PlayAlbum(*schema.Results)
	Volume(float64) float64
	Pause()
	FastForward()
	Seek(i int)
	Rewind()
	Queue() *schema.Results
	RemoveFromQueue([]int)
	NextSong(id string, f func(schema.Result))
	PlayProgress(id string, f func(schema.Progress))
	DownloadProgress(id string, f func(schema.Progress))
	Done(string)
	Close()
}
