package main

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/app/rdf_generator/cmd"
	"github.com/mkideal/cli"
	"os"
)

func main() {
	if err := cli.Root(cmd.Root,
		cli.Tree(cmd.DemoRDFCommand),
		cli.Tree(help),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var help = cli.HelpCommand("help:")
