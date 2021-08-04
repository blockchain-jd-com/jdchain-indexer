package worker

import (
	"encoding/base64"
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/value_indexer/schema"
	"github.com/tidwall/gjson"
	"io"
	"strconv"
	"time"
)

type Status int
type Command int

const (
	CommandClear   Command = 1
	CommandIndex   Command = 2
	CommandStop    Command = 3 //停止后可以重新开始
	CommandCleared Command = 4

	StatusDefault  Status = 0
	StatusIndexing Status = 1
	StatusStopped  Status = 2
	StatusClearing Status = 3
	StatusCleared  Status = 4
)

func (status Status) String() string {
	switch status {
	case StatusDefault:
		return "WorkerStatusDefault"
	case StatusIndexing:
		return "WorkerStatusIndexing"
	case StatusStopped:
		return "WorkerStatusStopped"
	case StatusClearing:
		return "WorkerStatusClearing"
	case StatusCleared:
		return "WorkerStatusCleared"

	}
	return ""
}

//ValueIndexWorker pull data from ledger and commit to db
type ValueIndexWorker struct {
	id                string
	dataSource        DataSource
	dataUpdater       DataUpdater
	currentStatus     Status
	nodeSchema        *schema.NodeSchema
	schemaStatus      SchemaIndexStatus
	builder           KVRDFBuilder
	chExternalCommand chan Command
	chInternalCommand chan Command
}

type DataUpdater interface {
	PushDelete(string) error
	PushUpdate(string) error
	UIDs(string) ([]string, error)
}

type DataSource interface {
	Read() (string, int, error)
	Stop()
}

type KVRDFBuilder interface {
	Build(src string) dgraph_helper.Mutations
}

func NewValueIndexWorker(id string, status SchemaIndexStatus, schema *schema.NodeSchema, builder KVRDFBuilder, updater DataUpdater) *ValueIndexWorker {
	worker := &ValueIndexWorker{
		dataUpdater:       updater,
		chExternalCommand: make(chan Command, 1024),
		chInternalCommand: make(chan Command, 1024),
	}

	worker.schemaStatus = status
	worker.nodeSchema = schema
	worker.builder = builder
	worker.id = id
	go worker.Run()
	return worker
}

func (worker *ValueIndexWorker) Stop() {
	worker.chExternalCommand <- CommandStop
}

func (worker *ValueIndexWorker) Clear() {
	//worker.schemaStatus = status
	worker.chExternalCommand <- CommandClear
	if worker.dataSource != nil {
		worker.dataSource.Stop()
		worker.dataSource = nil
	}
}

func (worker *ValueIndexWorker) Start(source DataSource) {
	if worker.dataSource == nil {
		worker.dataSource = source
	}
	worker.chExternalCommand <- CommandIndex
}

func (worker *ValueIndexWorker) startClearLoop() {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		//fmt.Println("try clear...")
		if worker.currentStatus != StatusClearing {
			continue
		}

		fmt.Println("start clear...")
		predict := worker.nodeSchema.PrimaryPredict()
		if len(predict) <= 0 {
			logger.Errorf("primary predict should not be empty")
			continue
		}
		uids, err := worker.dataUpdater.UIDs(predict)
		if err != nil {
			logger.Errorf("worker cannot get uids from datasource: %s", err)
			continue
		}
		var mutations dgraph_helper.Mutations
		for _, uid := range uids {
			mutations = mutations.Add(
				dgraph_helper.NewMutation(
					dgraph_helper.MutationItemUid(uid),
					dgraph_helper.MutationItemValue("*"),
					dgraph_helper.MutationPredict("*"),
				),
			)
		}
		if mutations.IsEmpty() == false {
			if err = worker.dataUpdater.PushDelete(mutations.Assembly()); err != nil {
				logger.Errorf("worker[%s] cannot push delete data to updater: %s", worker.id, err)
				continue
			}
		}

		status := worker.schemaStatus
		waitForClearStatus := NewSchemaIndexStatusCleared(status.Schema, status.uid)
		data := waitForClearStatus.UpdateStatusMutations().Assembly()
		logger.Infof("worker[%s] update status to cleared: ", worker.id)
		fmt.Println(data)

		if err := worker.dataUpdater.PushUpdate(data); err != nil {
			logger.Errorf("worker push update data failed: %s", err)
			return
		}

		worker.chInternalCommand <- CommandCleared
	}
}

func (worker *ValueIndexWorker) startIndexLoop() {
	for {
		if worker.currentStatus != StatusIndexing {
			time.Sleep(time.Second)
			continue
		}
		raw, h, err := worker.dataSource.Read()
		if err != nil {
			if err == io.EOF {
				time.Sleep(time.Second)
				continue
			} else {
				logger.Errorf("worker read data from datasource failed: %s", err)
				time.Sleep(5 * time.Second)
				continue
			}
		}
		rdfs := worker.builder.Build(raw)
		rdfs = rdfs.Add(worker.schemaStatus.ProgressMutations(int64(h))...)
		data := rdfs.Assembly()
		logger.Infof("index data: \n%s", data)
		if err := worker.dataUpdater.PushUpdate(data); err != nil {
			logger.Errorf("worker push data to updater failed: %s", err)
			time.Sleep(5 * time.Second)
			continue
		}
		logger.Infof("index data source success [%d]", h)
	}
}

func (worker *ValueIndexWorker) Run() {
	go worker.startIndexLoop()
	go worker.startClearLoop()
	fmt.Println("worker running...")
	for {
		select {
		case command := <-worker.chExternalCommand:
			switch command {
			case CommandStop:
				worker.currentStatus = StatusStopped
				logger.Infof("worker[%s] status now is %s", worker.id, worker.currentStatus.String())
			case CommandClear:
				worker.currentStatus = StatusClearing
				logger.Infof("worker[%s] status now is %s", worker.id, worker.currentStatus.String())
			case CommandIndex:
				worker.currentStatus = StatusIndexing
				logger.Infof("worker[%s] status now is %s", worker.id, worker.currentStatus.String())
			default:
				logger.Infof("unknown command for worker")
			}
		case command := <-worker.chInternalCommand:
			switch command {
			case CommandCleared:
				worker.currentStatus = StatusCleared
			default:
				logger.Infof("unknown command for worker")
			}
		}
	}
}

func NewKVSchemaBuilder(schemaUid string, schemaInfo SchemaInfo, ns *schema.NodeSchema) *KVSchemaBuilder {
	builder := &KVSchemaBuilder{
		schemaInfo: schemaInfo,
	}
	builder.schemaUid = schemaUid
	builder.schema = ns
	builder.schemaBuilder = schema.NewRDFBuilder(ns)
	return builder
}

type KVSchemaBuilder struct {
	schemaUid     string
	schemaInfo    SchemaInfo
	schema        *schema.NodeSchema
	schemaBuilder *schema.RDFBuilder
}

type WriteOperation struct {
	Key     string
	Value   string
	Version int64
	Time    int64
}

func (builder *KVSchemaBuilder) CreateMutations(wo WriteOperation, predictPrefix string) (mutations dgraph_helper.Mutations) {
	result, ok := builder.schemaBuilder.Build(wo.Value, nil)
	if ok == false {
		return
	}
	uid := result.UID
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(uid),
			dgraph_helper.MutationItemValue(wo.Key),
			dgraph_helper.MutationPredict(predictPrefix+"-key"),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(uid),
			dgraph_helper.MutationItemValue(strconv.Itoa(int(wo.Version))),
			dgraph_helper.MutationPredict(predictPrefix+"-version"),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(uid),
			dgraph_helper.MutationItemUid(builder.schemaUid),
			dgraph_helper.MutationPredict(predictPrefix+"-schema"),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(uid),
			dgraph_helper.MutationItemValue(strconv.Itoa(int(wo.Time))),
			dgraph_helper.MutationPredict(predictPrefix+"-time"),
		),
	)
	mutations = mutations.Add(result.Mutations...)
	return
}

func (builder *KVSchemaBuilder) Build(src string) (mutations dgraph_helper.Mutations) {
	data := gjson.Parse(src).Get("data")
	if data.Exists() == false {
		logger.Warnf("cannot find data to build: %s", src)
		return
	}
	var wos []WriteOperation

	for _, tx := range data.Array() {
		if tx.Get("result.executionState").String() != "SUCCESS" {
			continue
		}
		txTime := tx.Get("request.transactionContent.timestamp").Int()
		operations := tx.Get("request.transactionContent.operations")
		if operations.Exists() == false {
			continue
		}

		for _, op := range operations.Array() {
			if op.Get("accountAddress").String() != builder.schemaInfo.AssociateAccount {
				continue
			}

			wsets := op.Get("writeSet")
			if wsets.Exists() == false {
				continue
			}
			for _, ws := range wsets.Array() {
				value := ws.Get("value.bytes").String()
				decoded, err := base64.StdEncoding.DecodeString(value)
				if err != nil {
					logger.Warnf("decode %s, got error %s", value, err)
					continue
				}
				wos = append(wos, WriteOperation{
					Key:     ws.Get("key").String(),
					Value:   string(decoded),
					Version: ws.Get("expectedVersion").Int() + 1,
					Time:    txTime,
				})
			}
		}
	}

	for _, wo := range wos {
		mutations = mutations.Add(builder.CreateMutations(wo, builder.schema.LowerName())...)
	}
	return
}
