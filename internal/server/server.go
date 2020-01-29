package server

import (
	"log"
	"net"

	"bitbucket.org/cswank/mcli/internal/fetch"
	"bitbucket.org/cswank/mcli/internal/play"
	"bitbucket.org/cswank/mcli/internal/repo"
	"bitbucket.org/cswank/mcli/internal/schema"
	"bitbucket.org/cswank/mcli/internal/rpc"
	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

const (
	port = ":50051"
)

type client struct {
	play.Player
	fetch.Fetcher
}

type server struct {
	cli                    *client
	nextSongStream         rpc.Player_NextSongServer
	playProgressStream     rpc.Player_PlayProgressServer
	downloadProgressStream rpc.Player_DownloadProgressServer
	done                   chan bool
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

func (s *server) History(ctx context.Context, p *rpc.Page) (*rpc.Results, error) {
	r, err := s.cli.History(int(p.Page), int(p.PageSize), repo.Sort(p.Sort))
	return rpc.PBFromResults(r), err
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

func Start(p play.Player, f fetch.Fetcher) error {
	log.Println("rpc listening on ", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	// Creates a new gRPC server
	s := grpc.NewServer()
	RegisterPlayerServer(s, &server{
		cli:  &client{Player: p, Fetcher: f},
		done: make(chan bool),
	})

	// stop := make(chan os.Signal)
	// signal.Notify(stop, syscall.SIGTERM)
	// signal.Notify(stop, syscall.SIGINT)
	s.Serve(lis)

	// <-stop
	// fmt.Println("graceful stop")
	// s.GracefulStop()
	//cli.Close()
	return nil
}

func (s *server) Name(ctx context.Context, _ *rpc.Empty) (*rpc.String, error) {
	n := s.cli.Name()
	return &String{Value: n}, nil
}

func (s *server) Login(ctx context.Context, up *UsernamePassword) (*rpc.Empty, error) {
	err := s.cli.Login(up.Username, up.Passwrord)
	return &Empty{}, err
}

func (s *server) Ping(ctx context.Context, _ *Empty) (*rpc.Bool, error) {
	out := s.cli.Ping()
	return &Bool{Value: out}, nil
}

func (s *server) AlbumLink(ctx context.Context, _ *Empty) (*rpc.String, error) {
	s.cli.AlbumLink()
	return &String{}, nil
}

func (s *server) FindArtist(ctx context.Context, r *Request) (*rpc.Results, error) {
	out, err := s.cli.FindArtist(r.Term, int(r.N))
	return rpc.PBFromResults(out), err
}

func (s *server) FindAlbum(ctx context.Context, r *Request) (*rpc.Results, error) {
	out, err := s.cli.FindAlbum(r.Term, int(r.N))
	return PBFromResults(out), err
}

func (s *server) FindTrack(ctx context.Context, r *Request) (*rpc.Results, error) {
	out, err := s.cli.FindTrack(r.Term, int(r.N))
	return PBFromResults(out), err
}

func (s *server) GetAlbum(ctx context.Context, st *rpc.String) (*rpc.Results, error) {
	out, err := s.cli.GetAlbum(st.Value)
	return PBFromResults(out), err
}

func (s *server) GetTrack(ctx context.Context, st *rpc.String) (*rpc.String, error) {
	out, err := s.cli.GetTrack(st.Value)
	return &String{Value: out}, err
}

func (s *server) GetArtistAlbums(ctx context.Context, r *schema.Request) (*rpc.Results, error) {
	out, err := s.cli.GetArtistAlbums(r.Term, int(r.N))
	return PBFromResults(out), err
}

func (s *server) GetArtistTracks(ctx context.Context, r *schema.Request) (*rpc.Results, error) {
	out, err := s.cli.GetArtistTracks(r.Term, int(r.N))
	return PBFromResults(out), err
}

func (s *server) GetPlaylists(ctx context.Context, e *rpc.Empty) (*rpc.Results, error) {
	out, err := s.cli.GetPlaylists()
	return PBFromResults(out), err
}

func (s *server) GetPlaylist(ctx context.Context, r *rpc.Request) (*rpc.Results, error) {
	out, err := s.cli.GetPlaylist(r.Term, int(r.N))
	return PBFromResults(out), err
}
