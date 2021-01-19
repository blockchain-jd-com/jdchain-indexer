package query

import (
	"github.com/ssor/zlog"
)

var (
	logger = zlog.New("searcher", "query")

	MaxRecordsPerRequest = int64(1000)
)
