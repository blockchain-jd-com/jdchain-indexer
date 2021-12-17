package tasks

import (
	"fmt"
	"time"

	"git.jd.com/jd-blockchain/explorer/adaptor"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/app/rds_import/types"
)

type DataAccountTask struct {
	id       string
	apiHost  string
	ledger   string
	from     int64
	count    int64
	accounts []*types.DataAccounts
	kvs      []*types.DataAccountKVS
	err      error
}

func GetDataAccountCount(host, ledger string) (int64, error) {
	return adaptor.GetTotalAccountCountInLedgerFromServer(host, ledger)
}

func NewDataAccountTasks(host, ledger string, totalCount, taskSize int64) []*DataAccountTask {

	var dataAccountTasks []*DataAccountTask

	mod := totalCount % taskSize
	taskCount := totalCount / taskSize

	if mod != 0 {
		taskCount = taskCount + 1
	}

	for i := int64(1); i <= taskCount; i++ {
		from := (i - 1) * taskSize
		ct := &DataAccountTask{
			id:      fmt.Sprintf("data-account-task-%s-%d", ledger, i),
			apiHost: host,
			ledger:  ledger,
			from:    from,
			count:   taskSize,
		}
		dataAccountTasks = append(dataAccountTasks, ct)
	}

	return dataAccountTasks
}

func (dataAccountTask *DataAccountTask) ID() string {
	return dataAccountTask.id
}

func (dataAccountTask *DataAccountTask) Status() error {
	return dataAccountTask.err
}

func (dataAccountTask *DataAccountTask) Accounts() []*types.DataAccounts {
	return dataAccountTask.accounts
}

func (dataAccountTask *DataAccountTask) KVS() []*types.DataAccountKVS {
	return dataAccountTask.kvs
}

func (dataAccountTask *DataAccountTask) Do() error {

	var accounts []*types.DataAccounts
	var kvs []*types.DataAccountKVS

	dataAccountsList, err := adaptor.GetAccountsFromServer(dataAccountTask.apiHost, dataAccountTask.ledger, dataAccountTask.from, dataAccountTask.count)
	if err != nil {
		dataAccountTask.err = err
		return err
	}

	for _, dataAccount := range dataAccountsList {

		dataAccountInfo, err := adaptor.GetAccountInfoFromServer(dataAccountTask.apiHost, dataAccountTask.ledger, dataAccount.Address)
		if err != nil {
			dataAccountTask.err = err
			return err
		}

		role := dataAccountInfo.Get("permission.role").String()
		if role == "" {
			role = `["DEFAULT"]`
		}

		dataAccounts := types.DataAccounts{
			Ledger:                dataAccountTask.ledger,
			DataAccountAddress:    dataAccount.Address,
			DataAccountPubkey:     dataAccount.PublicKey,
			DataAccountRoles:      role,
			DataAccountPrivileges: dataAccountInfo.Get("permission.modeBits").String(),
			DataAccountCreator:    dataAccountInfo.Get("permission.owners").String(),
		}

		accounts = append(accounts, &dataAccounts)

		entriesCount, err := adaptor.GetTotalAccountEntriesCountFromServer(dataAccountTask.apiHost, dataAccountTask.ledger, dataAccount.Address)
		if err != nil {
			dataAccountTask.err = err
			return err
		}

		if entriesCount == 0 {
			continue
		}

		pageSize := int64(100)
		loops := calPage(entriesCount, pageSize)
		for i := int64(1); i <= loops; i++ {
			accountEntriesList, err := adaptor.GetAccountEntriesListFromServer(dataAccountTask.apiHost, dataAccountTask.ledger, dataAccount.Address, (i-1)*pageSize, pageSize)
			if err != nil {
				dataAccountTask.err = err
				return err
			}

			for _, accountEntry := range accountEntriesList {

				kv := types.DataAccountKVS{
					Ledger:             dataAccountTask.ledger,
					DataAccountAddress: dataAccount.Address,
					DataAccountKey:     accountEntry.Get("key").String(),
					DataAccountValue:   accountEntry.Get("value").String(),
					DataAccountType:    accountEntry.Get("type").String(),
					DataAccountVersion: int(accountEntry.Get("version").Int()),
					CreateTime:         time.Now(),
					UpdateTime:         time.Now(),
					State:              1,
				}

				kvs = append(kvs, &kv)
			}

		}

	}

	dataAccountTask.accounts = accounts
	dataAccountTask.kvs = kvs
	return nil
}

func calPage(totalCount, pageSize int64) int64 {
	mod := totalCount % pageSize
	count := totalCount / pageSize

	if mod != 0 {
		count = count + 1
	}

	return count
}
