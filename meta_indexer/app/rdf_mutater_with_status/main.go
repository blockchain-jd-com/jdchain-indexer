package main

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/cmds/drop_db"
	"git.jd.com/jd-blockchain/explorer/cmds/meta_index"
	"github.com/mkideal/cli"
	"os"
)

func main() {
	if err := cli.Root(Root,
		cli.Tree(CreateLevelTask),
		cli.Tree(meta_index.UpdateSchema),
		cli.Tree(drop_db.DropDB),
		cli.Tree(help),
	).Run(os.Args[1:]); err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}
}

var help = cli.HelpCommand("help:")
