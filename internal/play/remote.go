package play

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/cswank/mcli/internal/rpc"
	"github.com/cswank/mcli/internal/schema"
)

type Remote struct {
	client rpc.PlayerClient
}

func NewRemote(conn *grpc.ClientConn) *Remote {
	return &Remote{
		client: rpc.NewPlayerClient(conn),
	}
}

func (r Remote) Done(id string) {
	r.client.Done(context.Background(), &rpc.String{Value: id})
}

func (r Remote) Close() {}

func (r Remote) Play(rs schema.Result) {
	_, err := r.client.Play(context.Background(), rpc.PBFromResult(rs))
	if err != nil {
		log.Println(err)
	}
}

func (r Remote) PlayAlbum(rs *schema.Results) {
	_, err := r.client.PlayAlbum(context.Background(), rpc.PBFromResults(rs))
	if err != nil {
		log.Println(err)
	}
}

func (r Remote) Volume(v float64) float64 {
	f, err := r.client.Volume(context.Background(), &rpc.Float{Value: float32(v)})
	if err != nil {
		log.Println(err)
	}
	return float64(f.Value)
}

func (r Remote) Pause() {
	_, err := r.client.Pause(context.Background(), &rpc.Empty{})
	if err != nil {
		log.Println(err)
	}
}

func (r Remote) FastForward() {
	_, err := r.client.FastForward(context.Background(), &rpc.Empty{})
	if err != nil {
		log.Println(err)
	}
}

func (r Remote) Rewind() {
	_, err := r.client.Rewind(context.Background(), &rpc.Empty{})
	if err != nil {
		log.Println(err)
	}
}

func (r Remote) Queue() *schema.Results {
	out, err := r.client.Queue(context.Background(), &rpc.Empty{})
	if err != nil {
		log.Println(err)
		return nil
	}
	return rpc.ResultsFromPB(out)
}

func (r Remote) RemoveFromQueue(indices []int) {
	out := make([]int64, len(indices))
	for i, val := range indices {
		out[i] = int64(val)
	}
	_, err := r.client.RemoveFromQueue(context.Background(), &rpc.Ints{Value: out})
	if err != nil {
		log.Println(err)
	}
}

func (r Remote) NextSong(id string, f func(schema.Result)) {
	go func() {
		stream, err := r.client.NextSong(context.Background(), &rpc.String{Value: id})
		if err != nil {
			log.Fatal("could not get stream for next song", err)
		}
		for {
			ns, err := stream.Recv()
			if err == io.EOF {
				time.Sleep(time.Second)
			} else if err != nil {
				log.Println(err)
			} else {
				f(rpc.ResultFromPB(ns))
			}
		}
	}()
}

func (r Remote) PlayProgress(id string, f func(schema.Progress)) {
	go func() {
		stream, err := r.client.PlayProgress(context.Background(), &rpc.String{Value: id})
		if err != nil {
			log.Fatal("could not get stream for next song", err)
		}
		for {
			p, err := stream.Recv()
			if err == io.EOF {
				time.Sleep(time.Second)
			} else if err != nil {
				log.Println(err)
			} else {
				f(rpc.ProgressFromPB(p))
			}
		}
	}()
}

func (r Remote) DownloadProgress(id string, f func(schema.Progress)) {
	go func() {
		stream, err := r.client.DownloadProgress(context.Background(), &rpc.String{Value: id})
		if err != nil {
			log.Fatal("could not get stream for next song", err)
		}
		for {
			p, err := stream.Recv()
			if err == io.EOF {
				time.Sleep(time.Second)
			} else if err != nil {
				log.Println(err)
			} else {
				f(rpc.ProgressFromPB(p))
			}
		}
	}()
}
