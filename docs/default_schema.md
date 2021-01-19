## 基础数据索引

### nodes

| Name | Items | Type | Index |
| ------ | ------ | ------ | ------ |
| ledger | hash_id | string | @index(term,trigram) |
| block | height | int | @index(int) |
|  | time | int | @index(int)  |
|  | hash_id | string | @index(term,trigram)  |
| tx | hash_id | string | @index(term,trigram)  |
|  | execution_state | string | @index(exact) |
|  | block_height | int | @index(int) |
|  | time | dateTime | @index(hour) |
| kv | key | string | @index(term,trigram)  |
|  | version | int | @index(int) |
|  | tx_index | string | @index(exact) |
|  | data_account_address | string | @index(exact) |
| data_account | address | string | @index(term,trigram)  |
|  | public_key | string | @index(term,trigram)  |
| contract | address | string | @index(term,trigram)  |
|  | public_key | string | @index(term,trigram)  |
| contract_event | contract_event | string | @index(term,trigram)  |
|  | contract_address | string | @index(exact) |
|  | tx_index | string | @index(exact) |
| user | address | string | @index(term,trigram)  |
|  | public_key | string | @index(term,trigram)  |
| event_account | address | string | @index(term,trigram)  |
|  | public_key | string | @index(term,trigram)  |
| event | topic | string | @index(term,trigram)  |
|  | sequence | int | @index(int) |
|  | event_account_address | string | @index(exact) |
|  | tx_index | string | @index(exact) |

### links

| Name | Links |
| ------ | ------ |
| ledger-block | ledger-hash_id |
|  | block-hash_id |
| ledger-contract | ledger-hash_id |
|  | contract-address |
| ledger-data_account | ledger-hash_id |
|  | data_account-address |
| ledger-event_account | ledger-hash_id |
|  | event_account-address |
| ledger-user | ledger-hash_id |
|  | user-address |
| block-tx | block-hash_id |
|  | tx-hash_id |
| endpoint_user-tx | user-public_key |
|  | tx-hash_id |
| tx-contract | tx-hash_id |
|  | contract-address |
| tx-contract_event | tx-hash_id |
|  | contract_event-tx_index |
| tx-data_account | tx-hash_id |
|  | data_account-address |
| tx-event | tx-hash_id |
|  | event-tx_index |
| tx-event_account | tx-hash_id |
|  | event_account-address |
| tx-kv | tx-hash_id |
|  | kv-tx_index |
| tx-user | tx-hash_id |
|  | kv-user-address |