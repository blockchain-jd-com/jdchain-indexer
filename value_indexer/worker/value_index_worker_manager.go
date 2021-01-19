package worker

import (
	"git.jd.com/jd-blockchain/explorer/value_indexer/schema"
)

func NewValueIndexWorkerManager(host string, sync DataUpdater) *ValueIndexWorkerManager {
	manager := &ValueIndexWorkerManager{
		workers:  map[string]*ValueIndexWorker{},
		chEvent:  make(chan IndexEvent, 1024),
		host:     host,
		dataSync: sync,
	}
	go manager.run()
	return manager
}

type ValueIndexWorkerManager struct {
	listener *SchemaMonitor
	workers  map[string]*ValueIndexWorker
	chEvent  chan IndexEvent
	host     string
	dataSync DataUpdater
}

func (manager *ValueIndexWorkerManager) OnEvent(e IndexEvent) {
	manager.chEvent <- e
}

func (manager *ValueIndexWorkerManager) startWorkerClearing(schemaStatus SchemaIndexStatus) {
	info := schemaStatus.Schema
	id := info.ID
	if len(id) <= 0 {
		return
	}
	worker, ok := manager.workers[id]
	if ok == false {
		logger.Warnf("worker %s not found to stop, try to create worker", id)
		manager.createNewWorker(schemaStatus)
		worker, _ = manager.workers[id]
	}
	worker.Stop()
	worker.Clear()
}

func (manager *ValueIndexWorkerManager) stopWorkerIndexing(schemaStatus SchemaIndexStatus) {
	info := schemaStatus.Schema
	id := info.ID
	if len(id) <= 0 {
		return
	}
	worker, ok := manager.workers[id]
	if ok == false {
		logger.Warnf("worker %s not found to stop, try to create worker", id)
		manager.createNewWorker(schemaStatus)
		worker, _ = manager.workers[id]
	} else {
		delete(manager.workers, id)
	}
	worker.Stop()
}

func (manager *ValueIndexWorkerManager) startWorkerIndexing(schemaStatus SchemaIndexStatus) {
	info := schemaStatus.Schema
	id := info.ID
	if len(id) <= 0 {
		logger.Warnf("startWorkerIndexing failed for id is empty")
		return
	}
	worker, ok := manager.workers[id]
	if ok == false {
		logger.Warnf("startWorkerIndexing failed for no worker with id [%s] found, try to create worker", id)
		manager.createNewWorker(schemaStatus)
		worker, _ = manager.workers[id]
	}
	ledger := info.Ledger
	ledgerStream := NewLedgerDataSteam(manager.host, ledger, int(schemaStatus.Progress)+1, schemaStatus.to)

	worker.Start(ledgerStream)
}

func (manager *ValueIndexWorkerManager) createNewWorker(schemaStatus SchemaIndexStatus) {
	info := schemaStatus.Schema
	id := info.ID
	if len(id) <= 0 {
		logger.Warnf("createNewWorker failed for id is empty")
		return
	}
	_, ok := manager.workers[id]
	if ok {
		logger.Warnf("createNewWorker failed for id is already exits")
		return
	}
	ns, err := schema.NewSchemaParser().FirstNodeSchema(info.Content)
	if err != nil {
		logger.Errorf("create worker failed: %s", err)
		return
	}

	kvSchemaBuilder := NewKVSchemaBuilder(schemaStatus.uid, info, ns)
	worker := NewValueIndexWorker(info.ID, schemaStatus, ns, kvSchemaBuilder, manager.dataSync)
	manager.workers[id] = worker
}

func (manager *ValueIndexWorkerManager) run() {
	for {
		select {
		case e := <-manager.chEvent:
			switch e.Status {
			case SchemaStatusDefault:
				manager.createNewWorker(e.schemaStatus)
				logger.Infof("create index worker for schema %s at status %s", e.schemaStatus.Schema.ID, e.Status)

			case SchemaStatusRunning:
				manager.startWorkerIndexing(e.schemaStatus)
				logger.Infof("start index worker for schema %s at status %s", e.schemaStatus.Schema.ID, e.Status)
			case SchemaStatusStopped:
				manager.stopWorkerIndexing(e.schemaStatus)
				logger.Infof("create index worker for schema %s at status %s", e.schemaStatus.Schema.ID, e.Status)
			case SchemaStatusClearing:
				manager.startWorkerClearing(e.schemaStatus)
				logger.Infof(" worker start clearing for schema %s at status %s", e.schemaStatus.Schema.ID, e.Status)
			default:
				logger.Infof("no handler for %s in worker manager", e.Status)
			}
		}
	}
}
