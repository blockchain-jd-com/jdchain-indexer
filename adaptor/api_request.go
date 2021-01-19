package adaptor

import (
	"fmt"
	"github.com/imkira/go-interpol"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var (
	errStatusCodeNotOK = fmt.Errorf("status code is not 200")
	errResponseFailed  = fmt.Errorf("server response failed")
)

func buildUrl(host, apiName string, paras map[string]string) string {
	api, err := interpol.WithMap(apiList[apiName], paras)
	if err != nil {
		logger.Failedf("build url failed for %s", err)
		return ""
	}

	return host + api
}

func GetTxListInBlockRawFromServer(apiHost, ledgerID string, height, from, count int64) (raw []byte, debugInfo interface{}, e error) {
	url := buildUrl(apiHost, apiGetTxListOfBlock, map[string]string{
		"ledger": ledgerID,
		"height": strconv.FormatInt(height, 10),
		"from":   strconv.FormatInt(from, 10),
		"count":  strconv.FormatInt(count, 10),
	})
	debugInfo = url

	resp, err := http.Get(url)
	if err != nil {
		e = err
		logger.Failedf("failed getting response from [%s] for %s", url, err)
		return
	}
	statusCode := resp.StatusCode
	if http.StatusOK != statusCode {
		e = errStatusCodeNotOK
		return
	}
	defer resp.Body.Close()
	raw, e = ioutil.ReadAll(resp.Body)
	return
}

func GetAccountsFromServer(apiHost, ledgerID string) ([]DataAccountOperation, error) {
	paras := map[string]string{
		"ledger": ledgerID,
	}
	url := buildUrl(apiHost, apiGetAccountList, paras)
	raw, err := fetchRequestRawResult(url)
	if err != nil {
		return nil, err
	}

	result, err := parseToJson(raw, url)
	if err != nil {
		return nil, err
	}
	var accounts []DataAccountOperation
	for _, raw := range result.Array() {
		var account DataAccountOperation
		account.Address = raw.Get("address.value").String()
		account.PublicKey = raw.Get("pubKey").String()
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func GetContractsFromServer(apiHost, ledgerID string) ([]ContractDeployOperation, error) {
	paras := map[string]string{
		"ledger": ledgerID,
	}
	url := buildUrl(apiHost, apiGetContractList, paras)
	raw, err := fetchRequestRawResult(url)
	if err != nil {
		return nil, err
	}

	result, err := parseToJson(raw, url)
	if err != nil {
		return nil, err
	}
	var contracts []ContractDeployOperation
	for _, raw := range result.Array() {
		var contract ContractDeployOperation
		contract.Address = raw.Get("address.value").String()
		contract.PublicKey = raw.Get("pubKey").String()
		contracts = append(contracts, contract)
	}

	return contracts, nil
}

func GetUsersFromServer(apiHost, ledgerID string) ([]UserOperation, error) {
	paras := map[string]string{
		"ledger": ledgerID,
	}
	url := buildUrl(apiHost, apiGetUserList, paras)
	raw, err := fetchRequestRawResult(url)
	if err != nil {
		return nil, err
	}

	result, err := parseToJson(raw, url)
	if err != nil {
		return nil, err
	}
	var users []UserOperation
	for _, raw := range result.Array() {
		var user UserOperation
		user.Address = raw.Get("address.value").String()
		user.PublicKey = raw.Get("pubKey").String()
		users = append(users, user)
	}

	return users, nil
}

func GetTxListInBlockFromServer(apiHost, ledgerID string, height, from, count int64) ([]*Transaction, error) {
	paras := map[string]string{
		"ledger": ledgerID,
		"height": strconv.FormatInt(height, 10),
		"from":   strconv.FormatInt(from, 10),
		"count":  strconv.FormatInt(count, 10),
	}
	url := buildUrl(apiHost, apiGetTxListOfBlock, paras)
	start := time.Now()
	raw, err := fetchRequestRawResult(url)
	if err != nil {
		return nil, err
	}
	logGetRawData(height, start, time.Now())

	start = time.Now()
	result, err := parseToJson(raw, url)
	if err != nil {
		return nil, err
	}

	txs, err := parseTransactions(result, height, from), nil
	logParseRawData(height, int64(len(txs)), start, time.Now())

	return txs, err
}

func parseTransactions(result gjson.Result, blockHeight, startIndex int64) []*Transaction {
	var txs []*Transaction
	for _, txRaw := range result.Array() {
		tx := parseTransaction(txRaw, blockHeight, startIndex)
		txs = append(txs, tx)
		startIndex++
	}
	return txs
}

func parseTransaction(result gjson.Result, blockHeight, startIndex int64) *Transaction {
	var tx Transaction

	txHash := result.Get("result.transactionHash").String()
	tx.Hash = txHash
	tx.IndexInBlock = startIndex
	tx.ExecutionState = result.Get("result.executionState").String()
	tx.Time = result.Get("request.transactionContent.timestamp").Int()
	// 仅记录成功交易的操作内容
	if tx.IsSuccess() {
		operationsResult := result.Get("request.transactionContent.operations")
		tx.Contents = parseOperations(tx.Hash, operationsResult)
	}
	tx.NodePublicKey = parsePublicKeys(result.Get("request.nodeSignatures"))
	tx.EndpointPublicKey = parsePublicKeys(result.Get("request.endpointSignatures"))
	tx.BlockHeight = blockHeight
	return &tx
}

func parsePublicKeys(result gjson.Result) (list []string) {
	for _, nodeRaw := range result.Array() {
		list = append(list, nodeRaw.Get("pubKey").String())
	}
	return
}

func parseOperations(txHash string, result gjson.Result) []interface{} {
	var contents []interface{}
	operationIndex := int64(0)
	for _, contentResult := range result.Array() {
		operationIndex++
		ws := contentResult.Get("writeSet")
		if ws.Exists() {
			address := contentResult.Get("accountAddress.value").String()
			wo := newWriteOperation()
			for _, historyResult := range ws.Array() {
				key := historyResult.Get("key").String()
				valueType := historyResult.Get("value.type").String()
				value := historyResult.Get("value.bytes.value").String()
				version := historyResult.Get("expectedVersion").Int()
				wo.addHistory(txHash, operationIndex, address, WriteValueType(valueType), key, value, version)
			}
			contents = append(contents, wo)
			continue
		}

		userID := contentResult.Get("userID")
		if userID.Exists() {
			uo := newUserOperation(
				contentResult.Get("userID.address.value").String(),
				contentResult.Get("userID.pubKey").String(),
			)
			contents = append(contents, uo)
			continue
		}

		accountID := contentResult.Get("accountID")
		if accountID.Exists() {
			uo := newDataAccountOperation(
				contentResult.Get("accountID.address.value").String(),
				contentResult.Get("accountID.pubKey").String(),
			)
			contents = append(contents, uo)
			continue
		}

		contractID := contentResult.Get("contractID")
		if contractID.Exists() {
			uo := newContractDeployOperation(
				contentResult.Get("contractID.address.value").String(),
				contentResult.Get("contractID.pubKey").String(),
				contentResult.Get("chainCodeVersion").Int())
			contents = append(contents, uo)
			continue
		}

		contractAddress := contentResult.Get("contractAddress")
		if contractAddress.Exists() {
			uo := newContractEventOperation(
				txHash,
				operationIndex,
				contentResult.Get("contractAddress.value").String(),
				contentResult.Get("version").Int(),
				contentResult.Get("event").String(),
				contentResult.Get("args").Raw)
			contents = append(contents, uo)
			continue
		}

		eventAccountID := contentResult.Get("eventAccountID")
		if eventAccountID.Exists() {
			uo := newEventAccountOperation(
				contentResult.Get("eventAccountID.address.value").String(),
				contentResult.Get("eventAccountID.pubKey").String(),
			)
			contents = append(contents, uo)
			continue
		}

		events := contentResult.Get("events")
		if events.Exists() {
			eventAddress := contentResult.Get("eventAddress.value").String()
			eo := newEventOperation()
			for _, historyResult := range events.Array() {
				topic := historyResult.Get("name").String()
				valueType := historyResult.Get("content.type").String()
				content := historyResult.Get("content.bytes.value").String()
				sequence := historyResult.Get("sequence").Int()
				eo.addHistory(txHash, operationIndex, eventAddress, WriteValueType(valueType), topic, content, sequence)
			}
			contents = append(contents, eo)
			continue
		}
	}
	return contents
}

func getTotalContractCountInLedgerFromServer(apiHost, ledgerID string) (int64, error) {
	paras := map[string]string{
		"ledger": ledgerID,
	}
	url := buildUrl(apiHost, apiGetTotalContractCount, paras)
	raw, err := fetchRequestRawResult(url)
	if err != nil {
		return -1, err
	}
	result, err := parseToJson(raw, url)
	if err != nil {
		return -1, err
	}
	count := result.Int()
	return count, nil
}

func getTotalAccountCountInLedgerFromServer(apiHost, ledgerID string) (int64, error) {
	paras := map[string]string{
		"ledger": ledgerID,
	}
	url := buildUrl(apiHost, apiGetTotalAccountCount, paras)
	raw, err := fetchRequestRawResult(url)
	if err != nil {
		return -1, err
	}
	result, err := parseToJson(raw, url)
	if err != nil {
		return -1, err
	}
	count := result.Int()
	return count, nil
}

func getTotalUserCountInLedgerFromServer(apiHost, ledgerID string) (int64, error) {
	//result, _, err := doRequest(apiHost, apiGetUserTotalCount, map[string]string{
	//    "ledger": ledgerID,
	//})
	paras := map[string]string{
		"ledger": ledgerID,
	}
	url := buildUrl(apiHost, apiGetUserTotalCount, paras)
	raw, err := fetchRequestRawResult(url)
	if err != nil {
		return -1, err
	}
	result, err := parseToJson(raw, url)
	if err != nil {
		return -1, err
	}
	count := result.Int()
	return count, nil
}

//
//func getTxCountInBlockFromServer(apiHost, ledgerID string, Height int64) (int64, error) {
//    result, _, err := doRequest(apiHost, apiGetTxCountOfBlock, map[string]string{
//        "ledger": ledgerID,
//        "Height": strconv.FormatInt(Height, 10),
//    })
//    if err != nil {
//        return -1, err
//    }
//    count := result.Int()
//    return count, nil
//}

func GetBlockFromServer(apiHost, ledgerID string, height int64) (b *Block, e error) {
	paras := map[string]string{
		"ledger": ledgerID,
		"height": strconv.FormatInt(height, 10),
	}
	url := buildUrl(apiHost, apiGetBlockDetail, paras)
	raw, err := fetchRequestRawResult(url)
	if err != nil {
		e = err
		return
	}
	result, err := parseToJson(raw, url)
	if err != nil {
		e = err
		return
	}

	//spew.Dump(result.Summary())

	var block Block
	blockLedgerHash := result.Get("ledgerHash")
	if height == 0 && blockLedgerHash.Raw == "" {
		block.LedgerID = ledgerID
	} else {
		block.LedgerID = blockLedgerHash.String()
	}
	block.Hash = result.Get("hash").String()
	block.Height = result.Get("height").Int()
	block.Time = result.Get("timestamp").Int()
	block.PreviousHash = result.Get("previousHash").String()
	block.TransactionSetHash = result.Get("transactionSetHash").String()
	block.UserAccountSetHash = result.Get("userAccountSetHash").String()
	block.AdminAccountHash = result.Get("adminAccountHash").String()

	if block.LedgerID != ledgerID || block.Height != height {
		e = fmt.Errorf("request for block  is not equal with response [%d -> %d] [%s -> %s]",
			height, block.Height, ledgerID, block.LedgerID)
		return
	}
	b = &block
	return
}

func GetLedgersFromServer(apiHost string) (ledgers []*Ledger, e error) {
	//result, _, err := doRequest(apiHost, apiGetLedgers, nil)
	url := buildUrl(apiHost, apiGetLedgers, nil)
	raw, err := fetchRequestRawResult(url)
	if err != nil {
		e = err
		return
	}
	result, err := parseToJson(raw, url)
	if err != nil {
		e = err
		return
	}

	var ledgerIDs []string
	for _, v := range result.Array() {
		ledgerIDs = append(ledgerIDs, v.String())
	}

	for _, id := range ledgerIDs {
		ledger, err := GetLedgerDetailFromServer(apiHost, id)
		if err != nil {
			logger.Errorf("GetLedgerDetailFromServer failed: %s", err)
			e = err
			return
		}
		ledgers = append(ledgers, ledger)
	}
	return
}

//
//func getParticipants(apiHost, ledgerID string) (participants []*model.Admin, e error) {
//    //result, debugInfo, err := doRequest(apiHost, apiGetParticipants, map[string]string{"ledger": ledgerID})
//    paras := map[string]string{"ledger": ledgerID}
//    url := buildUrl(apiHost, apiGetParticipants, paras)
//    raw, err := fetchRequestRawResult(url)
//    if err != nil {
//        e = err
//        return
//    }
//    result, err := parseToJson(raw, url)
//    if err != nil {
//        e = err
//        return
//    }
//    for _, raw := range result.Array() {
//        var admin model.Admin
//        admin.PublicKey = model.PublicKey(raw.Get("pubKey").String())
//        admin.Name = raw.Get("name").String()
//        participants = append(participants, &admin)
//    }
//    if len(participants) < 4 {
//        logger.Errorf("no participants for ledger %s(%s)", ledgerID, url)
//    }
//    return
//}

func GetLedgerDetailFromServer(apiHost, ledgerID string) (ledger *Ledger, e error) {
	paras := map[string]string{"ledger": ledgerID}
	url := buildUrl(apiHost, apiGetLedgerDetail, paras)
	//result, _, err := doRequest(apiHost, apiGetLedgerDetail, map[string]string{"ledger": ledgerID})
	raw, err := fetchRequestRawResult(url)
	if err != nil {
		logger.Failedf("failed to get ledger[%s] detail for: %s", ledgerID, err)
		e = err
		return
	}
	result, err := parseToJson(raw, url)
	if err != nil {
		e = err
		return
	}
	newLedgerHash := result.Get("hash").String()
	if newLedgerHash != ledgerID {
		e = fmt.Errorf("ledger request is not equal with response")
		return
	}
	latestBlockHash := result.Get("latestBlockHash").String()

	latestBlockHeight := result.Get("latestBlockHeight").Int()
	ledger = &Ledger{
		Hash:      ledgerID,
		Height:    latestBlockHeight,
		BlockHash: latestBlockHash,
	}
	return
}

func fetchRequestRawResult(url string) (raw []byte, e error) {
	//url := buildUrl(apiHost, apiName, paras)
	resp, err := http.Get(url)
	if err != nil {
		e = err
		logger.Failedf("failed getting response from [%s] for %s", url, err)
		return
	}
	statusCode := resp.StatusCode
	if http.StatusOK != statusCode {
		e = errStatusCodeNotOK
		logger.Failedf("failed getting response from [%s] for status code %d", url, statusCode)
		return
	}
	defer resp.Body.Close()
	raw, e = ioutil.ReadAll(resp.Body)
	return
}

func parseToJson(raw []byte, tip interface{}) (result gjson.Result, e error) {
	body := string(raw)
	success := gjson.Get(body, "success").Bool()
	if success == false {
		logger.Failedf("failed getting response for body failed: %s", tip)
		e = errResponseFailed
		return
	}
	result = gjson.Get(body, "data")
	return
}
