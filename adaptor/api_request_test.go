package adaptor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	ledgerDefault *Ledger
	block         *Block
	apiHostInTest string
)

func TestGetLedger(t *testing.T) {
	ledgers, err := GetLedgersFromServer(apiHostInTest)
	assert.Nil(t, err)
	assert.True(t, len(ledgers) > 0)
	ledgerDefault = ledgers[0]
	assert.Len(t, ledgerDefault.Hash, 45)
	assert.True(t, ledgerDefault.Height > 0)
	t.Logf("get ledger [%s] success", ledgerDefault.Hash)
}

func TestGetBlock(t *testing.T) {
	b, err := GetBlockFromServer(apiHostInTest, ledgerDefault.Hash, ledgerDefault.Height)
	assert.Nil(t, err)
	assert.NotNil(t, b)
	block = b
	assert.Equal(t, b.LedgerID, ledgerDefault.Hash)
	t.Logf("get block (%d) success", b.Height)
}

//
//func TestGetTxCount(t *testing.T) {
//    c, err := getTxCountInBlockFromServer(apiHost, ledgerDefault.Hash, ledgerDefault.Height)
//    assert.Nil(t, err)
//    assert.True(t, c > 0)
//    block.TxCount = c
//    t.Logf("from block[%d] got [%d] txs", ledgerDefault.Height, c)
//}

func TestGetTxList(t *testing.T) {
	txs, err := GetTxListInBlockFromServer(apiHostInTest, ledgerDefault.Hash, ledgerDefault.Height, 0, block.TxCount)
	assert.Nil(t, err)
	assert.NotNil(t, txs)
	assert.Len(t, txs, int(block.TxCount))
	block.txs = txs
}

func TestGetUserCount(t *testing.T) {
	c, err := getTotalUserCountInLedgerFromServer(apiHostInTest, ledgerDefault.Hash)
	assert.Nil(t, err)
	assert.True(t, c > 0)
	ledgerDefault.userCount = c
	t.Logf("got [%d] users in ledger[%s]", ledgerDefault.userCount, ledgerDefault.Hash)
}

func TestGetAccountCount(t *testing.T) {
	c, err := getTotalAccountCountInLedgerFromServer(apiHostInTest, ledgerDefault.Hash)
	assert.Nil(t, err)
	assert.True(t, c > 0)
	ledgerDefault.userCount = c
	t.Logf("got [%d] account in ledger[%s]", ledgerDefault.userCount, ledgerDefault.Hash)
}

func TestGetContractCount(t *testing.T) {
	c, err := getTotalContractCountInLedgerFromServer(apiHostInTest, ledgerDefault.Hash)
	assert.Nil(t, err)
	assert.True(t, c > 0)
	ledgerDefault.userCount = c
	t.Logf("got [%d] contract in ledger[%s]", ledgerDefault.userCount, ledgerDefault.Hash)
}

//
//func TestGetParticipants(t *testing.T) {
//	participants, err := getParticipants(apiHostInTest, ledgerDefault.Hash)
//	assert.Nil(t, err)
//	assert.True(t, len(participants) > 3)
//}
