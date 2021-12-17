package tasks

import (
	"fmt"

	"git.jd.com/jd-blockchain/explorer/adaptor"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/app/rds_import/types"
)

type EventAccountTask struct {
	id       string
	apiHost  string
	ledger   string
	from     int64
	count    int64
	accounts []*types.EventAccounts
	events   []*types.EventAccountEvents
	err      error
}

func GetEventAccountCount(host, ledger string) (int64, error) {
	return adaptor.GetTotalEventAccountCountInLedgerFromServer(host, ledger)
}

func NewEventAccountTasks(host, ledger string, totalCount, taskSize int64) []*EventAccountTask {

	var eventAccountTasks []*EventAccountTask

	mod := totalCount % taskSize
	taskCount := totalCount / taskSize

	if mod != 0 {
		taskCount = taskCount + 1
	}

	for i := int64(1); i <= taskCount; i++ {
		from := (i - 1) * taskSize
		ct := &EventAccountTask{
			id:      fmt.Sprintf("event-account-task-%s-%d", ledger, i),
			apiHost: host,
			ledger:  ledger,
			from:    from,
			count:   taskSize,
		}
		eventAccountTasks = append(eventAccountTasks, ct)
	}

	return eventAccountTasks
}

func (eventAccountTask *EventAccountTask) ID() string {
	return eventAccountTask.id
}

func (eventAccountTask *EventAccountTask) Status() error {
	return eventAccountTask.err
}

func (eventAccountTask *EventAccountTask) Accounts() []*types.EventAccounts {
	return eventAccountTask.accounts
}

func (eventAccountTask *EventAccountTask) Events() []*types.EventAccountEvents {
	return eventAccountTask.events
}

func (eventAccountTask *EventAccountTask) Do() error {

	var accounts []*types.EventAccounts
	var events []*types.EventAccountEvents

	eventAccounts, err := adaptor.GetEventAccountListFromServer(eventAccountTask.apiHost, eventAccountTask.ledger, eventAccountTask.from, eventAccountTask.count)

	if err != nil {
		eventAccountTask.err = err
		return err
	}

	for _, eventAccount := range eventAccounts {

		eventAccountAddress := eventAccount.Get("address").String()

		info, err := adaptor.GetEventAccountInfoFromServer(eventAccountTask.apiHost, eventAccountTask.ledger, eventAccountAddress)
		if err != nil {
			eventAccountTask.err = err
			return err
		}

		role := info.Get("permission.role").String()
		if role == "" {
			role = `["DEFAULT"]`
		}

		eventAccount := types.EventAccounts{
			Ledger:                  eventAccountTask.ledger,
			EventAccountAddress:     eventAccountAddress,
			EventAccountPubkey:      info.Get("pubKey").String(),
			EventAccountRoles:       role,
			EventAccountPriviledges: info.Get("permission.modeBits").String(),
			EventAccountCreator:     info.Get("permission.owners").String(),
		}

		accounts = append(accounts, &eventAccount)

		eventCount, err := adaptor.GetTotalEventNameCountInLedgerFromServer(eventAccountTask.apiHost, eventAccountTask.ledger, eventAccountAddress)
		if err != nil {
			eventAccountTask.err = err
			return err
		}

		if eventCount == 0 {
			continue
		}

		eventNameList, err := adaptor.GetEventAccountEventNameListFromServer(eventAccountTask.apiHost, eventAccountTask.ledger, eventAccountAddress, 0, eventCount)
		if err != nil {
			eventAccountTask.err = err
			return err
		}

		for _, eventName := range eventNameList {

			name := eventName.Str
			eventInfos, err := adaptor.GetEventAccountEventNameInfoFromServer(eventAccountTask.apiHost, eventAccountTask.ledger, eventAccountAddress, name)
			if err != nil {
				eventAccountTask.err = err
				return err
			}

			for _, eventInfo := range eventInfos {
				event := types.EventAccountEvents{
					Ledger:               eventAccountTask.ledger,
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
	}

	eventAccountTask.accounts = accounts
	eventAccountTask.events = events
	return nil
}
