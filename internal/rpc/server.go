package rpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"bitbucket.org/cswank/mcli/internal/player"
	pb "bitbucket.org/cswank/mcli/internal/rpc/player"
	"golang.org/x/net/context"
)

const (
	port = ":50051"
)

type server struct {
	cli                    player.Player
	nextSongStream         pb.Player_NextSongServer
	playProgressStream     pb.Player_PlayProgressServer
	downloadProgressStream pb.Player_DownloadProgressServer
	done                   chan bool
}

func (s *server) Done(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	s.cli.Done()
	close(s.done)
	s.done = make(chan bool)
	return &pb.Empty{}, nil
}

func (s *server) Close(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *server) Play(ctx context.Context, in *pb.Result) (*pb.Empty, error) {
	s.cli.Play(ResultFromPB(in))
	return &pb.Empty{}, nil
}

func (s *server) PlayAlbum(ctx context.Context, in *pb.Results) (*pb.Empty, error) {
	s.cli.PlayAlbum(ResultsFromPB(in))
	return &pb.Empty{}, nil
}

func (s *server) Volume(ctx context.Context, r *pb.Float) (*pb.Float, error) {
	v := s.cli.Volume(float64(r.Value))
	return &pb.Float{Value: float32(v)}, nil
}

func (s *server) Pause(ctx context.Context, r *pb.Empty) (*pb.Empty, error) {
	s.cli.Pause()
	return &pb.Empty{}, nil
}

func (s *server) FastForward(ctx context.Context, r *pb.Empty) (*pb.Empty, error) {
	s.cli.FastForward()
	return &pb.Empty{}, nil
}

func (s *server) Rewind(ctx context.Context, r *pb.Empty) (*pb.Empty, error) {
	s.cli.Rewind()
	return &pb.Empty{}, nil
}

func (s *server) Queue(ctx context.Context, r *pb.Empty) (*pb.Results, error) {
	return PBFromResults(s.cli.Queue()), nil
}

func (s *server) RemoveFromQueue(ctx context.Context, r *pb.Int) (*pb.Results, error) {
	s.cli.RemoveFromQueue(int(r.Value))
	return PBFromResults(s.cli.Queue()), nil
}

func (s *server) NextSong(_ *pb.Empty, stream pb.Player_NextSongServer) error {
	s.nextSongStream = stream
	s.cli.NextSong(s.nextSong)
	<-s.done
	s.cli.NextSong(nil)
	return nil
}

func (s *server) PlayProgress(_ *pb.Empty, stream pb.Player_PlayProgressServer) error {
	s.playProgressStream = stream
	s.cli.PlayProgress(s.playProgress)
	<-s.done
	s.cli.PlayProgress(nil)
	return nil
}

func (s *server) DownloadProgress(_ *pb.Empty, stream pb.Player_DownloadProgressServer) error {
	s.downloadProgressStream = stream
	s.cli.DownloadProgress(s.downloadProgress)
	<-s.done
	s.cli.DownloadProgress(nil)
	return nil
}

func (s *server) History(ctx context.Context, p *pb.Page) (*pb.Results, error) {
	r, err := s.cli.History(int(p.Page), int(p.PageSize), player.Sort(p.Sort))
	return PBFromResults(r), err
}

func (s *server) nextSong(r player.Result) {
	if err := s.nextSongStream.Send(PBFromResult(r)); err != nil {
		log.Printf("could not stream result %v, err: %s", r, err)
	}
}

func (s *server) playProgress(p player.Progress) {
	if err := s.playProgressStream.Send(PBFromProgress(p)); err != nil {
		log.Printf("could not stream result %v, err: %s", p, err)
	}
}

func (s *server) downloadProgress(p player.Progress) {
	if err := s.downloadProgressStream.Send(PBFromProgress(p)); err != nil {
		log.Printf("could not stream result %v, err: %s", p, err)
	}
}

func Start() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	// Creates a new gRPC server
	s := grpc.NewServer()
	cli, err := player.NewTidal(nil)
	if err != nil {
		return err
	}

	pb.RegisterPlayerServer(s, &server{
		cli:  cli,
		done: make(chan bool),
	})

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)
	go s.Serve(lis)

	<-stop
	fmt.Println("graceful stop")
	s.GracefulStop()
	cli.Close()
	return nil
}

func ResultFromPB(r *pb.Result) player.Result {
	return player.Result{
		Service:   r.GetService(),
		Path:      r.GetPath(),
		PlayCount: int(r.GetPlaycount()),
		Track:     trackFromPB(r.GetTrack()),
		Album:     albumFromPB(r.GetAlbum()),
		Artist:    artistFromPB(r.GetArtist()),
		Playlist:  playlistFromPB(r.GetPlaylist()),
	}
}

func trackFromPB(t *pb.Result_Track) player.Track {
	if t == nil {
		return player.Track{}
	}
	return player.Track{
		ID:       t.GetId(),
		Title:    t.GetTitle(),
		Duration: int(t.GetDuration()),
	}
}

func playlistFromPB(p *pb.Result_Playlist) player.Album {
	if p == nil {
		return player.Album{}
	}
	return player.Album{
		ID:    p.GetId(),
		Title: p.GetTitle(),
	}
}

func albumFromPB(a *pb.Result_Album) player.Album {
	if a == nil {
		return player.Album{}
	}
	return player.Album{
		ID:    a.GetId(),
		Title: a.GetTitle(),
	}
}

func artistFromPB(a *pb.Result_Artist) player.Artist {
	if a == nil {
		return player.Artist{}
	}
	return player.Artist{
		ID:   a.GetId(),
		Name: a.GetName(),
	}
}

func pbFromTrack(t player.Track) *pb.Result_Track {
	return &pb.Result_Track{
		Id:       t.ID,
		Title:    t.Title,
		Duration: int64(t.Duration),
	}
}

func pbFromPlaylist(a player.Album) *pb.Result_Playlist {
	return &pb.Result_Playlist{
		Id:    a.ID,
		Title: a.Title,
	}
}

func pbFromAlbum(a player.Album) *pb.Result_Album {
	return &pb.Result_Album{
		Id:    a.ID,
		Title: a.Title,
	}
}

func pbFromArtist(a player.Artist) *pb.Result_Artist {
	return &pb.Result_Artist{
		Id:   a.ID,
		Name: a.Name,
	}
}

func ResultsFromPB(r *pb.Results) *player.Results {
	pbRes := r.GetResults()
	res := make([]player.Result, len(pbRes))
	for i := range pbRes {
		res[i] = ResultFromPB(pbRes[i])
	}
	return &player.Results{
		Type:    r.GetType(),
		Header:  r.GetHeader(),
		Fmt:     r.GetFmt(),
		Results: res,
	}
}

func PBFromResults(r *player.Results) *pb.Results {
	out := make([]*pb.Result, len(r.Results))
	for i := range r.Results {
		out[i] = PBFromResult(r.Results[i])
	}
	return &pb.Results{
		Header:  r.Header,
		Type:    r.Type,
		Fmt:     r.Fmt,
		Results: out,
	}
}

func PBFromResult(r player.Result) *pb.Result {
	return &pb.Result{
		Service:   r.Service,
		Path:      r.Path,
		Playcount: int64(r.PlayCount),
		Track:     pbFromTrack(r.Track),
		Album:     pbFromAlbum(r.Album),
		Artist:    pbFromArtist(r.Artist),
		Playlist:  pbFromPlaylist(r.Playlist),
	}
}

func PBFromProgress(p player.Progress) *pb.Progress {
	return &pb.Progress{
		N:     int64(p.N),
		Total: int64(p.Total),
	}
}

func ProgressFromPB(p *pb.Progress) player.Progress {
	return player.Progress{
		N:     int(p.GetN()),
		Total: int(p.GetTotal()),
	}
}
