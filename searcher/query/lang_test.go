package query

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestQueryLang_Assemble(t *testing.T) {
	query := newQueryLang(qlNameBlockByHeight, qlBlockQueryResultName)
	result := query.Assemble(map[string]interface{}{
		"from": "1",
		"to":   "2",
	})
	t.Log(result)
	assert.True(t, strings.HasPrefix(result, "{"))
}
