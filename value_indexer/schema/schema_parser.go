package schema

import (
	"fmt"
	"github.com/ssor/zlog"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
	"strings"
)

const (
	SchemaTypeNode = "node"
	SchemaTypeEdge = "edge"
)

type CommonSchema interface {
	Type() string
	LowerName() string
}

type CommonSchemas []CommonSchema

func (css CommonSchemas) FindNodeSchema(name string) (node *NodeSchema) {
	for _, cs := range css {
		if cs.Type() == SchemaTypeNode {
			schema := cs.(*NodeSchema)
			if strings.ToLower(name) == schema.LowerName() {
				return schema
			}
		}
	}
	return
}

func (css CommonSchemas) FindRelativeSchemas(name string) (result CommonSchemas) {
	for _, cs := range css {
		t := cs.Type()
		switch t {
		case SchemaTypeEdge:
			schema := cs.(*EdgeSchema)
			if schema.IsRelativeWith(name) {
				result = append(result, schema)
			} else {
			}
		case SchemaTypeNode:
			schema := cs.(*NodeSchema)
			if strings.ToLower(name) == schema.LowerName() {
				result = append(result, schema)
			}
		}
	}
	return
}

func NewSchemaParser() *SchemaParser {
	return &SchemaParser{}
}

type SchemaParser struct {
}

func (sp *SchemaParser) FirstNodeSchema(schema string) (ns *NodeSchema, e error) {
	css, err := sp.Parse(schema)
	if err != nil {
		e = err
		return
	}
	if len(css) <= 0 {
		return
	}
	for _, cs := range css {
		nodeSchema, ok := cs.(*NodeSchema)
		if ok == false {
			continue
		}
		ns = nodeSchema
		break
	}
	if ns == nil {
		e = fmt.Errorf("no node schema detected")
		return
	}
	return
}

func (sp *SchemaParser) Parse(schema string) (css CommonSchemas, e error) {
	doc, err := parser.ParseSchema(&ast.Source{Input: schema, Name: "spec"})
	if err != nil {
		e = err
		return
	}

	for _, def := range doc.Definitions {
		cs := sp.parseDefinition(def)
		if cs == nil {
			continue
		}
		css = append(css, cs)
	}
	return
}

func (sp *SchemaParser) parseDefinition(def *ast.Definition) (cs CommonSchema) {
	name := def.Name
	if strings.HasPrefix(strings.ToLower(name), SchemaTypeEdge) {
		if len(def.Fields) < 2 {
			zlog.Warnf("Edge Schema definition %s has less than 2 fields", def.Name)
			return
		}
		field0 := def.Fields[0]
		if len(field0.Arguments) <= 0 {
			zlog.Warnf("Edge Schema definition %s has no reference property", def.Name)
			return
		}
		from := NewEdgeSchemaField(field0.Name, field0.Arguments[0].Name)
		field1 := def.Fields[1]
		if len(field1.Arguments) <= 0 {
			zlog.Warnf("Edge Schema definition %s has no reference property", def.Name)
			return
		}
		to := NewEdgeSchemaField(field1.Name, field1.Arguments[0].Name)
		cs = NewEdgeSchema(name, from, to)
	} else {
		node := NewNodeSchema(def.Name)
		if len(def.Fields) <= 0 {
			zlog.Warnf("definition %s has no field", def.Name)
			return
		}
		//spew.Dump(def)
		for _, field := range def.Fields {
			var isPrimaryKey bool
			argPrimaryKey := field.Arguments.ForName("isPrimaryKey")
			if argPrimaryKey != nil && argPrimaryKey.DefaultValue.String() == "true" {
				isPrimaryKey = true
			}
			node.AddField(NewNodeSchemaField(field.Name, field.Type.String(), isPrimaryKey))
		}
		cs = node
	}
	return
}
