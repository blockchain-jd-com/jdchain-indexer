package level_task

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

var (
	LevelTaskSchemas = dgraph_helper.Schemas{
		dgraph_helper.NewSchemaIntIndex("level-task-level"),
		dgraph_helper.NewSchemaStringTermIndex("level-task-name"),
		dgraph_helper.NewSchemaStringTermIndex("level-task-ledger"),
		dgraph_helper.NewSchemaIntIndex("level-task-block"),
		dgraph_helper.NewSchemaString("level-task-content"),
	}
)

type LevelTask interface {
	CreateMutations() (mutations dgraph_helper.Mutations)
}
