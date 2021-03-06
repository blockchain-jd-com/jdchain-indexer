package handler

import (
	"errors"
	"git.jd.com/jd-blockchain/explorer/response"
	"git.jd.com/jd-blockchain/explorer/searcher/query"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func HandleQueryContractRange(c *gin.Context) {
	var obj struct {
		Ledgers string `form:"ledgers"`
		From    int64  `form:"from"`
		Count   int64  `form:"count"`
		IsDebug string `form:"debug"`
	}
	err := c.BindQuery(&obj)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}

	qe := query.NewQueryContractRange(parseLedgers(obj.Ledgers), obj.From, obj.Count)
	contracts, err := qe.DoQuery(dgClient)
	doQueryResponse(c, &QueryResult{Contracts: contracts.(query.Contracts)}, err, isDebugOn(obj.IsDebug), qe)

}

func HandleQueryContractCountByHash(c *gin.Context) {
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

	qe := query.NewQueryContractCountByHash(ledgers, obj.Keyword)
	count, err := qe.DoQuery(dgClient)
	doQueryResponse(c, count, err, isDebugOn(obj.IsDebug), qe)
}

func HandleQueryContractByHash(c *gin.Context) {
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

	if len(strings.TrimSpace(obj.Keyword)) < 20 {
		res := response.NewResponse(nil, errors.New("length of keyword must >= 20"))
		c.JSON(http.StatusOK, res)
	} else {
		qe := query.NewQueryContractByHash(ledgers, obj.Keyword, obj.From, obj.Count)
		contracts, err := qe.DoQuery(dgClient)
		doQueryResponse(c, &QueryResult{Contracts: contracts.(query.Contracts)}, err, isDebugOn(obj.IsDebug), qe)
	}
}
