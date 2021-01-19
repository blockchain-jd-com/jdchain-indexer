package worker

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/adaptor"
)

type LedgerData struct {
	Ledger    string
	BlockInfo *adaptor.Block
	Txs       []*adaptor.Transaction
}

func NewFetchTask(host, ledger string, height int64) *FetchTask {
	return &FetchTask{
		id:      fmt.Sprintf("%s-%d", ledger, height),
		apiHost: host,
		ledger:  ledger,
		Height:  height,
	}
}

type FetchTask struct {
	id      string
	apiHost string
	ledger  string
	Height  int64
	data    interface{}
}

func (fetcher *FetchTask) GetHeight() int64 {
	return fetcher.Height
}

func (fetcher *FetchTask) Data() interface{} {
	return fetcher.data
}

func (fetcher *FetchTask) Do() (e error) {
	logger.Debugf("start fetch block [%d]", fetcher.Height)

	block, err := adaptor.GetBlockFromServer(fetcher.apiHost, fetcher.ledger, fetcher.Height)
	if err != nil {
		e = err
		return
	}
	txs, err := adaptor.GetTxListInBlockFromServer(fetcher.apiHost, fetcher.ledger, fetcher.Height, 0, -1)
	if err != nil {
		e = err
		return
	}
	data := &LedgerData{
		Ledger:    fetcher.ledger,
		BlockInfo: block,
		Txs:       txs,
	}
	fetcher.data = data
	return
}
