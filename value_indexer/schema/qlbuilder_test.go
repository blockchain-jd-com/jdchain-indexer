package schema

import "testing"

func TestQLBuilder_Build(t *testing.T) {
	schema := NewNodeSchema("Product")

	schema.AddField(NewNodeSchemaField("field1", DataTypeString, true),
		NewNodeSchemaField("field2", DataTypeInt, false),
		NewNodeSchemaField("field3", DataTypeInt, false),
		NewNodeSchemaField("field4", DataTypeInt, false))

	builders := []*QLBuilder{
		NewQLBuilder(),
		NewQLBuilder(),
		NewQLBuilder(),
	}
	builders[1].SetRange("name1", "child-node", "range1")
	builders[2].SetRange("name2", "child-node", "range1", "range2")

	for _, builder := range builders {
		ql := builder.Build(schema, "abc")
		t.Log(ql)
	}
}
