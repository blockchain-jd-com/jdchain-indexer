package drop_db

import (
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"github.com/mkideal/cli"
)

// demoRDF command
type DropDBArgs struct {
	cli.Helper
	DgraphHost string `cli:"dgraph" usage:"dgraph server host" dft:"127.0.0.1:9080"`
}

var DropDB = &cli.Command{
	Name: "drop",
	Desc: "clear data in dgraph",
	Argv: func() interface{} { return new(DropDBArgs) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*DropDBArgs)
		return resetDb(argv.DgraphHost)
	},
}

func resetDb(host string) error {
	helper := dgraph_helper.NewHelper(host)
	return helper.DropDB()
}
