package server

import (
	"log"
	"net"
	"os"

	"github.com/cswank/mcli/internal/fetch"
	"github.com/cswank/mcli/internal/history"
	"github.com/cswank/mcli/internal/play"
	"github.com/cswank/mcli/internal/repo"
	"github.com/cswank/mcli/internal/rpc"
	"github.com/cswank/mcli/internal/schema"
	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

const (
	port = ":50051"
)

type client struct {
	play.Player
	fetch.Fetcher
	history.History
}

type server struct {
	rpc.UnsafePlayerServer
	rpc.UnsafeFetcherServer
	rpc.UnsafeDownloaderServer
	rpc.UnsafeHistoryServer
	pth                    string
	db                     *repo.SQLLite
	cli                    *client
	nextSongStream         rpc.Player_NextSongServer
	playProgressStream     rpc.Player_PlayProgressServer
	downloadProgressStream rpc.Player_DownloadProgressServer
	getTrackProgressStream rpc.Downloader_DownloadServer
	done                   chan bool
}

func Start(p play.Player, f fetch.Fetcher, h history.History, db *repo.SQLLite, pth string) error {
	log.Println("rpc listening on ", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	// Creates a new gRPC server
	srv := grpc.NewServer()

	s := &server{
		cli:  &client{Player: p, Fetcher: f, History: h},
		db:   db,
		pth:  pth,
		done: make(chan bool),
	}

	rpc.RegisterPlayerServer(srv, s)
	rpc.RegisterFetcherServer(srv, s)
	rpc.RegisterDownloaderServer(srv, s)
	rpc.RegisterHistoryServer(srv, s)

	srv.Serve(lis)
	return nil
}

func (s *server) Done(ctx context.Context, id *rpc.String) (*rpc.Empty, error) {
	s.cli.Done(id.Value)
	close(s.done)
	s.done = make(chan bool)
	return &rpc.Empty{}, nil
}

func (s *server) Close(ctx context.Context, _ *rpc.Empty) (*rpc.Empty, error) {
	return &rpc.Empty{}, nil
}

func (s *server) Play(ctx context.Context, in *rpc.Result) (*rpc.Empty, error) {
	s.cli.Play(rpc.ResultFromPB(in))
	return &rpc.Empty{}, nil
}

func (s *server) PlayAlbum(ctx context.Context, in *rpc.Results) (*rpc.Empty, error) {
	s.cli.PlayAlbum(rpc.ResultsFromPB(in))
	return &rpc.Empty{}, nil
}

func (s *server) Volume(ctx context.Context, r *rpc.Float) (*rpc.Float, error) {
	v := s.cli.Volume(float64(r.Value))
	return &rpc.Float{Value: float32(v)}, nil
}

func (s *server) Pause(ctx context.Context, r *rpc.Empty) (*rpc.Empty, error) {
	s.cli.Pause()
	return &rpc.Empty{}, nil
}

func (s *server) FastForward(ctx context.Context, r *rpc.Empty) (*rpc.Empty, error) {
	s.cli.FastForward()
	return &rpc.Empty{}, nil
}

func (s *server) Rewind(ctx context.Context, r *rpc.Empty) (*rpc.Empty, error) {
	s.cli.Rewind()
	return &rpc.Empty{}, nil
}

func (s *server) Queue(ctx context.Context, r *rpc.Empty) (*rpc.Results, error) {
	return rpc.PBFromResults(s.cli.Queue()), nil
}

func (s *server) RemoveFromQueue(ctx context.Context, r *rpc.Ints) (*rpc.Results, error) {
	out := make([]int, len(r.Value))
	for i, val := range r.Value {
		out[i] = int(val)
	}

	s.cli.RemoveFromQueue(out)
	return rpc.PBFromResults(s.cli.Queue()), nil
}

func (s *server) NextSong(id *rpc.String, stream rpc.Player_NextSongServer) error {
	s.nextSongStream = stream
	s.cli.NextSong(id.Value, s.nextSong)
	<-s.done
	s.cli.NextSong(id.Value, nil)
	return nil
}

func (s *server) PlayProgress(id *rpc.String, stream rpc.Player_PlayProgressServer) error {
	s.playProgressStream = stream
	s.cli.PlayProgress(id.Value, s.playProgress)
	<-s.done
	s.cli.PlayProgress(id.Value, nil)
	return nil
}

func (s *server) DownloadProgress(id *rpc.String, stream rpc.Player_DownloadProgressServer) error {
	s.downloadProgressStream = stream
	s.cli.DownloadProgress(id.Value, s.downloadProgress)
	<-s.done
	s.cli.DownloadProgress(id.Value, nil)
	return nil
}

func (s *server) Import(_ *rpc.Empty, stream rpc.Fetcher_ImportServer) error {
	s.playProgressStream = stream
	s.cli.Import(s.playProgress)
	<-s.done
	return nil
}

func (s *server) Fetch(ctx context.Context, p *rpc.Page) (*rpc.Results, error) {
	r, err := s.cli.Fetch(int(p.Page), int(p.PageSize), repo.Sort(p.Sort))
	return rpc.PBFromResults(r), err
}

func (s *server) Save(ctx context.Context, in *rpc.Result) (*rpc.Empty, error) {
	err := s.cli.Save(rpc.ResultFromPB(in))
	return &rpc.Empty{}, err
}

func (s *server) nextSong(r schema.Result) {
	if err := s.nextSongStream.Send(rpc.PBFromResult(r)); err != nil {
		log.Printf("could not stream result %v, err: %s", r, err)
	}
}

func (s *server) playProgress(p schema.Progress) {
	if err := s.playProgressStream.Send(rpc.PBFromProgress(p)); err != nil {
		log.Printf("could not stream result %v, err: %s", p, err)
	}
}

func (s *server) downloadProgress(p schema.Progress) {
	if err := s.downloadProgressStream.Send(rpc.PBFromProgress(p)); err != nil {
		log.Printf("could not stream result %v, err: %s", p, err)
	}
}

func (s *server) Name(ctx context.Context, _ *rpc.Empty) (*rpc.String, error) {
	n := s.cli.Name()
	return &rpc.String{Value: n}, nil
}

func (s *server) Login(ctx context.Context, up *rpc.UsernamePassword) (*rpc.Empty, error) {
	err := s.cli.Login(up.Username, up.Passwrord)
	return &rpc.Empty{}, err
}

func (s *server) Ping(ctx context.Context, _ *rpc.Empty) (*rpc.Bool, error) {
	out := s.cli.Ping()
	return &rpc.Bool{Value: out}, nil
}

func (s *server) AlbumLink(ctx context.Context, _ *rpc.Empty) (*rpc.String, error) {
	s.cli.AlbumLink()
	return &rpc.String{}, nil
}

func (s *server) FindArtist(ctx context.Context, r *rpc.Request) (*rpc.Results, error) {
	out, err := s.cli.FindArtist(r.Term, int(r.N))
	return rpc.PBFromResults(out), err
}

func (s *server) FindAlbum(ctx context.Context, r *rpc.Request) (*rpc.Results, error) {
	out, err := s.cli.FindAlbum(r.Term, int(r.N))
	return rpc.PBFromResults(out), err
}

func (s *server) FindTrack(ctx context.Context, r *rpc.Request) (*rpc.Results, error) {
	out, err := s.cli.FindTrack(r.Term, int(r.N))
	return rpc.PBFromResults(out), err
}

func (s *server) GetAlbum(ctx context.Context, r *rpc.Request) (*rpc.Results, error) {
	out, err := s.cli.GetAlbum(r.Id)
	return rpc.PBFromResults(out), err
}

// func (s *server) DownloadProgress(id *rpc.String, stream rpc.Player_DownloadProgressServer) error {
// 	s.downloadProgressStream = stream
// 	s.cli.DownloadProgress(id.Value, s.downloadProgress)
// 	<-s.done
// 	s.cli.DownloadProgress(id.Value, nil)
// 	return nil
// }

func (s *server) Download(req *rpc.Request, stream rpc.Downloader_DownloadServer) error {
	s.getTrackProgressStream = stream
	buf := make([]byte, 100000)

	pth, err := s.track(req.Id)
	if err != nil {
		return err
	}

	f, err := os.Open(pth)
	if err != nil {
		return err
	}

	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	l := int(fi.Size())
	var tot int
	for {
		n, err := f.Read(buf)
		tot += n
		s.getTrackProgressStream.Send(rpc.PBFromProgress(schema.Progress{N: tot, Total: l, Payload: buf[:n]}))
		if err != nil || n == 0 {
			break
		}
	}

	return nil
}

func (s *server) track(id int64) (string, error) {
	return s.db.Track(id)
}

func (s *server) GetArtistAlbums(ctx context.Context, r *rpc.Request) (*rpc.Results, error) {
	out, err := s.cli.GetArtistAlbums(r.Id, int(r.N))
	return rpc.PBFromResults(out), err
}

func (s *server) GetArtistTracks(ctx context.Context, r *rpc.Request) (*rpc.Results, error) {
	out, err := s.cli.GetArtistTracks(r.Id, int(r.N))
	return rpc.PBFromResults(out), err
}

func (s *server) GetPlaylists(ctx context.Context, e *rpc.Empty) (*rpc.Results, error) {
	out, err := s.cli.GetPlaylists()
	return rpc.PBFromResults(out), err
}

func (s *server) GetPlaylist(ctx context.Context, r *rpc.Request) (*rpc.Results, error) {
	out, err := s.cli.GetPlaylist(r.Id, int(r.N))
	return rpc.PBFromResults(out), err
}
