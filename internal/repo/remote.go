package repo

import (
	"context"

	"github.com/cswank/mcli/internal/rpc"
	"github.com/cswank/mcli/internal/schema"
	"google.golang.org/grpc"
)

type Remote struct {
	client rpc.HistoryClient
}

func NewRemote(conn *grpc.ClientConn) *Remote {
	return &Remote{
		client: rpc.NewHistoryClient(conn),
	}
}

func (r Remote) Fetch(page, pageSize int, sort Sort) (*schema.Results, error) {
	out, err := r.client.Fetch(context.Background(), &rpc.Page{Page: int64(page), PageSize: int64(pageSize), Sort: string(sort)})
	return rpc.ResultsFromPB(out), err
}

func (r Remote) Save(rs schema.Result) error {
	_, err := r.client.Save(context.Background(), rpc.PBFromResult(rs))
	return err
}
