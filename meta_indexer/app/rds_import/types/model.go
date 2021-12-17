package types

import "time"

func (Blocks) TableName() string {
	return "jdchain_blocks"
}

func (Contracts) TableName() string {
	return "jdchain_contracts"
}

func (DataAccounts) TableName() string {
	return "jdchain_data_accounts"
}

func (DataAccountKVS) TableName() string {
	return "jdchain_data_account_kvs"
}

func (EventAccounts) TableName() string {
	return "jdchain_event_accounts"
}

func (EventAccountEvents) TableName() string {
	return "jdchain_event_account_events"
}

func (Users) TableName() string {
	return "jdchain_users"
}

func (Txs) TableName() string {
	return "jdchain_txs"
}

// 区块信息表
type Blocks struct {
	Id                    int64     `gorm:"column:id"`                               //主键
	Ledger                string    `gorm:"column:ledger"`                           //账本
	BlockHeight           int64     `gorm:"column:block_height"`                     //区块高度
	BlockHash             string    `gorm:"column:block_hash"`                       //区块HASH
	PreBlockHash          string    `gorm:"column:pre_block_hash"`                   //前置区块HASH
	TxsSetHash            string    `gorm:"column:txs_set_hash"`                     //交易集HASH
	UsersSetHash          string    `gorm:"column:users_set_hash"`                   //用户集HASH
	ContractsSetHash      string    `gorm:"column:contracts_set_hash"`               //合约集HASH
	ConfigurationsSetHash string    `gorm:"column:configurations_set_hash"`          //配置集HASH
	DataAccountsSetHash   string    `gorm:"column:dataaccounts_set_hash"`            //数据账户HASH
	EventAccountsSetHash  string    `gorm:"column:eventaccounts_set_hash"`           //事件账户HASH
	BlockTimestamp        time.Time `gorm:"column:block_timestamp"`                  //区块创建时间
	CreateTime            time.Time `gorm:"column:create_time;autoCreateTime:milli"` //创建时间
	UpdateTime            time.Time `gorm:"column:update_time;autoUpdateTime:milli"` //更新时间
	State                 int8      `gorm:"column:state;default:1"`                  //数据状态
}

type Contracts struct {
	Id                  int64     `gorm:"column:id"`                               //主键
	Ledger              string    `gorm:"column:ledger"`                           //账本
	ContractAddress     string    `gorm:"column:contract_address"`                 //合约地址
	ContractPubkey      string    `gorm:"column:contract_pubkey"`                  //合约公钥
	ContractRoles       string    `gorm:"column:contract_roles"`                   //合约归属角色
	ContractPriviledges string    `gorm:"column:contract_priviledges"`             //合约权限
	ContractVersion     int8      `gorm:"column:contract_version"`                 //合约版本
	ContractStatus      string    `gorm:"column:contract_status"`                  //合约状态
	ContractCreator     string    `gorm:"column:contract_creator"`                 //合约创建者地址
	ContractContent     string    `gorm:"column:contract_content"`                 //合约内容
	CreateTime          time.Time `gorm:"column:create_time;autoCreateTime:milli"` //创建时间
	UpdateTime          time.Time `gorm:"column:update_time;autoUpdateTime:milli"` //更新时间
	State               int8      `gorm:"column:state;default:1"`                  //数据状态
}

type DataAccounts struct {
	Id                    int64     `gorm:"column:id"`                               //主键
	Ledger                string    `gorm:"column:ledger"`                           //账本
	DataAccountAddress    string    `gorm:"column:data_account_address"`             //数据账户地址
	DataAccountPubkey     string    `gorm:"column:data_account_pubkey"`              //数据账户公钥
	DataAccountRoles      string    `gorm:"column:data_account_roles"`               //数据账户角色
	DataAccountPrivileges string    `gorm:"column:data_account_privileges"`          //数据账户权限
	DataAccountCreator    string    `gorm:"column:data_account_creator"`             //数据账户创建者
	CreateTime            time.Time `gorm:"column:create_time;autoCreateTime:milli"` //创建时间
	UpdateTime            time.Time `gorm:"column:update_time;autoUpdateTime:milli"` //更新时间
	State                 int8      `gorm:"column:state;default:1"`                  //数据状态
}

type DataAccountKVS struct {
	Id                 int64     `gorm:"column:id"`                               //主键
	Ledger             string    `gorm:"column:ledger"`                           //账本
	DataAccountAddress string    `gorm:"column:data_account_address"`             //数据账户地址
	DataAccountKey     string    `gorm:"column:data_account_key"`                 //数据账户key
	DataAccountValue   string    `gorm:"column:data_account_value"`               //数据账户value
	DataAccountType    string    `gorm:"column:data_account_type"`                //数据账户类型
	DataAccountVersion int       `gorm:"column:data_account_version"`             //数据账户版本
	CreateTime         time.Time `gorm:"column:create_time;autoCreateTime:milli"` //创建时间
	UpdateTime         time.Time `gorm:"column:update_time;autoUpdateTime:milli"` //更新时间
	State              int8      `gorm:"column:state;default:1"`                  //数据状态
}

type EventAccounts struct {
	Id                      int64     `gorm:"column:id"`                               //主键
	Ledger                  string    `gorm:"column:ledger"`                           //账本
	EventAccountAddress     string    `gorm:"column:event_account_address"`            //事件账户地址
	EventAccountPubkey      string    `gorm:"column:event_account_pubkey"`             //事件账户公钥
	EventAccountRoles       string    `gorm:"column:event_account_roles"`              //事件账户归属角色
	EventAccountPriviledges string    `gorm:"column:event_account_priviledges"`        //事件账户权限
	EventAccountCreator     string    `gorm:"column:event_account_creator"`            //事件账户创建者
	CreateTime              time.Time `gorm:"column:create_time;autoCreateTime:milli"` //创建时间
	UpdateTime              time.Time `gorm:"column:update_time;autoUpdateTime:milli"` //更新时间
	State                   int8      `gorm:"column:state;default:1"`                  //数据状态
}

type EventAccountEvents struct {
	Id                   int64     `gorm:"column:id"`                               //主键
	Ledger               string    `gorm:"column:ledger"`                           //账本
	EventAccountAddress  string    `gorm:"column:event_account_address"`            //事件账户地址
	EventName            string    `gorm:"column:event_name"`                       //事件名称
	EventSequence        int32     `gorm:"column:event_sequence"`                   //事件序列
	EventTxHash          string    `gorm:"column:event_tx_hash"`                    //事件交易HASH
	EventBlockHeight     int64     `gorm:"column:event_block_height"`               //事件高度
	EventType            string    `gorm:"column:event_type"`                       //事件类型
	EventValue           string    `gorm:"column:event_value"`                      //事件值
	EventContractAddress string    `gorm:"column:event_contract_address"`           //事件合约地址
	CreateTime           time.Time `gorm:"column:create_time;autoCreateTime:milli"` //创建时间
	UpdateTime           time.Time `gorm:"column:update_time;autoUpdateTime:milli"` //更新时间
	State                int8      `gorm:"column:state;default:1"`                  //数据状态
}

type Users struct {
	Id               int64     `gorm:"column:id"`                               //主键
	Ledger           string    `gorm:"column:ledger"`                           //账本
	UserAddress      string    `gorm:"column:user_address"`                     //用户地址
	UserPubkey       string    `gorm:"column:user_pubkey"`                      //用户公钥
	UserKeyAlgorithm string    `gorm:"column:user_key_algorithm"`               //用户算法
	UserState        string    `gorm:"column:user_state"`                       //用户状态
	Roles            string    `gorm:"column:roles"`                            //用户角色
	Privileges       string    `gorm:"column:privileges"`                       //用户权限
	CreateTime       time.Time `gorm:"column:create_time;autoCreateTime:milli"` //创建时间
	UpdateTime       time.Time `gorm:"column:update_time;autoUpdateTime:milli"` //更新时间
	State            int8      `gorm:"column:state;default:1"`                  //数据状态
}

type Txs struct {
	Id                int64     `gorm:"column:id"`                               //主键
	Ledger            string    `gorm:"column:ledger"`                           //账本hash
	TxBlockHeight     int64     `gorm:"column:tx_block_height"`                  //区块高度
	TxIndex           int32     `gorm:"column:tx_index"`                         //交易顺序
	TxHash            string    `gorm:"column:tx_hash"`                          //交易HASH
	TxNodePubkeys     string    `gorm:"column:tx_node_pubkeys"`                  //交易节点签名公钥
	TxEndpointPubkeys string    `gorm:"column:tx_endpoint_pubkeys"`              //节点终端签名公钥
	TxContents        string    `gorm:"column:tx_contents"`                      //交易内容
	TxResponseState   int       `gorm:"column:tx_response_state"`                //交易执行结果状态
	TxResponseMsg     string    `gorm:"column:tx_response_msg"`                  //交易执行结果说明
	CreateTime        time.Time `gorm:"column:create_time;autoCreateTime:milli"` //创建时间
	UpdateTime        time.Time `gorm:"column:update_time;autoUpdateTime:milli"` //更新时间
	State             int8      `gorm:"column:state;default:1"`                  //数据状态
}
