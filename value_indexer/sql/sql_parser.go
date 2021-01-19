package sql

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/opcode"
	"github.com/pingcap/tidb/types/parser_driver"
	"strings"
)

func NewSqlParser() *Parser {
	return &Parser{}
}

type Parser struct {
}

func (sqlParser *Parser) Parse(src string) (from string, filter Filter, e error) {
	p := parser.New()

	stmtNodes, _, err := p.Parse(src, "", "")
	if err != nil {
		logger.Errorf("parse sql failed: %s", err)
		e = err
		return
	}

	spew.Config.Indent = "    "
	head, ok := stmtNodes[0].(*ast.SelectStmt)
	if ok == false {
		e = fmt.Errorf("invalid sql")
		return
	}
	from = head.From.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.L
	if tree, success := extractWhereExpressions(head.Where); success == true {
		filter = tree
	}
	return
}

func extractWhereExpressions(node ast.ExprNode) (express Filter, success bool) {
	if node == nil {
		return
	}
	switch t := node.(type) {
	case *ast.BinaryOperationExpr:
		left := t.L
		right := t.R
		columnNameExpr, isColumnName := left.(*ast.ColumnNameExpr)
		if isColumnName == false {
			leftExpress, ok1 := extractWhereExpressions(left)
			rightExpress, ok2 := extractWhereExpressions(right)
			if ok1 == false || ok2 == false {
				return
			}

			express = NewExpressTree(t.Op, leftExpress, rightExpress)
			success = true
			return
		}
		valueExpr, isValueExpr := right.(*driver.ValueExpr)
		if isValueExpr == false {
			logger.Warnf("column does not have value")
			return
		}
		v, err := valueExpr.Datum.ToString()
		if err != nil {
			logger.Warnf("get express value failed: %s", err)
			spew.Dump(valueExpr)
			return
		}
		express = NewExpress(t.Op, columnNameExpr.Name.Name.L, v)
		success = true
	case *ast.PatternInExpr:
		columnNameExpr, isColumnName := t.Expr.(*ast.ColumnNameExpr)
		if isColumnName == false {
			logger.Warnf("should be column at left in PatternInExpr")
			return
		}
		exprNodes := t.List
		var sb strings.Builder
		for _, en := range exprNodes {
			valueExpr, isValueExpr := en.(*driver.ValueExpr)
			if isValueExpr == false {
				logger.Warnf("column does not have value")
				return
			}
			v, err := valueExpr.Datum.ToString()
			if err != nil {
				logger.Warnf("get express value failed: %s", err)
				spew.Dump(valueExpr)
				continue
			}
			sb.WriteString("||")
			sb.WriteString(v)
		}
		express = NewExpress(opcode.In, columnNameExpr.Name.Name.L, sb.String())
		success = true
	case *ast.ParenthesesExpr:
		express, success = extractWhereExpressions(t.Expr)
	case *ast.BetweenExpr:
		columnNameExpr, isColumnName := t.Expr.(*ast.ColumnNameExpr)
		if isColumnName == false {
			logger.Warnf("should be column at left in PatternInExpr")
			return
		}
		valueExpr, isValueExpr := t.Left.(*driver.ValueExpr)
		if isValueExpr == false {
			logger.Warnf("between expr left should be value expr")
			return
		}
		leftValue, err := valueExpr.Datum.ToString()
		if err != nil {
			logger.Warnf("get express value failed: %s", err)
			spew.Dump(valueExpr)
			return
		}
		valueExpr, isValueExpr = t.Right.(*driver.ValueExpr)
		if isValueExpr == false {
			logger.Warnf("between expr right should be value expr")
			return
		}
		rightValue, err := valueExpr.Datum.ToString()
		if err != nil {
			logger.Warnf("get express value failed: %s", err)
			spew.Dump(valueExpr)
			return
		}
		leftExpress := NewExpress(opcode.GE, columnNameExpr.Name.Name.L, leftValue)
		rightExpress := NewExpress(opcode.LE, columnNameExpr.Name.Name.L, rightValue)
		express = NewExpressTree(opcode.LogicAnd, leftExpress, rightExpress)
		success = true
	case *ast.PatternLikeExpr:
		columnNameExpr, isColumnName := t.Expr.(*ast.ColumnNameExpr)
		if isColumnName == false {
			logger.Warnf("should be column at left in PatternInExpr")
			return
		}
		valueExpr, isValueExpr := t.Pattern.(*driver.ValueExpr)
		if isValueExpr == false {
			logger.Warnf("between expr left should be value expr")
			return
		}
		value, err := valueExpr.Datum.ToString()
		if err != nil {
			logger.Warnf("get express value failed: %s", err)
			spew.Dump(valueExpr)
			return
		}
		express = NewExpress(opcode.Like, columnNameExpr.Name.Name.L, value)
		express.(*Express).SetNot(t.Not)
		success = true
	default:
		logger.Warnf("sql parser has no handler for new type: input == nil -> %t", node == nil)
		spew.Dump(t)
	}
	return
}
