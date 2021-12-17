CREATE DATABASE `jdchain` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

DROP TABLE IF EXISTS `jdchain_blocks`;
DROP TABLE IF EXISTS `jdchain_contracts`;
DROP TABLE IF EXISTS `jdchain_data_account_kvs`;
DROP TABLE IF EXISTS `jdchain_data_accounts`;
DROP TABLE IF EXISTS `jdchain_event_account_events`;
DROP TABLE IF EXISTS `jdchain_event_accounts`;
DROP TABLE IF EXISTS `jdchain_txs`;
DROP TABLE IF EXISTS `jdchain_users`;

-- jdchain.jdchain_blocks definition

CREATE TABLE `jdchain_blocks` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ledger` varchar(46) NOT NULL COMMENT '账本',
  `block_height` bigint(20) NOT NULL COMMENT '区块高度',
  `block_hash` varchar(46) NOT NULL COMMENT '区块HASH',
  `pre_block_hash` varchar(46) DEFAULT NULL COMMENT '前置区块HASH',
  `txs_set_hash` varchar(46) DEFAULT NULL COMMENT '交易集HASH',
  `users_set_hash` varchar(46) DEFAULT NULL COMMENT '用户集HASH',
  `contracts_set_hash` varchar(46) DEFAULT NULL COMMENT '合约集HASH',
  `configurations_set_hash` varchar(46) DEFAULT NULL COMMENT '配置集HASH',
  `dataaccounts_set_hash` varchar(46) DEFAULT NULL COMMENT '数据账户HASH',
  `eventaccounts_set_hash` varchar(46) DEFAULT NULL COMMENT '事件账户HASH',
  `block_timestamp` timestamp NULL DEFAULT NULL COMMENT '区块创建时间',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `state` tinyint(4) DEFAULT '1' COMMENT '数据状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `I_U_LEDGER_BLOCK_HEIGHT` (`ledger`,`block_height`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='区块信息表';


-- jdchain.jdchain_contracts definition

CREATE TABLE `jdchain_contracts` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ledger` varchar(46) NOT NULL COMMENT '账本',
  `contract_address` varchar(46) NOT NULL COMMENT '合约地址',
  `contract_pubkey` varchar(48) NOT NULL COMMENT '合约公钥',
  `contract_roles` varchar(256) NOT NULL COMMENT '合约归属角色',
  `contract_priviledges` varchar(10) NOT NULL COMMENT '合约权限',
  `contract_version` int(11) NOT NULL COMMENT '合约版本',
  `contract_status` varchar(10) NOT NULL COMMENT '合约状态',
  `contract_creator` varchar(46) NOT NULL COMMENT '合约创建者地址',
  `contract_content` text COMMENT '合约内容',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `state` tinyint(4) DEFAULT '1' COMMENT '数据状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `I_U_LEDGER_CONTRACT_ADDR` (`ledger`,`contract_address`),
  KEY `I_CONTRACTS_ADDSS` (`contract_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='合约信息表';


-- jdchain.jdchain_data_account_kvs definition

CREATE TABLE `jdchain_data_account_kvs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ledger` varchar(46) NOT NULL COMMENT '账本',
  `data_account_address` varchar(46) NOT NULL COMMENT '数据账户地址',
  `data_account_key` varchar(512) NOT NULL COMMENT '数据账户key',
  `data_account_value` blob COMMENT '数据账户value',
  `data_account_type` varchar(10) NOT NULL COMMENT '数据账户类型',
  `data_account_version` int(11) NOT NULL COMMENT '数据账户版本',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `state` tinyint(4) DEFAULT '1' COMMENT '数据状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `I_LEDGER_ADDR_KEY_VER` (`ledger`,`data_account_address`,`data_account_key`,`data_account_version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据账户KV信息表';


-- jdchain.jdchain_data_accounts definition

CREATE TABLE `jdchain_data_accounts` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ledger` varchar(46) NOT NULL COMMENT '账本',
  `data_account_address` varchar(46) NOT NULL COMMENT '数据账户地址',
  `data_account_pubkey` varchar(48) NOT NULL COMMENT '数据账户公钥',
  `data_account_roles` varchar(256) NOT NULL COMMENT '数据账户角色',
  `data_account_privileges` varchar(10) NOT NULL COMMENT '数据账户权限',
  `data_account_creator` varchar(46) NOT NULL COMMENT '数据账户创建者',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `state` tinyint(4) DEFAULT '1' COMMENT '数据状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `I_U_LEDGER_DATA_ACCOUNT_ADDR` (`ledger`,`data_account_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据账户信息表';


-- jdchain.jdchain_event_account_events definition

CREATE TABLE `jdchain_event_account_events` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ledger` varchar(46) NOT NULL COMMENT '账本',
  `event_account_address` varchar(46) NOT NULL COMMENT '事件账户地址',
  `event_name` varchar(256) NOT NULL COMMENT '事件名称',
  `event_sequence` int(11) NOT NULL COMMENT '事件序列',
  `event_tx_hash` varchar(46) NOT NULL COMMENT '事件交易HASH',
  `event_block_height` bigint(20) NOT NULL COMMENT '事件高度',
  `event_type` varchar(32) NOT NULL COMMENT '事件类型',
  `event_value` varchar(1024) NOT NULL COMMENT '事件值',
  `event_contract_address` varchar(46) DEFAULT NULL COMMENT '事件合约地址',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `state` tinyint(4) DEFAULT '1' COMMENT '数据状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `I_LEDGER_EVENT_ADDR_NAME_SEQ` (`ledger`,`event_account_address`,`event_name`,`event_sequence`),
  KEY `I_LEDGER_EVENT_NAME` (`event_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件信息表';


-- jdchain.jdchain_event_accounts definition

CREATE TABLE `jdchain_event_accounts` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ledger` varchar(46) NOT NULL COMMENT '账本',
  `event_account_address` varchar(46) NOT NULL COMMENT '事件账户地址',
  `event_account_pubkey` varchar(48) NOT NULL COMMENT '事件账户公钥',
  `event_account_roles` varchar(256) NOT NULL COMMENT '事件账户归属角色',
  `event_account_priviledges` varchar(10) NOT NULL COMMENT '事件账户权限',
  `event_account_creator` varchar(46) NOT NULL COMMENT '事件账户创建者',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `state` tinyint(4) DEFAULT '1' COMMENT '数据状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `I_U_LEDGER_EVENT_ADDR` (`ledger`,`event_account_address`),
  KEY `I_EVENT_ADDR` (`event_account_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件账户信息表';


-- jdchain.jdchain_txs definition

CREATE TABLE `jdchain_txs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ledger` varchar(46) NOT NULL COMMENT '账本hash',
  `tx_block_height` bigint(20) NOT NULL COMMENT '区块高度',
  `tx_index` int(11) NOT NULL COMMENT '交易顺序',
  `tx_hash` varchar(46) NOT NULL COMMENT '交易HASH',
  `tx_node_pubkeys` varchar(1024) DEFAULT NULL COMMENT '交易节点签名公钥',
  `tx_endpoint_pubkeys` varchar(1024) DEFAULT NULL COMMENT '节点终端签名公钥',
  `tx_contents` text COMMENT '交易内容',
  `tx_response_state` int(11) DEFAULT NULL COMMENT '交易执行结果状态',
  `tx_response_msg` varchar(256) DEFAULT NULL COMMENT '交易执行结果说明',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `state` tinyint(4) DEFAULT '1' COMMENT '数据状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `I_U_TXS_LEDGER_HEIGHT_INDEX` (`ledger`,`tx_block_height`,`tx_index`),
  KEY `I_TXS_TX_HASH` (`tx_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交易信息表';


-- jdchain.jdchain_users definition

CREATE TABLE `jdchain_users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ledger` varchar(46) NOT NULL COMMENT '账本',
  `user_address` varchar(46) NOT NULL COMMENT '用户地址',
  `user_pubkey` varchar(48) NOT NULL COMMENT '用户公钥',
  `user_key_algorithm` varchar(16) NOT NULL COMMENT '用户算法',
  `user_state` varchar(32) NOT NULL COMMENT '用户状态',
  `roles` varchar(1024) NOT NULL COMMENT '用户角色',
  `privileges` varchar(1024) NOT NULL COMMENT '用户权限',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `state` tinyint(4) DEFAULT '1' COMMENT '状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `I_U_LEDGER_USER_ADDR` (`ledger`,`user_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息表';