package handler

import (
	"git.jd.com/jd-blockchain/explorer/response"
	"git.jd.com/jd-blockchain/explorer/searcher/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleQueryUserRange(c *gin.Context) {
	var obj struct {
		Ledgers string `form:"ledgers"`
		From    int64  `form:"fromIndex"`
		Count   int64  `form:"count"`
		IsDebug string `form:"debug"`
	}
	err := c.BindQuery(&obj)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}

	qe := query.NewQueryUsersRange(parseLedgers(obj.Ledgers), obj.From, obj.Count)
	users, err := qe.DoQuery(dgClient)
	doQueryResponse(c, &QueryResult{Users: users.(query.Users)}, err, isDebugOn(obj.IsDebug), qe)
}

func HandleQueryUserByHash(c *gin.Context) {
	var obj struct {
		Keyword string `form:"keyword"`
		From    int64  `form:"fromIndex"`
		Count   int64  `form:"count"`
		IsDebug string `form:"debug"`
	}
	err := c.BindQuery(&obj)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}

	ledgers := parseLedgers(c.Param("ledger"))
	qe := query.NewQueryUsersByHash(ledgers, obj.Keyword, obj.From, obj.Count)
	users, err := qe.DoQuery(dgClient)
	doQueryResponse(c, &QueryResult{Users: users.(query.Users)}, err, isDebugOn(obj.IsDebug), qe)
}

func HandleQueryUserCountByHash(c *gin.Context) {
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

	qe := query.NewQueryUserCountByHash(ledgers, obj.Keyword)
	count, err := qe.DoQuery(dgClient)
	doQueryResponse(c, count, err, isDebugOn(obj.IsDebug), qe)
}
