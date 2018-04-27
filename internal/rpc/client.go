package rpc

import (
	"io"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"bitbucket.org/cswank/mcli/internal/player"
	pb "bitbucket.org/cswank/mcli/internal/rpc/player"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.PlayerClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: pb.NewPlayerClient(conn),
	}, nil
}

func (c *Client) Done(id string) {
	c.client.Done(context.Background(), &pb.String{Value: id})
	c.conn.Close()
}

func (c *Client) Close() {}

func (c *Client) Play(r player.Result) {
	_, err := c.client.Play(context.Background(), PBFromResult(r))
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) PlayAlbum(r *player.Results) {
	_, err := c.client.PlayAlbum(context.Background(), PBFromResults(r))
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) Volume(v float64) float64 {
	f, err := c.client.Volume(context.Background(), &pb.Float{Value: float32(v)})
	if err != nil {
		log.Println(err)
	}
	return float64(f.Value)
}

func (c *Client) Pause() {
	_, err := c.client.Pause(context.Background(), &pb.Empty{})
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) FastForward() {
	_, err := c.client.FastForward(context.Background(), &pb.Empty{})
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) Rewind() {
	_, err := c.client.Rewind(context.Background(), &pb.Empty{})
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) Queue() *player.Results {
	out, err := c.client.Queue(context.Background(), &pb.Empty{})
	if err != nil {
		log.Println(err)
		return nil
	}
	return ResultsFromPB(out)
}

func (c *Client) RemoveFromQueue(indices []int) {
	out := make([]int64, len(indices))
	for i, val := range indices {
		out[i] = int64(val)
	}
	_, err := c.client.RemoveFromQueue(context.Background(), &pb.Ints{Value: out})
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) NextSong(id string, f func(player.Result)) {
	go func() {
		stream, err := c.client.NextSong(context.Background(), &pb.String{Value: id})
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

func (c *Client) PlayProgress(id string, f func(player.Progress)) {
	go func() {
		stream, err := c.client.PlayProgress(context.Background(), &pb.String{Value: id})
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

func (c *Client) DownloadProgress(id string, f func(player.Progress)) {
	go func() {
		stream, err := c.client.DownloadProgress(context.Background(), &pb.String{Value: id})
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

func (c *Client) History(page, pageSize int, sort player.Sort) (*player.Results, error) {
	out, err := c.client.History(context.Background(), &pb.Page{Page: int64(page), PageSize: int64(pageSize), Sort: string(sort)})
	return ResultsFromPB(out), err
}
