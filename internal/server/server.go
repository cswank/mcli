package server

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "bitbucket.org/cswank/mcli/internal/server/player"
	"golang.org/x/net/context"
)

const (
	port = ":50051"
)

type server struct {
}

func (s *server) Play(ctx context.Context, r *pb.Result) (*pb.Empty, error) {
	log.Println("play song", r)
	return &pb.Empty{}, nil
}

func (s *server) PlayAlbum(ctx context.Context, r *pb.Result) (*pb.Empty, error) {
	log.Println("play album", r)
	return &pb.Empty{}, nil
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
