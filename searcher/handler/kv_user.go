package handler

import (
	"git.jd.com/jd-blockchain/explorer/response"
	"git.jd.com/jd-blockchain/explorer/searcher/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleQueryKvEndpointUser(c *gin.Context) {
	var obj struct {
		Keyword string `form:"keyword"`
		Account string `form:"account"`
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
	qe := query.NewQueryKvEndpointUser(ledgers, obj.Account, obj.Keyword, obj.From, obj.Count)
	kvUsers, err := qe.DoQuery(dgClient)
	doQueryResponse(c, &QueryResult{KvUsers: kvUsers.(query.KvEndpointUsers)}, err, isDebugOn(obj.IsDebug), qe)
}
