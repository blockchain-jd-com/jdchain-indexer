package query

import (
	"github.com/dgraph-io/dgo/v200"
	"time"
)

type LedgerEventReceiver interface {
	AddStatus(ledger string)
}

func NewLedgerMonitor(eventReceiver LedgerEventReceiver, dgClient *dgo.Dgraph) *LedgerMonitor {
	return &LedgerMonitor{
		eventReceiver: eventReceiver,
		dgClient:      dgClient,
	}
}

type LedgerMonitor struct {
	eventReceiver LedgerEventReceiver
	dgClient      *dgo.Dgraph
}

func (monitor *LedgerMonitor) Run(seconds int) {
	if seconds <= 0 {
		seconds = 30
	}
	duration := time.Duration(seconds) * time.Second
	ticker := time.NewTicker(duration)
	go func() {
		for {
			<-ticker.C
			ledgers := getLedgers(monitor.dgClient)
			if len(ledgers) <= 0 {
				continue
			}

			if monitor.eventReceiver == nil {
				continue
			}
			for _, ledger := range ledgers {
				monitor.eventReceiver.AddStatus(ledger.HashID)
			}
		}
	}()
}

func getLedgers(dgClient *dgo.Dgraph) (ledgers Ledgers) {
	qe := NewLedgerQuery()
	result, e := qe.DoQuery(dgClient)
	if e != nil {
		logger.Warnf("query ledger from dgraph failed: %s", e)
		return
	}
	ledgers = result.(Ledgers)
	return
}
