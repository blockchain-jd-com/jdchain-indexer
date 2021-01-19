package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

func NewLedgerBlockLink(ledger, block string) *LedgerBlockLink {
	return &LedgerBlockLink{
		ledger: ledger,
		block:  block,
	}
}

type LedgerBlockLink struct {
	ledger    string
	ledgerUID string
	block     string
	blockUID  string
}

func (link *LedgerBlockLink) SetUidLeft(uid string) {
	link.ledgerUID = uid
}

func (link *LedgerBlockLink) SetUidRight(uid string) {
	link.blockUID = uid
}

func (link *LedgerBlockLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "ledger-hash_id"
	value = link.ledger
	ok = true
	return
}

func (link *LedgerBlockLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "block-hash_id"
	value = link.block
	ok = true
	return
}

func (link *LedgerBlockLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.ledgerUID),
			dgraph_helper.MutationItemUid(link.blockUID),
			dgraph_helper.MutationPredict("ledger-block"),
		),
	)
	return
}
