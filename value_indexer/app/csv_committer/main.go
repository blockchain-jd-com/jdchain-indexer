package main

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/cmds/drop_db"
	"git.jd.com/jd-blockchain/explorer/value_indexer/app/csv_committer/cmd"
	"github.com/mkideal/cli"
	"os"
)

func main() {
	if err := cli.Root(cmd.Root,
		cli.Tree(drop_db.DropDB),
		cli.Tree(help),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

var help = cli.HelpCommand("help:")
