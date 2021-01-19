package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

/*
 * Author: imuge
 * Date: 2020/11/5 下午5:56
 */

func NewTxContractEventLink(tx, event string) *TxContractEventLink {
	return &TxContractEventLink{
		tx: tx, event: event,
	}
}

type TxContractEventLink struct {
	tx       string
	txUID    string
	event    string
	eventUID string
}

func (link *TxContractEventLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.txUID),
			dgraph_helper.MutationItemUid(link.eventUID),
			dgraph_helper.MutationPredict("tx-contract_event"),
		),
	)
	return
}

func (link *TxContractEventLink) SetUidLeft(uid string) {
	link.txUID = uid
}

func (link *TxContractEventLink) SetUidRight(uid string) {
	link.eventUID = uid
}

func (link *TxContractEventLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "tx-hash_id"
	value = link.tx
	ok = true
	return
}

func (link *TxContractEventLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "contract_event-tx_index"
	value = link.event
	ok = true
	return
}
