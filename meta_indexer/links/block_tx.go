package links

import "git.jd.com/jd-blockchain/explorer/dgraph_helper"

func NewBlockTxLink(block, tx string) *BlockTxLink {
	return &BlockTxLink{
		tx: tx, block: block,
	}
}

type BlockTxLink struct {
	tx       string
	txUID    string
	block    string
	blockUID string
}

func (link *BlockTxLink) Mutations() (mutations dgraph_helper.Mutations) {
	mutations = mutations.Add(
		dgraph_helper.NewMutation(
			dgraph_helper.MutationItemUid(link.blockUID),
			dgraph_helper.MutationItemUid(link.txUID),
			dgraph_helper.MutationPredict("block-tx"),
		),
	)
	return
}

func (link *BlockTxLink) SetUidLeft(uid string) {
	link.blockUID = uid
}

func (link *BlockTxLink) SetUidRight(uid string) {
	link.txUID = uid
}

func (link *BlockTxLink) LeftQueryBy() (predict, value string, ok bool) {
	predict = "block-hash_id"
	value = link.block
	ok = true
	return
}

func (link *BlockTxLink) RightQueryBy() (predict, value string, ok bool) {
	predict = "tx-hash_id"
	value = link.tx
	ok = true
	return
}
