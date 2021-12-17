package cmd

import (
	"testing"

	"git.jd.com/jd-blockchain/explorer/adaptor"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/app/rds_import/tasks"
)

func Test_startServer(t *testing.T) {
	type args struct {
		cmd *ImportArgs
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startServer(tt.args.cmd)
		})
	}
}
func Test_userTask(t *testing.T) {

	host := "http://127.0.0.1:8080"
	ledger := "j5nfkEvfyHidqk9MHJZGFZxVbLBfy23M4TQwcqP6fFewkF"

	userCount, err := adaptor.GetTotalUserCountInLedgerFromServer(host, ledger)
	if err == nil {
		userTasks := tasks.NewUserTasks(host, ledger, userCount, 10)
		for _, userTask := range userTasks {
			userTask.Do()
		}
	}

}

func Test_eventTask(t *testing.T) {

	host := "http://127.0.0.1:8080"
	ledger := "j5nfkEvfyHidqk9MHJZGFZxVbLBfy23M4TQwcqP6fFewkF"

	eventAccount, err := adaptor.GetTotalEventAccountCountInLedgerFromServer(host, ledger)
	if err == nil {
		eventTasks := tasks.NewEventAccountTasks(host, ledger, eventAccount, 10)
		for _, eventTask := range eventTasks {
			eventTask.Do()
		}
	}

}

func Test_accountTask(t *testing.T) {

	host := "http://127.0.0.1:8080"
	ledger := "j5nfkEvfyHidqk9MHJZGFZxVbLBfy23M4TQwcqP6fFewkF"

	dataAccountCount, err := adaptor.GetTotalAccountCountInLedgerFromServer(host, ledger)
	if err == nil {
		dataAccountTasks := tasks.NewDataAccountTasks(host, ledger, dataAccountCount, 10)
		for _, dataAccountTask := range dataAccountTasks {
			dataAccountTask.Do()
		}
	}

}
