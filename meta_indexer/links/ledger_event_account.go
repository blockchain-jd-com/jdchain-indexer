package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

func NewLedgerEventAccountLink(ledger, address string) *LedgerEventAccountLink {
	return &LedgerEventAccountLink{
		ledger:       ledger,
		eventAccount: address,
	}
}

type LedgerEventAccountLink struct {
	ledger          string
	ledgerUID       string
	eventAccount    string
	eventAccountUID string
}

func (link *LedgerEventAccountLink) SetUidLeft(uid string) {
	link.ledgerUID = uid
}

func (link *LedgerEventAccountLink) SetUidRight(uid string) {
	link.eventAccountUID = uid
}

func (link *LedgerEventAccountLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "ledger-hash_id"
	value = link.ledger
	ok = true
	return
}

func (link *LedgerEventAccountLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "event_account-address"
	value = link.eventAccount
	ok = true
	return
}

func (link *LedgerEventAccountLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.ledgerUID),
			dgraph_helper.MutationItemUid(link.eventAccountUID),
			dgraph_helper.MutationPredict("ledger-event_account"),
		),
	)
	return
}
