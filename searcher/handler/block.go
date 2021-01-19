package handler

import (
	"git.jd.com/jd-blockchain/explorer/response"
	"git.jd.com/jd-blockchain/explorer/searcher/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleQueryBlockCountByHash(c *gin.Context) {
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
	qe := query.NewBlockCountQueryByKewword(ledgers, obj.Keyword)
	count, e := qe.DoQuery(dgClient)
	doQueryResponse(c, count, e, isDebugOn(obj.IsDebug), qe)
}

func HandleQueryBlockByHash(c *gin.Context) {
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
	qe := query.NewBlockQueryByKeyword(ledgers, obj.Keyword, obj.From, obj.Count)
	blocks, e := qe.DoQuery(dgClient)
	doQueryResponse(c, &QueryResult{Blocks: blocks.(query.Blocks)}, e, isDebugOn(obj.IsDebug), qe)
}

func HandleQueryBlockRange(c *gin.Context) {
	var obj struct {
		From    int64  `form:"from"`
		To      int64  `form:"to"`
		Ledgers string `form:"ledgers"`
		IsDebug string `form:"debug"`
	}
	err := c.BindQuery(&obj)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}

	qe := query.NewBlockQueryInRange(parseLedgers(obj.Ledgers), obj.From, obj.To)
	blocks, e := qe.DoQuery(dgClient)
	doQueryResponse(c, &QueryResult{Blocks: blocks.(query.Blocks)}, e, isDebugOn(obj.IsDebug), qe)
}
