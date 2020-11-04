package download

import (
	"context"
	"io"
	"log"

	"github.com/cswank/mcli/internal/rpc"
	"github.com/cswank/mcli/internal/schema"
	"google.golang.org/grpc"
)

type Remote struct {
	conn   *grpc.ClientConn
	client rpc.DownloaderClient
}

func NewRemote(conn *grpc.ClientConn) *Remote {
	return &Remote{
		conn:   conn,
		client: rpc.NewDownloaderClient(conn),
	}
}

func (r Remote) Download(id int64, w io.Writer, f func(pg schema.Progress)) {
	stream, err := r.client.Download(context.Background(), &rpc.Request{Id: id})
	log.Println("stream", err)
	if err != nil {
		log.Fatal("could not get stream for track", err)
	}
	for {
		p, err := stream.Recv()
		if err == io.EOF {
			return
		} else if err != nil {
			log.Println(err)
			return
		} else {
			_, err := w.Write(p.Payload)
			if err != nil {
				log.Println(err)
			}
			f(rpc.ProgressFromPB(p))
		}
	}
}
