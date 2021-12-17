package tasks

import (
	"encoding/json"
	"errors"
	"fmt"

	"git.jd.com/jd-blockchain/explorer/adaptor"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/app/rds_import/types"
	"github.com/tidwall/gjson"
)

type TxTask struct {
	id      string
	apiHost string
	ledger  string
	height  int64
	data    []*types.Txs
	err     error
}

func NewTxTask(host, ledger string, height int64) *TxTask {
	return &TxTask{
		id:      fmt.Sprintf("tx-task-%s-%d", ledger, height),
		apiHost: host,
		ledger:  ledger,
		height:  height,
	}
}

func (txTask *TxTask) ID() string {
	return txTask.id
}

func (txTask *TxTask) Status() error {
	return txTask.err
}

func (txTask *TxTask) Data() []*types.Txs {
	return txTask.data
}

func (txTask *TxTask) Do() error {
	count, err := adaptor.GetTxCountInBlockFromServer(txTask.apiHost, txTask.ledger, txTask.height)
	if err != nil {
		txTask.err = err
		return err
	}

	if count == 0 {
		return nil
	}

	raw, _, err := adaptor.GetTxListInBlockRawFromServer(txTask.apiHost, txTask.ledger, txTask.height, 0, count)
	if err != nil {
		txTask.err = err
		return err
	}

	body := string(raw)
	success := gjson.Get(body, "success").Bool()
	if success == false {
		err := errors.New("response faild")
		txTask.err = err
		return err
	}
	result := gjson.Get(body, "data")

	txs := parseTransactions(result, txTask.ledger, txTask.height, 0)

	txTask.data = txs
	return nil
}

func parseTransactions(result gjson.Result, ledger string, blockHeight, startIndex int64) []*types.Txs {
	var txs []*types.Txs
	for _, txRaw := range result.Array() {
		tx := parseTransaction(txRaw, ledger, blockHeight, startIndex)
		txs = append(txs, tx)
		startIndex++
	}
	return txs
}

func parseTransaction(result gjson.Result, ledger string, blockHeight, startIndex int64) *types.Txs {

	tx := types.Txs{
		Ledger:            ledger,
		TxBlockHeight:     blockHeight,
		TxIndex:           int32(startIndex + 1),
		TxHash:            result.Get("result.transactionHash").String(),
		TxNodePubkeys:     "",
		TxEndpointPubkeys: "",
		TxContents:        "",
		TxResponseState:   0,
		TxResponseMsg:     result.Get("result.executionState").String(),
	}

	nodePubkeys, _ := json.Marshal(parsePublicKeys(result.Get("request.nodeSignatures")))
	endpointPubkeys, _ := json.Marshal(parsePublicKeys(result.Get("request.endpointSignatures")))

	tx.TxNodePubkeys = string(nodePubkeys)
	tx.TxEndpointPubkeys = string(endpointPubkeys)

	operations := result.Get("request.transactionContent.operations").Array()
	derivedOperations := result.Get("result.derivedOperations").Array()
	operations = append(operations, derivedOperations...)

	var rawOperations []string
	for _, operation := range operations {
		rawOperations = append(rawOperations, operation.String())
	}

	contents, _ := json.Marshal(rawOperations)
	tx.TxContents = string(contents)

	if tx.TxResponseMsg == "SUCCESS" {
		tx.TxResponseState = 1
	}

	return &tx
}

func parsePublicKeys(result gjson.Result) (list []string) {
	for _, nodeRaw := range result.Array() {
		list = append(list, nodeRaw.Get("pubKey").String())
	}
	return
}
