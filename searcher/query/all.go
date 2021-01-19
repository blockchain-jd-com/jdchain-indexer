package query

import (
	"context"
	"github.com/dgraph-io/dgo/v200"
	"github.com/tidwall/gjson"
	"math"
	"strconv"
	"strings"
)

func NewSearchAllByKeyword(ledgers []string, keyword string) *SearchAllByKeyword {
	return &SearchAllByKeyword{
		lang: newQueryLangGroup().
			addQuery(newQueryLang(qlNameBlocksByHash, qlBlockQueryResultName)).
			addQuery(newQueryLang(qlNameUsersByHash, qlUserQueryResultName)).
			addQuery(newQueryLang(qlNameDatasetByHash, qlDatasetQueryResultName)).
			addQuery(newQueryLang(qlNameContractByHash, qlContractQueryResultName)).
			addQuery(newQueryLang(qlNameTxsByEndpointUser, qlTxQueryByEndpointUserResultName)).
			addQuery(newQueryLang(qlNameTxsByHash, qlTxQueryResultName)).
			addQuery(newQueryLang(qlNameEventAccountByHash, qlEventAccountQueryResultName)),
		args: map[string]interface{}{
			"keyword": keyword,
			"count":   strconv.FormatInt(math.MaxInt32, 10),
			"offset":  strconv.FormatInt(0, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

type SearchAllByKeyword struct {
	lang *QueryLangGroup
	args map[string]interface{}
}

func (query *SearchAllByKeyword) DoQuery(client *dgo.Dgraph) (blocks Blocks,
	txs Transactions, users Users, accounts Accounts, contracts Contracts, eventAccounts EventAccounts, e error) {
	return searchByKeyword(query.lang, query.args, client)
}

func (query *SearchAllByKeyword) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func searchByKeyword(ql *QueryLangGroup, args map[string]interface{}, client *dgo.Dgraph) (blocks Blocks,
	txs Transactions, users Users, accounts Accounts, contracts Contracts, eventAccounts EventAccounts, e error) {
	ctx := context.Background()
	resp, err := client.NewTxn().Query(ctx, ql.Assemble(args))
	if err != nil {
		e = err
		logger.Warn("While try to query failed: ", err)
		return
	}

	blocks = parseBlockInfo(gjson.Get(string(resp.Json), qlBlockQueryResultName))
	txs = txs.add(parseTxInfo(gjson.Get(string(resp.Json), qlTxQueryResultName))...)
	users = parseUserInfo(gjson.Get(string(resp.Json), qlUserQueryResultName))
	accounts = parseDatasetInfo(gjson.Get(string(resp.Json), qlDatasetQueryResultName))
	contracts = parseContractInfo(gjson.Get(string(resp.Json), qlContractQueryResultName))
	eventAccounts = parseEventAccountInfo(gjson.Get(string(resp.Json), qlEventAccountQueryResultName))

	return
}
