package main

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/cmds/task_monitor"
	"github.com/mkideal/cli"
	"os"
)

func main() {
	if err := cli.Root(task_monitor.Root).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
