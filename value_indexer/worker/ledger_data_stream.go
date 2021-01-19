package worker

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/adaptor"
	"github.com/RoseRocket/xerrs"
	"io"
	"time"
)

func NewLedgerDataSteam(host, ledger string, from, to int) *LedgerDataSteam {
	stream := &LedgerDataSteam{
		apiHost: host,
		ledger:  ledger,
		current: from,
		to:      to,
		chStop:  make(chan bool, 1),
	}
	if to < 0 {
		stream.autoDetect = true
	}
	go stream.run()

	return stream
}

type LedgerDataSteam struct {
	apiHost    string
	ledger     string
	autoDetect bool

	current int
	from    int
	to      int

	chStop chan bool
}

func (stream *LedgerDataSteam) run() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-stream.chStop:
			return
		case <-ticker.C:
			if stream.autoDetect == false {
				continue
			}
			ledger, err := adaptor.GetLedgerDetailFromServer(stream.apiHost, stream.ledger)
			if err != nil {
				logger.Errorf("get ledger latest info failed: %s", err)
				continue
			}
			if int(ledger.Height) > stream.to {
				stream.to = int(ledger.Height)
				logger.Infof("ledger %s stream height updated to %d", stream.ledger, stream.to)
			}
		}
	}
}

func (stream *LedgerDataSteam) Stop() {
	stream.chStop <- true
}

func (stream *LedgerDataSteam) Read() (values string, height int, e error) {
	if stream.current > stream.to {
		return "", stream.to, io.EOF
	}

	raw, _, err := adaptor.GetTxListInBlockRawFromServer(stream.apiHost, stream.ledger, int64(stream.current), 0, -1)
	if err != nil {
		e = xerrs.Mask(err, fmt.Errorf("get tx list from server failed"))
		return
	}
	values = string(raw)
	height = stream.current
	stream.current++
	return
}
