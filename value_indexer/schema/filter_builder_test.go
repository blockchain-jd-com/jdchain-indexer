package schema

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func NewFakeFilterSchema(name, typeName string, isIndex bool) FakeFilterSchema {
	return FakeFilterSchema{
		Name:     name,
		TypeName: typeName,
		isIndex:  isIndex,
	}
}

type FakeFilterSchema struct {
	Name     string
	TypeName string
	isIndex  bool
}

func (field FakeFilterSchema) LowerName() string {
	return strings.ToLower(field.Name)
}

func (field FakeFilterSchema) Type() string {
	return field.TypeName
}

func (field FakeFilterSchema) IsIndex() bool {
	return field.isIndex
}

func TestSchemaFields_BuildQLFilter(t *testing.T) {
	var fields []FilterSchema

	fields = append(fields, NewFakeFilterSchema("field1", DataTypeString, true),
		NewFakeFilterSchema("field2", DataTypeInt, true),
		NewFakeFilterSchema("field3", DataTypeInt, true),
		NewFakeFilterSchema("field4", DataTypeInt, false))

	filterBuilder := NewSearchFilterBuilder("fake")
	for _, f := range fields {
		filterBuilder.AddFields(f)
	}
	result := filterBuilder.Build("abc")
	countOr := strings.Count(result, "or")
	countKeyword := strings.Count(result, "abc")

	t.Log(result)

	assert.Equal(t, 0, countOr, result)
	assert.Equal(t, 1, countKeyword, result)

	filterBuilder = NewSearchFilterBuilder("fake")
	for _, f := range fields {
		filterBuilder.AddFields(f)
	}
	result = filterBuilder.Build("123")
	countOr = strings.Count(result, "or")
	countKeyword = strings.Count(result, "123")

	t.Log(result)

	assert.Equal(t, 2, countOr, result)
	assert.Equal(t, 3, countKeyword, result)
}
