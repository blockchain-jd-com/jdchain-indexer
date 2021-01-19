package meta_level_task

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/level_task"
	"github.com/tidwall/gjson"
	"strings"
)

var (
	DefaultMaxTaskCountEachUpdate = 300
)

func NewMetaInfoLevelTaskCreator(ledger string) *MetaInfoLevelTaskCreator {
	creator := &MetaInfoLevelTaskCreator{
		ledger:                 ledger,
		maxTaskCountEachUpdate: DefaultMaxTaskCountEachUpdate,
		cacheHeight:            -1,
	}
	return creator
}

type MetaInfoLevelTaskCreator struct {
	ledger                 string
	maxTaskCountEachUpdate int
	cacheHeight            int64
}

func (creator *MetaInfoLevelTaskCreator) RealtimeTaskStatus(dgraphHelper *dgraph_helper.Helper) (block int64, e error) {

	queryQL := fmt.Sprintf(`
    {
      nodes(func: eq(level-task-ledger,"%s"), orderdesc: level-task-block, first:1) {
        uid
        level-task-name
		level-task-level
		level-task-ledger
		level-task-block
      }
    }
    `, creator.ledger)
	raw, err := dgraphHelper.QueryObj(queryQL)
	if err != nil {
		e = err
		return
	}
	nodesResult := gjson.ParseBytes(raw).Get("nodes")
	if nodesResult.Exists() == false {
		logger.Warnf("query failed: %s", queryQL)
		e = fmt.Errorf("setup failed for invalid query response")
		return
	}
	nodesArray := nodesResult.Array()
	if len(nodesArray) > 0 {
		block = nodesArray[0].Get("level-task-block").Int()
		logger.Infof("metainfo task for ledger[%s] already updated -> %d", creator.ledger, block)
	} else {
		block = -1
	}
	return
}

func (creator *MetaInfoLevelTaskCreator) update(block int64, dgraphHelper *dgraph_helper.Helper) error {
	if creator.cacheHeight >= block {
		return nil
	}

	current, err := creator.RealtimeTaskStatus(dgraphHelper)
	if err != nil {
		return err
	}
	creator.cacheHeight = current
	if current >= block {
		return nil
	}

	var lts []level_task.LevelTask
	var count int
	for i := current + 1; i <= block; i++ {
		lt := CreateNewMetaInfoLevelTask("", creator.ledger, i, level_task.MetaInfoLevel1)
		lts = append(lts, lt)
		count++
		if count >= creator.maxTaskCountEachUpdate {
			break
		}
	}

	var builder strings.Builder
	for _, task := range lts {
		builder.WriteString(task.CreateMutations().Assembly())
	}

	//fmt.Println(builder.String())
	_, err = dgraphHelper.MutationRdfs([]byte(builder.String()))
	if err != nil {
		return err
	}

	return nil
}
