package main

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/cmds/drop_db"
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

	if err := cli.Root(value_index.Root,
		cli.Tree(drop_db.DropDB),
		cli.Tree(help),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

var help = cli.HelpCommand("help:")

func printEvn() {
	fmt.Println(strings.Repeat("-", 128))
	fmt.Println("build-time: ", BUILD_TIME)
	fmt.Println("branch:     ", BRANCH)
	fmt.Println("version:    ", VERSION)
	fmt.Println("go-version: ", GO_VERSION)
	fmt.Println(strings.Repeat("-", 128))
}
