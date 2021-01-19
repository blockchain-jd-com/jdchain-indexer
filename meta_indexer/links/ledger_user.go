package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

func NewLedgerUserLink(ledger, address string) *LedgerUserLink {
	return &LedgerUserLink{
		ledger: ledger,
		user:   address,
	}
}

type LedgerUserLink struct {
	ledger    string
	ledgerUID string
	user      string
	userUID   string
}

func (link *LedgerUserLink) SetUidLeft(uid string) {
	link.ledgerUID = uid
}

func (link *LedgerUserLink) SetUidRight(uid string) {
	link.userUID = uid
}

func (link *LedgerUserLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "ledger-hash_id"
	value = link.ledger
	ok = true
	return
}

func (link *LedgerUserLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "user-address"
	value = link.user
	ok = true
	return
}

func (link *LedgerUserLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.ledgerUID),
			dgraph_helper.MutationItemUid(link.userUID),
			dgraph_helper.MutationPredict("ledger-user"),
		),
	)
	return
}
