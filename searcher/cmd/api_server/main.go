package main

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/searcher/handler"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/cli"
	"github.com/ssor/zlog"
	"net/http"
	"os"
	"strings"
)

var (
	BRANCH     string
	VERSION    string
	BUILD_TIME string
	GO_VERSION string
)

type Args struct {
	cli.Helper
	Port       int    `cli:"p,port" usage:"listening port" dft:"8081"`
	DgraphHost string `cli:"dgraph" usage:"dgraph server host" dft:"127.0.0.1:9080"`
	CacheSize  int    `cli:"c,cache" usage:"cache size" dft:"10000"`
}

func main() {
	printEvn()

	cli.Run(new(Args), func(ctx *cli.Context) error {
		args := ctx.Argv().(*Args)
		zlog.WithFields(map[string]interface{}{
			"port":   args.Port,
			"dgraph": args.DgraphHost,
		}).Info("Run with args:")

		startServer(args.DgraphHost, args.Port, args.CacheSize)
		return nil
	})

}

func printEvn() {
	fmt.Println(strings.Repeat("-", 128))
	fmt.Println("build-time: ", BUILD_TIME)
	fmt.Println("branch:     ", BRANCH)
	fmt.Println("version:    ", VERSION)
	fmt.Println("go-version: ", GO_VERSION)
	fmt.Println(strings.Repeat("-", 128))
}

func startServer(dgraphHost string, listeningPort, cacheSize int) {
	handler.Init(dgraphHost)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/debug/info", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"build-time": BUILD_TIME,
			"branch":     BRANCH,
			"version":    VERSION,
			"go-version": GO_VERSION,
		})
	})
	r.GET("/ledgers/:ledger/all/search", handler.HandleSearch)

	r.GET("/ledgers/:ledger/contracts/search", handler.HandleQueryContractByHash)
	r.GET("/ledgers/:ledger/contracts/count/search", handler.HandleQueryContractCountByHash)

	r.GET("/ledgers/:ledger/blocks/search", handler.HandleQueryBlockByHash)
	r.GET("/ledgers/:ledger/blocks/count/search", handler.HandleQueryBlockCountByHash)

	r.GET("/ledgers/:ledger/txs/search", handler.HandleQueryTxByHash)
	r.GET("/ledgers/:ledger/txs/count/search", handler.HandleQueryTxCountByHash)

	r.GET("/ledgers/:ledger/users/txs/search", handler.HandleQueryTxByEndpointUser)
	r.GET("/ledgers/:ledger/users/txs/count/search", handler.HandleQueryTxCountByEndpoint)

	r.GET("/ledgers/:ledger/users/search", handler.HandleQueryUserByHash)
	r.GET("/ledgers/:ledger/users/count/search", handler.HandleQueryUserCountByHash)

	r.GET("/ledgers/:ledger/accounts/search", handler.HandleQueryDataAccountByHash)
	r.GET("/ledgers/:ledger/accounts/count/search", handler.HandleQueryDataAccountCountByHash)

	r.GET("/ledgers/:ledger/eventAccounts/search", handler.HandleQueryEventAccountByHash)
	r.GET("/ledgers/:ledger/eventAccounts/count/search", handler.HandleQueryEventAccountCountByHash)

	if err := r.Run(fmt.Sprintf("0.0.0.0:%d", listeningPort)); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
