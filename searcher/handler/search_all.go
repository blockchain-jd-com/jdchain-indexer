package handler

import (
	"bufio"
	"errors"
	"git.jd.com/jd-blockchain/explorer/response"
	"git.jd.com/jd-blockchain/explorer/searcher/query"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func isDebugOn(debug string) bool {
	return debug == "yes"
}

type QueryExecutor interface {
	OutputDebugInfo() interface{}
}

func doQueryResponse(c *gin.Context, data interface{}, e error, isDebug bool, qe QueryExecutor) {
	if isDebug {
		c.String(http.StatusOK, qe.OutputDebugInfo().(string))
	} else {
		res := response.NewResponse(data, e)
		c.JSON(http.StatusOK, res)
	}
}

type SearchAllResult struct {
	QueryResult *QueryResult
	Total       int64 `json:"total"`
}

func (result *SearchAllResult) ToJSON(writer *bufio.Writer) (e error) {

	e = result.QueryResult.ToJSON(writer)
	if e != nil {
		return e
	}

	if result.Total > 0 {
		writer.WriteString(",\"total\":" + strconv.FormatInt(result.Total, 10))
	} else {
		writer.WriteString("\"total\": 0")
	}

	return
}

func HandleSearch(c *gin.Context) {
	var obj struct {
		Keyword string `form:"keyword"`
		IsDebug string `form:"debug"`
	}
	err := c.BindQuery(&obj)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(paraError))
		return
	}

	ledger := c.Param("ledger")
	if len(strings.TrimSpace(obj.Keyword)) < 20 {
		res := response.NewResponse(nil, errors.New("length of keyword must >= 20"))
		c.JSON(http.StatusOK, res)
	} else {
		qe := query.NewSearchAllByKeyword([]string{ledger}, obj.Keyword)
		blocks, txs, users, accounts, contracts, eventAccounts, e := qe.DoQuery(dgClient)
		result := &SearchAllResult{
			QueryResult: &QueryResult{
				combine:       true,
				Blocks:        blocks,
				Txs:           txs,
				Users:         users,
				Accounts:      accounts,
				Contracts:     contracts,
				EventAccounts: eventAccounts,
			},
			Total: int64(len(blocks) + len(txs) + len(users) + len(accounts) + len(contracts) + len(eventAccounts)),
		}
		doQueryResponse(c, result, e, isDebugOn(obj.IsDebug), qe)
	}
}

func parseLedgers(all string) []string {
	return strings.Split(all, ",")
}
