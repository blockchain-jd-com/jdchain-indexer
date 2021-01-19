package level_task

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"github.com/RoseRocket/xerrs"
	"github.com/tidwall/gjson"
)

type LevelTaskParserMap map[string]LevelTaskParser

type LevelTaskParser func(uid, ledger string, block int64, level TaskLevel) LevelTask

func NewLevelTaskPuller(parsers LevelTaskParserMap) *LevelTaskPuller {
	return &LevelTaskPuller{
		parsers: parsers,
	}
}

type LevelTaskPuller struct {
	parsers LevelTaskParserMap
}

func (puller *LevelTaskPuller) Pull(dgraphHelper *dgraph_helper.Helper) (lts []LevelTask, e error) {
	queryQL := `
    {
      node(func:gt(level-task-level,0), first:100 ) {
        uid
        level-task-name
		level-task-level
		level-task-ledger
		level-task-block
      }
    }
    `
	raw, err := dgraphHelper.QueryObj(queryQL)
	if err != nil {
		e = xerrs.Mask(fmt.Errorf("do query failed to do query: %s", queryQL), err)
		return
	}

	//fmt.Println(string(pretty.Pretty(raw)))
	node := gjson.ParseBytes(raw).Get("node")
	if node.Exists() == false {
		return
	}
	node.ForEach(func(key, value gjson.Result) bool {
		uid := value.Get("uid").String()
		name := value.Get("level-task-name").String()
		level := value.Get("level-task-level").Int()
		ledger := value.Get("level-task-ledger").String()
		block := value.Get("level-task-block").Int()
		//content := value.Get("level-task-content").String()
		parser, ok := puller.parsers[name]
		if ok == false {
			logger.Warnf("no parser for task[%s]", name)
			return true
		}
		task := parser(uid, ledger, block, TaskLevel(int(level)))
		lts = append(lts, task)
		return true
	})

	return
}
