package download

import (
	"context"
	"io"
	"log"
	"time"

	"bitbucket.org/cswank/mcli/internal/rpc"
	"bitbucket.org/cswank/mcli/internal/schema"
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

func (r Remote) Download(id string, w io.Writer, f func(pg schema.Progress)) {
	go func() {
		stream, err := r.client.Download(context.Background(), &rpc.String{Value: id})
		if err != nil {
			log.Fatal("could not get stream for track", err)
		}
		for {
			p, err := stream.Recv()
			log.Println("recv", err)
			if err == io.EOF {
				time.Sleep(time.Second)
			} else if err != nil {
				log.Println(err)
			} else {
				_, err := w.Write(p.Payload)
				if err != nil {
					log.Println(err)
				}
				f(rpc.ProgressFromPB(p))
			}
		}
	}()
}
