package handler

import (
	"git.jd.com/jd-blockchain/explorer/response"
	"git.jd.com/jd-blockchain/explorer/searcher/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
 * Author: imuge
 * Date: 2020/10/22 上午11:40
 */

func HandleQueryEventAccountsRange(c *gin.Context) {
	var obj struct {
		Ledger  string `form:"ledger"`
		IsDebug string `form:"debug"`
		From    int64  `form:"from"`
		Count   int64  `form:"count"`
	}
	err := c.BindQuery(&obj)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}
	obj.Ledger = c.Param("ledger")

	qe := query.NewEventAccountRangeQuery(parseLedgers(obj.Ledger), obj.From, obj.Count)
	accounts, err := qe.DoQuery(dgClient)
	doQueryResponse(c, accounts, err, isDebugOn(obj.IsDebug), qe)
}

func HandleQueryEventAccountByHash(c *gin.Context) {
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
	qe := query.NewEventAccountHasKeyOrHasAddressQuery(ledgers, obj.Keyword, obj.From, obj.Count)
	accounts, err := qe.DoQuery(dgClient)
	doQueryResponse(c, accounts, err, isDebugOn(obj.IsDebug), qe)
}

func HandleQueryEventAccountCountByHash(c *gin.Context) {
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

	qe := query.NewQueryEventAccountCountByKeyword(ledgers, obj.Keyword)
	count, err := qe.DoQuery(dgClient)
	doQueryResponse(c, count, err, isDebugOn(obj.IsDebug), qe)
}
