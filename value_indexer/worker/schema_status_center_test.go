package worker

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	schemaSrc = `
    {
      "associate_account": "672GGY5G2rLQny4aZUcETEstK9uNfpGho6PVt5vyn67Kw",
      "content": "type Cast{     id(isIndex: Boolean = true, isPrimaryKey: Boolean = true):                    Int     cast_id(isIndex: Boolean = true):               Int     character(termIndex: Boolean = true):           String     credit_id(termIndex: Boolean = true):           String     gender(isIndex: Boolean = true):                Int     name(termIndex: Boolean = true):                String     order(isIndex: Boolean = true):                 Int } "
    }
    `
	fakeData = `
     {
        "node": [
            {
                "uid": "0xb4104a",
                "schemainfo-associate_account": "672GGY5G2rLQny4aZUcETEstK9uNfpGho6PVt5vyn67Kw",
                "schemainfo-id": "cast672GGY",
                "schemainfo-status": "0",
                "schemainfo-content": "type Cast{     id(isIndex: Boolean = true, isPrimaryKey: Boolean = true):                    Int     cast_id(isIndex: Boolean = true):               Int     character(termIndex: Boolean = true):           String     credit_id(termIndex: Boolean = true):           String     gender(isIndex: Boolean = true):                Int     name(termIndex: Boolean = true):                String     order(isIndex: Boolean = true):                 Int } "
            }
        ]
    }
    `
)

func TestNewDgraphDataSyncAdd(t *testing.T) {
	dgraphSync := NewFakeDataSync("")
	schemaSync := NewSchemaStatusCenter(dgraphSync)
	err := schemaSync.Prepare()
	assert.Nil(t, err)

	info, err := ParseSchemaInfo(schemaSrc)
	assert.Nil(t, err)

	err = schemaSync.Add(info)
	assert.Nil(t, err)
	err = schemaSync.Add(info)
	assert.NotNil(t, err)
}

func TestNewDgraphDataSyncStartStopDelete(t *testing.T) {
	dgraphSync := NewFakeDataSync(fakeData)
	schemaSync := NewSchemaStatusCenter(dgraphSync)
	err := schemaSync.Prepare()
	assert.Nil(t, err)
	assert.Equal(t, 1, schemaSync.IndexStatusList.Len(), spew.Sdump(schemaSync.IndexStatusList))

	err = schemaSync.Start("cast672GGY")
	assert.Nil(t, err)
	err = schemaSync.Stop("cast672GGY")
	assert.Nil(t, err)

	err = schemaSync.Delete("cast672GGY")
	assert.Nil(t, err)
}

func NewFakeDataSync(src string) *FakeDataSync {
	return &FakeDataSync{
		src: src,
	}
}

type FakeDataSync struct {
	src string
}

func (sync *FakeDataSync) PushDelete(data string) (e error) {
	fmt.Println("------------- delete --------------")
	fmt.Println(data)
	return nil
}

func (sync *FakeDataSync) PushUpdate(data string) (e error) {
	fmt.Println("------------- update --------------")
	fmt.Println(data)
	return nil
}

func (sync *FakeDataSync) Pull() (string, error) {
	return sync.src, nil
}

func (sync *FakeDataSync) SpecifiedSchemaUID(value string) (id string, e error) {
	return
}

func (sync *FakeDataSync) AlterSchema(schemas dgraph_helper.Schemas) error {
	fmt.Println("--------------AlterSchema-------------")
	fmt.Println(schemas.String())
	return nil
}
