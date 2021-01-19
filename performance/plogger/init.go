package plogger

const (
	MaxCacheSize = 1000
)

type RecordLogger interface {
	AddRecord(record []string)
}
