package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

/*
 * Author: imuge
 * Date: 2020/11/3 下午5:57
 */

func NewTxDataAccountLink(tx, address string) *TxDataAccountLink {
	return &TxDataAccountLink{
		tx: tx, dataAccount: address,
	}
}

type TxDataAccountLink struct {
	tx             string
	txUID          string
	dataAccount    string
	dataAccountUID string
}

func (link *TxDataAccountLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.txUID),
			dgraph_helper.MutationItemUid(link.dataAccountUID),
			dgraph_helper.MutationPredict("tx-data_account"),
		),
	)
	return
}

func (link *TxDataAccountLink) SetUidLeft(uid string) {
	link.txUID = uid
}

func (link *TxDataAccountLink) SetUidRight(uid string) {
	link.dataAccountUID = uid
}

func (link *TxDataAccountLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "tx-hash_id"
	value = link.tx
	ok = true
	return
}

func (link *TxDataAccountLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "data_account-address"
	value = link.dataAccount
	ok = true
	return
}
