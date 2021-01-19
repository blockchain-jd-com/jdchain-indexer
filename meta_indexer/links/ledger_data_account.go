package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

func NewLedgerDatasetLink(ledger, address string) *LedgerDatasetLink {
	return &LedgerDatasetLink{
		ledger:      ledger,
		dataAccount: address,
	}
}

type LedgerDatasetLink struct {
	ledger         string
	ledgerUID      string
	dataAccount    string
	dataAccountUID string
}

func (link *LedgerDatasetLink) SetUidLeft(uid string) {
	link.ledgerUID = uid
}

func (link *LedgerDatasetLink) SetUidRight(uid string) {
	link.dataAccountUID = uid
}

func (link *LedgerDatasetLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "ledger-hash_id"
	value = link.ledger
	ok = true
	return
}

func (link *LedgerDatasetLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "data_account-address"
	value = link.dataAccount
	ok = true
	return
}

func (link *LedgerDatasetLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.ledgerUID),
			dgraph_helper.MutationItemUid(link.dataAccountUID),
			dgraph_helper.MutationPredict("ledger-data_account"),
		),
	)
	return
}
