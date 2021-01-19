package sql

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlParser_Parse(t *testing.T) {
	spew.Config.Indent = "  "
	parser := Parser{}
	src := "select * from company where myid = 4403 AND myname in ('4404', '4405') or mycount in(4406, 4407)"
	from, tree, err := parser.Parse(src)
	assert.Nil(t, err)
	assert.Equal(t, from, "company")

	src2 := "select * from company where myid = 4403 AND (myname in ('4404', '4405') or mycount in(4406, 4407))"
	_, tree, err = parser.Parse(src2)
	assert.Nil(t, err)
	assert.Equal(t, "eq(myid, 4403) and eq(myname, [4404,4405]) or eq(mycount, [4406,4407])", tree.Expression("", nil))

	src3 := "select * from bank where id like '*001'"
	_, tree, err = parser.Parse(src3)
	assert.Nil(t, err)
	assert.Equal(t, "regexp(id, /^*001$/)", tree.Expression("", nil))
}
