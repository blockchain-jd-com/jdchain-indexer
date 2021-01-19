package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

/*
 * Author: imuge
 * Date: 2020/11/3 下午5:57
 */

func NewTxContractLink(tx, address string) *TxContractLink {
	return &TxContractLink{
		tx: tx, contract: address,
	}
}

type TxContractLink struct {
	tx          string
	txUID       string
	contract    string
	contractUID string
}

func (link *TxContractLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.txUID),
			dgraph_helper.MutationItemUid(link.contractUID),
			dgraph_helper.MutationPredict("tx-contract"),
		),
	)
	return
}

func (link *TxContractLink) SetUidLeft(uid string) {
	link.txUID = uid
}

func (link *TxContractLink) SetUidRight(uid string) {
	link.contractUID = uid
}

func (link *TxContractLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "tx-hash_id"
	value = link.tx
	ok = true
	return
}

func (link *TxContractLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "contract-address"
	value = link.contract
	ok = true
	return
}
