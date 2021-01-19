package query

import (
	"github.com/valyala/fasttemplate"
	"strings"
)

//type queryParser func(ql *QueryLang, args map[string]interface{}, client *dgo.Dgraph) (*QueryResult, error)

type QueryAssembler interface {
	Assemble(args map[string]interface{}) string
	AssembleForRead(args map[string]interface{}) string
	ResultNames() []string
}

func newQueryLangGroup() *QueryLangGroup {
	return &QueryLangGroup{}
}

type QueryLangGroup struct {
	queries []*QueryLang
}

func (group *QueryLangGroup) addQuery(query *QueryLang) *QueryLangGroup {
	group.queries = append(group.queries, query)
	return group
}

func (group *QueryLangGroup) ResultNames() (names []string) {
	for _, ql := range group.queries {
		names = append(names, ql.resultName)
	}
	return
}

func (group *QueryLangGroup) Assemble(args map[string]interface{}) string {
	if group == nil {
		logger.Warn("group is empty")
		return ""
	}

	var builder strings.Builder
	builder.WriteString("{")
	for _, ql := range group.queries {
		builder.WriteString(ql.innerAssemble(args))
	}
	builder.WriteString("}")
	return builder.String()
}

func (group *QueryLangGroup) AssembleForRead(args map[string]interface{}) string {
	if group == nil {
		logger.Warn("group is empty")
		return ""
	}

	var builder strings.Builder
	builder.WriteString("  {")
	for _, ql := range group.queries {
		builder.WriteString(ql.innerAssembleForRead(args))
	}
	builder.WriteString("}  ")
	return builder.String()
}

func newQueryLang(src, resultName string) *QueryLang {
	return &QueryLang{
		src:        src,
		resultName: resultName,
	}
}

type QueryLang struct {
	src        string
	resultName string
}

func (ql *QueryLang) innerAssembleForRead(args map[string]interface{}) string {
	t, err := fasttemplate.NewTemplate(ql.src, "[[", "]]")
	if err != nil {
		return ""
	}
	return strings.Replace(t.ExecuteString(args), "\n", "", -1)
}

func (ql *QueryLang) AssembleForRead(args map[string]interface{}) string {
	t, err := fasttemplate.NewTemplate("{"+ql.src+"}", "[[", "]]")
	if err != nil {
		return ""
	}
	return strings.Replace(t.ExecuteString(args), "\n", "", -1)
}

func (ql *QueryLang) Assemble(args map[string]interface{}) string {
	t, err := fasttemplate.NewTemplate("{"+ql.src+"}", "[[", "]]")
	if err != nil {
		return ""
	}
	l := t.ExecuteString(args)
	logger.Debugf("dql :\n %s", l)

	return l
}

func (ql *QueryLang) innerAssemble(args map[string]interface{}) string {
	t, err := fasttemplate.NewTemplate(ql.src, "[[", "]]")
	if err != nil {
		return ""
	}
	return t.ExecuteString(args)
}

func (ql *QueryLang) ResultNames() []string {
	return []string{ql.resultName}
}
