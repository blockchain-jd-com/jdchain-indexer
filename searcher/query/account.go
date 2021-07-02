package query

import (
	"context"
	"github.com/dgraph-io/dgo/v200"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

var (
	qlDatasetCountQueryResultName = "query_dataset_count"
	qlNameDatasetCountByHash      = `
        var(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-data_account @filter(regexp(data_account-address, /\S*[[keyword]]\S*/) or regexp(data_account-public_key,    /\S*[[keyword]]\S*/  ) )
            {
                counts as count(data_account-public_key)
            }
        }

		query_dataset_count() {
            data_account-count: sum(val(counts))
	   }
    `
	qlDatasetQueryResultName = "query_dataset"
	qlNameDatasetByHash      = `
        query_dataset(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-data_account @filter(regexp(data_account-address, /\S*[[keyword]]\S*/) or regexp(data_account-public_key,    /\S*[[keyword]]\S*/  ) ) (first:[[count]], offset:[[offset]])
            {
                dataset-address:data_account-address
                dataset-public_key:data_account-public_key
            }
        }
    `
	qlNameDatasetRange = `
        query_dataset(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-data_account (first:[[count]], offset:[[offset]])
            {
                dataset-address:data_account-address
                dataset-public_key:data_account-public_key
            }
        }
    `
)

func NewQueryDataAccountCountByKeyword(ledgers []string, keyword string) *DatasetCountQuery {
	query := newQueryLang(qlNameDatasetCountByHash, qlDatasetCountQueryResultName)
	return &DatasetCountQuery{
		lang: query,
		args: map[string]interface{}{
			"keyword": keyword,
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

type DatasetCountQuery struct {
	lang *QueryLang
	args map[string]interface{}
}

func (query *DatasetCountQuery) DoQuery(client *dgo.Dgraph) (interface{}, error) {
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

	count := result.Array()[0].Get("data_account-count").Int()
	//if count <= 0 {
	//    logger.Debugf(string(resp.Json))
	//}
	return count, nil
}

func (query *DatasetCountQuery) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func NewDatasetRangeQuery(ledgers []string, from, count int64) *DatasetQuery {
	if count > MaxRecordsPerRequest {
		count = MaxRecordsPerRequest
	}
	return &DatasetQuery{
		lang: newQueryLang(qlNameDatasetRange, qlDatasetQueryResultName),
		args: map[string]interface{}{
			"count":   strconv.FormatInt(count, 10),
			"offset":  strconv.FormatInt(from, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

func NewDatasetHasKeyOrHasAddressQuery(ledgers []string, keyword string, from, count int64) *DatasetQuery {
	if count > MaxRecordsPerRequest {
		count = MaxRecordsPerRequest
	}
	return &DatasetQuery{
		lang: newQueryLang(qlNameDatasetByHash, qlDatasetQueryResultName),
		args: map[string]interface{}{
			"keyword": keyword,
			"count":   strconv.FormatInt(count, 10),
			"offset":  strconv.FormatInt(from, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

func NewDatasetByHashAddressQuery(ledgers []string, keyword string) *DatasetQuery {
	return &DatasetQuery{
		lang: newQueryLang(qlNameDatasetByHash, qlDatasetQueryResultName),
		args: map[string]interface{}{
			"keyword": keyword,
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

type DatasetQuery struct {
	lang QueryAssembler
	args map[string]interface{}
}

func (query *DatasetQuery) DoQuery(client *dgo.Dgraph) (interface{}, error) {
	return doDatasetQuery(query.lang, query.args, client)
}

func (query *DatasetQuery) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func doDatasetQuery(ql QueryAssembler, args map[string]interface{}, client *dgo.Dgraph) (accounts Accounts, e error) {
	ctx := context.Background()
	resp, err := client.NewTxn().Query(ctx, ql.Assemble(args))
	if err != nil {
		logger.Warn("While try to query failed: ", err)
		logger.Infof("\n%s\n", ql.Assemble(args))
		return nil, err
	}
	for _, resultName := range ql.ResultNames() {
		accounts = accounts.Append(parseDatasetInfo(gjson.Get(string(resp.Json), resultName))...)
	}
	if len(accounts) <= 0 {
		//logger.Debugf("While try to query empty")
		//logger.Debugf("\n%s\n", ql.Assemble(args))
		accounts = Accounts{}
	}

	return
}

func parseDatasetInfo(result gjson.Result) (sets Accounts) {
	for _, setRaw := range result.Array() {
		var set Account
		set.Address = setRaw.Get("dataset-address").String()
		set.PublicKey = setRaw.Get("dataset-public_key").String()
		sets = append(sets, &set)
	}
	return
}
