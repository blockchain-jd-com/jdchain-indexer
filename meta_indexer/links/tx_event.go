package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

/*
 * Author: imuge
 * Date: 2020/11/5 上午11:04
 */

func NewTxEventLink(tx, event string) *TxEventLink {
	return &TxEventLink{
		tx: tx, event: event,
	}
}

type TxEventLink struct {
	tx       string
	txUID    string
	event    string
	eventUID string
}

func (link *TxEventLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.txUID),
			dgraph_helper.MutationItemUid(link.eventUID),
			dgraph_helper.MutationPredict("tx-event"),
		),
	)
	return
}

func (link *TxEventLink) SetUidLeft(uid string) {
	link.txUID = uid
}

func (link *TxEventLink) SetUidRight(uid string) {
	link.eventUID = uid
}

func (link *TxEventLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "tx-hash_id"
	value = link.tx
	ok = true
	return
}

func (link *TxEventLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "event-tx_index"
	value = link.event
	ok = true
	return
}
