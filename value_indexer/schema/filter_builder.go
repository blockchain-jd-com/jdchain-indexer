package schema

import (
	"github.com/valyala/fasttemplate"
	"strconv"
	"strings"
)

func NewSearchFilterBuilder(prefix string) *SearchFilterBuilder {
	return &SearchFilterBuilder{
		predictionPrefix: prefix,
	}
}

type SearchFilterBuilder struct {
	predictionPrefix string
	fields           []FilterSchema
}

type FilterSchema interface {
	Type() string
	LowerName() string
}

func (builder *SearchFilterBuilder) AddFields(field FilterSchema) {
	builder.fields = append(builder.fields, field)
}

func (builder *SearchFilterBuilder) Build(keyword string) string {
	return builder.buildFields(keyword)
}

func (builder *SearchFilterBuilder) buildFields(keyword string) string {
	var filters []string
	for _, field := range builder.fields {
		filter := builder.build(field, keyword)
		if len(filter) <= 0 {
			continue
		}
		filters = append(filters, filter)
	}
	return strings.Join(filters, " or ")
}

func (builder *SearchFilterBuilder) build(field FilterSchema, keyword string) string {
	templateRegexp := `regexp([[predict_name]], /\S*[[keyword]]\S*/)`
	templateEq := `eq([[predict_name]], [[keyword]])`
	tRegexp, err := fasttemplate.NewTemplate(templateRegexp, "[[", "]]")
	if err != nil {
		return ""
	}

	tEq, err := fasttemplate.NewTemplate(templateEq, "[[", "]]")
	if err != nil {
		return ""
	}
	switch field.Type() {
	case DataTypeString:
		return tRegexp.ExecuteString(map[string]interface{}{
			"predict_name": builder.predictionPrefix + "-" + field.LowerName(),
			"keyword":      keyword,
		})
	case DataTypeInt, DataTypeFloat:
		_, err := strconv.ParseFloat(keyword, 64)
		if err == nil {
			return tEq.ExecuteString(map[string]interface{}{
				"predict_name": builder.predictionPrefix + "-" + field.LowerName(),
				"keyword":      keyword,
			})
		}
	}
	return ""
}
