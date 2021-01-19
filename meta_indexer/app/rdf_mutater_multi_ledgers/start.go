package main

import (
	"git.jd.com/jd-blockchain/explorer/chain"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/meta_level_task"
	"git.jd.com/jd-blockchain/explorer/performance"
	"github.com/mkideal/cli"
	"os"
	"os/signal"
	"syscall"
)

var Root = &cli.Command{
	Name: "ledger-rdf",
	Desc: "fetch data from ledger, and generate RDF mutations, and then commit to dgraph",
	Argv: func() interface{} { return new(LedgerRDFArgs) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*LedgerRDFArgs)
		startLedgerServer(argv)
		return nil
	},
}

// root command
type LedgerRDFArgs struct {
	cli.Helper
	ApiHost       string `cli:"api" usage:"api server host, like http://127.0.0.1:8080" dft:""`
	DgraphHost    string `cli:"dgraph" usage:"dgraph server host" dft:"127.0.0.1:9080"`
	ProfileCPU    bool   `cli:"cpu" usage:"profile cpu" dft:"false"`
	ProfileMemory bool   `cli:"memory" usage:"profile memory" dft:"false"`
}

var (
	dgraphHelper *dgraph_helper.Helper
)

func startLedgerServer(args *LedgerRDFArgs) {
	if args.ProfileCPU {
		performance.StartCpuProfile()
	}

	dgraphHelper = dgraph_helper.NewHelper(args.DgraphHost)

	chainMonitor := chain.NewChainMonitor(args.ApiHost)
	creatorManager := meta_level_task.NewMetaInfoLevelTaskCreatorManager(args.ApiHost, dgraphHelper)

	chainMonitor.AddListeners(creatorManager)

	chainMonitor.Run(0)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	if args.ProfileCPU {
		performance.StopCpuProfile()
	}
	if args.ProfileMemory {
		performance.StartMemoryProfile()
	}
}
