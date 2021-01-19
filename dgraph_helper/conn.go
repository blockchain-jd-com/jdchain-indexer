package dgraph_helper

import (
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
)

func CreateDgClient(host string) (*dgo.Dgraph, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		logger.Warn("While trying to dial gRPC failed: ", err)
		return nil, err
	}

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)
	return dg, nil
}
