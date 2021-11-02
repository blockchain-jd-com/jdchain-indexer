package meta_level_task

import (
	"container/list"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/event"
	"git.jd.com/jd-blockchain/explorer/level_task"
	"git.jd.com/jd-blockchain/explorer/worker"
	"github.com/davecgh/go-spew/spew"
	"time"
)

//type LevelTaskHandlerMap map[string]

func NewLevelTaskHandler(host, ledger string, dgraphHelper *dgraph_helper.Helper, rdfSaver func(string) error) *LevelTaskHandler {
	handler := &LevelTaskHandler{
		ledgerHost:   host,
		ledger:       ledger,
		rdfSaver:     rdfSaver,
		lowWaterLine: 5,
	}
	cache := dgraph_helper.NewUidLruCache(dgraphHelper, -1)
	handler.cache = cache
	worker4Write := worker.NewConcurrentWorker("writer", 8)
	worker4Write.AddListeners(handler)
	handler.writeWorker = worker4Write
	handler.tasks = list.New()
	handler.chTryHandleTask = make(chan bool, 1024)
	handler.chNewTask = make(chan level_task.LevelTask, 1024)
	return handler
}

type LevelTaskHandler struct {
	ledger       string
	ledgerHost   string
	lowWaterLine int

	cache       *dgraph_helper.UidLruCache
	rdfSaver    func(string) error
	writeWorker *worker.ConcurrentWorker
	tasks       *list.List

	chNewTask       chan level_task.LevelTask
	chTryHandleTask chan bool
}

func (handler *LevelTaskHandler) IsTaskLow() bool {
	return handler.tasks.Len() < handler.lowWaterLine
}

func (handler *LevelTaskHandler) AddLevelTask(tasks ...level_task.LevelTask) {
	for _, task := range tasks {
		handler.chNewTask <- task
	}
}

func (handler *LevelTaskHandler) Setup() *LevelTaskHandler {
	go handler.run()
	return handler
}

func (handler *LevelTaskHandler) run() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case lt := <-handler.chNewTask:
			handler.tasks.PushBack(lt)
			handler.addWriteRDFTask()
		case <-handler.chTryHandleTask:
			handler.addWriteRDFTask()
		case <-ticker.C:
			handler.addWriteRDFTask()
		}
	}
}

func (handler *LevelTaskHandler) addWriteRDFTask() {
	count := 0
	ele := handler.tasks.Front()
	for {
		if ele == nil {
			break
		}
		switch t := ele.Value.(type) {
		case *MetaInfoLevelTask:
			task := newMetaInfoLevelTaskHandler(t, handler.ledgerHost, handler.cache, handler.rdfSaver)
			if handler.writeWorker.AddTask(task) == false {
				break
			}
			count++
			temp := ele.Next()
			handler.tasks.Remove(ele)
			ele = temp
		default:
			logger.Warnf("unknown level task: \n%s", spew.Sdump(ele.Value))
		}
	}
	logger.Debugf("addWriteRDFTask: %d ", count)
}

func (handler *LevelTaskHandler) EventReceived(e event.Event) bool {
	sponsor := e.GetSponsor()
	switch e.GetName() {
	case event.EventWorkerNoTask:
		logger.Infof("no task in worker [%s] now", sponsor)
	case event.EventWorkerTaskComplete:
		if sponsor == "writer" {
			//handler.addWriteRDFTask()
			handler.chTryHandleTask <- true
		}
	default:
		logger.Infof("no handler for event [%s]", e.GetName())
	}
	return true
}

func (handler *LevelTaskHandler) ID() string {
	return "handler"
}
