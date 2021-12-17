package tasks

import (
	"fmt"

	"git.jd.com/jd-blockchain/explorer/adaptor"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/app/rds_import/types"
	"github.com/tidwall/gjson"
)

type ContractTask struct {
	id      string
	apiHost string
	ledger  string
	from    int64
	count   int64
	data    []*types.Contracts
	err     error
}

func GetContractCount(host, ledger string) (int64, error) {
	return adaptor.GetTotalContractCountInLedgerFromServer(host, ledger)
}

func NewContractTasks(host, ledger string, totalCount, taskSize int64) []*ContractTask {

	var contractTasks []*ContractTask

	mod := totalCount % taskSize
	taskCount := totalCount / taskSize

	if mod != 0 {
		taskCount = taskCount + 1
	}

	for i := int64(1); i <= taskCount; i++ {
		from := (i - 1) * taskSize
		ct := &ContractTask{
			id:      fmt.Sprintf("contract-task-%s-%d", ledger, i),
			apiHost: host,
			ledger:  ledger,
			from:    from,
			count:   taskSize,
		}
		contractTasks = append(contractTasks, ct)
	}

	return contractTasks
}

func (contractTask *ContractTask) ID() string {
	return contractTask.id
}

func (contractTask *ContractTask) Status() error {
	return contractTask.err
}

func (contractTask *ContractTask) Data() []*types.Contracts {
	return contractTask.data
}

func (contractTask *ContractTask) Do() error {

	contractAaddresses, err := adaptor.GetContractsFromServer(contractTask.apiHost, contractTask.ledger, contractTask.from, contractTask.count)

	if err != nil {
		contractTask.err = err
		return err
	}

	var contracts []*types.Contracts

	for _, address := range contractAaddresses {
		result, err := adaptor.GetContractDetailFromServer(contractTask.apiHost, contractTask.ledger, address.Address)
		if err != nil {
			contractTask.err = err
			return err
		}

		gResult := result.(gjson.Result)

		role := gResult.Get("permission.role").String()
		if role == "" {
			role = `["DEFAULT"]`
		}

		contract := &types.Contracts{
			Ledger:              contractTask.ledger,
			ContractAddress:     address.Address,
			ContractPubkey:      address.PublicKey,
			ContractRoles:       role,
			ContractPriviledges: gResult.Get("permission.modeBits").String(),
			ContractVersion:     int8(gResult.Get("chainCodeVersion").Int()),
			ContractStatus:      gResult.Get("state").String(),
			ContractCreator:     gResult.Get("permission.owners").String(),
			ContractContent:     gResult.Get("chainCode").String(),
		}

		contracts = append(contracts, contract)
	}

	contractTask.data = contracts
	return nil
}
