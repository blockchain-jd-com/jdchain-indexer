package handler

import (
	"git.jd.com/jd-blockchain/explorer/response"
	"git.jd.com/jd-blockchain/explorer/searcher/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleQueryTxCountByHash(c *gin.Context) {
	var obj struct {
		Keyword string `form:"keyword"`
		IsDebug string `form:"debug"`
	}
	err := c.BindQuery(&obj)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}

	ledgers := parseLedgers(c.Param("ledger"))
	qe := query.NewQueryTxCountByKeyword(ledgers, obj.Keyword)
	count, e := qe.DoQuery(dgClient)
	doQueryResponse(c, count, e, isDebugOn(obj.IsDebug), qe)
}

func HandleQueryTxByHash(c *gin.Context) {
	var obj struct {
		Keyword string `form:"keyword"`
		From    int64  `form:"from"`
		Count   int64  `form:"count"`
		IsDebug string `form:"debug"`
	}
	err := c.BindQuery(&obj)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}

	ledgers := parseLedgers(c.Param("ledger"))
	qe := query.NewQueryTxByKeyword(ledgers, obj.Keyword, obj.From, obj.Count)
	txs, err := qe.DoQuery(dgClient)
	doQueryResponse(c, &QueryResult{Txs: txs.(query.Transactions)}, err, isDebugOn(obj.IsDebug), qe)
}

func HandleQueryTxCountByEndpoint(c *gin.Context) {
	var obj struct {
		Keyword string `form:"keyword"`
		IsDebug string `form:"debug"`
	}
	err := c.BindQuery(&obj)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}

	ledgers := parseLedgers(c.Param("ledger"))
	qe := query.NewQueryTxCountByEndpointUser(ledgers, obj.Keyword)
	count, e := qe.DoQuery(dgClient)
	doQueryResponse(c, count, e, isDebugOn(obj.IsDebug), qe)
}

func HandleQueryTxByEndpointUser(c *gin.Context) {
	var obj struct {
		Keyword string `form:"keyword"`
		From    int64  `form:"from"`
		Count   int64  `form:"count"`
		IsDebug string `form:"debug"`
	}
	err := c.BindQuery(&obj)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}

	ledgers := parseLedgers(c.Param("ledger"))
	qe := query.NewQueryTxByEndpoint(ledgers, obj.Keyword, obj.From, obj.Count)
	txs, err := qe.DoQuery(dgClient)
	doQueryResponse(c, &QueryResult{Txs: txs.(query.Transactions)}, err, isDebugOn(obj.IsDebug), qe)
}

func HandleQueryTxRange(c *gin.Context) {
	var obj struct {
		Height  int64  `form:"height"`
		From    int64  `form:"from"`
		Count   int64  `form:"count"`
		Ledgers string `form:"ledgers"`
		IsDebug string `form:"debug"`
	}
	err := c.BindQuery(&obj)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}

	if obj.Height < 0 {
		c.JSON(http.StatusOK, response.NewFailedResponse("区块高度 height 不能小于 0"))
		return
	}
	qe := query.NewQueryTxRangeInBlock(parseLedgers(obj.Ledgers), obj.Height, obj.From, obj.Count)
	txs, err := qe.DoQuery(dgClient)
	doQueryResponse(c, &QueryResult{Txs: txs.(query.Transactions)}, err, isDebugOn(obj.IsDebug), qe)
}
