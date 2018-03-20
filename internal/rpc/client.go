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

func (c *Client) Done() {
	c.client.Done(context.Background(), &pb.Empty{})
	c.conn.Close()
}

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

func (c *Client) Volume(v float64) {
	_, err := c.client.Volume(context.Background(), &pb.Float{Value: float32(v)})
	if err != nil {
		log.Println(err)
	}
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

func (c *Client) RemoveFromQueue(i int) {
	_, err := c.client.RemoveFromQueue(context.Background(), &pb.Int{Value: int64(i)})
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) NextSong(f func(player.Result)) {
	go func() {
		stream, err := c.client.NextSong(context.Background(), &pb.Empty{})
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

func (c *Client) PlayProgress(f func(player.Progress)) {
	go func() {
		stream, err := c.client.PlayProgress(context.Background(), &pb.Empty{})
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

func (c *Client) DownloadProgress(f func(player.Progress)) {
	go func() {
		stream, err := c.client.DownloadProgress(context.Background(), &pb.Empty{})
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

func (c *Client) History(page, pageSize int) (*player.Results, error) {
	out, err := c.client.History(context.Background(), &pb.Page{Page: int64(page), PageSize: int64(pageSize)})
	return ResultsFromPB(out), err
}
