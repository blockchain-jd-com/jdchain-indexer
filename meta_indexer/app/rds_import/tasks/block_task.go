package tasks

import (
	"fmt"
	"time"

	"git.jd.com/jd-blockchain/explorer/adaptor"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/app/rds_import/types"
)

type BlockTask struct {
	id      string
	apiHost string
	ledger  string
	height  int64
	data    *types.Blocks
	err     error
}

func NewBlockTask(host, ledger string, height int64) *BlockTask {
	return &BlockTask{
		id:      fmt.Sprintf("block-task-%s-%d", ledger, height),
		apiHost: host,
		ledger:  ledger,
		height:  height,
	}
}

func (blockTask *BlockTask) ID() string {
	return blockTask.id
}

func (blockTask *BlockTask) Status() error {
	return blockTask.err
}

func (blockTask *BlockTask) Data() *types.Blocks {
	return blockTask.data
}

func (blockTask *BlockTask) Do() error {

	block, err := adaptor.GetBlockFromServer(blockTask.apiHost, blockTask.ledger, blockTask.height)
	if err != nil {
		blockTask.err = err
		return err
	}

	dbBlock := types.Blocks{
		Ledger:                blockTask.ledger,
		BlockHeight:           blockTask.height,
		BlockHash:             block.Hash,
		PreBlockHash:          block.PreviousHash,
		TxsSetHash:            block.TransactionSetHash,
		UsersSetHash:          block.UserAccountSetHash,
		ContractsSetHash:      block.ContractAccountSetHash,
		ConfigurationsSetHash: block.AdminAccountHash,
		DataAccountsSetHash:   block.DataAccountSetHash,
		EventAccountsSetHash:  block.UserEventSetHash,
		BlockTimestamp:        time.UnixMilli(block.Time),
	}

	blockTask.data = &dbBlock

	return nil
}
