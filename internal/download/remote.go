package download

import (
	"context"
	"io"

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

func (r Remote) Download(id string, w io.Writer, f func(pg schema.Progress)) error {
	stream, err := r.client.Download(context.Background(), &rpc.String{Value: id})
	if err != nil {
		return err
	}

	for {
		p, err := stream.Recv()
		if err != nil {
			return isEOF(err)
		} else {
			_, err := w.Write(p.Payload)
			if err != nil {
				return err
			}
			f(rpc.ProgressFromPB(p))
		}
	}
}
