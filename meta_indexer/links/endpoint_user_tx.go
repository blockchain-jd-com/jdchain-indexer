package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

/*
 * Author: imuge
 * Date: 2020/10/30 下午5:16
 */

func NewEndpointUserTxLink(userPublicKey, txHash string) *EndpointUserTxLink {
	return &EndpointUserTxLink{
		user: userPublicKey,
		tx:   txHash,
	}
}

// 用户-交易 关联关系
type EndpointUserTxLink struct {
	user    string
	userUID string
	tx      string
	txUID   string
}

func (link *EndpointUserTxLink) SetUidLeft(uid string) {
	link.userUID = uid
}

func (link *EndpointUserTxLink) SetUidRight(uid string) {
	link.txUID = uid
}

func (link *EndpointUserTxLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "user-public_key"
	value = link.user
	ok = true
	return
}

func (link *EndpointUserTxLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "tx-hash_id"
	value = link.tx
	ok = true
	return
}

func (link *EndpointUserTxLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.userUID),
			dgraph_helper.MutationItemUid(link.txUID),
			dgraph_helper.MutationPredict("endpoint_user-tx"),
		),
	)
	return
}
