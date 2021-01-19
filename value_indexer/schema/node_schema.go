package schema

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"strings"
)

func NewNodeSchema(name string) *NodeSchema {
	return &NodeSchema{
		Name:   name,
		fields: NodeSchemaFields{},
	}
}

type NodeSchema struct {
	Name   string
	fields NodeSchemaFields
}

func (schema *NodeSchema) Type() string {
	return SchemaTypeNode
}

func (schema *NodeSchema) FindField(name string) (field NodeSchemaField, ok bool) {
	for _, field := range schema.fields {
		if field.LowerName() == strings.ToLower(name) {
			return field, true
		}
	}
	return NodeSchemaField{}, false
}
func (schema *NodeSchema) ListTypeFields() (fields NodeSchemaFields) {
	for _, field := range schema.fields {
		if field.IsListType() {
			fields = append(fields, field)
		}
	}
	return
}

func (schema *NodeSchema) PrimaryPredict() (s string) {
	primaryField, ok := schema.PrimaryField()
	if ok == false {
		return
	}
	return fmt.Sprintf("%s-%s", schema.LowerName(), primaryField.LowerName())
}

func (schema *NodeSchema) PrimaryField() (NodeSchemaField, bool) {
	for _, field := range schema.fields {
		if field.IsPrimaryKey() {
			return field, true
		}
	}
	return NodeSchemaField{}, false
}

func (schema *NodeSchema) LowerName() string {
	return strings.ToLower(schema.Name)
}

func (schema *NodeSchema) FormatCacheQueryKey(v string) string {
	return fmt.Sprintf("%s-%s", schema.LowerName(), v)
}

func (schema *NodeSchema) queryName() string {
	return strings.ToLower(schema.Name) + "s"
}

func (schema *NodeSchema) Fields() NodeSchemaFields {
	return schema.fields
}

func (schema *NodeSchema) AddField(fields ...NodeSchemaField) {
	schema.fields = append(schema.fields, fields...)
}

type NodeSchemas []*NodeSchema

func (schemas NodeSchemas) FindByName(name string) *NodeSchema {
	for _, schema := range schemas {
		if strings.ToLower(schema.Name) == strings.ToLower(name) {
			return schema
		}
	}
	return nil
}

type NodeSchemaFields []NodeSchemaField

func NewNodeSchemaField(name, typeName string, isPrimaryKey bool) NodeSchemaField {
	return NodeSchemaField{
		Name:         name,
		TypeName:     typeName,
		isPrimaryKey: isPrimaryKey,
	}
}

type NodeSchemaField struct {
	Name         string
	TypeName     string
	isPrimaryKey bool
}

func (field NodeSchemaField) ToSchemaRDF(prefix string) (schema *dgraph_helper.Schema) {
	switch strings.ToLower(field.TypeName) {
	case DataTypeString:
		schema = dgraph_helper.NewSchemaStringTrigramIndex(fmt.Sprintf("%s-%s", prefix, field.Name))
	case DataTypeInt, DataTypeID:
		schema = dgraph_helper.NewSchemaIntIndex(fmt.Sprintf("%s-%s", prefix, field.Name))
	case DataTypeFloat:
		schema = dgraph_helper.NewSchemaFloatIndex(fmt.Sprintf("%s-%s", prefix, field.Name))
	case DataTypeBool:
		schema = dgraph_helper.NewSchemaBoolIndex(fmt.Sprintf("%s-%s", prefix, field.Name))
	case DataTypeDateTime:
		schema = dgraph_helper.NewSchemaDateIndex(fmt.Sprintf("%s-%s", prefix, field.Name))
	case DataTypeUid:
		schema = dgraph_helper.NewSchemaUidIndex(fmt.Sprintf("%s-%s", prefix, field.Name))
	case DataTypeUids:
		schema = dgraph_helper.NewSchemaUidsIndex(fmt.Sprintf("%s-%s", prefix, field.Name))
	default:
		if strings.HasPrefix(field.TypeName, "[") {

		} else {
			fmt.Println("cannot recognize type: ", field.TypeName)
		}
	}
	return
}

func (field NodeSchemaField) IsListType() bool {
	return strings.HasPrefix(field.TypeName, "[") &&
		!strings.EqualFold(field.TypeName, DataTypeUid) &&
		!strings.EqualFold(field.TypeName, DataTypeUids)
}

func (field NodeSchemaField) IsUidType() bool {
	return strings.EqualFold(field.TypeName, DataTypeUid)
}

func (field NodeSchemaField) IsUidsType() bool {
	return strings.EqualFold(field.TypeName, DataTypeUids)
}

func (field NodeSchemaField) IsPrimaryKey() bool {
	return field.isPrimaryKey
}

func (field NodeSchemaField) Type() string {
	return field.TypeName
}

func (field NodeSchemaField) LowerName() string {
	return strings.ToLower(field.Name)
}

func (field NodeSchemaField) FormatPredict(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, field.LowerName())
}
