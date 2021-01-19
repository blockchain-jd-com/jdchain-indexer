package worker

import (
	"fmt"
	"time"
)

func NewIndexEvent(status SchemaStatus, schemaInfo SchemaIndexStatus, from, to int64) IndexEvent {
	return IndexEvent{
		//Name:       name,
		Status:       status,
		schemaStatus: schemaInfo,
		From:         from,
		To:           to,
	}
}

type IndexEvent struct {
	Status       SchemaStatus
	schemaStatus SchemaIndexStatus
	From, To     int64
}

func (e IndexEvent) String() string {
	return fmt.Sprintf("event: schema %s -> %s", e.schemaStatus.Schema.ID, e.schemaStatus.Schema.ID)
}

type SchemaStatusSource interface {
	GetSchemaStatus() (SchemaIndexStatusList, error)
}

func NewSchemaMonitor(source SchemaStatusSource) *SchemaMonitor {
	monitor := &SchemaMonitor{
		source:   source,
		chRemind: make(chan bool, 1024),
	}
	go monitor.run()
	return monitor
}

type SchemaStatusChangeEventListener interface {
	OnEvent(e IndexEvent)
}

func (monitor *SchemaMonitor) run() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			monitor.Refresh()
		case <-monitor.chRemind:
			monitor.Refresh()
		}
	}
}

func (monitor *SchemaMonitor) Remind() {
	monitor.chRemind <- true
}

func (monitor *SchemaMonitor) Refresh() {

	// 全部重置为 SchemaIndexStatusVersionRemoved
	verStatusList := monitor.verStatusList.Map(func(status VersionedSchemaIndexStatus) VersionedSchemaIndexStatus {
		return VersionedSchemaIndexStatus{
			schemaIndexStatus: status.schemaIndexStatus,
			version:           SchemaIndexStatusVersionRemoved,
		}
	})

	// 载入数据库中数据
	statusList, err := monitor.source.GetSchemaStatus()
	if err != nil {
		logger.Errorf("get schema list failed: %s", err)
		return
	}

	for _, ss := range statusList {
		vss := verStatusList.Filter(
			func(status VersionedSchemaIndexStatus) bool {
				return status.schemaIndexStatus.Schema.ID == ss.Schema.ID
			})
		// 在 verStatusList 中不存在则加入 并设置状态为 SchemaIndexStatusVersionUpdated
		if vss.Len() <= 0 {
			verStatusList = verStatusList.Append(
				VersionedSchemaIndexStatus{
					schemaIndexStatus: ss,
					version:           SchemaIndexStatusVersionUpdated,
				})
			continue
		}

		// 在 verStatusList 存在且状态不一致的 重置状态为 SchemaIndexStatusVersionUpdated
		vs := vss[0]
		if vs.schemaIndexStatus.Status != ss.Status {
			verStatusList = verStatusList.
				FilterNot(
					func(status VersionedSchemaIndexStatus) bool {
						return status.schemaIndexStatus.Schema.ID == ss.Schema.ID
					}).
				Append(
					VersionedSchemaIndexStatus{
						schemaIndexStatus: ss,
						version:           SchemaIndexStatusVersionUpdated,
					})
			continue
		}

		// 在 verStatusList 中存在 且状态一致 为 SchemaStatusCleared 的 设置为 SchemaIndexStatusVersionUpdated
		if ss.Status == SchemaStatusCleared {
			verStatusList = verStatusList.
				FilterNot(
					func(status VersionedSchemaIndexStatus) bool {
						return status.schemaIndexStatus.Schema.ID == ss.Schema.ID
					}).
				Append(
					VersionedSchemaIndexStatus{
						schemaIndexStatus: ss,
						version:           SchemaIndexStatusVersionUpdated,
					})
			continue
		}

		// 在 verStatusList 中存在 且状态一致 不为 SchemaStatusCleared 的 设置为 SchemaIndexStatusVersionDefault
		verStatusList = verStatusList.
			FilterNot(
				func(status VersionedSchemaIndexStatus) bool {
					return status.schemaIndexStatus.Schema.ID == ss.Schema.ID
				}).
			Append(
				VersionedSchemaIndexStatus{
					schemaIndexStatus: ss,
					version:           SchemaIndexStatusVersionDefault,
				})

	}
	updated := verStatusList.Filter(func(status VersionedSchemaIndexStatus) bool {
		return status.version == SchemaIndexStatusVersionUpdated
	})
	monitor.PublishEvent(updated)

	monitor.verStatusList = verStatusList.FilterNot(func(status VersionedSchemaIndexStatus) bool {
		return status.version == SchemaIndexStatusVersionRemoved
	})
}

func (monitor *SchemaMonitor) PublishEvent(updated VersionedSchemaIndexStatusList) {
	for _, vs := range updated {
		e := NewIndexEvent(vs.schemaIndexStatus.Status, vs.schemaIndexStatus, 0, 0)
		monitor.NotifyListener(e)
		logger.Infof("event notified to listener: %s", e.String())
	}
}

func (monitor *SchemaMonitor) NotifyListener(e IndexEvent) {
	for _, listener := range monitor.eventListeners {
		listener.OnEvent(e)
	}
}

func (monitor *SchemaMonitor) AddListener(listeners ...SchemaStatusChangeEventListener) {
	monitor.eventListeners = append(monitor.eventListeners, listeners...)
}

type SchemaMonitor struct {
	source         SchemaStatusSource
	eventListeners []SchemaStatusChangeEventListener
	verStatusList  VersionedSchemaIndexStatusList
	chRemind       chan bool
}
