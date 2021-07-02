package query

import (
	"context"
	"github.com/dgraph-io/dgo/v200"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

var (
	qlUserCountQueryResultName = "query_users_count"
	qlNameUsersCountByHash     = `
      var(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-user @filter(regexp(user-address, /\S*[[keyword]]\S*/) or regexp(user-public_key,    /\S*[[keyword]]\S*/  ))
            {
                counts as count(user-public_key)
            }
        }

	query_users_count() {
            user-count: sum(val(counts))
	   }
    `

	qlUserQueryResultName = "query_users"
	qlNameUsersByHash     = `
      query_users(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-user @filter(regexp(user-address, /\S*[[keyword]]\S*/) or regexp(user-public_key,    /\S*[[keyword]]\S*/  )) (first:[[count]], offset:[[offset]])
            {
                user-address:user-address
                user-public_key:user-public_key
            }
        }
    `
	qlNameUsersRange = `
      query_users(func: has(ledger-hash_id))
            @filter(anyofterms(ledger-hash_id, "[[ledgers]]")) @normalize
        {
            ledger-user (first:[[count]], offset:[[offset]])
            {
                user-address:user-address
                user-public_key:user-public_key
            }
        }
    `
)

func NewQueryUserCountByHash(ledgers []string, keyword string) *UserCountQuery {
	query := newQueryLang(qlNameUsersCountByHash, qlUserCountQueryResultName)
	return &UserCountQuery{
		lang: query,
		args: map[string]interface{}{
			"keyword": keyword,
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

type UserCountQuery UserQuery

func (query *UserCountQuery) DoQuery(client *dgo.Dgraph) (interface{}, error) {
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

	count := result.Array()[0].Get("user-count").Int()
	return count, nil
}

func (query *UserCountQuery) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func NewQueryUsersByHash(ledgers []string, keyword string, from, count int64) *UserQuery {
	if count > MaxRecordsPerRequest {
		count = MaxRecordsPerRequest
	}
	query := newQueryLang(qlNameUsersByHash, qlUserQueryResultName)
	return &UserQuery{
		lang: query,
		args: map[string]interface{}{
			"keyword": keyword,
			"count":   strconv.FormatInt(count, 10),
			"offset":  strconv.FormatInt(from, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

func NewQueryUsersRange(ledgers []string, from, count int64) *UserQuery {
	if count > MaxRecordsPerRequest {
		count = MaxRecordsPerRequest
	}
	query := newQueryLang(qlNameUsersRange, qlUserQueryResultName)
	return &UserQuery{
		lang: query,
		args: map[string]interface{}{
			"count":   strconv.FormatInt(count, 10),
			"offset":  strconv.FormatInt(from, 10),
			"ledgers": strings.Join(ledgers, " "),
		},
	}
}

type UserQuery struct {
	lang *QueryLang
	args map[string]interface{}
}

func (query *UserQuery) DoQuery(client *dgo.Dgraph) (interface{}, error) {
	return doUsersQuery(query.lang, query.args, client)
}

func (query *UserQuery) OutputDebugInfo() interface{} {
	return query.lang.AssembleForRead(query.args)
}

func doUsersQuery(ql *QueryLang, args map[string]interface{}, client *dgo.Dgraph) (users Users, e error) {
	ctx := context.Background()
	resp, err := client.NewTxn().Query(ctx, ql.Assemble(args))
	if err != nil {
		logger.Warn("While try to query failed: ", err)
		logger.Infof("\n%s\n", ql.Assemble(args))
		return nil, err
	}

	users = parseUserInfo(gjson.Get(string(resp.Json), ql.resultName))
	if len(users) <= 0 {
		users = Users{}
	}
	return
}

func parseUserInfo(result gjson.Result) (users Users) {
	for _, userRaw := range result.Array() {
		var user User
		user.Address = userRaw.Get("user-address").String()
		user.PublicKey = userRaw.Get("user-public_key").String()
		users = append(users, &user)
	}
	return
}
