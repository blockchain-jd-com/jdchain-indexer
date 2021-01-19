package chain

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/adaptor"
	"git.jd.com/jd-blockchain/explorer/event"
	"strings"
	"time"
)

const (
	EventLedgerStatus = "ledger-status"
)

func NewLedgerStatus(ledger, apiHost string, h int64) *LedgerStatus {
	return &LedgerStatus{
		Ledger:  ledger,
		Height:  h,
		ApiHost: apiHost,
	}
}

type LedgerStatus struct {
	ApiHost string
	Ledger  string
	Height  int64
}

func NewChainMonitor(apiHost string) *Monitor {
	return &Monitor{
		apiHost: apiHost,
		Manager: event.NewManager(),
	}
}

type Monitor struct {
	*event.Manager
	apiHost string
}

func (monitor *Monitor) Run(duration time.Duration) {
	if duration <= time.Duration(0) {
		duration = 15 * time.Second
	}
	ticker := time.NewTicker(duration)
	go func() {
		for {
			<-ticker.C
			ledgers, err := adaptor.GetLedgersFromServer(monitor.apiHost)
			if err != nil {
				continue
			}
			if len(ledgers) <= 0 {
				continue
			}

			var ledgersInfo []string
			for _, ledger := range ledgers {
				info := fmt.Sprintf("-  %s  [height: %d]", ledger.Hash, ledger.Height)
				ledgersInfo = append(ledgersInfo, info)
				e := event.NewCommonEvent(EventLedgerStatus, NewLedgerStatus(ledger.Hash, monitor.apiHost, ledger.Height), "monitor")
				monitor.Notify(e)
			}
			logger.Info("ledgers list:")
			logger.Info(strings.Join(ledgersInfo, "\n"))
		}
	}()
}
