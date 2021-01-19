package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

/*
 * Author: imuge
 * Date: 2020/11/3 下午5:57
 */

func NewTxUserLink(tx, user string) *TxUserLink {
	return &TxUserLink{
		tx: tx, user: user,
	}
}

type TxUserLink struct {
	tx      string
	txUID   string
	user    string
	userUID string
}

func (link *TxUserLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.txUID),
			dgraph_helper.MutationItemUid(link.userUID),
			dgraph_helper.MutationPredict("tx-user"),
		),
	)
	return
}

func (link *TxUserLink) SetUidLeft(uid string) {
	link.txUID = uid
}

func (link *TxUserLink) SetUidRight(uid string) {
	link.userUID = uid
}

func (link *TxUserLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "tx-hash_id"
	value = link.tx
	ok = true
	return
}

func (link *TxUserLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "user-address"
	value = link.user
	ok = true
	return
}
