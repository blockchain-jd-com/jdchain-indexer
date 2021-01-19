package dgraph_helper

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/tidwall/gjson"
	"strings"
)

type Helper struct {
	client *dgo.Dgraph
}

func NewHelper(host string) *Helper {
	c, err := CreateDgClient(host)
	if err != nil {
		logger.Failedf("create dgraph client failed")
		return nil
	}
	return &Helper{
		client: c,
	}
}

func (helper *Helper) DeleteRdfs(rdfs []byte) (e error) {
	mu := &api.Mutation{
		CommitNow: true,
	}

	mu.DelNquads = rdfs
	ctx := context.Background()
	_, err := helper.client.NewTxn().Mutate(ctx, mu)
	if err != nil {
		logger.Warn("While trying to delete failed: ", err)
		e = err
		return
	}
	return
}

func (helper *Helper) MutationRdfs(rdfs []byte) (uids map[string]string, e error) {
	mu := &api.Mutation{
		CommitNow: true,
	}

	mu.SetNquads = rdfs
	ctx := context.Background()
	assigned, err := helper.client.NewTxn().Mutate(ctx, mu)
	if err != nil {
		logger.Warn("While trying to mutate failed: ", err)
		e = err
		return
	}
	uids = assigned.Uids
	return
}

func (helper *Helper) QueryNode(qb Queryable) (json string, e error) {
	name, value := qb.QueryBy()
	query := fmt.Sprintf(`{
        query_node(func: eq(%s, "%s")) {
            uid
        }
    }`, name, value)

	raw, err := helper.QueryObj(query)
	if err != nil {
		e = err
		return
	}
	json = gjson.Get(string(raw), "query_node").String()
	return
}

type UidQueryable interface {
	UniqueMutationName() string
	QueryBy() (string, string)
}

func (helper *Helper) QueryUids(uidQueryable UidQueryable) (uids map[string]string, e error) {
	tmpl := `
    {
        query_nodes(func: has(%s)) {
            uid
            %s
        }
    }
    `
	predict, _ := uidQueryable.QueryBy()
	query := fmt.Sprintf(tmpl, predict, predict)
	resultRaw, err := helper.QueryObj(query)
	if err != nil {
		e = err
		return
	}
	uids = make(map[string]string)
	list := gjson.Parse(string(resultRaw)).Get("query_nodes").Array()
	if len(list) <= 0 {
		return
	}
	for _, obj := range list {
		uids[uidQueryable.UniqueMutationName()] = obj.Get("uid").String()
	}
	return
}

// QueryObj do query from dgraph
func (helper *Helper) QueryObj(query string) (json []byte, e error) {
	ctx := context.Background()
	resp, err := helper.client.NewTxn().Query(ctx, query)
	if err != nil {
		logger.Warn("While try to query failed: ", err)
		logger.Warn(strings.Repeat("-", 128))
		logger.Info(query)
		logger.Warn(strings.Repeat("-", 128))
		e = err
		return
	}
	json = resp.Json
	return
}

// QueryObjWithVars do query from dgraph with paras
func (helper *Helper) QueryObjWithVars(query string, variables map[string]string) (json []byte, e error) {
	ctx := context.Background()
	resp, err := helper.client.NewTxn().QueryWithVars(ctx, query, variables)
	if err != nil {
		logger.Warn("While try to query failed: ", err)
		logger.Warn(strings.Repeat("-", 128))
		logger.Info(query)
		logger.Warn(strings.Repeat("-", 128))
		e = err
		return
	}
	json = resp.Json
	return
}

func (helper *Helper) QueryUID(predict, value string) (uid string, exists bool, e error) {
	query := fmt.Sprintf(`{
        query_uid(func: eq(%s, "%s")) {
            uid
        }
    }`, predict, value)
	ctx := context.Background()
	resp, err := helper.client.NewTxn().Query(ctx, query)
	if err != nil {
		logger.Warn("While try to query failed: ", err)
		spew.Dump(query)
		e = err
		return
	}

	result := gjson.Get(string(resp.Json), "query_uid.0.uid")
	if result.Exists() {
		uid = result.String()
		exists = true
	}
	return
}

func (helper *Helper) Alter(schemes Schemas) (e error) {
	op := &api.Operation{}
	op.Schema = schemes.String()
	ctx := context.Background()
	err := helper.client.Alter(ctx, op)
	if err != nil {
		logger.Warn("While try to alter scheme failed: ", err)
		e = err
		return
	}
	return
}

func (helper *Helper) DropDB() error {
	op := api.Operation{
		DropAll: true,
	}
	ctx := context.Background()
	if err := helper.client.Alter(ctx, &op); err != nil {
		logger.Failed(err)
		return err
	}

	logger.Success("dgraph drop success")
	return nil
}
