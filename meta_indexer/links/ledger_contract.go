package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

func NewLedgerContractLink(ledger, address string) *LedgerContractLink {
	return &LedgerContractLink{
		ledger:   ledger,
		contract: address,
	}
}

type LedgerContractLink struct {
	ledger      string
	ledgerUID   string
	contract    string
	contractUID string
}

func (link *LedgerContractLink) SetUidLeft(uid string) {
	link.ledgerUID = uid
}

func (link *LedgerContractLink) SetUidRight(uid string) {
	link.contractUID = uid
}

func (link *LedgerContractLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "ledger-hash_id"
	value = link.ledger
	ok = true
	return
}

func (link *LedgerContractLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "contract-address"
	value = link.contract
	ok = true
	return
}

func (link *LedgerContractLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.ledgerUID),
			dgraph_helper.MutationItemUid(link.contractUID),
			dgraph_helper.MutationPredict("ledger-contract"),
		),
	)
	return
}
