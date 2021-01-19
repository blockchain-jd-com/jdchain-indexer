package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

/*
 * Author: imuge
 * Date: 2020/11/3 下午5:57
 */

func NewTxEventAccountLink(tx, address string) *TxEventAccountLink {
	return &TxEventAccountLink{
		tx: tx, eventAccount: address,
	}
}

type TxEventAccountLink struct {
	tx              string
	txUID           string
	eventAccount    string
	eventAccountUID string
}

func (link *TxEventAccountLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.txUID),
			dgraph_helper.MutationItemUid(link.eventAccountUID),
			dgraph_helper.MutationPredict("tx-event_account"),
		),
	)
	return
}

func (link *TxEventAccountLink) SetUidLeft(uid string) {
	link.txUID = uid
}

func (link *TxEventAccountLink) SetUidRight(uid string) {
	link.eventAccountUID = uid
}

func (link *TxEventAccountLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "tx-hash_id"
	value = link.tx
	ok = true
	return
}

func (link *TxEventAccountLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "event_account-address"
	value = link.eventAccount
	ok = true
	return
}
