package adaptor

import (
	"fmt"
	"strconv"

	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
)

const Success = "SUCCESS"

type Ledger struct {
	Uid               string
	Hash              string
	BlockHash         string
	Height            int64
	userCount         int64
	contractCount     int64
	accountCount      int64
	eventAccountCount int64
}

func NewLedger(hash string) *Ledger {
	return &Ledger{
		Hash: hash,
	}
}

func (l *Ledger) copy(src *Ledger) {
	l.Hash = src.Hash
	l.BlockHash = src.BlockHash
	l.Height = src.Height
	l.userCount = src.userCount
	l.contractCount = src.contractCount
	l.accountCount = src.accountCount
	l.eventAccountCount = src.eventAccountCount
}

func (l *Ledger) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(ModelTypeLedger),
			dgraph_helper.MutationItemValue(l.Hash),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeLedger, "hash_id")),
		),
	)

	return
}

func (ledger *Ledger) QueryBy() (string, string) {
	return PredictTo(ModelTypeLedger, "hash_id"), ledger.Hash
}

type Block struct {
	Hash                   string
	Height                 int64
	Time                   int64
	LedgerID               string
	PreviousHash           string
	TransactionSetHash     string
	UserAccountSetHash     string
	AdminAccountHash       string
	ContractAccountSetHash string
	DataAccountSetHash     string
	UserEventSetHash       string
	TxCount                int64
	txs                    []*Transaction
}

func (block *Block) UniqueMutationName() string {
	return fmt.Sprintf("%s-%s", ModelTypeBlock, block.Hash)
}

func (block *Block) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(block.UniqueMutationName()),
			dgraph_helper.MutationItemValue(block.Hash),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeBlock, "hash_id")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(block.UniqueMutationName()),
			dgraph_helper.MutationItemValue(strconv.FormatInt(block.Height, 10)),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeBlock, "height")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(block.UniqueMutationName()),
			dgraph_helper.MutationItemValue(strconv.FormatInt(block.Time, 10)),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeBlock, "time")),
		),
	)
	return
}

type Transaction struct {
	Hash              string
	IndexInBlock      int64
	BlockHeight       int64
	ExecutionState    string
	Time              int64
	NodePublicKey     []string
	EndpointPublicKey []string
	Contents          []interface{}
}

func (tx Transaction) IsSuccess() bool {
	return tx.ExecutionState == Success
}

func (tx *Transaction) UniqueMutationName() string {
	return fmt.Sprintf("%s-%s", ModelTypeTx, tx.Hash)
}

func (tx *Transaction) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(tx.UniqueMutationName()),
			dgraph_helper.MutationItemValue(tx.Hash),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeTx, "hash_id")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(tx.UniqueMutationName()),
			dgraph_helper.MutationItemValue(tx.ExecutionState),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeTx, "execution_state")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(tx.UniqueMutationName()),
			dgraph_helper.MutationItemValue(strconv.FormatInt(tx.BlockHeight, 10)),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeTx, "block_height")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(tx.UniqueMutationName()),
			dgraph_helper.MutationItemValue(strconv.FormatInt(tx.Time, 10)),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeTx, "time")),
		),
	)
	return
}

func newContractDeployOperation(address, pubKey string, version int64, lang string) *ContractDeployOperation {
	if len(lang) == 0 {
		lang = "Java"
	}
	return &ContractDeployOperation{
		Address:   address,
		PublicKey: pubKey,
		Version:   version,
		Lang:      lang,
	}
}

type ContractDeployOperation struct {
	Address   string
	PublicKey string
	Version   int64
	Lang   	  string
}

func (contract *ContractDeployOperation) UniqueMutationName() string {
	return fmt.Sprintf("%s-%s-%d", ModelTypeContract, contract.Address, contract.Version)
}

func (contract *ContractDeployOperation) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(contract.UniqueMutationName()),
			dgraph_helper.MutationItemValue(contract.Address),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeContract, "address")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(contract.UniqueMutationName()),
			dgraph_helper.MutationItemValue(contract.PublicKey),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeContract, "public_key")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(contract.UniqueMutationName()),
			dgraph_helper.MutationItemValue(strconv.FormatInt(contract.Version, 10)),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeContract, "version")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(contract.UniqueMutationName()),
			dgraph_helper.MutationItemValue(contract.Lang),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeContract, "lang")),
		),
	)
	return
}

func newDataAccountOperation(address, pubKey string) *DataAccountOperation {
	return &DataAccountOperation{
		Address:   address,
		PublicKey: pubKey,
	}
}

type DataAccountOperation struct {
	Address   string
	PublicKey string
}

func (ds *DataAccountOperation) UniqueMutationName() string {
	return fmt.Sprintf("%s-%s", ModelTypeDataAccount, ds.Address)
}

func (ds *DataAccountOperation) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(ds.UniqueMutationName()),
			dgraph_helper.MutationItemValue(ds.Address),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeDataAccount, "address")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(ds.UniqueMutationName()),
			dgraph_helper.MutationItemValue(ds.PublicKey),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeDataAccount, "public_key")),
		),
	)

	return
}

func newUserOperation(address, pubKey string) *UserOperation {
	return &UserOperation{
		Address:   address,
		PublicKey: pubKey,
	}
}

type UserOperation struct {
	Address   string
	PublicKey string
}

func (user *UserOperation) UniqueMutationName() string {
	return fmt.Sprintf("%s-%s", ModelTypeUser, user.Address)
}

func (user *UserOperation) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(user.UniqueMutationName()),
			dgraph_helper.MutationItemValue(user.Address),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeUser, "address")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(user.UniqueMutationName()),
			dgraph_helper.MutationItemValue(user.PublicKey),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeUser, "public_key")),
		),
	)
	return
}

func newWriteOperation() *KVSetOperation {
	return &KVSetOperation{
		Histories: []*WriteHistory{},
	}
}

type WriteValueType string

const (
	WriteTEXT    = "TEXT"
	WriteBYTES   = "BYTES"
	WriteJSON    = "JSON"
	WriteINT64   = "INT64"
	WriteINT16   = "INT16"
	WriteINT32   = "INT32"
	WriteBoolean = "BOOLEAN"
)

func newWriteHistory(txHash string, operationIndex int64, accountAddress string, valueType WriteValueType, key, value string, version int64) *WriteHistory {
	return &WriteHistory{
		TxHash:         txHash,
		OperationIndex: operationIndex,
		DataSetAddress: accountAddress,
		Key:            key,
		Value:          value,
		Version:        version,
		Type:           valueType,
	}
}

type WriteHistory struct {
	TxHash         string
	OperationIndex int64
	DataSetAddress string
	Type           WriteValueType
	Key            string
	Value          string
	Version        int64
}

func (wh *WriteHistory) UniqueMutationName() string {
	return fmt.Sprintf("%s-%s-%d", ModelTypeKV, wh.TxHash, wh.OperationIndex)
}

func (wh *WriteHistory) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(wh.UniqueMutationName()),
			dgraph_helper.MutationItemValue(wh.DataSetAddress),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeKV, "data_account_address")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(wh.UniqueMutationName()),
			dgraph_helper.MutationItemValue(wh.Key),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeKV, "key")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(wh.UniqueMutationName()),
			dgraph_helper.MutationItemValue(strconv.FormatInt(wh.Version, 10)),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeKV, "version")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(wh.UniqueMutationName()),
			dgraph_helper.MutationItemValue(wh.UniqueMutationName()),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeKV, "tx_index")),
		),
	)
	return
}

type KVSetOperation struct {
	Histories []*WriteHistory
}

func (wo *KVSetOperation) addHistory(txHash string, operationIndex int64, accountAddress string, valueType WriteValueType, key, value string, version int64) {
	wo.Histories = append(wo.Histories, newWriteHistory(txHash, operationIndex, accountAddress, valueType, key, value, version))
}

func (wo *KVSetOperation) Mutations() (mutations dgraph_helper.Mutations) {
	for _, history := range wo.Histories {
		mutations = mutations.Add(history.Mutations()...)
	}
	return
}

// event account
func newEventAccountOperation(address, pubKey string) *EventAccountOperation {
	return &EventAccountOperation{
		Address:   address,
		PublicKey: pubKey,
	}
}

type EventAccountOperation struct {
	Address   string
	PublicKey string
}

func (es *EventAccountOperation) UniqueMutationName() string {
	return fmt.Sprintf("%s-%s", ModelTypeEventAccount, es.Address)
}

func (es *EventAccountOperation) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(es.UniqueMutationName()),
			dgraph_helper.MutationItemValue(es.Address),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeEventAccount, "address")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(es.UniqueMutationName()),
			dgraph_helper.MutationItemValue(es.PublicKey),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeEventAccount, "public_key")),
		),
	)

	return
}

func newEventHistory(txHash string, operationIndex int64, eventAddress string, valueType WriteValueType, topic, content string, version int64) *EventHistory {
	return &EventHistory{
		TxHash:         txHash,
		OperationIndex: operationIndex,
		EventAddress:   eventAddress,
		Topic:          topic,
		Content:        content,
		Sequence:       version,
		Type:           valueType,
	}
}

type EventHistory struct {
	TxHash         string
	OperationIndex int64
	EventAddress   string
	Type           WriteValueType
	Topic          string
	Content        string
	Sequence       int64
}

func (eh *EventHistory) UniqueMutationName() string {
	return fmt.Sprintf("%s-%s-%d", ModelTypeEventSet, eh.TxHash, eh.OperationIndex)
}

func (eh *EventHistory) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(eh.UniqueMutationName()),
			dgraph_helper.MutationItemValue(eh.Topic),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeEventSet, "event_account_address")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(eh.UniqueMutationName()),
			dgraph_helper.MutationItemValue(eh.Topic),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeEventSet, "topic")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(eh.UniqueMutationName()),
			dgraph_helper.MutationItemValue(strconv.FormatInt(eh.Sequence, 10)),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeEventSet, "sequence")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(eh.UniqueMutationName()),
			dgraph_helper.MutationItemValue(eh.UniqueMutationName()),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeEventSet, "tx_index")),
		),
	)
	return
}

// event dataSet
func newEventOperation() *EventOperation {
	return &EventOperation{
		Histories: []*EventHistory{},
	}
}

type EventOperation struct {
	Histories []*EventHistory
}

func (wo *EventOperation) addHistory(txHash string, operationIndex int64, accountAddress string, valueType WriteValueType, topic, content string, sequence int64) {
	wo.Histories = append(wo.Histories, newEventHistory(txHash, operationIndex, accountAddress, valueType, topic, content, sequence))
}

func (wo *EventOperation) Mutations() (mutations dgraph_helper.Mutations) {
	for _, history := range wo.Histories {
		mutations = mutations.Add(history.Mutations()...)
	}
	return
}

// contract call
func newContractEventOperation(txHash string, operationIndex int64, contractAddress string, version int64, event, args string) *ContractEventOperation {
	return &ContractEventOperation{
		TxHash:          txHash,
		OperationIndex:  operationIndex,
		ContractAddress: contractAddress,
		Version:         version,
		Args:            args,
		Event:           event,
	}
}

type ContractEventOperation struct {
	TxHash          string
	OperationIndex  int64
	ContractAddress string
	Version         int64
	Event           string
	Args            string
}

func (ceo *ContractEventOperation) UniqueMutationName() string {
	return fmt.Sprintf("%s-%s-%d", ModelTypeContractEvent, ceo.TxHash, ceo.OperationIndex)
}

func (ceo *ContractEventOperation) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(ceo.UniqueMutationName()),
			dgraph_helper.MutationItemValue(ceo.Event),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeContractEvent, "contract_address")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(ceo.UniqueMutationName()),
			dgraph_helper.MutationItemValue(ceo.Event),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeContractEvent, "event")),
		),
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemEmpty(ceo.UniqueMutationName()),
			dgraph_helper.MutationItemValue(ceo.UniqueMutationName()),
			dgraph_helper.MutationPredict(PredictTo(ModelTypeContractEvent, "tx_index")),
		),
	)
	return
}
