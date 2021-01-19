package schema

import (
	"fmt"
	"github.com/valyala/fasttemplate"
	"strings"
)

func NewQLBuilder() *QLBuilder {
	return &QLBuilder{}
}

type QLBuilder struct {
	Name           string
	linkedNodeName string
	Ranges         []string
}

func (builder *QLBuilder) SetRange(name, linkedNodeName string, ranges ...string) {
	builder.Name = name
	builder.linkedNodeName = linkedNodeName
	builder.Ranges = ranges
}

func (builder *QLBuilder) Build(schema *NodeSchema, keyword string) string {
	if len(builder.Ranges) > 0 {
		return builder.buildInRange(schema, keyword)
	} else {
		return builder.buildNoRange(schema, keyword)
	}
}

func (builder *QLBuilder) buildNoRange(schema *NodeSchema, keyword string) string {
	template := `{
        [[query_name]](func: has([[main_predict]])) 
            @filter([[filters]])
            {
                uid
                [[fields]]
            }
        }
    `
	filterBuilder := NewSearchFilterBuilder(strings.ToLower(schema.Name))
	for _, field := range schema.fields {
		filterBuilder.AddFields(field)
	}
	filters := filterBuilder.Build(keyword)
	fields := ""
	for i := 0; i < len(schema.Fields()); i++ {
		fields += schema.Fields()[i].FormatPredict(schema.LowerName()) + "\n"
	}
	t, err := fasttemplate.NewTemplate(template, "[[", "]]")
	if err != nil {
		return ""
	}
	return t.ExecuteString(map[string]interface{}{
		"filters":      filters,
		"query_name":   schema.queryName(),
		"main_predict": fmt.Sprintf("%s-%s", schema.LowerName(), schema.fields[0].LowerName()),
		"fields":       fields,
	})
}

func (builder *QLBuilder) buildInRange(schema *NodeSchema, keyword string) string {
	template := `{
        [[query_name]](func: anyofterms([[range_name]], "[[ranges]]")) @normalize
            {
                [[linked_node]] @filter([[filters]]){
                    uid
                    [[fields]]
                }
            }
        }
    `
	filterBuilder := NewSearchFilterBuilder(strings.ToLower(schema.Name))
	for _, field := range schema.fields {
		filterBuilder.AddFields(field)
	}
	filters := filterBuilder.Build(keyword)
	fields := ""
	for i := 0; i < len(schema.Fields()); i++ {
		fields += schema.Fields()[i].FormatPredict(schema.LowerName()) + "\n"
	}
	t, err := fasttemplate.NewTemplate(template, "[[", "]]")
	if err != nil {
		return ""
	}
	return t.ExecuteString(map[string]interface{}{
		"range_name":   builder.Name,
		"ranges":       strings.Join(builder.Ranges, " "),
		"linked_node":  builder.linkedNodeName,
		"filters":      filters,
		"query_name":   schema.queryName(),
		"main_predict": fmt.Sprintf("%s-%s", schema.LowerName(), schema.fields[0].LowerName()),
		"fields":       fields,
	})
}
