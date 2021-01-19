package query

import (
	"context"
	"github.com/dgraph-io/dgo/v200"
	"github.com/tidwall/gjson"
)

var (
	qlWriteKeyQueryResultName = "query_write_key"

	qlNameWriteKeyByKeyword = `
        query_write_key(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-dataset 
            {
                dataset-write_operation @filter(regexp(write_operation-key, /\S*[[keyword]]\S*/))
                {
                    write_operation-key:write_operation-key
                    write_operation-version:write_operation-version
                    write_operation-value:write_operation-value
                }
            }
        }
    `
)

type WriteKeyQuery struct {
	lang *QueryLang
	args map[string]interface{}
}

func (query *WriteKeyQuery) DoQuery(client *dgo.Dgraph) (interface{}, error) {
	return doWriteKeyQuery(query.lang, query.args, client)
}

func (query *WriteKeyQuery) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func doWriteKeyQuery(ql *QueryLang, args map[string]interface{}, client *dgo.Dgraph) (kvs WriteKvs, e error) {
	ctx := context.Background()
	resp, err := client.NewTxn().Query(ctx, ql.Assemble(args))
	if err != nil {
		logger.Warn("While try to query failed: ", err)
		logger.Infof("\n%s\n", ql.Assemble(args))
		return nil, err
	}

	kvs = parseWriteKV(gjson.Get(string(resp.Json), ql.resultName))
	return
}

func parseWriteKV(result gjson.Result) (kvs WriteKvs) {
	for _, raw := range result.Array() {
		var kv WriteKV
		kv.Key = raw.Get("write_operation-key").String()
		kv.Value = raw.Get("write_operation-value").String()
		kv.Version = raw.Get("write_operation-version").Int()
		kvs = append(kvs, &kv)
	}
	return
}
