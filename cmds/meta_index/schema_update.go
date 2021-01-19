package meta_index

import (
	"git.jd.com/jd-blockchain/explorer/adaptor"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/level_task"
	"github.com/mkideal/cli"
	"github.com/ssor/zlog"
)

type SchemaUpdateArgs struct {
	cli.Helper
	DgraphHost string `cli:"dgraph" usage:"dgraph server host" dft:"127.0.0.1:9080"`
}

var UpdateSchema = &cli.Command{
	Name: "schema-update",
	Desc: "update schema in dgraph",
	Argv: func() interface{} { return new(SchemaUpdateArgs) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*SchemaUpdateArgs)
		return alterSchema(argv.DgraphHost)
	},
}

func alterSchema(host string) error {
	var schemas dgraph_helper.Schemas
	schemas = schemas.Add(adaptor.MetaSchemas...)
	schemas = schemas.Add(level_task.LevelTaskSchemas...)

	helper := dgraph_helper.NewHelper(host)
	err := helper.Alter(schemas)
	if err != nil {
		zlog.Failedf("alter schema failed: %s", err)
		return err
	}
	zlog.Success("alter schema success")
	return nil
}
