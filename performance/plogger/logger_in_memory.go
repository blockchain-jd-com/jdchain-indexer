package plogger

import (
	"container/list"
)

func NewCSVLoggerInMemory(headers []string) *CSVLoggerInMemory {
	logger := &CSVLoggerInMemory{
		headers: headers,
		caches:  list.New(),
	}

	return logger
}

type CSVLoggerInMemory struct {
	headers []string
	caches  *list.List
}

func (logger *CSVLoggerInMemory) resetCache() {
	logger.caches = list.New()
}

func (logger *CSVLoggerInMemory) AddRecord(record []string) {
	logger.caches.PushBack(record)
	cacheSize := logger.caches.Len()
	if cacheSize > MaxCacheSize {
		logger.caches.Remove(logger.caches.Front())
	}
}
