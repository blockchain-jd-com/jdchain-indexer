package query

import (
	"context"
	"github.com/dgraph-io/dgo/v200"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

var (
	qlTxQueryResultName = "query_txs"
	qlNameTxsByHash     = `
		query_txs(func: has(tx-execution_state), orderasc:tx-block_height, first:[[count]], offset:[[offset]]) 
			@filter(regexp(tx-hash_id, /\S*[[keyword]]\S*/)) @normalize @cascade {
			tx-execution_state:tx-execution_state
			tx-time:tx-time
			tx-block_height:tx-block_height
			tx-hash_id:tx-hash_id
			~block-tx {
				~ledger-block @filter(anyofterms(ledger-hash_id, "[[ledgers]]"))
			}		
		}
    `
	qlNameTxsRange = `
		query_txs(func: has(tx-execution_state), orderasc:tx-block_height, first:[[count]], offset:[[offset]]) @normalize @cascade {
			tx-execution_state:tx-execution_state
			tx-time:tx-time
			tx-block_height:tx-block_height
			tx-hash_id:tx-hash_id
			~block-tx {
				~ledger-block @filter(anyofterms(ledger-hash_id, "[[ledgers]]"))
			}		
		}
    `

	qlTxCountQueryResultName = "query_txs_count"
	qlNameTxsCountByHash     = `
       var(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
       {
            ledger-block {
                block-tx @filter(regexp(tx-hash_id, /\S*[[keyword]]\S*/)){
                    counts as count(tx-hash_id)
                }
            }
       }
       
       query_txs_count() {
            tx-count: sum(val(counts))
	   }
    `

	qlTxInBlockQueryResultName = "query_txs_in_block"
	qlNameTxsInBlock           = `
       query_txs_in_block(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
       {
            ledger-block 
            {
                block-tx @filter(eq(tx-block_height, [[height]])) (first:[[count]], offset:[[offset]])
                {
                    tx-execution_state:tx-execution_state
					tx-time:tx-time
                    tx-hash_id:tx-hash_id
                    tx-block_height:tx-block_height
                }
            }
       }
    `

	// query txs by endpoint user
	qlTxQueryByEndpointUserResultName = "query_txs_by_endpoint_user"
	qlNameTxsByEndpointUser           = `
		query_txs_by_endpoint_user(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
       {
            ledger-user @filter(regexp(user-public_key, /\S*[[endpointUser]]\S*/) or regexp(user-address, /\S*[[endpointUser]]\S*/))
			{
				endpoint_user-tx (orderasc: tx-block_height, first:[[count]], offset:[[offset]])
				{
					tx-execution_state:tx-execution_state
					tx-time:tx-time
					tx-hash_id:tx-hash_id
					tx-block_height:tx-block_height
				}
			}
       }
    `

	qlTxCountQueryByEndpointUserResultName = "query_txs_by_endpoint_user_count"
	qlNameTxsCountEndpointUser             = `
       var(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
       {
            ledger-user @filter(regexp(user-public_key, /\S*[[endpointUser]]\S*/) or regexp(user-address, /\S*[[endpointUser]]\S*/))
		   {
				counts as count(endpoint_user-tx)
		   }
       }
       query_txs_by_endpoint_user_count() {
            tx-count: sum(val(counts))
	   }
    `
)

func NewQueryTxCountByKeyword(ledgers []string, keyword string) *QueryTxCount {
	query := newQueryLang(qlNameTxsCountByHash, qlTxCountQueryResultName)
	return &QueryTxCount{
		lang: query,
		args: map[string]interface{}{
			"keyword": keyword,
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

type QueryTxCount QueryTx

func (query *QueryTxCount) DoQuery(client *dgo.Dgraph) (interface{}, error) {
	ctx := context.Background()
	resp, err := client.NewTxn().Query(ctx, query.lang.Assemble(query.args))
	if err != nil {
		logger.Warn("While try to query failed: ", err)
		logger.Infof("\n%s\n", query.lang.Assemble(query.args))
		return nil, err
	}

	result := gjson.Get(string(resp.Json), query.lang.resultName)
	if len(result.Array()) != 1 {
		return 0, nil
	}

	count := result.Array()[0].Get("tx-count").Int()
	return count, nil
}

func (query *QueryTxCount) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func NewQueryTxByKeyword(ledgers []string, keyword string, from, count int64) *QueryTx {
	if count > MaxRecordsPerRequest {
		count = MaxRecordsPerRequest
	}
	query := newQueryLang(qlNameTxsByHash, qlTxQueryResultName)
	return &QueryTx{
		lang: query,
		args: map[string]interface{}{
			"keyword": keyword,
			"count":   strconv.FormatInt(count, 10),
			"offset":  strconv.FormatInt(from, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

func NewQueryTxRange(ledgers []string, from, to int64) *QueryTx {
	count, offset := mapFromToToOffset(from, to, -1)
	query := newQueryLang(qlNameTxsRange, qlTxQueryResultName)
	return &QueryTx{
		lang: query,
		args: map[string]interface{}{
			"count":   strconv.FormatInt(count, 10),
			"offset":  strconv.FormatInt(offset, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

func NewQueryTxRangeInBlock(ledgers []string, height, from, count int64) *QueryTx {
	if count > MaxRecordsPerRequest {
		count = MaxRecordsPerRequest
	}
	query := newQueryLang(qlNameTxsInBlock, qlTxInBlockQueryResultName)
	return &QueryTx{
		lang: query,
		args: map[string]interface{}{
			"height":  strconv.FormatInt(height, 10),
			"count":   strconv.FormatInt(count, 10),
			"offset":  strconv.FormatInt(from, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

func NewQueryTxCountByEndpointUser(ledgers []string, endpointUser string) *QueryTxCount {
	query := newQueryLang(qlNameTxsCountEndpointUser, qlTxCountQueryByEndpointUserResultName)
	return &QueryTxCount{
		lang: query,
		args: map[string]interface{}{
			"endpointUser": endpointUser,
			"ledgers":      strings.Join(ledgers, " "),
		},
	}
}

func NewQueryTxByEndpoint(ledgers []string, endpointUser string, from, count int64) *QueryTx {
	if count > MaxRecordsPerRequest {
		count = MaxRecordsPerRequest
	}
	query := newQueryLang(qlNameTxsByEndpointUser, qlTxQueryByEndpointUserResultName)
	return &QueryTx{
		lang: query,
		args: map[string]interface{}{
			"endpointUser": endpointUser,
			"count":        strconv.FormatInt(count, 10),
			"offset":       strconv.FormatInt(from, 10),
			"ledgers":      strings.Join(ledgers, " "),
		},
	}
}

type QueryTx struct {
	lang *QueryLang
	args map[string]interface{}
}

func (query *QueryTx) DoQuery(client *dgo.Dgraph) (interface{}, error) {
	return doTxQuery(query.lang, query.args, client)
}

func (query *QueryTx) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func doTxQuery(ql *QueryLang, args map[string]interface{}, client *dgo.Dgraph) (txs Transactions, e error) {
	ctx := context.Background()
	resp, err := client.NewTxn().Query(ctx, ql.Assemble(args))
	if err != nil {
		logger.Warn("While try to query failed: ", err)
		logger.Infof("\n%s\n", ql.Assemble(args))
		return nil, err
	}

	txs = parseTxInfo(gjson.Get(string(resp.Json), ql.resultName))
	if len(txs) <= 0 {
		txs = []*Transaction{}
	}
	return
}

func parseTxInfo(result gjson.Result) (txs Transactions) {
	for _, txr := range result.Array() {
		var tx Transaction
		tx.HashID = txr.Get("tx-hash_id").String()
		tx.BlockHeight = txr.Get("tx-block_height").Int()
		tx.Time = txr.Get("tx-time").Int()
		tx.ExecutionState = txr.Get("tx-execution_state").String()
		txs = append(txs, &tx)
	}
	return
}
