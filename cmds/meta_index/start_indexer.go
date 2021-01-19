package meta_index

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/chain"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/meta_level_task"
	"github.com/mkideal/cli"
	"github.com/ssor/zlog"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
)

var ConvertStart = &cli.Command{
	Name: "ledger-rdf",
	Desc: "index meta data of ledger",
	Argv: func() interface{} { return new(LedgerRDFArg) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*LedgerRDFArg)
		StartLedgerServer(argv)
		return nil
	},
}

// root command
type LedgerRDFArg struct {
	cli.Helper
	ApiHost    string `cli:"*ledger-host" usage:"api server host, like http://127.0.0.1:8080" dft:""`
	DgraphHost string `cli:"dgraph" usage:"dgraph server host" dft:"127.0.0.1:9080"`
	Production bool   `cli:"production" usage:"if use production mode" dft:"false"`
}

func StartLedgerServer(argv *LedgerRDFArg) {
	apiHost, dgraphHost := argv.ApiHost, argv.DgraphHost
	production := argv.Production

	if production {
		fmt.Println("** Production ** mode")

		logFile := &lumberjack.Logger{
			Filename:   "ledger-rdf.log",
			MaxSize:    30, // megabytes
			MaxBackups: 10,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		}
		zlog.SetLevel(zlog.InfoLevel)
		zlog.SetOutput(logFile)
	}

	httpPrefix := "http://"
	if strings.HasPrefix(apiHost, httpPrefix) == false {
		apiHost = httpPrefix + apiHost
	}

	if strings.HasPrefix(dgraphHost, "http://") {
		dgraphHost = strings.Replace(dgraphHost, httpPrefix, "", 1)
	}
	dgraphHelper := dgraph_helper.NewHelper(dgraphHost)

	chainMonitor := chain.NewChainMonitor(apiHost)
	creatorManager := meta_level_task.NewMetaInfoLevelTaskCreatorManager(apiHost, dgraphHelper)

	chainMonitor.AddListeners(creatorManager)

	chainMonitor.Run(0)

	select {}
}
