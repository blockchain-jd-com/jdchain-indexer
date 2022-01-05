package query

import (
	"context"
	"github.com/dgraph-io/dgo/v200"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

var (
	qlEndpointUserByKvResultName = "query_kv_users"
	qlNameEndpointUserByKv       = `
	query_kv_users(func: has(ledger-hash_id)) @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize {
		ledger-tx {
			tx-kv @filter(eq(kv-data_account_address, "[[account]]") and regexp(kv-key, /\S*[[keyword]]\S*/)) (orderasc: kv-version, first:[[count]], offset:[[offset]]) {
				kv-key:kv-key
				kv-version:kv-version
				kv-data_account_address:kv-data_account_address
				~tx-kv {
					~endpoint_user-tx {
						user-address:user-address
					}
					tx-hash_id:tx-hash_id
					tx-block_height:tx-block_height
				}
			}
		}
	}
    `
)

func NewQueryKvEndpointUser(ledgers []string, account, keyword string, from, count int64) *QueryKvEndpointUser {
	if count > MaxRecordsPerRequest {
		count = MaxRecordsPerRequest
	}
	query := newQueryLang(qlNameEndpointUserByKv, qlEndpointUserByKvResultName)
	return &QueryKvEndpointUser{
		lang: query,
		args: map[string]interface{}{
			"account": account,
			"keyword": keyword,
			"count":   strconv.FormatInt(count, 10),
			"offset":  strconv.FormatInt(from, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

type QueryKvEndpointUser struct {
	lang *QueryLang
	args map[string]interface{}
}

func (query *QueryKvEndpointUser) DoQuery(client *dgo.Dgraph) (interface{}, error) {
	return doKvEndpointUserQuery(query.lang, query.args, client)
}

func (query *QueryKvEndpointUser) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func doKvEndpointUserQuery(ql *QueryLang, args map[string]interface{}, client *dgo.Dgraph) (kvusers KvEndpointUsers, e error) {
	ctx := context.Background()
	resp, err := client.NewTxn().Query(ctx, ql.Assemble(args))
	if err != nil {
		logger.Warn("While try to query failed: ", err)
		logger.Infof("\n%s\n", ql.Assemble(args))
		return nil, err
	}

	kvusers = parseKvEndpointUserInfo(gjson.Get(string(resp.Json), ql.resultName))
	if len(kvusers) <= 0 {
		kvusers = []KvEndpointUser{}
	}
	return
}

func parseKvEndpointUserInfo(result gjson.Result) (kvusers KvEndpointUsers) {
	for _, kus := range result.Array() {
		var ku KvEndpointUser
		ku.Account = kus.Get("kv-data_account_address").String()
		ku.Key = kus.Get("kv-key").String()
		ku.Version = kus.Get("kv-version").Int()
		ku.BlockHeight = kus.Get("tx-block_height").Int()
		ku.User = kus.Get("user-address").String()
		ku.Tx = kus.Get("tx-hash_id").String()
		kvusers = append(kvusers, ku)
	}
	return
}
