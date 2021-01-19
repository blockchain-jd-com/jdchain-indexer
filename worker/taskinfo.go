package worker

type WorkStatus int

const (
	WorkStatusNew       WorkStatus = 0
	WorkStatusOnProcess WorkStatus = 1
	WorkStatusFailed    WorkStatus = 2
	WorkStatusComplete  WorkStatus = 3
)

func newTaskInfo(task Task) *TaskInfo {
	ti := &TaskInfo{
		status: WorkStatusNew,
		task:   task,
	}
	return ti
}

type TaskInfo struct {
	task   Task
	status WorkStatus
}

func (ti *TaskInfo) do() {
	e := ti.task.Do()
	if e != nil {
		ti.status = WorkStatusFailed
		return
	}
	ti.status = WorkStatusComplete
	return
}

func (ti *TaskInfo) isCompleted() bool {
	return ti.status == WorkStatusComplete
}

func (ti *TaskInfo) isFailed() bool {
	return ti.status == WorkStatusFailed
}

func (ti *TaskInfo) isUncompleted() bool {
	return ti.status != WorkStatusComplete
}

func (ti *TaskInfo) needStartNow() bool {
	switch ti.status {
	case WorkStatusNew:
		return true
	case WorkStatusOnProcess:
		return false
	case WorkStatusFailed:
		return false
	case WorkStatusComplete:
		return false
	default:
		return false
	}
}
