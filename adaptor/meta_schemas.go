package adaptor

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
)

type ModelType string

const (
	ModelTypeLedger        ModelType = "ledger"         // 账本
	ModelTypeBlock         ModelType = "block"          // 区块
	ModelTypeTx            ModelType = "tx"             // 交易
	ModelTypeKV            ModelType = "kv"             // KV
	ModelTypeDataAccount   ModelType = "data_account"   // 数据账户
	ModelTypeContract      ModelType = "contract"       // 合约账户
	ModelTypeContractEvent ModelType = "contract_event" // 合约调用
	ModelTypeUser          ModelType = "user"           // 用户
	ModelTypeEventAccount  ModelType = "event_account"  // 事件账户
	ModelTypeEventSet      ModelType = "event"          // 事件集

	ModelTypeLedgerBlockLink    ModelType = "ledger-block"     // 账本-区块
	ModelTypeBlockTxLink        ModelType = "block-tx"         // 区块-交易
	ModelTypeLedgerTxLink       ModelType = "ledger-tx"        // 账本-交易
	ModelTypeTxKvLink           ModelType = "tx-kv"            // 交易-键值
	ModelTypeEndpointUserTxLink ModelType = "endpoint_user-tx" // 终端用户-交易
)

func PredictTo(t ModelType, name string) string {
	return fmt.Sprintf("%s-%s", t, name)
}

var (
	MetaSchemas = (dgraph_helper.Schemas{}).
		Add(dgraph_helper.NewSchemaIntIndex(PredictTo(ModelTypeBlock, "height"))).
		Add(dgraph_helper.NewSchemaIntIndex(PredictTo(ModelTypeBlock, "time"))).
		Add(dgraph_helper.NewSchemaStringTrigramIndex(PredictTo(ModelTypeBlock, "hash_id"))).
		Add(dgraph_helper.NewSchemaStringTrigramIndex(PredictTo(ModelTypeContract, "public_key"))).
		Add(dgraph_helper.NewSchemaStringTrigramIndex(PredictTo(ModelTypeContract, "address"))).
		Add(dgraph_helper.NewSchemaStringExactIndex(PredictTo(ModelTypeContractEvent, "contract_address"))).
		Add(dgraph_helper.NewSchemaStringTrigramIndex(PredictTo(ModelTypeContractEvent, "contract_event"))).
		Add(dgraph_helper.NewSchemaStringExactIndex(PredictTo(ModelTypeContractEvent, "tx_index"))).
		Add(dgraph_helper.NewSchemaStringTrigramIndex(PredictTo(ModelTypeDataAccount, "public_key"))).
		Add(dgraph_helper.NewSchemaStringTrigramIndex(PredictTo(ModelTypeDataAccount, "address"))).
		Add(dgraph_helper.NewSchemaStringTrigramIndex(PredictTo(ModelTypeLedger, "hash_id"))).
		Add(dgraph_helper.NewSchemaStringTrigramIndex(PredictTo(ModelTypeTx, "hash_id"))).
		Add(dgraph_helper.NewSchemaStringExactIndex(PredictTo(ModelTypeTx, "execution_state"))).
		Add(dgraph_helper.NewSchemaIntIndex(PredictTo(ModelTypeTx, "block_height"))).
		Add(dgraph_helper.NewSchemaIntIndex(PredictTo(ModelTypeTx, "time"))).
		Add(dgraph_helper.NewSchemaStringTrigramIndex(PredictTo(ModelTypeUser, "public_key"))).
		Add(dgraph_helper.NewSchemaStringTrigramIndex(PredictTo(ModelTypeUser, "address"))).
		Add(dgraph_helper.NewSchemaStringExactIndex(PredictTo(ModelTypeKV, "data_account_address"))).
		Add(dgraph_helper.NewSchemaIntIndex(PredictTo(ModelTypeKV, "version"))).
		Add(dgraph_helper.NewSchemaStringExactIndex(PredictTo(ModelTypeKV, "tx_index"))).
		Add(dgraph_helper.NewSchemaStringTrigramIndex(PredictTo(ModelTypeKV, "key"))).
		Add(dgraph_helper.NewSchemaStringTrigramIndex(PredictTo(ModelTypeEventAccount, "public_key"))).
		Add(dgraph_helper.NewSchemaStringTrigramIndex(PredictTo(ModelTypeEventAccount, "address"))).
		Add(dgraph_helper.NewSchemaStringExactIndex(PredictTo(ModelTypeEventSet, "tx_index"))).
		Add(dgraph_helper.NewSchemaStringExactIndex(PredictTo(ModelTypeEventSet, "event_account_address"))).
		Add(dgraph_helper.NewSchemaStringTrigramIndex(PredictTo(ModelTypeEventSet, "topic"))).
		Add(dgraph_helper.NewSchemaIntIndex(PredictTo(ModelTypeEventSet, "sequence"))).
		Add(dgraph_helper.NewSchemaUidsIndex(string(ModelTypeLedgerBlockLink))).
		Add(dgraph_helper.NewSchemaUidsIndex(string(ModelTypeLedgerTxLink))).
		Add(dgraph_helper.NewSchemaUidsIndex(string(ModelTypeBlockTxLink))).
		Add(dgraph_helper.NewSchemaUidsIndex(string(ModelTypeTxKvLink))).
		Add(dgraph_helper.NewSchemaUidsIndex(string(ModelTypeEndpointUserTxLink)))
)
