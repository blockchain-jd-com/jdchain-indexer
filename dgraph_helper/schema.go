package dgraph_helper

import "fmt"

func NewSchemaIntIndex(name string) *Schema {
	return newSchema(name, "int", "@index(int)")
}

func NewSchemaStringExactIndex(name string) *Schema {
	return newSchema(name, "string", "@index(exact)")
}

func NewSchemaStringTermIndex(name string) *Schema {
	return newSchema(name, "string", "@index(term)")
}

func NewSchemaStringTrigramIndex(name string) *Schema {
	return newSchema(name, "string", "@index(term,trigram)")
}

func NewSchemaString(name string) *Schema {
	return newSchema(name, "string", "")
}

func NewSchemaFloatIndex(name string) *Schema {
	return newSchema(name, "float", "@index(float)")
}

func NewSchemaBoolIndex(name string) *Schema {
	return newSchema(name, "bool", "@index(bool)")
}

func NewSchemaDateIndex(name string) *Schema {
	return newSchema(name, "dateTime", "@index(hour)")
}

func NewSchemaUidIndex(name string) *Schema {
	return newSchema(name, "uid", "@reverse")
}

func NewSchemaUidsIndex(name string) *Schema {
	return newSchema(name, "[uid]", "@count @reverse")
}

func newSchema(name, schemeType, others string) *Schema {
	return &Schema{
		Name:   name,
		Type:   schemeType,
		Others: others,
	}
}

type Schema struct {
	Name   string
	Type   string
	Others string
}

func (schema *Schema) String() string {
	return fmt.Sprintf("%s: %s %s . \r\n", schema.Name, schema.Type, schema.Others)
}

func (schema *Schema) Equal(src *Schema) bool {
	if schema.Name == src.Name &&
		schema.Type == src.Type &&
		schema.Others == src.Others {
		return true
	}
	return false
}

type Schemas []*Schema

func (schemas Schemas) Add(newSchemes ...*Schema) Schemas {
	list := schemas
	for _, ns := range newSchemes {
		if schemas.Exists(ns) {
			continue
		}
		list = append(list, ns)
	}
	return list
}

func (schemas Schemas) Exists(s *Schema) bool {
	for _, scheme := range schemas {
		if s.Equal(scheme) {
			return true
		}
	}
	return false
}

func (schemas Schemas) String() string {
	all := ""
	for _, s := range schemas {
		all = all + s.String()
	}
	return all
}
