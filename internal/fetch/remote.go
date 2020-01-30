package fetch

import (
	"context"

	"bitbucket.org/cswank/mcli/internal/rpc"
	"bitbucket.org/cswank/mcli/internal/schema"
	"google.golang.org/grpc"
)

type Remote struct {
	conn   *grpc.ClientConn
	client rpc.FetcherClient
}

func NewRemote() (*Remote, error) {
	return &Remote{}, nil
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
	out, _ := r.client.Ping(context.Background(), &rpc.Empty{})
	return out.Value
}

func (r *Remote) AlbumLink() string {
	return ""
}

func (r Remote) FindArtist(term string, n int) (*schema.Results, error) {
	out, err := r.client.FindArtist(context.Background(), &rpc.Request{Term: term, N: int64(n)})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) FindAlbum(term string, n int) (*schema.Results, error) {
	out, err := r.client.FindAlbum(context.Background(), &rpc.Request{Term: term, N: int64(n)})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) FindTrack(term string, n int) (*schema.Results, error) {
	out, err := r.client.FindTrack(context.Background(), &rpc.Request{Term: term, N: int64(n)})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) GetAlbum(id string) (*schema.Results, error) {
	out, err := r.client.GetAlbum(context.Background(), &rpc.String{Value: id})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) GetTrack(id string) (string, error) {
	out, err := r.client.GetTrack(context.Background(), &rpc.String{Value: id})
	return out.Value, err
}

func (r Remote) GetArtistAlbums(id string, n int) (*schema.Results, error) {
	out, err := r.client.GetArtistAlbums(context.Background(), &rpc.Request{Term: id, N: int64(n)})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) GetArtistTracks(id string, n int) (*schema.Results, error) {
	out, err := r.client.GetArtistTracks(context.Background(), &rpc.Request{Term: id, N: int64(n)})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) GetPlaylists() (*schema.Results, error) {
	out, err := r.client.GetPlaylists(context.Background(), &rpc.Empty{})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) GetPlaylist(id string, n int) (*schema.Results, error) {
	out, err := r.client.GetPlaylist(context.Background(), &rpc.Request{Term: id, N: int64(n)})
	return rpc.ResultsFromPB(out), err
}
