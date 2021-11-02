package handler

import (
	"git.jd.com/jd-blockchain/explorer/response"
	"git.jd.com/jd-blockchain/explorer/searcher/query"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

func HandleQueryTxCountByTime(c *gin.Context) {
	ledgers := parseLedgers(c.Param("ledger"))
	from, err := strconv.ParseInt(c.Param("from"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}
	to, err := strconv.ParseInt(c.Param("to"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}
	debug := c.Query("debug")
	qe := query.NewQueryTxCountByTime(ledgers, from, to)
	count, e := qe.DoQuery(dgClient)
	doQueryResponse(c, count, e, isDebugOn(debug), qe)
}

func HandleQueryTxByTime(c *gin.Context) {
	ledgers := parseLedgers(c.Param("ledger"))
	from, err := strconv.ParseInt(c.Param("from"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}
	to, err := strconv.ParseInt(c.Param("to"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}
	count, err := strconv.ParseInt(c.DefaultQuery("count", "1000"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}
	if count > 1000 {
		count = 1000
	}
	debug := c.Query("debug")
	qe := query.NewQueryTxByTime(ledgers, from, to, count)
	txs, e := qe.DoQuery(dgClient)
	doQueryResponse(c, txs, e, isDebugOn(debug), qe)
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
		Ledger  string `form:"ledger"`
		IsDebug string `form:"debug"`
	}
	err := c.BindQuery(&obj)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}
	obj.Ledger = c.Param("ledger")

	if obj.Height < 0 {
		c.JSON(http.StatusOK, response.NewFailedResponse("区块高度 height 不能小于 0"))
		return
	}
	qe := query.NewQueryTxRangeInBlock(parseLedgers(obj.Ledger), obj.Height, obj.From, obj.Count)
	txs, err := qe.DoQuery(dgClient)
	doQueryResponse(c, &QueryResult{Txs: txs.(query.Transactions)}, err, isDebugOn(obj.IsDebug), qe)
}
