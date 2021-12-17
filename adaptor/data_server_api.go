package adaptor

var (
	//apiHost = ""

	apiGetLedgers                    = "/ledgers"
	apiGetLedgerDetail               = "/ledgers/{ledger}"
	apiGetParticipants               = "/ledgers/{ledger}/participants"
	apiGetBlockDetail                = "/ledgers/{ledger}/blocks/height/{height}"
	apiGetTxListOfBlock              = "/ledgers/{ledger}/blocks/height/{height}/txs/additional-txs?fromIndex={from}&count={count}"
	apiGetUserTotalCount             = "/ledgers/{ledger}/users/count"
	apiGetUserList                   = "/ledgers/{ledger}/users?fromIndex={from}&count={count}"
	apiGetTotalAccountCount          = "/ledgers/{ledger}/accounts/count"
	apiGetAccountList                = "/ledgers/{ledger}/accounts?fromIndex={from}&count={count}"
	apiGetAccountInfo                = "/ledgers/{ledger}/accounts/address/{address}"
	apiGetAccountEntriesCount        = "/ledgers/{ledger}/accounts/address/{address}/entries/count"
	apiGetAccountEntriesList         = "/ledgers/{ledger}/accounts/address/{address}/entries?fromIndex={from}&count={count}"
	apiGetTotalContractCount         = "/ledgers/{ledger}/contracts/count"
	apiGetContractAccountList        = "/ledgers/{ledger}/contracts?fromIndex={from}&count={count}"
	apiGetContractDetail             = "/ledgers/{ledger}/contracts/address/{address}"
	apiGetTxCountOfBlock             = "/ledgers/{ledger}/blocks/height/{height}/txs/additional-count"
	apiGetUserAuthorization          = "/ledgers/{ledger}/authorization/user/{address}"
	apiGetUserDetail                 = "/ledgers/{ledger}/users/address/{address}"
	apiGetTotalEventAccountCount     = "/ledgers/{ledger}/events/user/accounts/count"
	apiGetEventAccountList           = "/ledgers/{ledger}/events/user/accounts?fromIndex={from}&count={count}"
	apiGetEventAccountInfo           = "/ledgers/{ledger}/events/user/accounts/{address}"
	apiGetEventAccountEventNameCount = "/ledgers/{ledger}/events/user/accounts/{address}/names/count"
	apiGetEventAccountEventNameList  = "/ledgers/{ledger}/events/user/accounts/{address}/names?fromIndex={from}&count={count}"
	apiGetEventAccountEventNameInfo  = "/ledgers/{ledger}/events/user/accounts/{address}/names/{eventName}"

	apiList = map[string]string{
		apiGetLedgers:                    apiGetLedgers,
		apiGetLedgerDetail:               apiGetLedgerDetail,
		apiGetParticipants:               apiGetParticipants,
		apiGetBlockDetail:                apiGetBlockDetail,
		apiGetTxListOfBlock:              apiGetTxListOfBlock,
		apiGetUserTotalCount:             apiGetUserTotalCount,
		apiGetUserList:                   apiGetUserList,
		apiGetTotalAccountCount:          apiGetTotalAccountCount,
		apiGetAccountList:                apiGetAccountList,
		apiGetTotalContractCount:         apiGetTotalContractCount,
		apiGetContractAccountList:        apiGetContractAccountList,
		apiGetContractDetail:             apiGetContractDetail,
		apiGetTxCountOfBlock:             apiGetTxCountOfBlock,
		apiGetUserDetail:                 apiGetUserDetail,
		apiGetUserAuthorization:          apiGetUserAuthorization,
		apiGetAccountInfo:                apiGetAccountInfo,
		apiGetAccountEntriesCount:        apiGetAccountEntriesCount,
		apiGetAccountEntriesList:         apiGetAccountEntriesList,
		apiGetEventAccountList:           apiGetEventAccountList,
		apiGetEventAccountInfo:           apiGetEventAccountInfo,
		apiGetEventAccountEventNameCount: apiGetEventAccountEventNameCount,
		apiGetEventAccountEventNameList:  apiGetEventAccountEventNameList,
		apiGetTotalEventAccountCount:     apiGetTotalEventAccountCount,
		apiGetEventAccountEventNameInfo:  apiGetEventAccountEventNameInfo,
	}
)
