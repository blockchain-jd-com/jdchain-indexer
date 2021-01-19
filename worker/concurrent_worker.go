package worker

import (
	"container/list"
	"git.jd.com/jd-blockchain/explorer/event"
	"git.jd.com/jd-blockchain/explorer/performance"
	"time"
)

type Task interface {
	Do() error
}

func NewConcurrentWorker(id string, maxConcurrentWorkerCount int) *ConcurrentWorker {
	worker := &ConcurrentWorker{
		id:                         id,
		tasks:                      list.New(),
		eventNewTaskAddedListener:  make(chan int, 1024),
		eventTaskCompletedListener: make(chan bool, 1024),
		maxConcurrentWorkerCount:   maxConcurrentWorkerCount,
		maxTotalWorkCount:          64,
		performanceCounter:         performance.NewTimeCounter(-1),
		Manager:                    event.NewManager(),
	}
	worker.run()
	return worker
}

type ConcurrentWorker struct {
	id                       string
	workingWorkerCount       int
	maxConcurrentWorkerCount int
	maxTotalWorkCount        int
	tasks                    *list.List
	*event.Manager
	eventNewTaskAddedListener  chan int
	eventTaskCompletedListener chan bool
	performanceCounter         *performance.TimeCounter
}

func (cw *ConcurrentWorker) AddTask(tasks ...Task) bool {
	return cw.addTask(tasks...)
}

func (cw *ConcurrentWorker) addTask(tasks ...Task) bool {
	if cw.tasks.Len()+len(tasks) > cw.maxTotalWorkCount {
		return false
	}

	for _, task := range tasks {
		cw.tasks.PushBack(newTaskInfo(task))
	}
	cw.eventNewTaskAddedListener <- 0
	return true
}

func (cw *ConcurrentWorker) run() {
	go cw.startWorkerResultListening()
}

func (cw *ConcurrentWorker) startWorkerResultListening() {
	timer := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timer.C:
			if cw.workingWorkerCount <= 0 {
				cw.Notify(event.NewCommonEvent(event.EventWorkerNoTask, cw.id, cw.id))
			}
			cw.printPerformance()

		case <-cw.eventNewTaskAddedListener:
			cw.tryDoTask()

		case <-cw.eventTaskCompletedListener:
			cw.performanceCounter.AddTick()

			if cw.workingWorkerCount > 0 {
				cw.workingWorkerCount--
			}

			elements, TasksCompleted := cw.completedTasks()
			if len(TasksCompleted) > 0 {
				cw.Notify(event.NewCommonEvent(event.EventWorkerTaskComplete, TasksCompleted, cw.id))
				cw.removeTasks(elements)
			}

			elements, TasksFailed := cw.failedTasks()
			if len(TasksFailed) > 0 {
				cw.removeTasks(elements)
			}

			if cw.tryDoTask() == false {
				cw.Notify(event.NewCommonEvent(event.EventWorkerNoTask, cw.id, cw.id))
			}
		}
	}
}

func (cw *ConcurrentWorker) completedTasks() (elements []*list.Element, tasks []Task) {
	ele := cw.tasks.Front()
	for {
		if ele == nil {
			break
		}
		ti := ele.Value.(*TaskInfo)
		if ti.isCompleted() {
			elements = append(elements, ele)
			tasks = append(tasks, ti.task)
		}
		ele = ele.Next()
	}
	return
}

func (cw *ConcurrentWorker) failedTasks() (elements []*list.Element, tasks []Task) {
	ele := cw.tasks.Front()
	for {
		if ele == nil {
			break
		}
		ti := ele.Value.(*TaskInfo)
		if ti.isFailed() {
			elements = append(elements, ele)
			tasks = append(tasks, ti.task)
		}
		ele = ele.Next()
	}
	return
}

func (cw *ConcurrentWorker) removeTasks(elements []*list.Element) {
	for _, ele := range elements {
		cw.tasks.Remove(ele)
	}
	return
}

func (cw *ConcurrentWorker) printPerformance() {
	count, timeCost := cw.performanceCounter.Summary()
	if count <= 0 {
		return
	}
	logger.Infof("performance:  [%d] task cost [%s]", count, timeCost)
}

func (cw *ConcurrentWorker) tryDoTask() (hasWorkToDo bool) {
	if cw.workingWorkerCount >= cw.maxConcurrentWorkerCount {
		return true
	}

	wi := cw.firstTaskToStart()
	if wi == nil {
		return
	}

	cw.workingWorkerCount++
	//logger.Debugf("try to start new work [%s] -> current concurrent task: [%d]", wi.task.ID(), group.workingWorkerCount)
	wi.status = WorkStatusOnProcess
	go cw.doTask(wi)
	hasWorkToDo = true
	return
}

func (cw *ConcurrentWorker) doTask(taskInfo *TaskInfo) {
	taskInfo.do()
	cw.eventTaskCompletedListener <- true
}

func (cw *ConcurrentWorker) firstTaskToStart() (ti *TaskInfo) {
	ele := cw.tasks.Front()
	for {
		if ele == nil {
			break
		}

		wi := ele.Value.(*TaskInfo)
		if wi.needStartNow() {
			ti = wi
			break
		}
		ele = ele.Next()
	}

	return
}
