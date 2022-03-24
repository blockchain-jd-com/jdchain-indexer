package adaptor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/imkira/go-interpol"
	"github.com/tidwall/gjson"
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

func GetAccountsFromServer(apiHost, ledgerID string, from, count int64) ([]DataAccountOperation, error) {
	paras := map[string]string{
		"ledger": ledgerID,
		"from":   strconv.FormatInt(from, 10),
		"count":  strconv.FormatInt(count, 10),
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
		account.Address = raw.Get("address").String()
		account.PublicKey = raw.Get("pubKey").String()
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func GetContractsFromServer(apiHost, ledgerID string, from, count int64) ([]ContractDeployOperation, error) {
	paras := map[string]string{
		"ledger": ledgerID,
		"from":   strconv.FormatInt(from, 10),
		"count":  strconv.FormatInt(count, 10),
	}
	url := buildUrl(apiHost, apiGetContractAccountList, paras)
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
		contract.Address = raw.Get("address").String()
		contract.PublicKey = raw.Get("pubKey").String()
		contracts = append(contracts, contract)
	}

	return contracts, nil
}

func GetContractDetailFromServer(apiHost, ledgerId, address string) (interface{}, error) {
	paras := map[string]string{
		"ledger":  ledgerId,
		"address": address,
	}

	url := buildUrl(apiHost, apiGetContractDetail, paras)
	raw, err := fetchRequestRawResult(url)
	if err != nil {
		return nil, err
	}

	result, err := parseToJson(raw, url)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetUsersFromServer(apiHost, ledgerID string, from, count int64) ([]UserOperation, error) {
	paras := map[string]string{
		"ledger": ledgerID,
		"from":   strconv.FormatInt(from, 10),
		"count":  strconv.FormatInt(count, 10),
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
		user.Address = raw.Get("address").String()
		user.PublicKey = raw.Get("pubKey").String()
		users = append(users, user)
	}

	return users, nil
}

func GetUserInfoFromServer(apiHost, ledgerID, address string) (result gjson.Result, e error) {
	paras := map[string]string{
		"ledger":  ledgerID,
		"address": address,
	}
	url := buildUrl(apiHost, apiGetUserDetail, paras)
	return requestFromServer(url, paras)
}

func GetUserAuthorizationFromServer(apiHost, ledgerID, address string) (result gjson.Result, e error) {
	paras := map[string]string{
		"ledger":  ledgerID,
		"address": address,
	}
	url := buildUrl(apiHost, apiGetUserAuthorization, paras)
	return requestFromServer(url, paras)
}

func GetEventAccountListFromServer(apiHost, ledgerID string, from, count int64) (result []gjson.Result, e error) {
	paras := map[string]string{
		"ledger": ledgerID,
		"from":   strconv.FormatInt(from, 10),
		"count":  strconv.FormatInt(count, 10),
	}
	url := buildUrl(apiHost, apiGetEventAccountList, paras)
	r, err := requestFromServer(url, paras)
	if err != nil {
		return nil, err
	}

	return r.Array(), nil
}

func GetEventAccountInfoFromServer(apiHost, ledgerID, address string) (result gjson.Result, e error) {
	paras := map[string]string{
		"ledger":  ledgerID,
		"address": address,
	}
	url := buildUrl(apiHost, apiGetEventAccountInfo, paras)
	return requestFromServer(url, paras)
}

func GetEventAccountEventNameListFromServer(apiHost, ledgerID, address string, from, count int64) (result []gjson.Result, e error) {
	paras := map[string]string{
		"ledger":  ledgerID,
		"address": address,
		"from":    strconv.FormatInt(from, 10),
		"count":   strconv.FormatInt(count, 10),
	}
	url := buildUrl(apiHost, apiGetEventAccountEventNameList, paras)
	r, err := requestFromServer(url, paras)
	if err != nil {
		return nil, err
	}

	return r.Array(), nil
}

func GetEventAccountEventNameInfoFromServer(apiHost, ledgerID, address, eventName string) (result []gjson.Result, e error) {
	paras := map[string]string{
		"ledger":    ledgerID,
		"address":   address,
		"eventName": eventName,
	}
	url := buildUrl(apiHost, apiGetEventAccountEventNameInfo, paras)
	r, err := requestFromServer(url, paras)
	if err != nil {
		return nil, err
	}

	return r.Array(), nil
}

func GetAccountEntriesListFromServer(apiHost, ledgerID, address string, from, count int64) (result []gjson.Result, e error) {
	paras := map[string]string{
		"ledger":  ledgerID,
		"address": address,
		"from":    strconv.FormatInt(from, 10),
		"count":   strconv.FormatInt(count, 10),
	}
	url := buildUrl(apiHost, apiGetAccountEntriesList, paras)
	r, err := requestFromServer(url, paras)
	if err != nil {
		return nil, err
	}

	return r.Array(), nil
}

func GetAccountInfoFromServer(apiHost, ledgerID, address string) (result gjson.Result, e error) {
	paras := map[string]string{
		"ledger":  ledgerID,
		"address": address,
	}
	url := buildUrl(apiHost, apiGetAccountInfo, paras)
	return requestFromServer(url, paras)
}

func requestFromServer(url string, para map[string]string) (result gjson.Result, e error) {
	raw, err := fetchRequestRawResult(url)
	if err != nil {
		return gjson.Result{}, err
	}

	return parseToJson(raw, url)
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
		if operationsResult.Exists() {
			operations := operationsResult.Array()
			derivedOperations := result.Get("result.derivedOperations")
			if derivedOperations.Exists() {
				operations = append(operations, derivedOperations.Array()...)
			}
			tx.Contents = parseOperations(tx.Hash, operations)
		}
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

func parseOperations(txHash string, operations []gjson.Result) []interface{} {
	var contents []interface{}
	operationIndex := int64(0)
	for _, contentResult := range operations {
		operationIndex++
		opType := contentResult.Get("@type").String()
		switch opType {
		case "com.jd.blockchain.ledger.DataAccountKVSetOperation": // KV 写入
			ws := contentResult.Get("writeSet")
			address := contentResult.Get("accountAddress").String()
			wo := newWriteOperation()
			for _, historyResult := range ws.Array() {
				key := historyResult.Get("key").String()
				valueType := historyResult.Get("value.type").String()
				value := historyResult.Get("value.bytes").String()
				version := historyResult.Get("expectedVersion").Int()
				wo.addHistory(txHash, operationIndex, address, WriteValueType(valueType), key, value, version)
			}
			contents = append(contents, wo)
		case "com.jd.blockchain.ledger.UserRegisterOperation": // 用户注册
			uo := newUserOperation(
				contentResult.Get("userID.address").String(),
				contentResult.Get("userID.pubKey").String(),
			)
			contents = append(contents, uo)
		case "com.jd.blockchain.ledger.DataAccountRegisterOperation": // 数据账户注册
			uo := newDataAccountOperation(
				contentResult.Get("accountID.address").String(),
				contentResult.Get("accountID.pubKey").String(),
			)
			contents = append(contents, uo)
		case "com.jd.blockchain.ledger.ContractCodeDeployOperation": // 合约部署 TODO 合约升级操作
			uo := newContractDeployOperation(
				contentResult.Get("contractID.address").String(),
				contentResult.Get("contractID.pubKey").String(),
				contentResult.Get("chainCodeVersion").Int(),
				contentResult.Get("lang").String())
			contents = append(contents, uo)
		case "com.jd.blockchain.ledger.ContractEventSendOperation": // 合约调用
			uo := newContractEventOperation(
				txHash,
				operationIndex,
				contentResult.Get("contractAddress").String(),
				contentResult.Get("version").Int(),
				contentResult.Get("event").String(),
				contentResult.Get("args").Raw)
			contents = append(contents, uo)
		case "com.jd.blockchain.ledger.EventAccountRegisterOperation": // 事件账户注册
			uo := newEventAccountOperation(
				contentResult.Get("eventAccountID.address").String(),
				contentResult.Get("eventAccountID.pubKey").String(),
			)
			contents = append(contents, uo)
		case "com.jd.blockchain.ledger.EventPublishOperation": // 事件发布
			events := contentResult.Get("events")
			eventAddress := contentResult.Get("eventAddress").String()
			eo := newEventOperation()
			for _, historyResult := range events.Array() {
				topic := historyResult.Get("name").String()
				valueType := historyResult.Get("content.type").String()
				content := historyResult.Get("content.bytes").String()
				sequence := historyResult.Get("sequence").Int()
				eo.addHistory(txHash, operationIndex, eventAddress, WriteValueType(valueType), topic, content, sequence)
			}
			contents = append(contents, eo)
		}
	}
	return contents
}

func GetTotalContractCountInLedgerFromServer(apiHost, ledgerID string) (int64, error) {
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

func GetTotalAccountCountInLedgerFromServer(apiHost, ledgerID string) (int64, error) {
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

func GetTotalUserCountInLedgerFromServer(apiHost, ledgerID string) (int64, error) {
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

func GetTotalEventAccountCountInLedgerFromServer(apiHost, ledgerID string) (int64, error) {
	paras := map[string]string{
		"ledger": ledgerID,
	}
	url := buildUrl(apiHost, apiGetTotalEventAccountCount, paras)
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

func GetTotalAccountEntriesCountFromServer(apiHost, ledgerID, address string) (int64, error) {
	paras := map[string]string{
		"ledger":  ledgerID,
		"address": address,
	}
	url := buildUrl(apiHost, apiGetAccountEntriesCount, paras)
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

func GetTotalEventNameCountInLedgerFromServer(apiHost, ledgerID, address string) (int64, error) {
	paras := map[string]string{
		"ledger":  ledgerID,
		"address": address,
	}
	url := buildUrl(apiHost, apiGetEventAccountEventNameCount, paras)
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

func GetTxCountInBlockFromServer(apiHost, ledgerID string, height int64) (int64, error) {

	paras := map[string]string{
		"ledger": ledgerID,
		"height": strconv.FormatInt(height, 10),
	}

	url := buildUrl(apiHost, apiGetTxCountOfBlock, paras)
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
	block.ContractAccountSetHash = result.Get("contractAccountSetHash").String()
	block.DataAccountSetHash = result.Get("dataAccountSetHash").String()
	block.UserEventSetHash = result.Get("userEventSetHash").String()

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
