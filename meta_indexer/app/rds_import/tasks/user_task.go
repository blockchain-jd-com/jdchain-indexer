package tasks

import (
	"encoding/json"
	"fmt"

	"git.jd.com/jd-blockchain/explorer/adaptor"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/app/rds_import/types"
)

type UserTask struct {
	id      string
	apiHost string
	ledger  string
	from    int64
	count   int64
	data    []*types.Users
	err     error
}

func GetUsersCount(host, ledger string) (int64, error) {
	return adaptor.GetTotalUserCountInLedgerFromServer(host, ledger)
}

func NewUserTasks(host, ledger string, totalCount, taskSize int64) []*UserTask {

	var userTasks []*UserTask

	mod := totalCount % taskSize
	taskCount := totalCount / taskSize

	if mod != 0 {
		taskCount = taskCount + 1
	}

	for i := int64(1); i <= taskCount; i++ {
		from := (i - 1) * taskSize
		ct := &UserTask{
			id:      fmt.Sprintf("user-task-%s-%d", ledger, i),
			apiHost: host,
			ledger:  ledger,
			from:    from,
			count:   taskSize,
		}
		userTasks = append(userTasks, ct)
	}

	return userTasks
}

func (userTask *UserTask) ID() string {
	return userTask.id
}

func (userTask *UserTask) Status() error {
	return userTask.err
}

func (userTask *UserTask) Data() []*types.Users {
	return userTask.data
}

func (userTask *UserTask) Do() error {

	users, err := adaptor.GetUsersFromServer(userTask.apiHost, userTask.ledger, userTask.from, userTask.count)
	if err != nil {
		userTask.err = err
		return err
	}

	var typeUsers []*types.Users

	for _, user := range users {

		typeUser := types.Users{
			Ledger:           userTask.ledger,
			UserAddress:      user.Address,
			UserPubkey:       user.PublicKey,
			UserKeyAlgorithm: "ED25519",
			UserState:        "",
			Roles:            "",
			Privileges:       "",
		}

		userInfo, err := adaptor.GetUserInfoFromServer(userTask.apiHost, userTask.ledger, user.Address)
		if err != nil {
			userTask.err = err
			return err
		}

		typeUser.UserState = userInfo.Get("state").String()

		userAuth, err := adaptor.GetUserAuthorizationFromServer(userTask.apiHost, userTask.ledger, user.Address)
		if err != nil {
			userTask.err = err
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

		typeUsers = append(typeUsers, &typeUser)
	}

	userTask.data = typeUsers
	return nil
}
