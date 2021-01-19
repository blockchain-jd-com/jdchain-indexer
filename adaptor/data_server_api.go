package adaptor

var (
	//apiHost = ""

	apiGetLedgers            = "/ledgers"
	apiGetLedgerDetail       = "/ledgers/{ledger}"
	apiGetParticipants       = "/ledgers/{ledger}/participants"
	apiGetBlockDetail        = "/ledgers/{ledger}/blocks/height/{height}"
	apiGetTxListOfBlock      = "/ledgers/{ledger}/blocks/height/{height}/txs/additional-txs?fromIndex={from}&count={count}"
	apiGetUserTotalCount     = "/ledgers/{ledger}/users/count"
	apiGetUserList           = "/ledgers/{ledger}/users"
	apiGetTotalAccountCount  = "/ledgers/{ledger}/accounts/count"
	apiGetAccountList        = "/ledgers/{ledger}/accounts"
	apiGetTotalContractCount = "/ledgers/{ledger}/contracts/count"
	apiGetContractList       = "/ledgers/{ledger}/contracts"
	//apiGetTxCountOfBlock     = "/ledgers/{ledger}/blocks/Height/{Height}/txs/additional-count"

	apiList = map[string]string{
		apiGetLedgers:            apiGetLedgers,
		apiGetLedgerDetail:       apiGetLedgerDetail,
		apiGetParticipants:       apiGetParticipants,
		apiGetBlockDetail:        apiGetBlockDetail,
		apiGetTxListOfBlock:      apiGetTxListOfBlock,
		apiGetUserTotalCount:     apiGetUserTotalCount,
		apiGetUserList:           apiGetUserList,
		apiGetTotalAccountCount:  apiGetTotalAccountCount,
		apiGetAccountList:        apiGetAccountList,
		apiGetTotalContractCount: apiGetTotalContractCount,
		apiGetContractList:       apiGetContractList,
		//apiGetTxCountOfBlock:     apiGetTxCountOfBlock,
	}
)
