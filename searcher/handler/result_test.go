package handler

import (
	"bufio"
	"bytes"
	"git.jd.com/jd-blockchain/explorer/searcher/query"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestQueryResult_MarshalJSON(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	writer := bufio.NewWriter(buf)

	result := &QueryResult{}
	err := result.ToJSON(writer)
	assert.Nil(t, err)
	err = writer.Flush()
	assert.Nil(t, err)
	t.Logf("empty result: %s", buf.String())
	assert.Len(t, buf.String(), 2, buf.String())

	buf.Reset()

	result.Contracts = make([]query.Contract, 0)
	result.Blocks = make(query.Blocks, 0)
	err = result.ToJSON(writer)
	assert.Nil(t, err)
	err = writer.Flush()
	assert.Nil(t, err)
	t.Logf("0 contracts: %s", buf.String())
	assert.True(t, strings.Contains(buf.String(), "contracts"), buf.String())

}
