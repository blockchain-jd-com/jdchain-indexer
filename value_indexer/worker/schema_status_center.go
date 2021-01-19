package worker

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/value_indexer/schema"
	"github.com/tidwall/gjson"
	"sync"
)

var (
	errNotPrepared = fmt.Errorf("schema sync not prepared")
)

type DataSync interface {
	PushUpdate(string) error
	PushDelete(string) error
	Pull() (string, error)
	SpecifiedSchemaUID(value string) (id string, e error)
	AlterSchema(schemas dgraph_helper.Schemas) error
}

func NewSchemaStatusCenter(dataSync DataSync) *SchemaStatusCenter {
	return &SchemaStatusCenter{
		locker:   sync.Mutex{},
		dataSync: dataSync,
	}
}

type SchemaStatusCenter struct {
	IndexStatusList SchemaIndexStatusList
	locker          sync.Mutex
	prepared        bool
	dataSync        DataSync
}

func (center *SchemaStatusCenter) OnEvent(e IndexEvent) {
	if e.Status == SchemaStatusCleared {
		logger.Infof("schema %s is cleared, will be deleted", e.schemaStatus.Schema.ID)
		err := center.delete(e.schemaStatus.Schema.ID)
		if err != nil {
			logger.Errorf("delete schema %s from db failed: %s", e.schemaStatus.Schema.ID, err)
		}
	}
}

func (center *SchemaStatusCenter) PullSync() error {
	indexStatusList, err := center.GetSchemaStatus()

	if err != nil {
		return err
	}

	logger.Infof("%d status updated", indexStatusList.Len())
	center.IndexStatusList = indexStatusList
	return nil
}

func (center *SchemaStatusCenter) GetSchemaStatus() (list SchemaIndexStatusList, e error) {
	s, err := center.dataSync.Pull()
	if err != nil {
		e = err
		return
	}

	result := gjson.Parse(s)
	nodesResult := result.Get("node")
	if nodesResult.Exists() {
		for _, node := range nodesResult.Array() {
			id := node.Get("schemainfo-id").String()
			if len(id) <= 0 {
				continue
			}
			uid := node.Get("uid").String()
			account := node.Get("schemainfo-associate_account").String()
			ledger := node.Get("schemainfo-ledger").String()
			status := node.Get("schemainfo-status").Int()
			content := node.Get("schemainfo-content").String()
			progress := node.Get("schemainfo-progress").Int()
			si, err := NewSchemaInfo(id, ledger, account, content)
			if err != nil {
				logger.Errorf("cannot sync schema info with db, schema invalid: %s", content)
				e = err
				return
			}
			list = list.Append(NewSchemaIndexStatus(uid, si, SchemaStatus(status), progress))
		}
	} else {
		logger.Warnf("no node found from db: \n%s", s)
	}
	return
}

func (center *SchemaStatusCenter) Prepare() error {
	center.locker.Lock()
	defer center.locker.Unlock()

	err := center.PullSync()
	if err != nil {
		return err
	}
	center.prepared = true
	return nil
}

func (center *SchemaStatusCenter) Find(name string) *schema.NodeSchema {
	statusList := center.IndexStatusList.Filter(func(status SchemaIndexStatus) bool {
		return name == status.Schema.nodeSchema.LowerName()
	})
	if statusList.Len() <= 0 {
		logger.Errorf("no such Schema named %s", name)
		return nil

	}
	return statusList[0].Schema.nodeSchema
}

func (center *SchemaStatusCenter) Stop(id string) (e error) {
	center.locker.Lock()
	defer center.locker.Unlock()

	status, err := center.isValid(id)
	if err != nil {
		e = err
		return
	}

	stoppedStatus := NewSchemaIndexStatusStopped(status.Schema, status.uid)
	data := stoppedStatus.UpdateStatusMutations().Assembly()
	logger.Info("update node status to stopped")
	fmt.Println(data)
	if err := center.dataSync.PushUpdate(data); err != nil {
		e = err
		return
	}

	if err := center.PullSync(); err != nil {
		e = err
		return
	}

	return
}

func (center *SchemaStatusCenter) Start(id string) (e error) {
	center.locker.Lock()
	defer center.locker.Unlock()

	status, err := center.isValid(id)
	if err != nil {
		e = err
		return
	}

	runningStatus := NewSchemaIndexStatusRunning(status.Schema, status.uid)
	data := runningStatus.UpdateStatusMutations().Assembly()
	logger.Info("update status to running: ")
	fmt.Println(data)
	if err := center.dataSync.PushUpdate(data); err != nil {
		e = err
		return
	}

	if err := center.PullSync(); err != nil {
		e = err
		return
	}

	return
}

func (center *SchemaStatusCenter) isValid(id string) (status SchemaIndexStatus, e error) {
	if center.prepared == false {
		e = errNotPrepared
		return
	}

	statusList := center.IndexStatusList.Filter(func(status SchemaIndexStatus) bool {
		return id == status.Schema.ID
	})
	if statusList.Len() <= 0 {
		e = fmt.Errorf("no such Schema")
		return

	}
	return statusList[0], nil
}

func (center *SchemaStatusCenter) Delete(id string) error {
	center.locker.Lock()
	defer center.locker.Unlock()

	status, err := center.isValid(id)
	if err != nil {
		return err
	}
	if status.IsRunning() {
		return fmt.Errorf("schema is running, stop first")
	}

	waitForClearStatus := NewSchemaIndexStatusClearing(status.Schema, status.uid)
	data := waitForClearStatus.UpdateStatusMutations().Assembly()
	logger.Info("update status to clearing: ")
	logger.Info(data)

	if err := center.dataSync.PushUpdate(data); err != nil {
		return err
	}

	if err := center.PullSync(); err != nil {
		return err
	}

	return nil
}

func (center *SchemaStatusCenter) delete(id string) error {
	center.locker.Lock()
	defer center.locker.Unlock()

	status, err := center.isValid(id)
	if err != nil {
		return err
	}
	if status.IsRunning() {
		return fmt.Errorf("schema is running, stop first")
	}

	data := status.DeleteMutations().Assembly()
	logger.Info("delete status node:")
	logger.Info(data)

	if err := center.dataSync.PushDelete(data); err != nil {
		return err
	}

	if err := center.PullSync(); err != nil {
		return err
	}

	return nil
}

func (center *SchemaStatusCenter) Add(info SchemaInfo) error {
	center.locker.Lock()
	defer center.locker.Unlock()
	if center.prepared == false {
		return errNotPrepared
	}

	id := info.ID
	exists := center.IndexStatusList.Any(func(value SchemaIndexStatus) bool {
		return id == value.Schema.ID
	})
	if exists {
		return fmt.Errorf("already exists")
	}

	logger.Info(info.String())

	metaSchemaBuilder := schema.NewSchemaMetaBuilder(info.nodeSchema)

	err := center.dataSync.AlterSchema(metaSchemaBuilder.Build())
	if err != nil {
		logger.Errorf("alter schema for %s failed: %s \n %s", info.ID, err, metaSchemaBuilder.Build().String())
		return fmt.Errorf("cannot create schema for invalid format")
	}

	logger.Infof("add schema for schema %s: ", info.nodeSchema.LowerName())
	logger.Info(metaSchemaBuilder.Build().String())

	status := NewSchemaIndexStatusDefault(info, "")

	data := status.CreateMutations().Assembly()
	logger.Info("create status node: ")
	fmt.Println(data)

	if err := center.dataSync.PushUpdate(data); err != nil {
		logger.Errorf("commit status failed: %s", err)
		return fmt.Errorf("commit schema failed")
	}

	if err := center.PullSync(); err != nil {
		return err
	}

	return nil
}

func (center *SchemaStatusCenter) Update(info SchemaInfo) error {
	center.locker.Lock()
	defer center.locker.Unlock()
	if center.prepared == false {
		return errNotPrepared
	}

	status, err := center.isValid(info.ID)
	if err != nil {
		return err
	}
	if status.IsRunning() {
		return fmt.Errorf("schema is running, stop first")
	}

	logger.Info(info.String())

	metaSchemaBuilder := schema.NewSchemaMetaBuilder(info.nodeSchema)

	err = center.dataSync.AlterSchema(metaSchemaBuilder.Build())
	if err != nil {
		logger.Errorf("alter schema for %s failed: %s \n %s", info.ID, err, metaSchemaBuilder.Build().String())
		return fmt.Errorf("cannot create schema for invalid format")
	}

	status.Schema.Content = info.Content
	data := status.UpdateContentMutations().Assembly()
	logger.Info("update node content")
	fmt.Println(data)
	if err = center.dataSync.PushUpdate(data); err != nil {
		return err
	}

	if err := center.PullSync(); err != nil {
		return err
	}

	return nil
}
