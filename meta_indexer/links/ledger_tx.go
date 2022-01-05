package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

func NewLedgerTxLink(ledger, tx string) *LedgerTxLink {
	return &LedgerTxLink{
		ledger: ledger,
		tx:     tx,
	}
}

type LedgerTxLink struct {
	ledger    string
	ledgerUID string
	tx        string
	txUID     string
}

func (link *LedgerTxLink) SetUidLeft(uid string) {
	link.ledgerUID = uid
}

func (link *LedgerTxLink) SetUidRight(uid string) {
	link.txUID = uid
}

func (link *LedgerTxLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "ledger-hash_id"
	value = link.ledger
	ok = true
	return
}

func (link *LedgerTxLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "tx-hash_id"
	value = link.tx
	ok = true
	return
}

func (link *LedgerTxLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.ledgerUID),
			dgraph_helper.MutationItemUid(link.txUID),
			dgraph_helper.MutationPredict("ledger-tx"),
		),
	)
	return
}
