package main

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/level_task"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/meta_level_task"
	"git.jd.com/jd-blockchain/explorer/performance"
	"github.com/mkideal/cli"
	"github.com/ssor/zlog"
	"os"
	"os/signal"
	"strings"
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

var CreateLevelTask = &cli.Command{
	Name: "task",
	Desc: "create meta info level task",
	Argv: func() interface{} { return new(LedgerRDFArgs) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*LedgerRDFArgs)
		startWriteTasks(argv)
		return nil
	},
}

// root command
type LedgerRDFArgs struct {
	cli.Helper
	ApiHost       string `cli:"api" usage:"api server host, like http://127.0.0.1:8080" dft:""`
	DgraphHost    string `cli:"dgraph" usage:"dgraph server host" dft:"127.0.0.1:9080"`
	Ledger        string `cli:"l,ledger" usage:"ledger to search in" dft:""`
	BlockFrom     int64  `cli:"from" usage:"start from block" dft:"1"`
	BlockTo       int64  `cli:"to" usage:"stop at block" dft:"100"`
	ProfileCPU    bool   `cli:"cpu" usage:"profile cpu" dft:"false"`
	ProfileMemory bool   `cli:"memory" usage:"profile memory" dft:"false"`
}

var (
	dgraphHelper *dgraph_helper.Helper
)

func startWriteTasks(args *LedgerRDFArgs) {

	if args.ProfileCPU {
		performance.StartCpuProfile()
	}

	dgraphHelper = dgraph_helper.NewHelper(args.DgraphHost)

	createTasks(args.Ledger, args.BlockFrom, args.BlockTo)

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

func startLedgerServer(args *LedgerRDFArgs) {

	if args.ProfileCPU {
		performance.StartCpuProfile()
	}

	dgraphHelper = dgraph_helper.NewHelper(args.DgraphHost)

	if err := prepareLedgerNode(args.Ledger, dgraphHelper); err != nil {
		return
	}

	startTaskMonitor(args.ApiHost, args.Ledger, dgraphHelper)

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

func startTaskMonitor(apiHost, ledger string, dgraphHelper *dgraph_helper.Helper) {
	handler := meta_level_task.NewLevelTaskHandler(apiHost, ledger, dgraphHelper, func(s string) error {
		_, err := dgraphHelper.MutationRdfs([]byte(s))
		return err
	}).Setup()

	parsers := level_task.LevelTaskParserMap{
		meta_level_task.TaskMetaTaskLevelName: meta_level_task.ParseMetaInfoLevelTask,
	}

	level_task.NewLevelTaskMonitor(handler, parsers, dgraphHelper).Setup()
}

func createTasks(ledger string, from, to int64) {
	var lts []level_task.LevelTask
	for i := from; i <= to; i++ {
		lt := meta_level_task.CreateNewMetaInfoLevelTask("", ledger, i, level_task.MetaInfoLevel1)
		lts = append(lts, lt)
	}
	if err := writeRawLevelTasks(lts, dgraphHelper); err != nil {
		panic(err)
	}
}

func writeRawLevelTasks(tasks []level_task.LevelTask, dgraphHelper *dgraph_helper.Helper) error {
	var builder strings.Builder
	for _, task := range tasks {
		builder.WriteString(task.CreateMutations().Assembly())
	}
	fmt.Println(builder.String())
	_, err := dgraphHelper.MutationRdfs([]byte(builder.String()))
	return err
}

func prepareLedgerNode(hash string, dgraphHelper *dgraph_helper.Helper) error {
	_, _, err := meta_level_task.PrepareLedgerNode(dgraphHelper, hash)
	if err != nil {
		zlog.Warnf("prepare ledger [%s] failed: %s", hash, err)
		return err
	}
	return nil
}
