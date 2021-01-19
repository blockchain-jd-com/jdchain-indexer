package links

import (
	"git.jd.com/jd-blockchain/explorer/adaptor"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/worker"
	"github.com/davecgh/go-spew/spew"
	"github.com/ssor/zlog"
)

func ToCommonNode(ld *worker.LedgerData) (all []dgraph_helper.MutationData) {
	all = append(all, ld.BlockInfo)

	for _, tx := range ld.Txs {
		all = append(all, tx)
		for _, content := range tx.Contents {
			switch op := content.(type) {
			case *adaptor.KVSetOperation:
				all = append(all, op)
			case *adaptor.UserOperation:
				all = append(all, op)
			case *adaptor.DataAccountOperation:
				all = append(all, op)
			case *adaptor.ContractDeployOperation:
				all = append(all, op)
			case *adaptor.EventAccountOperation:
				all = append(all, op)
			case *adaptor.EventOperation:
				all = append(all, op)
			case *adaptor.ContractEventOperation:
				all = append(all, op)
			default:
				zlog.Warn("do not support content type now ")
				spew.Dump(content)
			}
		}
	}
	return
}

func ToLinks(ld *worker.LedgerData) (all []dgraph_helper.MutationData) {

	ledger := ld.Ledger
	blockInfo := ld.BlockInfo

	ledgerBlockLink := NewLedgerBlockLink(ledger, blockInfo.Hash)
	all = append(all, ledgerBlockLink)

	for _, tx := range ld.Txs {
		blockTxLink := NewBlockTxLink(blockInfo.Hash, tx.Hash)
		all = append(all, blockTxLink)

		for _, endpoint := range tx.EndpointPublicKey {
			all = append(all, NewEndpointUserTxLink(endpoint, tx.Hash))
		}
	}

	for _, tx := range ld.Txs {
		for _, content := range tx.Contents {
			switch op := content.(type) {
			case *adaptor.KVSetOperation:
				for _, history := range op.Histories {
					link := NewTxKVSetLink(tx.Hash, history.UniqueMutationName())
					all = append(all, link)
				}
			case *adaptor.UserOperation:
				link := NewLedgerUserLink(ledger, op.Address)
				all = append(all, link)
				tulink := NewTxUserLink(tx.Hash, op.Address)
				all = append(all, tulink)
			case *adaptor.DataAccountOperation:
				link := NewLedgerDatasetLink(ledger, op.Address)
				all = append(all, link)
				tdulink := NewTxDataAccountLink(tx.Hash, op.Address)
				all = append(all, tdulink)
			case *adaptor.ContractDeployOperation:
				link := NewLedgerContractLink(ledger, op.Address)
				all = append(all, link)
				tclink := NewTxContractLink(tx.Hash, op.Address)
				all = append(all, tclink)
			case *adaptor.ContractEventOperation:
				link := NewTxContractEventLink(tx.Hash, op.UniqueMutationName())
				all = append(all, link)
			case *adaptor.EventAccountOperation:
				link := NewLedgerEventAccountLink(ledger, op.Address)
				all = append(all, link)
				tealink := NewTxEventAccountLink(tx.Hash, op.Address)
				all = append(all, tealink)
			case *adaptor.EventOperation:
				for _, history := range op.Histories {
					link := NewTxEventLink(tx.Hash, history.UniqueMutationName())
					all = append(all, link)
				}
			default:
				zlog.Warn("do not support content type now ")
				spew.Dump(content)
			}
		}
	}
	return
}
