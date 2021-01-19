package sql

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/value_indexer/schema"
	"github.com/pingcap/parser/opcode"
	"strings"
)

func NewExpress(op opcode.Op, column, value string) *Express {
	return &Express{
		column: column,
		op:     op,
		value:  value,
	}
}

type Express struct {
	column string
	op     opcode.Op
	value  string
	not    bool
}

func (expr *Express) SetNot(b bool) {
	expr.not = b
}

func (expr *Express) Expression(prefix string, nodeSchema *schema.NodeSchema) (s string) {
	var key string
	if len(prefix) > 0 {
		key = fmt.Sprintf("%s-%s", prefix, expr.column)
	} else {
		key = expr.column
	}
	switch expr.op {
	case opcode.EQ:
		s = fmt.Sprintf("eq(%s, %s)", key, expr.value)
	case opcode.NE:
		s = fmt.Sprintf("(not eq(%s, %s))", key, expr.value)
	case opcode.In:
		trimEmpty := strings.TrimLeft(expr.value, "||")
		if nodeSchema == nil {
			v := strings.Join(strings.Split(trimEmpty, "||"), ",")
			s = fmt.Sprintf("eq(%s, [%s])", key, v)
		} else {
			field, ok := nodeSchema.FindField(expr.column)
			if ok == false {
				return
			}
			switch field.Type() {
			case schema.DataTypeString:
				keywords := strings.Split(trimEmpty, "||")
				for _, kw := range keywords {
					if len(s) <= 0 {
						s = fmt.Sprintf(`regexp(%s, /^%s$/)`, key, kw)
					} else {
						s = fmt.Sprintf(`%s or regexp(%s, /^%s$/)`, s, key, kw)
					}
				}
				s = fmt.Sprintf("(%s)", s)
				//s = fmt.Sprintf(`anyofterms(%s, "%s")`, key, v)
			case schema.DataTypeInt, schema.DataTypeID, schema.DataTypeFloat:
				v := strings.Join(strings.Split(trimEmpty, "||"), ",")
				s = fmt.Sprintf("eq(%s, [%s])", key, v)
			case schema.DataTypeBool:
			case schema.DataTypeDateTime:
			default:
			}
		}
	case opcode.LE:
		s = fmt.Sprintf("le(%s, %s)", key, expr.value)
	case opcode.LT:
		s = fmt.Sprintf("lt(%s, %s)", key, expr.value)
	case opcode.GE:
		s = fmt.Sprintf("ge(%s, %s)", key, expr.value)
	case opcode.GT:
		s = fmt.Sprintf("gt(%s, %s)", key, expr.value)
	case opcode.Like:
		keyword := strings.Replace(expr.value, "%", "", -1)
		isHasWordBeforeKeyword := strings.HasPrefix(expr.value, "%")
		isHasWordAfterKeyword := strings.HasSuffix(expr.value, "%")

		if isHasWordBeforeKeyword && isHasWordAfterKeyword {
			// like SELECT * FROM Persons WHERE City LIKE '%Ne%'
			s = fmt.Sprintf(`regexp(%s, /\S*%s\S*/)`, key, keyword)
		} else if isHasWordBeforeKeyword {
			// like SELECT * FROM Persons WHERE City LIKE '%Ne'
			s = fmt.Sprintf(`regexp(%s, /\S*%s$/)`, key, keyword)
		} else if isHasWordAfterKeyword {
			// like SELECT * FROM Persons WHERE City LIKE 'Ne%'
			s = fmt.Sprintf(`regexp(%s, /^%s\S*/)`, key, keyword)
		} else {
			// like SELECT * FROM Persons WHERE City LIKE 'Ne'
			s = fmt.Sprintf(`regexp(%s, /^%s$/)`, key, keyword)
		}
		if expr.not {
			s = fmt.Sprintf(`(not %s )`, s)
		}
	default:
		logger.Warnf("no handler for %s in Express", expr.op)
	}
	return
}

func NewExpressTree(op opcode.Op, left, right Filter) *ExpressTree {
	return &ExpressTree{
		op:    op,
		left:  left,
		right: right,
	}
}

type ExpressTree struct {
	op          opcode.Op
	left, right Filter
}

func (expr ExpressTree) Expression(prefix string, nodeSchema *schema.NodeSchema) (s string) {
	var left, right string
	if expr.left != nil {
		left = expr.left.Expression(prefix, nodeSchema)
	}
	if expr.right != nil {
		right = expr.right.Expression(prefix, nodeSchema)
	}

	if len(left) <= 0 || len(right) <= 0 {
		s = left + right
		return
	}
	switch expr.op {
	case opcode.LogicAnd:
		s = fmt.Sprintf("%s and %s", left, right)
	case opcode.LogicOr:
		s = fmt.Sprintf("%s or %s", left, right)
	}
	return
}
