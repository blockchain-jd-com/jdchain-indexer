package query

import (
	"context"
	"github.com/dgraph-io/dgo/v200"
	"github.com/tidwall/gjson"
)

var (
	qlLedgerQueryResultName = "query_ledgers"
	qlNameLedgers           = `
		query_ledgers(func: has(ledger-hash_id)) @normalize
		{
			uid
		     status-ledger:ledger-hash_id
             status-height:count(ledger-block)
        }
    `
)

func NewLedgerQuery() *LedgerQuery {
	return &LedgerQuery{
		lang: newQueryLang(qlNameLedgers, qlLedgerQueryResultName),
		args: map[string]interface{}{},
	}
}

type LedgerQuery struct {
	lang *QueryLang
	args map[string]interface{}
}

func (query *LedgerQuery) DoQuery(client *dgo.Dgraph) (interface{}, error) {
	return doLedgerQuery(query.lang, query.args, client)
}

func (query *LedgerQuery) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func doLedgerQuery(ql *QueryLang, args map[string]interface{}, client *dgo.Dgraph) (ledgers Ledgers, e error) {
	ctx := context.Background()
	resp, err := client.NewTxn().Query(ctx, ql.Assemble(args))
	if err != nil {
		logger.Warnf("query ledgers failed: %s", err)
		logger.Info(ql.Assemble(args))
		return nil, err
	}

	ledgers = parseLedgerInfo(gjson.Get(string(resp.Json), ql.resultName))
	if len(ledgers) <= 0 {
		logger.Debugf("query ledgers empty")
		logger.Debugf(ql.Assemble(args))
	}
	return
}

func parseLedgerInfo(result gjson.Result) (ledgers []*LedgerInfo) {
	for _, ledger := range result.Array() {
		var li LedgerInfo
		li.HashID = ledger.Get("status-ledger").String()
		li.Height = int(ledger.Get("status-height").Int())
		ledgers = append(ledgers, &li)
	}
	return
}
