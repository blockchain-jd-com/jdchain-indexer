package adaptor

import (
	"strings"
)

func IsExistsInTx(host, ledger, keyword string, maxHeight int64) (bool, error) {
	cs := NewContentSearcher(host, ledger, keyword, 1, maxHeight)
	b, err := cs.Search()
	if err != nil {
		return false, err
	}

	return b, nil
}

func NewContentSearcher(apiHost, ledger, keyword string, from, to int64) *ContentSearcher {
	cs := &ContentSearcher{
		keyword: keyword,
		apiHost: apiHost,
		ledger:  ledger,
		from:    from,
		to:      to,
	}
	return cs
}

type ContentSearcher struct {
	from, to int64
	keyword  string
	apiHost  string
	ledger   string
}

func (searcher *ContentSearcher) Search() (isFound bool, e error) {
	for i := searcher.from; i <= searcher.to; i++ {
		ok, err := searcher.SearchInBlock(i)
		if err != nil {
			logger.Errorf("stopped")
			return
		}
		isFound = ok
		if ok {
			logger.Infof("find keyword [%s] in block [%d]", searcher.keyword, i)
			return
		}
	}
	return
}

func (searcher *ContentSearcher) SearchInBlock(blockHeight int64) (isFound bool, e error) {
	//zlog.Debugf("start search ledger: [%s] from height: [%d] ", searcher.ledger, blockHeight)
	body, debugInfo, err := GetTxListInBlockRawFromServer(searcher.apiHost, searcher.ledger, blockHeight, 0, 100)
	if err != nil {
		e = err
		logger.WithStruct(debugInfo).Errorf("fetch block failed: %s", err)
		return
	}

	if strings.Contains(string(body), searcher.keyword) {
		isFound = true
	}
	return
}
