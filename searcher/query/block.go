package query

import (
	"context"
	"github.com/dgraph-io/dgo/v200"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

var (
	qlBlockCountQueryResultName = "query_blocks_count"
	qlNameBlocksCountByHash     = `
       var(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-block @filter(regexp(block-hash_id, /\S*[[keyword]]\S*/)) {
               counts as count(block-hash_id)
            } 
        }
		query_blocks_count() {
            block-count: sum(val(counts))
	   }
    `
	qlBlockQueryResultName = "query_blocks"
	qlNameBlocksByHash     = `
       query_blocks(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-block @filter(regexp(block-hash_id, /\S*[[keyword]]\S*/)) (orderasc:block-height, first:[[count]], offset:[[offset]]) {
                block-hash_id:block-hash_id
                block-height:block-height
				block-time:block-time
                tx-count:count(block-tx)
            } 
        }
    `
	qlNameBlockByHeight = `
       query_blocks(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-block(orderasc: block-height) @filter(ge(block-height, [[from]]) and le(block-height, [[to]]))  {
                block-hash_id:block-hash_id
                block-height:block-height
                block-time:block-time
                tx-count:count(block-tx)
            } 
        }
    `
)

func NewBlockCountQueryByKewword(ledgers []string, keyword string) *BlockCountQuery {
	return &BlockCountQuery{
		lang: newQueryLang(qlNameBlocksCountByHash, qlBlockCountQueryResultName),
		args: map[string]interface{}{
			"keyword": keyword,
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

type BlockCountQuery BlockQuery

func (query *BlockCountQuery) DoQuery(client *dgo.Dgraph) (int64, error) {
	ctx := context.Background()
	resp, err := client.NewTxn().Query(ctx, query.lang.Assemble(query.args))
	if err != nil {
		logger.Warnf("query blocks failed: %s", err)
		logger.Info(query.lang.Assemble(query.args))
		return 0, err
	}

	result := gjson.Get(string(resp.Json), query.lang.resultName)
	if len(result.Array()) <= 0 {
		//logger.Debugf("query blocks failed: %s", err)
		//logger.Debugf(query.lang.Assemble(query.args))
		//logger.Debugf(string(resp.Json))
		return 0, nil
	}
	count := result.Array()[0].Get("block-count").Int()
	return count, nil
}

func (query *BlockCountQuery) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func NewBlockQueryByKeyword(ledgers []string, keyword string, from, count int64) *BlockQuery {
	if count > MaxRecordsPerRequest {
		count = MaxRecordsPerRequest
	}
	return &BlockQuery{
		lang: newQueryLang(qlNameBlocksByHash, qlBlockQueryResultName),
		args: map[string]interface{}{
			"keyword": keyword,
			"count":   strconv.FormatInt(count, 10),
			"offset":  strconv.FormatInt(from, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

func NewBlockQueryInRange(ledgers []string, from, to int64) *BlockQuery {
	if from+to > MaxRecordsPerRequest {
		to = MaxRecordsPerRequest - from
	}
	return &BlockQuery{
		lang: newQueryLang(qlNameBlockByHeight, qlBlockQueryResultName),
		args: map[string]interface{}{
			"from":    strconv.FormatInt(from, 10),
			"to":      strconv.FormatInt(to, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

type BlockQuery struct {
	lang *QueryLang
	args map[string]interface{}
}

func (query *BlockQuery) DoQuery(client *dgo.Dgraph) (interface{}, error) {
	return doBlockQuery(query.lang, query.args, client)
}

func (query *BlockQuery) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func doBlockQuery(ql *QueryLang, args map[string]interface{}, client *dgo.Dgraph) (blocks Blocks, e error) {
	ctx := context.Background()
	resp, err := client.NewTxn().Query(ctx, ql.Assemble(args))
	if err != nil {
		logger.Warnf("query blocks failed: %s", err)
		logger.Info(ql.Assemble(args))
		return nil, err
	}

	blocks = parseBlockInfo(gjson.Get(string(resp.Json), ql.resultName))
	if len(blocks) <= 0 {
		//logger.Debugf("query blocks empty")
		//logger.Debugf(ql.Assemble(args))

		blocks = []*BlockInfo{}
	}
	return
}

func parseBlockInfo(result gjson.Result) (blocks Blocks) {
	for _, block := range result.Array() {
		var bi BlockInfo
		bi.HashID = block.Get("block-hash_id").String()
		bi.Height = int(block.Get("block-height").Int())
		bi.Time = block.Get("tx-count").Int()
		blocks = append(blocks, &bi)
	}
	return
}
