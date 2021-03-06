package fetch

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/cswank/mcli/internal/rpc"
	"github.com/cswank/mcli/internal/schema"
	"google.golang.org/grpc"
)

type Remote struct {
	conn   *grpc.ClientConn
	client rpc.FetcherClient
}

func NewRemote(conn *grpc.ClientConn) *Remote {
	return &Remote{
		conn:   conn,
		client: rpc.NewFetcherClient(conn),
	}
}

func (r *Remote) Name() string {
	out, _ := r.client.Name(context.Background(), &rpc.Empty{})
	return out.String()
}

func (r *Remote) Login(u, pw string) error {
	l := &rpc.UsernamePassword{
		Username:  u,
		Passwrord: pw,
	}
	_, err := r.client.Login(context.Background(), l)
	return err
}

func (r *Remote) Ping() bool {
	out, err := r.client.Ping(context.Background(), &rpc.Empty{})
	if err != nil {
		log.Println(err)
		return false
	}
	return out.Value
}

func (r *Remote) AlbumLink() string {
	return ""
}

func (r Remote) FindArtist(term string, p, ps int) (*schema.Results, error) {
	out, err := r.client.FindArtist(context.Background(), &rpc.Request{Term: term, Page: &rpc.Page{Page: int64(p), PageSize: int64(ps)}})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) FindAlbum(term string, p, ps int) (*schema.Results, error) {
	out, err := r.client.FindAlbum(context.Background(), &rpc.Request{Term: term, Page: &rpc.Page{Page: int64(p), PageSize: int64(ps)}})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) FindTrack(term string, p, ps int) (*schema.Results, error) {
	out, err := r.client.FindTrack(context.Background(), &rpc.Request{Term: term, Page: &rpc.Page{Page: int64(p), PageSize: int64(ps)}})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) GetAlbum(id int64) (*schema.Results, error) {
	out, err := r.client.GetAlbum(context.Background(), &rpc.Request{Id: id})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) GetArtistAlbums(id int64, p, ps int) (*schema.Results, error) {
	out, err := r.client.GetArtistAlbums(context.Background(), &rpc.Request{Id: id, Page: &rpc.Page{Page: int64(p), PageSize: int64(ps)}})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) GetArtistTracks(id int64, p, ps int) (*schema.Results, error) {
	out, err := r.client.GetArtistTracks(context.Background(), &rpc.Request{Id: id, Page: &rpc.Page{Page: int64(p), PageSize: int64(ps)}})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) GetPlaylists() (*schema.Results, error) {
	out, err := r.client.GetPlaylists(context.Background(), &rpc.Empty{})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) GetPlaylist(id int64, p, ps int) (*schema.Results, error) {
	out, err := r.client.GetPlaylist(context.Background(), &rpc.Request{Id: id, Page: &rpc.Page{Page: int64(p), PageSize: int64(ps)}})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) Import(fn func(schema.Progress)) error {
	go func() {
		stream, err := r.client.Import(context.Background(), &rpc.Empty{})
		if err != nil {
			log.Fatal("could not get stream for importing new songs", err)
		}
		for {
			p, err := stream.Recv()
			if err == io.EOF {
				time.Sleep(time.Second)
			} else if err != nil {
				log.Println(err)
			} else {
				fn(rpc.ProgressFromPB(p))
			}
		}
	}()

	return nil
}
