package tasks

import (
	"encoding/json"
	"fmt"
	"time"

	"git.jd.com/jd-blockchain/explorer/adaptor"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/app/rds_import/types"
	"github.com/tidwall/gjson"
)

type AddressInfoTask struct {
	id             string
	apiHost        string
	ledger         string
	address        string
	addressType    int8 // 1 user  2 dataaccount 3 eventaccount 4 contract
	UserData       *types.Users
	ContractData   *types.Contracts
	DataAccount    *types.DataAccounts
	DataAccountKV  []*types.DataAccountKVS
	EventAccount   *types.EventAccounts
	EventAccountKV []*types.EventAccountEvents
	err            error
}

func NewAddressInfoTask(host, ledger string, address string, addressType int8) *AddressInfoTask {

	return &AddressInfoTask{
		id:          fmt.Sprintf("address-info-task-%s-%s-%d", ledger, address, addressType),
		apiHost:     host,
		ledger:      ledger,
		address:     address,
		addressType: addressType,
	}
}

func (addressInfoTask *AddressInfoTask) ID() string {
	return addressInfoTask.id
}

func (addressInfoTask *AddressInfoTask) Status() error {
	return addressInfoTask.err
}

func (addressInfoTask *AddressInfoTask) Data() []*types.Users {
	return nil
}

func (addressInfoTask *AddressInfoTask) Do() error {

	if addressInfoTask.addressType == 1 {
		return processUserInfo(addressInfoTask.apiHost, addressInfoTask.ledger, addressInfoTask.address, addressInfoTask)
	}

	if addressInfoTask.addressType == 2 {
		return processDataAccountInfo(addressInfoTask.apiHost, addressInfoTask.ledger, addressInfoTask.address, addressInfoTask)
	}

	if addressInfoTask.addressType == 3 {
		return processEventAccountInfo(addressInfoTask.apiHost, addressInfoTask.ledger, addressInfoTask.address, addressInfoTask)
	}

	if addressInfoTask.addressType == 4 {
		return processContractInfo(addressInfoTask.apiHost, addressInfoTask.ledger, addressInfoTask.address, addressInfoTask)
	}

	return nil
}

func processContractInfo(host, ledger, address string, addressInfoTask *AddressInfoTask) error {
	result, err := adaptor.GetContractDetailFromServer(host, ledger, address)
	if err != nil {
		addressInfoTask.err = err
		return err
	}

	gResult := result.(gjson.Result)

	role := gResult.Get("permission.role").String()
	if role == "" {
		role = `["DEFAULT"]`
	}

	addressInfoTask.ContractData = &types.Contracts{
		Ledger:              ledger,
		ContractAddress:     address,
		ContractPubkey:      gResult.Get("pubKey").String(),
		ContractRoles:       role,
		ContractPriviledges: gResult.Get("permission.modeBits").String(),
		ContractVersion:     int8(gResult.Get("chainCodeVersion").Int()),
		ContractStatus:      gResult.Get("state").String(),
		ContractCreator:     gResult.Get("permission.owners").String(),
		ContractContent:     gResult.Get("chainCode").String(),
	}

	return nil
}

func processEventAccountInfo(host, ledger, address string, addressInfoTask *AddressInfoTask) error {

	eventAccountAddress := address

	info, err := adaptor.GetEventAccountInfoFromServer(host, ledger, eventAccountAddress)
	if err != nil {
		addressInfoTask.err = err
		return err
	}

	role := info.Get("permission.role").String()
	if role == "" {
		role = `["DEFAULT"]`
	}

	eventAccount := types.EventAccounts{
		Ledger:                  ledger,
		EventAccountAddress:     eventAccountAddress,
		EventAccountPubkey:      info.Get("iD.pubKey").String(),
		EventAccountRoles:       role,
		EventAccountPriviledges: info.Get("permission.modeBits").String(),
		EventAccountCreator:     info.Get("permission.owners").String(),
	}

	addressInfoTask.EventAccount = &eventAccount

	eventCount, err := adaptor.GetTotalEventNameCountInLedgerFromServer(host, ledger, eventAccountAddress)
	if err != nil {
		addressInfoTask.err = err
		return err
	}

	if eventCount == 0 {
		return nil
	}

	eventNameList, err := adaptor.GetEventAccountEventNameListFromServer(host, ledger, eventAccountAddress, 0, eventCount)
	if err != nil {
		addressInfoTask.err = err
		return err
	}

	var events []*types.EventAccountEvents

	for _, eventName := range eventNameList {

		name := eventName.Str
		eventInfos, err := adaptor.GetEventAccountEventNameInfoFromServer(host, ledger, eventAccountAddress, name)
		if err != nil {
			addressInfoTask.err = err
			return err
		}

		for _, eventInfo := range eventInfos {
			event := types.EventAccountEvents{
				Ledger:               ledger,
				EventAccountAddress:  eventAccountAddress,
				EventName:            name,
				EventSequence:        int32(eventInfo.Get("sequence").Int()),
				EventTxHash:          eventInfo.Get("transactionSource").String(),
				EventBlockHeight:     eventInfo.Get("blockHeight").Int(),
				EventType:            eventInfo.Get("content.type").String(),
				EventValue:           eventInfo.Get("content.bytes").String(),
				EventContractAddress: eventInfo.Get("contractSource").String(),
			}

			events = append(events, &event)
		}
	}

	addressInfoTask.EventAccountKV = events
	return nil
}

func processDataAccountInfo(host, ledger, address string, addressInfoTask *AddressInfoTask) error {

	var kvs []*types.DataAccountKVS

	dataAccountInfo, err := adaptor.GetAccountInfoFromServer(host, ledger, address)
	if err != nil {
		addressInfoTask.err = err
		return err
	}

	role := dataAccountInfo.Get("permission.role").String()
	if role == "" {
		role = `["DEFAULT"]`
	}

	dataAccounts := types.DataAccounts{
		Ledger:                ledger,
		DataAccountAddress:    address,
		DataAccountPubkey:     dataAccountInfo.Get("iD.pubKey").String(),
		DataAccountRoles:      role,
		DataAccountPrivileges: dataAccountInfo.Get("permission.modeBits").String(),
		DataAccountCreator:    dataAccountInfo.Get("permission.owners").String(),
	}

	addressInfoTask.DataAccount = &dataAccounts

	entriesCount, err := adaptor.GetTotalAccountEntriesCountFromServer(host, ledger, address)
	if err != nil {
		addressInfoTask.err = err
		return err
	}

	if entriesCount == 0 {
		return nil
	}

	pageSize := int64(100)
	loops := calPage(entriesCount, pageSize)
	for i := int64(1); i <= loops; i++ {
		accountEntriesList, err := adaptor.GetAccountEntriesListFromServer(host, ledger, address, (i-1)*pageSize, pageSize)
		if err != nil {
			addressInfoTask.err = err
			return err
		}

		for _, accountEntry := range accountEntriesList {

			kv := types.DataAccountKVS{
				Ledger:             ledger,
				DataAccountAddress: address,
				DataAccountKey:     accountEntry.Get("key").String(),
				DataAccountValue:   accountEntry.Get("value").String(),
				DataAccountType:    accountEntry.Get("type").String(),
				DataAccountVersion: int(accountEntry.Get("version").Int()),
				CreateTime:         time.Now(),
				UpdateTime:         time.Now(),
				State:              1,
			}

			kvs = append(kvs, &kv)
		}

	}

	addressInfoTask.DataAccountKV = kvs
	return nil
}

func processUserInfo(host, ledger, address string, addressInfoTask *AddressInfoTask) error {
	typeUser := types.Users{
		Ledger:           ledger,
		UserAddress:      address,
		UserPubkey:       "",
		UserKeyAlgorithm: "ED25519",
		UserState:        "",
		Roles:            "",
		Privileges:       "",
	}

	userInfo, err := adaptor.GetUserInfoFromServer(host, ledger, address)
	if err != nil {
		addressInfoTask.err = err
		return err
	}

	typeUser.UserPubkey = userInfo.Get("iD.pubKey").String()
	typeUser.UserState = userInfo.Get("state").String()

	userAuth, err := adaptor.GetUserAuthorizationFromServer(host, ledger, address)
	if err != nil {
		addressInfoTask.err = err
		return err
	}

	typeUser.Roles = userAuth.Get("userRole").String()

	var transactionPrivileges []string
	var ledgerPrivileges []string

	for _, result := range userAuth.Get("transactionPrivilegesBitset.privilege").Array() {
		transactionPrivileges = append(transactionPrivileges, result.String())
	}

	for _, result := range userAuth.Get("ledgerPrivilegesBitset.privilege").Array() {
		ledgerPrivileges = append(ledgerPrivileges, result.String())
	}

	privileges, _ := json.Marshal(map[string][]string{
		"transactionPrivileges": transactionPrivileges,
		"ledgerPrivileges":      ledgerPrivileges,
	})
	typeUser.Privileges = string(privileges)

	addressInfoTask.UserData = &typeUser
	return nil
}
