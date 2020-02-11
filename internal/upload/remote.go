package upload

import (
	"context"
	"io"

	"bitbucket.org/cswank/mcli/internal/rpc"
	"bitbucket.org/cswank/mcli/internal/schema"
	"google.golang.org/grpc"
)

type Remote struct {
	client rpc.UploaderClient
}

func NewRemote(conn *grpc.ClientConn) *Remote {
	return &Remote{
		client: rpc.NewUploaderClient(conn),
	}
}

func (r Remote) Upload(u schema.Upload, rd io.Reader, f func(pg schema.Progress)) error {
	stream, err := r.client.Send(context.Background())
	if err != nil {
		return err
	}

	buf := make([]byte, 10000)
	for {
		n, err := rd.Read(buf)
		if err != nil {
			return isEOF(err)
		} else {
			u.Payload = buf[:n]
			stream.Send(rpc.PBFromUpload(u))
			if err != nil {
				return err
			}
			f(u.Progress)
		}
	}
}

func isEOF(err error) error {
	if err == io.EOF {
		return nil
	}
	return err
}
