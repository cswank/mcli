package server

import (
	"net"

	"google.golang.org/grpc"

	"bitbucket.org/cswank/mcli/internal/player"
	pb "bitbucket.org/cswank/mcli/internal/server/player"
	"golang.org/x/net/context"
)

const (
	port = ":50051"
)

type server struct {
	cli player.Player
}

func (s *server) Play(ctx context.Context, in *pb.Result) (*pb.Empty, error) {
	s.cli.Play(ResultFromPB(in))
	return &pb.Empty{}, nil
}

func (s *server) PlayAlbum(ctx context.Context, in *pb.Results) (*pb.Empty, error) {
	s.cli.PlayAlbum(ResultsFromPB(in))
	return &pb.Empty{}, nil
}

func (s *server) Volume(ctx context.Context, r *pb.Float) (*pb.Empty, error) {
	s.cli.Volume(float64(r.Value))
	return &pb.Empty{}, nil
}

func (s *server) Pause(ctx context.Context, r *pb.Empty) (*pb.Empty, error) {
	s.cli.Pause()
	return &pb.Empty{}, nil
}

func (s *server) FastForward(ctx context.Context, r *pb.Empty) (*pb.Empty, error) {
	s.cli.FastForward()
	return &pb.Empty{}, nil
}

func (s *server) Queue(ctx context.Context, r *pb.Empty) (*pb.Results, error) {
	res := s.cli.Queue()
	return &pb.Results{}, nil
}

func Start() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterPlayerServer(s, &server{})
	s.Serve(lis)
	return nil
}

func ResultFromPB(r *pb.Result) player.Result {
	return player.Result{}
}

func ResultsFromPB(r *pb.Results) *player.Results {
	return &player.Results{}
}
