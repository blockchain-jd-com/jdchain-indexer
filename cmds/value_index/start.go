package value_index

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/value_indexer/worker"
	"github.com/mkideal/cli"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"

	//"github.com/ssor/dgraph_memory_loader"
	"github.com/ssor/zlog"
)

var Root = &cli.Command{
	Name: "data",
	Desc: "index data in txs of ledger",
	Argv: func() interface{} { return new(ServerArgs) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*ServerArgs)
		return StartLedgerServer(argv)
	},
}

type ServerArgs struct {
	cli.Helper
	ApiHost string `cli:"*ledger-host" usage:"api server host, like http://127.0.0.1:8080" dft:""`
	// TODO from to
	//From       int64  `cli:"from" usage:"from block height" dft:"0"`
	//To         int64  `cli:"to" usage:"to block height" dft:"0"`
	Port       int    `cli:"p,port" usage:"server listening port" dft:"8082"`
	DgraphHost string `cli:"dgraph" usage:"dgraph server host" dft:"127.0.0.1:9080"`
	Production bool   `cli:"production" usage:"if use production mode" dft:"false"`
}

var (
	dgraphHelper  *dgraph_helper.Helper
	schemaCenter  *worker.SchemaStatusCenter
	workerManager *worker.ValueIndexWorkerManager
	monitor       *worker.SchemaMonitor
)

func StartLedgerServer(args *ServerArgs) error {
	production := args.Production

	if production {
		fmt.Println("** Production ** mode")

		logFile := &lumberjack.Logger{
			Filename:   "value_index.log",
			MaxSize:    30, // megabytes
			MaxBackups: 10,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		}
		zlog.SetOutput(logFile)
	}

	dgraphHelper = dgraph_helper.NewHelper(args.DgraphHost)
	dataSync := worker.NewDgraphDataSync(dgraphHelper)

	err := dataSync.AlterSchema(worker.SchemaIndexStatus{}.MetaSchemes())
	if err != nil {
		zlog.Errorf("initialize schema failed: %s", err)
		return err
	}

	schemaCenter = worker.NewSchemaStatusCenter(dataSync)
	err = schemaCenter.Prepare()
	if err != nil {
		return err
	}
	apiHost := args.ApiHost
	httpPrefix := "http://"
	if strings.HasPrefix(apiHost, httpPrefix) == false {
		apiHost = httpPrefix + apiHost
	}
	monitor = worker.NewSchemaMonitor(schemaCenter)
	workerManager = worker.NewValueIndexWorkerManager(apiHost, dataSync)
	monitor.AddListener(workerManager, schemaCenter)

	router := initRouter()
	return router.Run(fmt.Sprintf("0.0.0.0:%d", args.Port))
}
