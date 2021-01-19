package meta_level_task

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/adaptor"
	"git.jd.com/jd-blockchain/explorer/chain"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/event"
	"git.jd.com/jd-blockchain/explorer/level_task"
	"github.com/davecgh/go-spew/spew"
	"github.com/tidwall/gjson"
)

func NewMetaInfoLevelTaskCreatorManager(apiHost string, dgraphHelper *dgraph_helper.Helper) *MetaInfoLevelTaskCreatorManager {
	return &MetaInfoLevelTaskCreatorManager{
		dgraphHelper: dgraphHelper,
		apiHost:      apiHost,
	}
}

type MetaInfoLevelTaskCreatorManager struct {
	dgraphHelper *dgraph_helper.Helper
	creators     []*MetaInfoLevelTaskCreator
	apiHost      string
}

func (manager *MetaInfoLevelTaskCreatorManager) ID() string {
	return "MetaInfoLevelTaskCreators"
}

func (manager *MetaInfoLevelTaskCreatorManager) EventReceived(e event.Event) bool {
	data := e.GetData()
	if data == nil {
		return true
	}
	switch t := data.(type) {
	case *chain.LedgerStatus:
		manager.updateTaskCreator(t.Ledger, t.Height)
	default:
		logger.Warnf("no handler in MetaInfoLevelTaskCreators for: \n%s", spew.Sdump(data))
	}
	return true
}

func (manager *MetaInfoLevelTaskCreatorManager) updateTaskCreator(ledger string, height int64) {
	creator := manager.findCreator(ledger)
	if creator != nil {
		err := creator.update(height, manager.dgraphHelper)
		if err != nil {
			logger.Errorf("creator update failed: %s", err)
			return
		}
		return
	}

	_, _, err := PrepareLedgerNode(manager.dgraphHelper, ledger)
	if err != nil {
		logger.Warnf("prepare ledger [%s] failed: %s", ledger, err)
		return
	}

	creator = NewMetaInfoLevelTaskCreator(ledger)

	manager.creators = append(manager.creators, creator)
	err = creator.update(height, manager.dgraphHelper)
	if err != nil {
		logger.Errorf("creator update failed: %s", err)
		return
	}
	startTaskMonitor(manager.apiHost, ledger, manager.dgraphHelper)
}

func startTaskMonitor(apiHost, ledger string, dgraphHelper *dgraph_helper.Helper) {
	handler := NewLevelTaskHandler(apiHost, ledger, dgraphHelper, func(s string) error {
		_, err := dgraphHelper.MutationRdfs([]byte(s))
		return err
	}).Setup()

	parsers := level_task.LevelTaskParserMap{
		TaskMetaTaskLevelName: ParseMetaInfoLevelTask,
	}

	level_task.NewLevelTaskMonitor(handler, parsers, dgraphHelper).Setup()
}

func (manager *MetaInfoLevelTaskCreatorManager) findCreator(ledger string) *MetaInfoLevelTaskCreator {
	for _, creator := range manager.creators {
		if creator.ledger == ledger {
			return creator
		}
	}
	return nil
}

func PrepareLedgerNode(helper *dgraph_helper.Helper, ledger string) (ledgerNode *adaptor.Ledger, alreadyExists bool, e error) {
	ledgerExists, node, err := isLedgerNodeExists(helper, ledger)
	if err != nil {
		e = err
		return
	}

	if ledgerExists {
		ledgerNode = node
		alreadyExists = true
		logger.Warnf("detect ledger(%s) node already exists ", ledger)
		return
	} else {
		logger.Infof("detect ledger(%s) node NOT exists, need to create", ledger)
	}

	ledgerNode, e = initLedgerNode(helper, ledger)
	return
}

func initLedgerNode(helper *dgraph_helper.Helper, ledger string) (node *adaptor.Ledger, e error) {
	node = adaptor.NewLedger(ledger)
	uids, err := helper.MutationRdfs([]byte(node.Mutations().Assembly()))
	if err != nil {
		e = err
		logger.Failedf("mutate ledger failed: %s", err)
		return
	}
	uid, ok := uids[string(adaptor.ModelTypeLedger)]
	if ok == false {
		spew.Dump(uids)
		return nil, fmt.Errorf("create ledger[%s] but node uid return", ledger)
	}
	node.Uid = uid
	logger.Infof("ledger[%s] node[%s] created success", ledger, uid)
	return
}

func isLedgerNodeExists(helper *dgraph_helper.Helper, ledger string) (ok bool, node *adaptor.Ledger, e error) {
	newLedger := adaptor.NewLedger(ledger)
	result, err := helper.QueryNode(newLedger)
	if err != nil {
		e = err
		return
	}
	nodes := gjson.Parse(result).Array()
	if len(nodes) <= 0 {
		return
	}
	uid := nodes[0].Get("uid").String()
	if len(uid) > 0 {
		ok = true
		newLedger.Uid = uid
		node = newLedger
	}
	return
}
