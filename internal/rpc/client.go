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
	flac   player.Client
}

func NewClient(addr string, opts ...func(*Client) error) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn:   conn,
		client: pb.NewPlayerClient(conn),
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func LocalPlay(local bool, addr string) func(c *Client) error {
	return func(c *Client) error {
		if !local {
			return nil
		}

		f, err := player.NewDisk(nil, addr)
		c.flac = f
		return err
	}
}

func (c *Client) Done(id string) {
	c.client.Done(context.Background(), &pb.String{Value: id})
	c.conn.Close()
}

func (c *Client) Close() {}

func (c *Client) Play(r player.Result) {
	c.flac.Play(r)
}

func (c *Client) PlayAlbum(r *player.Results) {
	if c.flac != nil {
		c.flac.PlayAlbum(r)
	} else {
		_, err := c.client.PlayAlbum(context.Background(), PBFromResults(r))
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *Client) Volume(v float64) float64 {
	if c.flac != nil {
		return c.flac.Volume(v)
	}

	f, err := c.client.Volume(context.Background(), &pb.Float{Value: float32(v)})
	if err != nil {
		log.Println(err)
	}
	return float64(f.Value)
}

func (c *Client) Pause() {
	if c.flac != nil {
		c.flac.Pause()
	} else {
		_, err := c.client.Pause(context.Background(), &pb.Empty{})
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *Client) FastForward() {
	if c.flac != nil {
		c.flac.FastForward()
	} else {
		_, err := c.client.FastForward(context.Background(), &pb.Empty{})
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *Client) Rewind() {
	if c.flac != nil {
		c.flac.FastForward()
	} else {
		_, err := c.client.Rewind(context.Background(), &pb.Empty{})
		if err != nil {
			log.Println(err)
		}
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
	if c.flac != nil {
		c.flac.NextSong(id, f)
	} else {
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
}

func (c *Client) PlayProgress(id string, f func(player.Progress)) {
	if c.flac != nil {
		c.flac.PlayProgress(id, f)
	} else {
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
}

func (c *Client) DownloadProgress(id string, f func(player.Progress)) {
	if c.flac != nil {
		c.flac.DownloadProgress(id, f)
	} else {
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
}

func (c *Client) History(page, pageSize int, sort player.Sort) (*player.Results, error) {
	if c.flac != nil {
		return c.flac.History(page, pageSize, sort)
	}
	out, err := c.client.History(context.Background(), &pb.Page{Page: int64(page), PageSize: int64(pageSize), Sort: string(sort)})
	return ResultsFromPB(out), err
}

func (c *Client) Name() string {
	out, _ := c.client.Name(context.Background(), &pb.Empty{})
	return out.String()
}

func (c *Client) Login(u, pw string) error {
	l := &pb.UsernamePassword{
		Username:  u,
		Passwrord: pw,
	}
	_, err := c.client.Login(context.Background(), l)
	return err
}

func (c *Client) Ping() bool {
	out, _ := c.client.Ping(context.Background(), &pb.Empty{})
	return out.Value
}

func (c *Client) AlbumLink() string {
	return ""
}

func (c *Client) FindArtist(term string, n int) (*player.Results, error) {
	out, err := c.client.FindArtist(context.Background(), &pb.Request{Term: term, N: int64(n)})
	return ResultsFromPB(out), err
}

func (c *Client) FindAlbum(term string, n int) (*player.Results, error) {
	out, err := c.client.FindAlbum(context.Background(), &pb.Request{Term: term, N: int64(n)})
	return ResultsFromPB(out), err
}

func (c *Client) FindTrack(term string, n int) (*player.Results, error) {
	out, err := c.client.FindTrack(context.Background(), &pb.Request{Term: term, N: int64(n)})
	return ResultsFromPB(out), err
}

func (c *Client) GetAlbum(id string) (*player.Results, error) {
	out, err := c.client.GetAlbum(context.Background(), &pb.String{Value: id})
	return ResultsFromPB(out), err
}

func (c *Client) GetTrack(id string) (string, error) {
	out, err := c.client.GetTrack(context.Background(), &pb.String{Value: id})
	return out.Value, err
}

func (c *Client) GetArtistAlbums(id string, n int) (*player.Results, error) {
	out, err := c.client.GetArtistAlbums(context.Background(), &pb.Request{Term: id, N: int64(n)})
	return ResultsFromPB(out), err
}

func (c *Client) GetArtistTracks(id string, n int) (*player.Results, error) {
	out, err := c.client.GetArtistTracks(context.Background(), &pb.Request{Term: id, N: int64(n)})
	return ResultsFromPB(out), err
}

func (c *Client) GetPlaylists() (*player.Results, error) {
	out, err := c.client.GetPlaylists(context.Background(), &pb.Empty{})
	return ResultsFromPB(out), err
}

func (c *Client) GetPlaylist(id string, n int) (*player.Results, error) {
	out, err := c.client.GetPlaylist(context.Background(), &pb.Request{Term: id, N: int64(n)})
	return ResultsFromPB(out), err
}
