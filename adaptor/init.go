package adaptor

import (
	"git.jd.com/jd-blockchain/explorer/performance/plogger"
	"github.com/ssor/zlog"
	"strconv"
	"time"
)

var (
	logger = zlog.New("indexer", "adaptor")
)

var (
	plogGetRawData   plogger.RecordLogger = plogger.NewCSVLoggerInMemory([]string{"block", "time"})
	plogParseRawData plogger.RecordLogger = plogger.NewCSVLoggerInMemory([]string{"block", "tx-count", "time"})
)

func logParseRawData(height, txCount int64, start, end time.Time) {
	logger.Debugf("parse data: h: %d tx: %d time: %v", height, txCount, end.Sub(start))
	plogParseRawData.AddRecord([]string{
		strconv.FormatInt(height, 10),
		strconv.FormatInt(txCount, 10),
		strconv.FormatInt(end.Sub(start).Nanoseconds()/1000/1000, 10),
	})
}

func logGetRawData(height int64, start, end time.Time) {
	logger.Debugf("fetch raw data: %s at height [%d]", end.Sub(start), height)
	plogGetRawData.AddRecord([]string{
		strconv.FormatInt(height, 10),
		strconv.FormatInt(end.Sub(start).Nanoseconds()/1000/1000, 10),
	})
}
