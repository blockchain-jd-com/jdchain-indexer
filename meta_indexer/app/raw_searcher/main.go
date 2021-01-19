package main

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/adaptor"
	"github.com/mkideal/cli"
	"github.com/ssor/zlog"
	"strings"
)

type Args struct {
	cli.Helper
	ApiHost   string `cli:"api" usage:"api server host, like http://127.0.0.1:8080" dft:""`
	Keyword   string `cli:"w,keyword" usage:"keyword to search for" dft:""`
	Ledger    string `cli:"l,ledger" usage:"ledger to search in" dft:""`
	BlockFrom int64  `cli:"from" usage:"start search from block" dft:"1"`
	BlockTo   int64  `cli:"to" usage:"start search to block" dft:"100"`
}

func main() {
	cli.Run(new(Args), func(ctx *cli.Context) error {
		args := ctx.Argv().(*Args)
		zlog.WithFields(map[string]interface{}{
			"api":     args.ApiHost,
			"ledger":  args.Ledger,
			"from":    args.BlockFrom,
			"to":      args.BlockTo,
			"kewword": args.Keyword,
		}).Info("Run with args:")

		if len(args.Keyword) <= 0 {
			return fmt.Errorf("keyword is empty")
		}

		startupServer(args.ApiHost, args.Keyword, args.Ledger, args.BlockFrom, args.BlockTo)

		//waitForStopSignal()
		return nil
	})
}

func startupServer(apiHost string, keyword string, ledger string, from, to int64) {
	fu := &ContentSearcher{
		keyword: keyword,
		apiHost: apiHost,
	}
	for i := from; i <= to; i++ {
		ok, err := fu.Search(ledger, i)
		if err != nil {
			zlog.Errorf("stopped")
			return
		}
		if ok {
			zlog.Infof("find keyword in block [%d]", i)
			return
		}
	}
	zlog.Infof("no result found for keyword in blocks[%d - %d]", from, to)
}

type ContentSearcher struct {
	keyword string
	apiHost string
}

func (searcher *ContentSearcher) Search(ledger string, blockHeight int64) (isFound bool, e error) {
	zlog.Infof("start search ledger: [%s] from height: [%d] ", ledger, blockHeight)
	body, debugInfo, err := adaptor.GetTxListInBlockRawFromServer(searcher.apiHost, ledger, blockHeight, 0, 100)
	if err != nil {
		e = err
		zlog.WithStruct(debugInfo).Errorf("fetch block failed: %s", err)
		return
	}

	if strings.Contains(string(body), searcher.keyword) {
		isFound = true
	}
	return
}
