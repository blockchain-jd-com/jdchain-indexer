package query

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v200"
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"github.com/ssor/zlog"
	"strings"
	"sync"
	"time"
)

const (
	SearchTypeBlock        = "**block**"
	SearchTypeTx           = "**tx**"
	SearchTypeContract     = "**contract**"
	SearchTypeUser         = "**user**"
	SearchTypeAccount      = "**account**"
	SearchTypeEventAccount = "**event_account**"
)

func NewDocCacher(engine *riot.Engine, cacheSize int, dgClient *dgo.Dgraph) *DocCacher {
	updater := &DocCacher{
		indexEngine:          engine,
		LedgerDocCacheStatus: sync.Map{},
		cacheSize:            cacheSize,
		dgClient:             dgClient,
	}
	updater.run()
	return updater
}

type DocCacher struct {
	dgClient             *dgo.Dgraph
	cacheSize            int
	indexEngine          *riot.Engine
	LedgerDocCacheStatus sync.Map
}

func (cacher *DocCacher) AddStatus(ledger string) {
	_, ok := cacher.LedgerDocCacheStatus.Load(ledger)
	if ok == false {
		cacher.LedgerDocCacheStatus.Store(ledger, newDocCacheStatus(ledger, cacher, cacher.cacheSize))
	}
}

type IndexObject interface {
	GetHashID() string
}

func (cacher *DocCacher) Index(obj IndexObject, ledger, indexType string) {
	labels := stringSplit(obj.GetHashID())
	bs, err := json.Marshal(obj)
	if err != nil {
		zlog.Warnf("marshal obj failed: %s", err)
		return
	}

	content := append(labels, ledger, indexType)
	data := types.DocData{
		Content: strings.Join(content, " "),
		Attri:   []string{indexType, string(bs)},
	}
	cacher.indexEngine.Index(obj.GetHashID(), data)
	cacher.indexEngine.Flush()
}

func (cacher *DocCacher) allStatusOK(size int) bool {
	allOK := true
	cacher.LedgerDocCacheStatus.Range(func(key, value interface{}) bool {
		b := value.(*DocCacheStatus).Ok(size)
		if b == false {
			allOK = false
			return false
		}
		return true
	})
	return allOK
}

func (cacher *DocCacher) updateStatus(size int) (e error) {
	cacher.LedgerDocCacheStatus.Range(func(key, value interface{}) bool {
		err := value.(*DocCacheStatus).Update(cacher.dgClient)
		if err != nil {
			e = err
			return false
		}
		//logger.Debugf(value.(*DocCacheStatus).Summary())
		return true
	})
	return
}

func (cacher *DocCacher) run() {
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		for {
			<-ticker.C
			if cacher.allStatusOK(cacher.cacheSize) {
				break
			}

			err := cacher.updateStatus(cacher.cacheSize)
			if err != nil {
				zlog.Warnf("update cache status of ledger [%s] failed: %s", err)
				continue
			}
		}
	}()
}

type IndexEngine interface {
	Index(obj IndexObject, ledger, indexType string)
}

func newDocCacheStatus(ledger string, engine IndexEngine, cacheSize int) *DocCacheStatus {
	return &DocCacheStatus{
		ledger:      ledger,
		indexEngine: engine,
		cacheSize:   cacheSize,
	}
}

type DocCacheStatus struct {
	ledger                 string
	indexEngine            IndexEngine
	cacheSize              int
	blockCacheCount        int
	txCacheCount           int
	contractCacheCount     int
	accountCacheCount      int
	userCacheCount         int
	eventAccountCacheCount int
}

func (status *DocCacheStatus) Ok(size int) bool {
	if status.blockCacheCount < size ||
		status.txCacheCount < size ||
		status.contractCacheCount < size ||
		status.accountCacheCount < size ||
		status.userCacheCount < size ||
		status.eventAccountCacheCount < size {
		return false
	}
	return true
}

func (status *DocCacheStatus) Summary() string {
	return fmt.Sprintf(`
     ledger cache status:   %s
        block:    			%d
        tx:       			%d
        contract: 			%d
        account:  			%d
        user:     			%d
        event_account:      %d
`,
		status.ledger, status.blockCacheCount, status.txCacheCount, status.contractCacheCount, status.accountCacheCount, status.userCacheCount, status.eventAccountCacheCount)
}

func (status *DocCacheStatus) Update(dgClient *dgo.Dgraph) error {
	if status.blockCacheCount < status.cacheSize {
		count, err := status.updateBlockCache(dgClient)
		if err != nil {
			return err
		}
		status.blockCacheCount += count
	}
	if status.txCacheCount < status.cacheSize {
		count, err := status.updateTxCache(dgClient)
		if err != nil {
			return err
		}
		status.txCacheCount += count
	}
	if status.contractCacheCount < status.cacheSize {
		count, err := status.updateContractCache(dgClient)
		if err != nil {
			return err
		}
		status.contractCacheCount += count
	}
	if status.accountCacheCount < status.cacheSize {
		count, err := status.updateAccountCache(dgClient)
		if err != nil {
			return err
		}
		status.accountCacheCount += count
	}
	if status.userCacheCount < status.cacheSize {
		count, err := status.updateUserCache(dgClient)
		if err != nil {
			return err
		}
		status.userCacheCount += count
	}
	if status.eventAccountCacheCount < status.cacheSize {
		count, err := status.updateEventAccountCache(dgClient)
		if err != nil {
			return err
		}
		status.eventAccountCacheCount += count
	}
	return nil
}

func (status *DocCacheStatus) updateTxCache(dgClient *dgo.Dgraph) (count int, e error) {
	txs, err := txListInLedger(status.ledger, int64(status.txCacheCount), int64(status.cacheSize-1), dgClient)
	if err != nil {
		return 0, err
	}
	for _, tx := range txs {
		status.indexEngine.Index(tx, status.ledger, SearchTypeTx)
	}
	return len(txs), nil
}

func (status *DocCacheStatus) updateBlockCache(dgClient *dgo.Dgraph) (count int, e error) {
	blocks, err := blockListInLedger(status.ledger, int64(status.blockCacheCount), int64(status.cacheSize-1), dgClient)
	if err != nil {
		return 0, err
	}
	for _, block := range blocks {
		status.indexEngine.Index(block, status.ledger, SearchTypeBlock)
	}
	return len(blocks), nil
}

func (status *DocCacheStatus) updateContractCache(dgClient *dgo.Dgraph) (count int, e error) {
	contracts, err := contractListInLedger(status.ledger, int64(status.contractCacheCount), int64(status.cacheSize-1), dgClient)
	if err != nil {
		return 0, err
	}
	for _, contract := range contracts {
		status.indexEngine.Index(contract, status.ledger, SearchTypeContract)
	}
	return len(contracts), nil
}

func (status *DocCacheStatus) updateUserCache(dgClient *dgo.Dgraph) (count int, e error) {
	users, err := userListInLedger(status.ledger, int64(status.userCacheCount), int64(status.cacheSize)-1, dgClient)
	if err != nil {
		return 0, err
	}
	for _, user := range users {
		status.indexEngine.Index(user, status.ledger, SearchTypeUser)
	}
	return len(users), nil
}

func (status *DocCacheStatus) updateAccountCache(dgClient *dgo.Dgraph) (count int, e error) {
	accounts, err := accountListInLedger(status.ledger, int64(status.accountCacheCount), int64(status.cacheSize)-1, dgClient)
	if err != nil {
		return 0, err
	}
	for _, account := range accounts {
		status.indexEngine.Index(account, status.ledger, SearchTypeAccount)
	}
	return len(accounts), nil
}

func (status *DocCacheStatus) updateEventAccountCache(dgClient *dgo.Dgraph) (count int, e error) {
	accounts, err := eventAccountListInLedger(status.ledger, int64(status.eventAccountCacheCount), int64(status.cacheSize)-1, dgClient)
	if err != nil {
		return 0, err
	}
	for _, account := range accounts {
		status.indexEngine.Index(account, status.ledger, SearchTypeEventAccount)
	}
	return len(accounts), nil
}

func blockListInLedger(ledger string, from, to int64, dgClient *dgo.Dgraph) (blocks Blocks, e error) {
	query := NewBlockQueryInRange([]string{ledger}, from, to)
	result, err := query.DoQuery(dgClient)
	if err != nil {
		zlog.Errorf("query block failed: %s", err)
		return nil, err
	}
	return result.(Blocks), nil
}

func txListInLedger(ledger string, from, to int64, dgClient *dgo.Dgraph) (txs Transactions, e error) {
	query := NewQueryTxRange([]string{ledger}, from, to)
	result, err := query.DoQuery(dgClient)
	if err != nil {
		zlog.Errorf("query block failed: %s", err)
		return nil, err
	}
	return result.(Transactions), nil
}

func contractListInLedger(ledger string, from, to int64, dgClient *dgo.Dgraph) (contracts Contracts, e error) {
	query := NewQueryContractRange([]string{ledger}, from, to)
	result, err := query.DoQuery(dgClient)
	if err != nil {
		zlog.Errorf("query contract failed: %s", err)
		return nil, err
	}
	return result.(Contracts), nil
}

func userListInLedger(ledger string, from, to int64, dgClient *dgo.Dgraph) (users Users, e error) {
	query := NewQueryUsersRange([]string{ledger}, from, to)
	result, err := query.DoQuery(dgClient)
	if err != nil {
		zlog.Errorf("query user failed: %s", err)
		return nil, err
	}
	return result.(Users), nil
}

func accountListInLedger(ledger string, from, to int64, dgClient *dgo.Dgraph) (accounts Accounts, e error) {
	query := NewDatasetRangeQuery([]string{ledger}, from, to)
	result, err := query.DoQuery(dgClient)
	if err != nil {
		zlog.Errorf("query account failed: %s", err)
		return nil, err
	}
	return result.(Accounts), nil
}

func eventAccountListInLedger(ledger string, from, to int64, dgClient *dgo.Dgraph) (accounts EventAccounts, e error) {
	query := NewEventAccountRangeQuery([]string{ledger}, from, to)
	result, err := query.DoQuery(dgClient)
	if err != nil {
		zlog.Errorf("query account failed: %s", err)
		return nil, err
	}
	return result.(EventAccounts), nil
}
