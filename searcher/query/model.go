package query

type Transactions []*Transaction

func (txs Transactions) exists(tx *Transaction) (ok bool) {
	for _, txTemp := range txs {
		if tx.HashID == txTemp.HashID {
			ok = true
			return
		}
	}
	return
}

func (txs Transactions) add(txsIn ...*Transaction) Transactions {
	txsNew := txs
	for _, txTemp := range txsIn {
		if txsNew.exists(txTemp) == false {
			txsNew = append(txsNew, txTemp)
		}
	}
	return txsNew
}

type Ledgers []*LedgerInfo

type LedgerInfo struct {
	HashID string `json:"hash"`
	Height int    `json:"height"`
}

type Blocks []*BlockInfo

type BlockInfo struct {
	HashID string `json:"hash"`
	Height int    `json:"height"`
	Time   int64  `json:"time"`
}

func (b *BlockInfo) GetHashID() string {
	return b.HashID
}

type Transaction struct {
	HashID string `json:"hash"`
	//IndexInBlock   int64  `json:"index_in_block"`
	BlockHeight    int64  `json:"block_height"`
	ExecutionState string `json:"execution_state"`
	Time           int64  `json:"time"`
	//NodePublicKey     string                 `json:"node_public_key"`
	//EndpointPublicKey string                 `json:"endpoint_public_key"`
	//Block             *BlockInfo             `json:"block_detail"`
	//KVPairs           map[string]interface{} `json:"kv_pairs"`
}

func (t *Transaction) GetHashID() string {
	return t.HashID
}

type Content struct {
	KVPairs map[string]interface{} `json:"kv_pairs"`
	Tx      *Transaction           `json:"tx"`
}

type WriteKvs []*WriteKV

type WriteKV struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Version int64  `json:"version"`
}

type Users []*User

type User struct {
	Address   string `json:"address"`
	PublicKey string `json:"pubKey"`
}

func (u *User) GetHashID() string {
	return u.PublicKey + u.Address
}

type Accounts []*Account

func (dss Accounts) Exists(dsIn *Account) bool {
	for _, ds := range dss {
		if ds.IsEqualWith(dsIn) {
			return true
		}
	}
	return false
}

func (dss Accounts) Append(dssIn ...*Account) (out Accounts) {
	out = dss
	for _, dsIn := range dssIn {
		if out.Exists(dsIn) {
			continue
		}
		out = append(out, dsIn)
	}
	return
}

type Account struct {
	Address   string `json:"address"`
	PublicKey string `json:"pubKey"`
}

func (ds *Account) IsEqualWith(dsIn *Account) bool {
	if len(ds.Address) > 0 {
		return ds.Address == dsIn.Address
	} else if len(ds.PublicKey) > 0 {
		return ds.PublicKey == dsIn.PublicKey
	} else if len(dsIn.PublicKey) > 0 || len(dsIn.Address) > 0 {
		return false
	}
	return true
}

func (ds *Account) GetHashID() string {
	return ds.PublicKey + ds.Address
}

type Contracts []Contract

type Contract struct {
	Address   string `json:"address"`
	PublicKey string `json:"pubKey"`
}

func (c Contract) GetHashID() string {
	return c.PublicKey + c.Address
}

type EventAccounts []*EventAccount

func (dss EventAccounts) Exists(dsIn *EventAccount) bool {
	for _, ds := range dss {
		if ds.IsEqualWith(dsIn) {
			return true
		}
	}
	return false
}

func (dss EventAccounts) Append(dssIn ...*EventAccount) (out EventAccounts) {
	out = dss
	for _, dsIn := range dssIn {
		if out.Exists(dsIn) {
			continue
		}
		out = append(out, dsIn)
	}
	return
}

type EventAccount struct {
	Address   string `json:"address"`
	PublicKey string `json:"pubKey"`
}

func (ds *EventAccount) IsEqualWith(dsIn *EventAccount) bool {
	if len(ds.Address) > 0 {
		return ds.Address == dsIn.Address
	} else if len(ds.PublicKey) > 0 {
		return ds.PublicKey == dsIn.PublicKey
	} else if len(dsIn.PublicKey) > 0 || len(dsIn.Address) > 0 {
		return false
	}
	return true
}

func (ds *EventAccount) GetHashID() string {
	return ds.PublicKey + ds.Address
}

type Events []*Event

type Event struct {
	Topic    string `json:"topic"`
	Content  string `json:"content"`
	Sequence int64  `json:"sequence"`
}
