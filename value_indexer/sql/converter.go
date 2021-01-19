package sql

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/value_indexer/schema"
	"github.com/valyala/fasttemplate"
)

var (
	baseTemplate = `
    {
      nodes(func: has([[node-primary-predict]]))
      {
        [[fields]]
      }
    }
    `
	baseTemplateWithFilter = `
    {
      nodes(func: has([[node-primary-predict]]))
      @filter( [[filter-strings]] )
      {
        [[fields]]
      }
    }
    `
)

func NewConverter(schemaDef *schema.NodeSchema) *Converter {
	c := &Converter{
		schemaDef: schemaDef,
	}
	twithFilter, err := fasttemplate.NewTemplate(baseTemplateWithFilter, "[[", "]]")
	if err != nil {
		logger.Errorf("create fasttemplate failed: %s", err)
		return nil
	}
	c.templateWithFilter = twithFilter

	t, err := fasttemplate.NewTemplate(baseTemplate, "[[", "]]")
	if err != nil {
		logger.Errorf("create fasttemplate failed: %s", err)
		return nil
	}
	c.template = t
	return c
}

type Converter struct {
	schemaDef          *schema.NodeSchema
	template           *fasttemplate.Template
	templateWithFilter *fasttemplate.Template
}

type Filter interface {
	Expression(s string, nodeSchema *schema.NodeSchema) string
}

func (converter *Converter) Do(from string, filter Filter) (s string, ok bool) {
	schemaName := converter.schemaDef.LowerName()
	if schemaName != from {
		return
	}

	// add key/version/time
	fields := "uid\n" + fmt.Sprintf("%s-%s\n", schemaName, "key") + fmt.Sprintf("%s-%s\n", schemaName, "version") + fmt.Sprintf("%s-%s\n", schemaName, "time")
	for i := 0; i < len(converter.schemaDef.Fields()); i++ {
		fields += converter.schemaDef.Fields()[i].FormatPredict(schemaName) + "\n"
	}

	if filter == nil {
		s = converter.template.ExecuteString(map[string]interface{}{
			"node-primary-predict": converter.schemaDef.PrimaryPredict(),
			"fields":               fields,
		})
	} else {
		s = converter.templateWithFilter.ExecuteString(map[string]interface{}{
			"node-primary-predict": converter.schemaDef.PrimaryPredict(),
			"filter-strings":       filter.Expression(from, converter.schemaDef),
			"fields":               fields,
		})
	}
	ok = true
	return
}
