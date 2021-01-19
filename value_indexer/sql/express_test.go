package sql

import (
	"github.com/pingcap/parser/opcode"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewExpressTree(t *testing.T) {
	left := NewExpress(opcode.EQ, "c1", "v1")
	tree := NewExpressTree(opcode.LogicAnd, left, nil)
	assert.Equal(t, "eq(c1, v1)", tree.Expression("", nil))

	right1 := NewExpress(opcode.EQ, "c2", "v2")
	right2 := NewExpress(opcode.In, "c3", "v3")
	assert.Equal(t, "eq(c2, v2)", right1.Expression("", nil))
	assert.Equal(t, "eq(c3, [v3])", right2.Expression("", nil))

	tree = NewExpressTree(opcode.LogicOr, right1, right2)
	assert.Equal(t, "eq(c2, v2) or eq(c3, [v3])", tree.Expression("", nil))

	deepTree := NewExpressTree(opcode.LogicAnd, left, tree)
	assert.Equal(t, "eq(c1, v1) and eq(c2, v2) or eq(c3, [v3])", deepTree.Expression("", nil))
}
