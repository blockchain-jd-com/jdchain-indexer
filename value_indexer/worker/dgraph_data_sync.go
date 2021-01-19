package worker

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"github.com/RoseRocket/xerrs"
	"github.com/tidwall/gjson"
)

var (
	getUIDsWithPredictSrc = `
    {
      node(func: has(%s)) {
            uid
        }
    }
    `
	srcQuerySpecifiedSchema = `
    {
      node(func: eq(schemainfo-id, "%s")) {
        uid
      }
    }
    `

	srcQuerySchemas = `
    {
      node(func: has(schemainfo-id)) {
        uid
		schemainfo-id
        schemainfo-associate_account
		schemainfo-ledger
		schemainfo-status
		schemainfo-content
		schemainfo-progress
      }
    }
    `
)

func NewDgraphDataSync(dgraphHelper *dgraph_helper.Helper) *DgraphDataSync {
	return &DgraphDataSync{
		helper: dgraphHelper,
	}
}

type DgraphDataSync struct {
	helper *dgraph_helper.Helper
}

func (sync *DgraphDataSync) PushDelete(data string) (e error) {
	return sync.helper.DeleteRdfs([]byte(data))
}

func (sync *DgraphDataSync) PushUpdate(data string) (e error) {
	if len(data) <= 0 {
		return nil
	}
	_, err := sync.helper.MutationRdfs([]byte(data))
	return err
}

func (sync *DgraphDataSync) SpecifiedSchemaUID(value string) (id string, e error) {
	query := fmt.Sprintf(srcQuerySpecifiedSchema, value)
	raw, err := sync.helper.QueryObj(query)
	if err != nil {
		return "", err
	}
	result := gjson.ParseBytes(raw)
	nodesResult := result.Get("node")
	if nodesResult.Exists() {
		for _, node := range nodesResult.Array() {
			uid := node.Get("uid").String()
			if len(uid) > 0 {
				id = uid
				return
			}
		}
	} else {
		logger.Warnf("no node found from db: \n%s", raw)
	}
	return
}

func (sync *DgraphDataSync) UIDs(predict string) (ids []string, e error) {
	query := fmt.Sprintf(getUIDsWithPredictSrc, predict)
	logger.Infof("UIDs -> \n %s", query)
	raw, err := sync.helper.QueryObj(query)
	if err != nil {
		return nil, err
	}
	result := gjson.ParseBytes(raw)
	nodesResult := result.Get("node")
	if nodesResult.Exists() {
		for _, node := range nodesResult.Array() {
			uid := node.Get("uid").String()
			if len(uid) > 0 {
				ids = append(ids, uid)
			}
		}
	} else {
		logger.Warnf("no node found from db: \n%s", raw)
	}
	return
}

func (sync *DgraphDataSync) AlterSchema(schemas dgraph_helper.Schemas) error {
	if e := sync.helper.Alter(schemas); e != nil {
		logger.Errorf("alter schema error: %s", e)
		return xerrs.Mask(e, fmt.Errorf("alter schema failed"))
	}
	return nil
}

func (sync *DgraphDataSync) Pull() (string, error) {
	bs, err := sync.helper.QueryObj(srcQuerySchemas)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
