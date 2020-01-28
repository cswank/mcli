package rpc

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "bitbucket.org/cswank/mcli/internal/play/rpc"
	"bitbucket.org/cswank/mcli/internal/player"
)

type RPC struct {
	conn   *grpc.ClientConn
	client pb.PlayerClient
}

func NewRPC(addr string) (*RPC, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &RPC{
		conn:   conn,
		client: pb.NewPlayerClient(conn),
	}, nil
}

func (c *Client) Done(id string) {
	c.client.Done(context.Background(), &pb.String{Value: id})
	c.conn.Close()
}

func (r RPC) Close() {}

func (r RPC) Play(r player.Result) {
	_, err := r.client.PlayPlay(context.Background(), PBFromResult(r))
	if err != nil {
		log.Println(err)
	}
}

func (r RPC) PlayAlbum(r *player.Results) {
	_, err := r.client.PlayAlbum(context.Background(), PBFromResults(r))
	if err != nil {
		log.Println(err)
	}
}

func (r RPC) Volume(v float64) float64 {
	f, err := r.client.Volume(context.Background(), &pb.Float{Value: float32(v)})
	if err != nil {
		log.Println(err)
	}
	return float64(f.Value)
}

func (r RPC) Pause() {
	_, err := r.client.Pause(context.Background(), &pb.Empty{})
	if err != nil {
		log.Println(err)
	}
}

func (r *RPC) FastForward() {
	_, err := r.client.FastForward(context.Background(), &pb.Empty{})
	if err != nil {
		log.Println(err)
	}
}

func (r *RPC) Rewind() {
	_, err := r.client.Rewind(context.Background(), &pb.Empty{})
	if err != nil {
		log.Println(err)
	}
}

func (r *RPC) Queue() *player.Results {
	out, err := r.client.Queue(context.Background(), &pb.Empty{})
	if err != nil {
		log.Println(err)
		return nil
	}
	return ResultsFromPB(out)
}

func (r RPC) RemoveFromQueue(indices []int) {
	out := make([]int64, len(indices))
	for i, val := range indices {
		out[i] = int64(val)
	}
	_, err := r.client.RemoveFromQueue(context.Background(), &pb.Ints{Value: out})
	if err != nil {
		log.Println(err)
	}
}

func (r RPC) NextSong(id string, f func(player.Result)) {
	go func() {
		stream, err := r.client.NextSong(context.Background(), &pb.String{Value: id})
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
				f(ResultFromPB(ns))
			}
		}
	}()
}

func (r RPC) PlayProgress(id string, f func(player.Progress)) {
	go func() {
		stream, err := r.client.PlayProgress(context.Background(), &pb.String{Value: id})
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
				f(ProgressFromPB(p))
			}
		}
	}()
}

func (r RPC) DownloadProgress(id string, f func(player.Progress)) {
	go func() {
		stream, err := r.client.DownloadProgress(context.Background(), &pb.String{Value: id})
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
				f(ProgressFromPB(p))
			}
		}
	}()
}

func (r RPC) History(page, pageSize int, sort player.Sort) (*player.Results, error) {
	out, err := r.client.History(context.Background(), &pb.Page{Page: int64(page), PageSize: int64(pageSize), Sort: string(sort)})
	return ResultsFromPB(out), err
}
