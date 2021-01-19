package main

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/cmds/drop_db"
	"git.jd.com/jd-blockchain/explorer/cmds/meta_index"
	"git.jd.com/jd-blockchain/explorer/cmds/value_index"
	"github.com/mkideal/cli"
	"os"
	"strings"
)

var (
	BRANCH     string
	VERSION    string
	BUILD_TIME string
	GO_VERSION string
)

func main() {
	printEvn()

	if err := cli.Root(argus,
		cli.Tree(helpRoot),
		cli.Tree(meta_index.ConvertStart),
		cli.Tree(meta_index.ApiServer),
		cli.Tree(meta_index.UpdateSchema),
		cli.Tree(value_index.Root),
		cli.Tree(drop_db.DropDB),
		//cli.Tree(task_monitor.Root),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// argus
var argus = &cli.Command{
	Name: "argus",
	Desc: "start all servers in argus",
	Argv: func() interface{} { return new(LedgerRDFArg) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*LedgerRDFArg)
		// start ledger-rdf
		go meta_index.StartLedgerServer(&meta_index.LedgerRDFArg{
			Helper:     argv.Helper,
			ApiHost:    argv.LedgerHost,
			DgraphHost: argv.DgraphHost,
			Production: argv.Production,
		})

		// start api-server
		go meta_index.StartApiServer(argv.DgraphHost, argv.ApiHost, argv.ApiPort, argv.Production)

		// start data
		go value_index.StartLedgerServer(&value_index.ServerArgs{
			Helper:     argv.Helper,
			ApiHost:    argv.LedgerHost,
			Port:       argv.SchemaPort,
			DgraphHost: argv.DgraphHost,
			Production: argv.Production,
		})

		// start task
		//go task_monitor.StartLedgerServer(&task_monitor.ServerArgs{
		//	Helper:     argv.Helper,
		//	Port:       argv.TaskPort,
		//	DgraphHost: argv.DgraphHost,
		//})

		select {}

		return nil
	},
}

// root command
type LedgerRDFArg struct {
	cli.Helper
	DgraphHost string `cli:"dgraph" usage:"dgraph server host" dft:"127.0.0.1:9080"`
	Production bool   `cli:"production" usage:"if use production mode" dft:"false"`
	// argus for ledger-rdf
	LedgerHost string `cli:"*ledger-host" usage:"api server host, like http://127.0.0.1:8080" dft:""`
	// argus for api-server
	ApiHost string `cli:"ao,api-host" usage:"argus ledger api listening host" dft:"0.0.0.0"`
	ApiPort int    `cli:"ap,api-port" usage:"argus ledger api listening port" dft:"10001"`
	// argus for data
	SchemaPort int `cli:"sp,schema-port" usage:"argus schema server listening port" dft:"8082"`
	// argus for task
	//TaskPort int `cli:"tp,task-port" usage:"argus value indexing task server listening port" dft:"10005"`
}

// help
var helpRoot = &cli.Command{
	Name: "help",
	Desc: "make index for data from jd-chain",
	Argv: func() interface{} { return new(HelpArgs) },
	Fn: func(ctx *cli.Context) error {
		return nil
	},
}

type HelpArgs struct {
	cli.Helper
}

func printEvn() {
	fmt.Println(strings.Repeat("-", 128))
	fmt.Println("build-time: ", BUILD_TIME)
	fmt.Println("branch:     ", BRANCH)
	fmt.Println("version:    ", VERSION)
	fmt.Println("go-version: ", GO_VERSION)
	fmt.Println(strings.Repeat("-", 128))
}
