package performance

import (
	"container/list"
	"time"
)

type TimeCounter struct {
	count   int
	records *list.List
}

func NewTimeCounter(cacheCount int) *TimeCounter {
	timer := &TimeCounter{
		records: list.New(),
	}
	if cacheCount <= 0 {
		timer.count = 100
	} else {
		timer.count = cacheCount
	}
	return timer
}

func (counter *TimeCounter) Summary() (count int, cost time.Duration) {
	if counter.records.Len() <= 0 {
		return
	}
	start := counter.records.Front().Value.(time.Time)
	end := counter.records.Back().Value.(time.Time)
	elapse := end.Sub(start)
	return counter.records.Len(), elapse
}

func (counter *TimeCounter) AddTick() {
	counter.records.PushBack(time.Now())

	l := counter.records.Len()
	if l > counter.count {
		counter.records.Remove(counter.records.Front())
	}
}
