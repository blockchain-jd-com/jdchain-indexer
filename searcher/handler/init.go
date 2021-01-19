package handler

import (
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"github.com/dgraph-io/dgo/v200"
	"github.com/ssor/zlog"
)

var (
	logger = zlog.New("searcher", "handler")
	//searcherEngine = &riot.Engine{}
	dgClient *dgo.Dgraph
)

func Init(host string) {
	initDgraphClient(host)
	//initDataCacher(cacheSize)
	//searcherEngine.Init(types.EngineOpts{
	//	Using:             3,
	//	NotUseGse: true,
	//})
}

//func initDataCacher(cacheSize int) {
//	query.NewAutoDocCacher(searcherEngine, cacheSize, dgClient)
//}

func initDgraphClient(host string) {
	client, err := dgraph_helper.CreateDgClient(host)
	if err != nil {
		panic(err)
	}
	dgClient = client
}
