package schema

import (
	"fmt"
	"strings"
)

func NewEdgeSchema(name string, from, to EdgeSchemaField) *EdgeSchema {
	return &EdgeSchema{
		name: name,
		from: from,
		to:   to,
	}
}

type EdgeSchemas []*EdgeSchema

func (schemas EdgeSchemas) Find(name string) (l EdgeSchemas) {
	for _, schema := range schemas {
		if schema.IsRelativeWith(name) {
			l = append(l, schema)
		}
	}
	return
}

type EdgeSchema struct {
	name     string
	from, to EdgeSchemaField
}

type NQuad struct {
	Subject string
	Predict string
	Object  string
}

func (schema *EdgeSchema) ToNQuads(fromUID string, toValues []string, uidCache UIDCache) (l []NQuad) {
	for _, to := range toValues {
		predictForward := fmt.Sprintf("%s-%s", strings.ToLower(schema.from.ObjectName), strings.ToLower(schema.to.ObjectName))
		predictBackward := fmt.Sprintf("%s-%s", strings.ToLower(schema.to.ObjectName), strings.ToLower(schema.from.ObjectName))

		toUID := fmt.Sprintf("%s-%s", strings.ToLower(schema.to.ObjectName), to)
		if uidCache != nil {
			v, ok := uidCache.GetUidInCache(toUID)
			if ok {
				toUID = v
			}
		}

		l = append(l, NQuad{Subject: fromUID, Predict: predictForward, Object: toUID})
		l = append(l, NQuad{Subject: toUID, Predict: predictBackward, Object: fromUID})
	}
	return
}

func (schema *EdgeSchema) LowerName() string {
	return strings.ToLower(schema.name)
}

func (schema *EdgeSchema) FormatCacheQueryKey(v string) string {
	return fmt.Sprintf("%s-%s", schema.LowerName(), v)
}

func (schema *EdgeSchema) IsRelativeWith(name string) bool {
	return schema.from.IsEqualWith(name) || schema.to.IsEqualWith(name)
}

func (schema *EdgeSchema) Type() string {
	return SchemaTypeEdge
}

func NewEdgeSchemaField(obj, property string) EdgeSchemaField {
	return EdgeSchemaField{
		ObjectName:   obj,
		PropertyName: property,
	}
}

type EdgeSchemaField struct {
	ObjectName   string
	PropertyName string
}

func (field EdgeSchemaField) IsEqualWith(name string) bool {
	return strings.ToLower(field.ObjectName) == strings.ToLower(name)
}
