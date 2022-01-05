package meta_index

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/searcher/handler"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/cli"
	"github.com/ssor/zlog"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

var ApiServer = &cli.Command{
	Name: "api-server",
	Desc: "server api for blockchain data search",
	Argv: func() interface{} { return new(ApiServerArgs) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*ApiServerArgs)
		StartApiServer(argv.DgraphHost, argv.Host, argv.Port, argv.Production)
		return nil
	},
}

type ApiServerArgs struct {
	cli.Helper
	Host       string `cli:"o,host" usage:"listening host" dft:"0.0.0.0"`
	Port       int    `cli:"p,port" usage:"listening port" dft:"10001"`
	DgraphHost string `cli:"dgraph" usage:"dgraph server host" dft:"127.0.0.1:9080"`
	Production bool   `cli:"production" usage:"if use production mode" dft:"false"`
}

func StartApiServer(dgraphHost string, listeningHost string, listeningPort int, production bool) {

	if production {
		fmt.Println("** Production ** mode")

		logFile := &lumberjack.Logger{
			Filename:   "api-server.log",
			MaxSize:    30, // megabytes
			MaxBackups: 10,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		}
		zlog.SetOutput(logFile)
	}

	httpPrefix := "http://"

	if strings.HasPrefix(dgraphHost, "http://") {
		dgraphHost = strings.Replace(dgraphHost, httpPrefix, "", 1)
	}
	if strings.HasPrefix(listeningHost, "http://") {
		listeningHost = strings.Replace(listeningHost, httpPrefix, "", 1)
	}

	handler.Init(dgraphHost)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ledgers/:ledger/all/search", handler.HandleSearch)

	r.GET("/ledgers/:ledger/contracts/search", handler.HandleQueryContractByHash)
	r.GET("/ledgers/:ledger/contracts/count/search", handler.HandleQueryContractCountByHash)

	r.GET("/ledgers/:ledger/blocks/search", handler.HandleQueryBlockByHash)
	r.GET("/ledgers/:ledger/blocks/count/search", handler.HandleQueryBlockCountByHash)

	r.GET("/ledgers/:ledger/txs/search", handler.HandleQueryTxByHash)
	r.GET("/ledgers/:ledger/txs/count/search", handler.HandleQueryTxCountByHash)

	// 按时间查询
	r.GET("/ledgers/:ledger/txs/count/from/:from/to/:to", handler.HandleQueryTxCountByTime)
	r.GET("/ledgers/:ledger/txs/from/:from/to/:to", handler.HandleQueryTxByTime)

	r.GET("/ledgers/:ledger/users/txs/search", handler.HandleQueryTxByEndpointUser)
	r.GET("/ledgers/:ledger/users/txs/count/search", handler.HandleQueryTxCountByEndpoint)

	r.GET("/ledgers/:ledger/users/search", handler.HandleQueryUserByHash)
	r.GET("/ledgers/:ledger/users/count/search", handler.HandleQueryUserCountByHash)

	r.GET("/ledgers/:ledger/accounts/search", handler.HandleQueryDataAccountByHash)
	r.GET("/ledgers/:ledger/accounts/count/search", handler.HandleQueryDataAccountCountByHash)

	r.GET("/ledgers/:ledger/eventAccounts/search", handler.HandleQueryEventAccountByHash)
	r.GET("/ledgers/:ledger/eventAccounts/count/search", handler.HandleQueryEventAccountCountByHash)

	r.GET("/ledgers/:ledger/kvs/users/search", handler.HandleQueryKvEndpointUser)

	err := r.Run(fmt.Sprintf("%s:%d", listeningHost, listeningPort))
	if err != nil {
		os.Exit(1)
	}
}
