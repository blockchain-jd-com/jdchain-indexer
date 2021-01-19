package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

/*
 * Author: imuge
 * Date: 2020/11/5 上午11:04
 */

func NewTxKVSetLink(tx, kv string) *TxKVSetLink {
	return &TxKVSetLink{
		tx: tx, kv: kv,
	}
}

type TxKVSetLink struct {
	tx    string
	txUID string
	kv    string
	kvUID string
}

func (link *TxKVSetLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.txUID),
			dgraph_helper.MutationItemUid(link.kvUID),
			dgraph_helper.MutationPredict("tx-kv"),
		),
	)
	return
}

func (link *TxKVSetLink) SetUidLeft(uid string) {
	link.txUID = uid
}

func (link *TxKVSetLink) SetUidRight(uid string) {
	link.kvUID = uid
}

func (link *TxKVSetLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "tx-hash_id"
	value = link.tx
	ok = true
	return
}

func (link *TxKVSetLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "kv-tx_index"
	value = link.kv
	ok = true
	return
}
