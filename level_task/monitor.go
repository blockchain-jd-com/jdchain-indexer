package level_task

import (
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"github.com/RoseRocket/xerrs"
	"time"
)

type LevelTaskReceiver interface {
	AddLevelTask(tasks ...LevelTask)
	IsTaskLow() bool
}

func NewLevelTaskMonitor(taskReceiver LevelTaskReceiver, parsers LevelTaskParserMap, dgraphHelper *dgraph_helper.Helper) *LevelTaskMonitor {
	monitor := &LevelTaskMonitor{
		dgraphHelper: dgraphHelper,
		taskReceiver: taskReceiver,
	}
	monitor.puller = NewLevelTaskPuller(parsers)
	return monitor
}

type LevelTaskMonitor struct {
	puller       *LevelTaskPuller
	dgraphHelper *dgraph_helper.Helper
	taskReceiver LevelTaskReceiver
}

func (monitor *LevelTaskMonitor) Setup() *LevelTaskMonitor {
	go monitor.run()
	return monitor
}

func (monitor *LevelTaskMonitor) run() {
	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-ticker.C:
			if monitor.taskReceiver.IsTaskLow() {
				tasks, err := monitor.puller.Pull(monitor.dgraphHelper)
				if err != nil {
					logger.Warnf("task monitor pull task failed: \n%s", xerrs.Details(err, 5))
				}
				monitor.taskReceiver.AddLevelTask(tasks...)
			}
		}
	}
}
