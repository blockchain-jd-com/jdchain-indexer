package sql

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/value_indexer/schema"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestConverter_Do(t *testing.T) {
	src := `type Company{
                id(isPrimaryKey: Boolean = true):    Int
                name:                                String 
            }`
	ns, err := schema.NewSchemaParser().FirstNodeSchema(src)
	assert.Nil(t, err)

	converter := NewConverter(ns)
	parser := Parser{}

	testCases := []struct {
		sql           string
		isFilterEmpty bool
		graphQL       string
	}{
		{
			sql:           "select * from company",
			isFilterEmpty: true,
			graphQL: `
            {
                nodes(func: has(company-id))
                {
                    expand(_all_)
                }
            }
            `,
		},
		{
			sql:           "select * from company where id = 4403 or id in(26281, 26282)",
			isFilterEmpty: false,
			graphQL: `
            {
                nodes(func: has(company-id)) @filter( eq(company-id, 4403) or eq(company-id, [26281,26282]) )
                {
                    expand(_all_)
                }
            }
            `,
		},
		{
			sql:           "select * from company where id in(26281, 26282) or name in ('Columbia', 'Lightstorm')",
			isFilterEmpty: false,
			graphQL: `
            {
              nodes(func: has(company-id))
              @filter( eq(company-id, [26281,26282]) or (regexp(company-name, /^Columbia$/) or regexp(company-name, /^Lightstorm$/)) )
              {
                expand(_all_)
              }
            }
            `,
		},
		{
			sql:           "select * from company where id between 26281 and 26282",
			isFilterEmpty: false,
			graphQL: `
            {
              nodes(func: has(company-id))
              @filter( ge(company-id, 26281) and le(company-id, 26282) )
              {
                expand(_all_)
              }
            }
            `,
		},
		{
			sql:           "select * from company where id > 26281 and id < 26282",
			isFilterEmpty: false,
			graphQL: `
            {
              nodes(func: has(company-id))
              @filter( gt(company-id, 26281) and lt(company-id, 26282) )
              {
                expand(_all_)
              }
            }
            `,
		},
		{
			sql:           "select * from company where id >= 26281 and id <= 26282",
			isFilterEmpty: false,
			graphQL: `
            {
              nodes(func: has(company-id))
              @filter( ge(company-id, 26281) and le(company-id, 26282) )
              {
                expand(_all_)
              }
            }
            `,
		},
		{
			sql:           "select * from company where id <> 26281 and id <> 26282",
			isFilterEmpty: false,
			graphQL: `
            {
              nodes(func: has(company-id))
              @filter( (not eq(company-id, 26281)) and (not eq(company-id, 26282)) )
              {
                expand(_all_)
              }
            }
            `,
		},
		{
			sql:           "select * from company where name like 'Columbi%'",
			isFilterEmpty: false,
			graphQL: `
            {
              nodes(func: has(company-id))
              @filter( regexp(company-name, /^Columbi\S*/) )
              {
                expand(_all_)
              }
            }
            `,
		},
		{
			sql:           "select * from company where name not like 'Columbi%'",
			isFilterEmpty: false,
			graphQL: `
            {
              nodes(func: has(company-id))
              @filter( (not regexp(company-name, /^Columbi\S*/)) )
              {
                expand(_all_)
              }
            }
            `,
		},
		{
			sql:           "select * from company where name like '%olumbia'",
			isFilterEmpty: false,
			graphQL: `
            {
              nodes(func: has(company-id))
              @filter( regexp(company-name, /\S*olumbia$/) )
              {
                expand(_all_)
              }
            }
            `,
		},
		{
			sql:           "select * from company where name like '%olumbi%'",
			isFilterEmpty: false,
			graphQL: `
            {
              nodes(func: has(company-id))
              @filter( regexp(company-name, /\S*olumbi\S*/) )
              {
                expand(_all_)
              }
            }
            `,
		},
	}

	for _, tc := range testCases {
		from, tree, e := parser.Parse(tc.sql)
		assert.Nil(t, e)
		//assert.Equal(t, e, tc.isFilterEmpty, spew.Sdump(tc))
		//if tc.isFilterEmpty == false {
		//	assert.NotNil(t, tree)
		//}
		graphql, success := converter.Do(from, tree)
		assert.True(t, success)
		errMsg := fmt.Sprintf(`
SQL:
   %s
Expect:
   %s
Actual:
   %s 
Tree:
   %s
            `, tc.sql, tc.graphQL, graphql, spew.Sdump(tree))
		assert.Equal(t, pureString(tc.graphQL), pureString(graphql), errMsg)
		//t.Log(graphql)
	}
}

func pureString(src string) string {
	s := strings.Replace(src, " ", "", -1)
	return strings.Replace(s, "\n", "", -1)
}
