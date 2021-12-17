package main

import (
	"fmt"
	"os"

	"git.jd.com/jd-blockchain/explorer/meta_indexer/app/rds_import/cmd"
	"github.com/mkideal/cli"
)

func main() {
	if err := cli.Root(cmd.Import,
		cli.Tree(help),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var help = cli.HelpCommand("help:")
