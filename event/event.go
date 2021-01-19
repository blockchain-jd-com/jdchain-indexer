package event

type Event interface {
	GetSponsor() string
	GetName() string
	GetData() interface{}
}

const (
	EventWorkerTaskComplete = "worker-task-complete"
	EventWorkerNoTask       = "worker-no-task"
)

func NewCommonEvent(name string, data interface{}, sponsors ...string) *CommonEvent {
	e := &CommonEvent{
		name: name,
		data: data,
	}
	if len(sponsors) > 0 {
		e.sponsor = sponsors[0]
	}
	return e
}

type CommonEvent struct {
	data    interface{}
	name    string
	sponsor string
}

func (event *CommonEvent) GetData() interface{} {
	return event.data
}

func (event *CommonEvent) GetSponsor() string {
	return event.sponsor
}

func (event *CommonEvent) GetName() string {
	return event.name
}
