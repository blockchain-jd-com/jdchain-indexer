package meta_level_task

import (
	"crypto/md5"
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/level_task"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/links"
	"git.jd.com/jd-blockchain/explorer/worker"
	"github.com/RoseRocket/xerrs"
	"strconv"
	"strings"
)

const (
	TaskMetaTaskLevelName = "meta-level-task"
)

func ParseMetaInfoLevelTask(uid, ledger string, block int64, level level_task.TaskLevel) level_task.LevelTask {
	return CreateNewMetaInfoLevelTask(uid, ledger, block, level)
}

func CreateNewMetaInfoLevelTask(uid, ledger string, block int64, level level_task.TaskLevel) *MetaInfoLevelTask {
	lt := &MetaInfoLevelTask{
		uid:    uid,
		Ledger: ledger,
		Block:  block,
		level:  level,
	}
	data := md5.Sum([]byte(lt.SerialString()))
	lt.contentHash = fmt.Sprintf("%x", data)
	return lt
}

type MetaInfoLevelTask struct {
	uid         string
	Ledger      string
	Block       int64
	level       level_task.TaskLevel
	contentHash string
}

func (lt *MetaInfoLevelTask) UpdateMutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(lt.uid),
			dgraph_helper.MutationItemValue(lt.level.Upgrade().String()),
			dgraph_helper.MutationPredict("level-task-level"),
		),
	)
	return
}

func (lt *MetaInfoLevelTask) CreateMutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(lt.contentHash),
			dgraph_helper.MutationItemValue(lt.level.String()),
			dgraph_helper.MutationPredict("level-task-level"),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(lt.contentHash),
			dgraph_helper.MutationItemValue(TaskMetaTaskLevelName),
			dgraph_helper.MutationPredict("level-task-name"),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(lt.contentHash),
			dgraph_helper.MutationItemValue(lt.Ledger),
			dgraph_helper.MutationPredict("level-task-ledger"),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(lt.contentHash),
			dgraph_helper.MutationItemValue(strconv.FormatInt(lt.Block, 10)),
			dgraph_helper.MutationPredict("level-task-block"),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(lt.contentHash),
			dgraph_helper.MutationItemValue(lt.SerialString()),
			dgraph_helper.MutationPredict("level-task-content"),
		),
	)
	return
}

func (lt *MetaInfoLevelTask) SerialString() string {
	return fmt.Sprintf("%s|%s|%d", TaskMetaTaskLevelName, lt.Ledger, lt.Block)
}

func newMetaInfoLevelTaskHandler(task *MetaInfoLevelTask, host string, cache *dgraph_helper.UidLruCache, rdfSaver func(string) error) *MetaInfoLevelTaskHandler {
	return &MetaInfoLevelTaskHandler{
		rdfSaver:   rdfSaver,
		cache:      cache,
		task:       task,
		ledgerHost: host,
	}
}

type MetaInfoLevelTaskHandler struct {
	cache      *dgraph_helper.UidLruCache
	rdfSaver   func(string) error
	task       *MetaInfoLevelTask
	ledgerHost string
	datas      []dgraph_helper.MutationData
}

func (handler *MetaInfoLevelTaskHandler) Do() error {
	task := handler.task
	fetchTask := worker.NewFetchTask(handler.ledgerHost, task.Ledger, task.Block)
	err := fetchTask.Do()
	if err != nil {
		return xerrs.Mask(fmt.Errorf("fetch data failed"), err)
	}

	ld := fetchTask.Data().(*worker.LedgerData)
	var rdfs string
	switch task.level {
	case level_task.MetaInfoLevel1:
		var builder strings.Builder
		cns := links.ToCommonNode(ld)
		for _, cn := range cns {
			rdfs, err := dgraph_helper.AssembleMutationDatas(handler.cache, cn)
			if err != nil {
				logger.Errorf("assemble failed: %s", err)
				break
			}
			builder.WriteString(rdfs)
		}
		rdfs = builder.String()
	case level_task.MetaInfoLevel2:
		modelLinks := links.ToLinks(ld)
		rdfs, err = dgraph_helper.AssembleMutationDatas(handler.cache, modelLinks...)
		if err != nil {
			logger.Errorf("assemble failed: %s", err)
			return err
		}
	}

	rdfs = fmt.Sprintf("%s\n%s", rdfs, handler.task.UpdateMutations().Assembly())
	fmt.Println(rdfs)
	if err := handler.rdfSaver(rdfs); err != nil {
		return err
	}
	return nil
}
