package handler

import (
	"git.jd.com/jd-blockchain/explorer/response"
	"git.jd.com/jd-blockchain/explorer/searcher/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleQueryLedgerRange(c *gin.Context) {
	var obj struct {
		IsDebug string `form:"debug"`
	}
	err := c.BindQuery(&obj)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}

	qe := query.NewLedgerQuery()
	ledgers, e := qe.DoQuery(dgClient)
	doQueryResponse(c, &QueryResult{Ledgers: ledgers.(query.Ledgers)}, e, isDebugOn(obj.IsDebug), qe)
}
