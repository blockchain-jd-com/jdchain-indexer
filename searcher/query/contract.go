package query

import (
	"context"
	"github.com/dgraph-io/dgo/v200"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

var (
	qlContractCountQueryResultName = "query_contract_count"
	qlNameContractCountByHash      = `
        var(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-contract @filter(regexp(contract-address, /\S*[[keyword]]\S*/) or regexp(contract-public_key,    /\S*[[keyword]]\S*/  ) )
            {
                counts as count(contract-public_key)
            }
        }
		query_contract_count() {
            contract-count: sum(val(counts))
	   }
    `
	qlContractQueryResultName = "query_contract"
	qlNameContractByHash      = `
        query_contract(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-contract @filter(regexp(contract-address, /\S*[[keyword]]\S*/) or regexp(contract-public_key,    /\S*[[keyword]]\S*/  ) ) (first:[[count]], offset:[[offset]])
            {
                contract-address:contract-address
                contract-public_key:contract-public_key
            }
        }
    `
	qlNameContractRange = `
        query_contract(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-contract (first:[[count]], offset:[[offset]])
            {
                contract-address:contract-address
                contract-public_key:contract-public_key
            }
        }
    `
)

func NewQueryContractCountByHash(ledgers []string, keyword string) *QueryContractCount {
	query := newQueryLang(qlNameContractCountByHash, qlContractCountQueryResultName)
	return &QueryContractCount{
		lang: query,
		args: map[string]interface{}{
			"keyword": keyword,
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

type QueryContractCount QueryContract

func (query *QueryContractCount) DoQuery(client *dgo.Dgraph) (interface{}, error) {
	ctx := context.Background()
	resp, err := client.NewTxn().Query(ctx, query.lang.Assemble(query.args))
	if err != nil {
		logger.Warn("While try to query failed: ", err)
		logger.Infof("\n%s\n", query.lang.Assemble(query.args))
		return nil, err
	}

	result := gjson.Get(string(resp.Json), query.lang.resultName)
	if len(result.Array()) <= 0 {
		return 0, nil
	}

	count := result.Array()[0].Get("contract-count").Int()
	return count, nil
}

func (query *QueryContractCount) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

type QueryContract struct {
	lang *QueryLang
	args map[string]interface{}
}

func (query *QueryContract) DoQuery(client *dgo.Dgraph) (interface{}, error) {
	return doContractQuery(query.lang, query.args, client)
}

func (query *QueryContract) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func NewQueryContractRange(ledgers []string, from, count int64) *QueryContract {
	if count > MaxRecordsPerRequest {
		count = MaxRecordsPerRequest
	}
	query := newQueryLang(qlNameContractRange, qlContractQueryResultName)
	return &QueryContract{
		lang: query,
		args: map[string]interface{}{
			"count":   strconv.FormatInt(count, 10),
			"offset":  strconv.FormatInt(from, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

func NewQueryContractByHash(ledgers []string, keyword string, from, count int64) *QueryContract {
	if count > MaxRecordsPerRequest {
		count = MaxRecordsPerRequest
	}
	query := newQueryLang(qlNameContractByHash, qlContractQueryResultName)
	return &QueryContract{
		lang: query,
		args: map[string]interface{}{
			"keyword": keyword,
			"count":   strconv.FormatInt(count, 10),
			"offset":  strconv.FormatInt(from, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

func doContractQuery(ql *QueryLang, args map[string]interface{}, client *dgo.Dgraph) (contracts Contracts, e error) {
	ctx := context.Background()
	resp, err := client.NewTxn().Query(ctx, ql.Assemble(args))
	if err != nil {
		logger.WithJsonRaw([]byte(ql.Assemble(args))).Warn("While try to query failed: ", err)
		return nil, err
	}
	contracts = parseContractInfo(gjson.Get(string(resp.Json), ql.resultName))
	if len(contracts) <= 0 {
		contracts = Contracts{}
	}
	return
}

func parseContractInfo(result gjson.Result) (contracts Contracts) {
	for _, raw := range result.Array() {
		var contract Contract
		contract.Address = AddressValue{raw.Get("contract-address").String()}
		contract.PublicKey = raw.Get("contract-public_key").String()
		contracts = append(contracts, contract)
	}
	return
}
