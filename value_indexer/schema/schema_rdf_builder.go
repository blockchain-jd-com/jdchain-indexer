package schema

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"strings"
)

type SchemaSource interface {
	Fields() NodeSchemaFields
	LowerName() string
}

func NewSchemaMetaBuilder(schema SchemaSource) *SchemaMetaBuilder {
	return &SchemaMetaBuilder{
		schema: schema,
	}
}

type SchemaMetaBuilder struct {
	schema SchemaSource
}

func (builder *SchemaMetaBuilder) Build() (schemas dgraph_helper.Schemas) {
	name := builder.schema.LowerName()
	fields := builder.schema.Fields()
	if !strings.HasPrefix(name, SchemaTypeEdge) {
		for _, field := range fields {
			schema := field.ToSchemaRDF(name)
			if schema == nil {
				continue
			}
			schemas = schemas.Add(schema)
		}

		// key/version/time
		schemas = schemas.Add(dgraph_helper.NewSchemaStringExactIndex(fmt.Sprintf("%s-%s", name, "key")))
		schemas = schemas.Add(dgraph_helper.NewSchemaIntIndex(fmt.Sprintf("%s-%s", name, "version")))
		schemas = schemas.Add(dgraph_helper.NewSchemaIntIndex(fmt.Sprintf("%s-%s", name, "time")))

		// links schema status uid
		schemas = schemas.Add(dgraph_helper.NewSchemaUidIndex(name + "-schema"))
	}
	return
}
