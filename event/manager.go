package event

import (
	"container/list"
)

type Listener interface {
	EventReceived(event Event) bool
	ID() string
}

type ListenerObj struct {
	Listener
	messagePool *list.List
	poolSize    int
}

func (obj *ListenerObj) Notify(event Event) {
	obj.cacheEvent(event)
	ele := obj.messagePool.Front()
	for {
		if ele == nil {
			break
		}

		v := ele.Value
		b := obj.EventReceived(v.(Event))
		if b {
			tmp := ele
			ele = ele.Next()
			obj.messagePool.Remove(tmp)
		} else {
			break
		}
	}
}

func (obj *ListenerObj) cacheEvent(event Event) {
	obj.messagePool.PushBack(event)
	if obj.messagePool.Len() > obj.poolSize {
		obj.messagePool.Remove(obj.messagePool.Front())
	}
}

func newListenerObj(listener Listener) *ListenerObj {
	return &ListenerObj{
		Listener:    listener,
		messagePool: list.New(),
		poolSize:    1024,
	}
}

func NewManager() *Manager {
	m := &Manager{
		listeners: list.New(),
	}
	return m
}

type Manager struct {
	listeners *list.List
}

func (em *Manager) Notify(event Event) {
	if em.listeners.Len() > 0 {
		ele := em.listeners.Front()
		for {
			if ele == nil {
				break
			}
			listener := ele.Value.(*ListenerObj)
			listener.Notify(event)

			ele = ele.Next()
		}
	}
}

func (em *Manager) RemoveListeners(listeners ...Listener) {
	for _, listenerOut := range listeners {
		ele := em.listeners.Front()
		for {
			if ele == nil {
				break
			}

			listener := ele.Value.(Listener)
			if listener.ID() == listenerOut.ID() {
				tmp := ele
				ele = ele.Next()
				em.listeners.Remove(tmp)
			} else {
				ele = ele.Next()
			}
		}
	}
}

func (em *Manager) AddListeners(listeners ...Listener) {
	for _, listener := range listeners {
		em.listeners.PushBack(newListenerObj(listener))
	}
}
