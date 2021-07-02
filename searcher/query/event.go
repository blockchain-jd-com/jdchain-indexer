package query

import (
	"context"
	"github.com/dgraph-io/dgo/v200"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

var (
	qlEventAccountCountQueryResultName = "query_event_account_count"
	qlNameEventAccountCountByHash      = `
        var(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-event_account @filter(regexp(event_account-address, /\S*[[keyword]]\S*/) or regexp(event_account-public_key,    /\S*[[keyword]]\S*/  ) )
            {
               counts as count(event_account-public_key)
            }
        }
		query_event_account_count() {
            event_account-count: sum(val(counts))
	   }
    `
	qlEventAccountQueryResultName = "query_event_accounts"
	qlNameEventAccountByHash      = `
        query_event_accounts(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-event_account @filter(regexp(event_account-address, /\S*[[keyword]]\S*/) or regexp(event_account-public_key,    /\S*[[keyword]]\S*/  ) ) (first:[[count]], offset:[[offset]])
            {
                event_account-address:event_account-address
                event_account-public_key:event_account-public_key
            }
        }
    `
	qlNameEventAccountRange = `
        query_event_accounts(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-event_account (first:[[count]], offset:[[offset]])
            {
                event_account-address:event_account-address
                event_account-public_key:event_account-public_key
            }
        }
    `
)

func NewQueryEventAccountCountByKeyword(ledgers []string, keyword string) *EventAccountCountQuery {
	query := newQueryLang(qlNameEventAccountCountByHash, qlEventAccountCountQueryResultName)
	return &EventAccountCountQuery{
		lang: query,
		args: map[string]interface{}{
			"keyword": keyword,
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

type EventAccountCountQuery struct {
	lang *QueryLang
	args map[string]interface{}
}

func (query *EventAccountCountQuery) DoQuery(client *dgo.Dgraph) (interface{}, error) {
	ctx := context.Background()
	resp, err := client.NewTxn().Query(ctx, query.lang.Assemble(query.args))
	if err != nil {
		logger.Warn("While try to query failed: ", err)
		logger.Infof("\n%s\n", query.lang.Assemble(query.args))
		return nil, err
	}

	result := gjson.Get(string(resp.Json), query.lang.resultName)
	if len(result.Array()) <= 0 {
		//logger.Debugf("While try to query empty")
		//logger.Debugf("\n%s\n", query.lang.Assemble(query.args))
		return 0, nil
	}

	count := result.Array()[0].Get("event_account-count").Int()
	return count, nil
}

func (query *EventAccountCountQuery) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func NewEventAccountRangeQuery(ledgers []string, from, count int64) *EventAccountQuery {
	if count > MaxRecordsPerRequest {
		count = MaxRecordsPerRequest
	}
	return &EventAccountQuery{
		lang: newQueryLang(qlNameEventAccountRange, qlEventAccountQueryResultName),
		args: map[string]interface{}{
			"count":   strconv.FormatInt(count, 10),
			"offset":  strconv.FormatInt(from, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

func NewEventAccountHasKeyOrHasAddressQuery(ledgers []string, keyword string, from, count int64) *EventAccountQuery {
	if count > MaxRecordsPerRequest {
		count = MaxRecordsPerRequest
	}
	return &EventAccountQuery{
		lang: newQueryLangGroup().
			addQuery(newQueryLang(qlNameEventAccountByHash, qlEventAccountQueryResultName)),
		args: map[string]interface{}{
			"keyword": keyword,
			"count":   strconv.FormatInt(count, 10),
			"offset":  strconv.FormatInt(from, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

type EventAccountQuery struct {
	lang QueryAssembler
	args map[string]interface{}
}

func (query *EventAccountQuery) DoQuery(client *dgo.Dgraph) (interface{}, error) {
	return doEventAccountQuery(query.lang, query.args, client)
}

func (query *EventAccountQuery) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func doEventAccountQuery(ql QueryAssembler, args map[string]interface{}, client *dgo.Dgraph) (accounts EventAccounts, e error) {
	ctx := context.Background()
	resp, err := client.NewTxn().Query(ctx, ql.Assemble(args))
	if err != nil {
		logger.Warn("While try to query failed: ", err)
		logger.Infof("\n%s\n", ql.Assemble(args))
		return nil, err
	}
	for _, resultName := range ql.ResultNames() {
		accounts = accounts.Append(parseEventAccountInfo(gjson.Get(string(resp.Json), resultName))...)
	}
	if len(accounts) <= 0 {
		//logger.Debugf("While try to query empty")
		//logger.Debugf("\n%s\n", ql.Assemble(args))
		accounts = EventAccounts{}
	}

	return
}

func parseEventAccountInfo(result gjson.Result) (sets EventAccounts) {
	for _, setRaw := range result.Array() {
		var set EventAccount
		set.Address = setRaw.Get("event_account-address").String()
		set.PublicKey = setRaw.Get("event_account-public_key").String()
		sets = append(sets, &set)
	}
	return
}
