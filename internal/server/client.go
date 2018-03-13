package server

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "bitbucket.org/cswank/mcli/internal/server/player"
)

func StartClient(addr string) error {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	// Creates a new CustomerClient
	client := pb.NewPlayerClient(conn)
	log.Println("client", client)

	r := &pb.Result{
		Service: "tidal",
		Path:    "/x",
		Artist: &pb.Result_Artist{
			Id:    "x",
			Title: "fake",
		},
		Album: &pb.Result_Album{
			Id:    "x",
			Title: "fake",
		},
		Track: &pb.Result_Track{
			Id:       "x",
			Title:    "fake",
			Duration: 1,
		},
		Playlist: &pb.Result_Playlist{
			Id:    "x",
			Title: "fake",
		},
	}

	resp, err := client.Play(context.Background(), r)
	log.Println("result", resp, err)
	return err
}
