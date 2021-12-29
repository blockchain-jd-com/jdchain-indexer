package handler

import (
	"bufio"
	"encoding/json"
	"git.jd.com/jd-blockchain/explorer/searcher/query"
)

type QueryResult struct {
	combine       bool
	Ledgers       query.Ledgers         `json:"Ledgers,omitempty"`
	Blocks        query.Blocks          `json:"blocks,omitempty"`
	Txs           query.Transactions    `json:"txs,omitempty"`
	Users         query.Users           `json:"users,omitempty"`
	Accounts      query.Accounts        `json:"accounts,omitempty"`
	Contracts     query.Contracts       `json:"contracts,omitempty"`
	KVs           query.WriteKvs        `json:"kvs,omitempty"`
	EventAccounts query.EventAccounts   `json:"event_accounts,omitempty"`
	Events        query.Events          `json:"events,omitempty"`
	KvUsers       query.KvEndpointUsers `json:"kvusers,omitempty"`
}

func (result *QueryResult) ToJSON(writer *bufio.Writer) (e error) {
	if result == nil {
		return
	}

	if result.combine {
		return result.toCombineJSON(writer)
	} else {
		return result.toSingleJSON(writer)
	}
}

func (result *QueryResult) toCombineJSON(writer *bufio.Writer) (e error) {

	_, e = writer.WriteString("{")

	needComma := false

	if result.Ledgers != nil {
		_, e = writer.WriteString(`"ledgers":`)
		if len(result.Ledgers) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.Ledgers)
		}
		needComma = true
	}

	if result.Blocks != nil {
		if needComma {
			_, e = writer.WriteString(",")
		}
		_, e = writer.WriteString(` "blocks":`)
		if len(result.Blocks) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.Blocks)
		}
		needComma = true
	}

	if result.Txs != nil {
		if needComma {
			_, e = writer.WriteString(",")
		}
		_, e = writer.WriteString(` "txs":`)
		if len(result.Txs) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.Txs)
		}
		needComma = true
	}

	if result.Users != nil {
		if needComma {
			_, e = writer.WriteString(",")
		}
		_, e = writer.WriteString(` "users":`)
		if len(result.Users) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.Users)
		}
		needComma = true
	}

	if result.Accounts != nil {
		if needComma {
			_, e = writer.WriteString(",")
		}
		_, e = writer.WriteString(` "accounts":`)
		if len(result.Accounts) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.Accounts)
		}
		needComma = true
	}

	if result.Contracts != nil {
		if needComma {
			_, e = writer.WriteString(",")
		}
		_, e = writer.WriteString(` "contracts":`)
		if len(result.Contracts) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.Contracts)
		}
		needComma = true
	}

	if result.EventAccounts != nil {
		if needComma {
			_, e = writer.WriteString(",")
		}
		_, e = writer.WriteString(` "event_accounts":`)
		if len(result.EventAccounts) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.EventAccounts)
		}
		needComma = true
	}

	if result.KvUsers != nil {
		if needComma {
			_, e = writer.WriteString(",")
		}
		_, e = writer.WriteString(` "kv_users":`)
		if len(result.KvUsers) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.KvUsers)
		}
		needComma = true
	}

	_, e = writer.WriteString("}")

	return e
}

func (result *QueryResult) toSingleJSON(writer *bufio.Writer) (e error) {

	if result.Ledgers != nil {
		if len(result.Ledgers) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.Ledgers)
		}
	} else if result.Blocks != nil {
		if len(result.Blocks) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.Blocks)
		}
	} else if result.Txs != nil {
		if len(result.Txs) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.Txs)
		}
	} else if result.Users != nil {
		if len(result.Users) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.Users)
		}
	} else if result.Accounts != nil {
		if len(result.Accounts) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.Accounts)
		}
	} else if result.Contracts != nil {
		if len(result.Contracts) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.Contracts)
		}
	} else if result.EventAccounts != nil {
		if len(result.EventAccounts) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.EventAccounts)
		}
	} else if result.KvUsers != nil {
		if len(result.KvUsers) <= 0 {
			_, e = writer.WriteString("[]")
		} else {
			e = result.writeProp(writer, result.KvUsers)
		}
	}

	return e
}

func (result *QueryResult) writeProp(writer *bufio.Writer, prop interface{}) (e error) {
	bs, err := json.Marshal(prop)
	if err != nil {
		e = err
		return
	}
	_, e = writer.Write(bs)
	return
}
