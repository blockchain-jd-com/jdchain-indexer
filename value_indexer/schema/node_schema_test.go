package schema

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSchemaSimple1(t *testing.T) {
	schemaSrc := `
    type City{
        geonameid(isIndex: Boolean = true, isPrimaryKey: Boolean = true):          Int
        name(isIndex: Boolean = true):               String
        subcountrys:         [String]
    }`
	css, err := NewSchemaParser().Parse(schemaSrc)
	assert.Nil(t, err)
	assert.Len(t, css, 1)
	//t.Log(spew.Sdump(css))
	nodes, _ := SplitByType(css...)
	builder := NewSchemaMetaBuilder(nodes[0])
	spew.Dump(nodes[0])
	schemas := builder.Build()
	t.Log("\n", schemas.String())
}

func TestParseSchemaSimple2(t *testing.T) {
	schemaSrcs := []string{`
    type EdgeCinemaCompany{
        Cinema(id: Int): EdgeFrom
        Company(id: Int): EdgeTo
    }
    `, `
    type CinemaCompany2{
		id(isPrimaryKey: Boolean = true):String
		dest:String
		items:String
		source:String
	}
    `}
	for _, schemaSrc := range schemaSrcs {
		css, err := NewSchemaParser().Parse(schemaSrc)
		assert.Nil(t, err, schemaSrc)
		assert.Len(t, css, 1)
	}
	//t.Log(spew.Sdump(css))

	//builder := NewSchemaRDFBuilder(css[0])
	//schemas := builder.Build()
	//t.Log("\n", schemas.String())
	/*
		        type cc-fin01-01{
				id(isPrimaryKey: Boolean = true):String
				dest:String
				items:String
				source:String
			}
	*/
}

func TestNodeSchema_PrimaryPredict(t *testing.T) {
	src := `type Company{
            id(isIndex: Boolean = true, isPrimaryKey: Boolean = true):                   Int
            name(termIndex: Boolean = true):               String }`
	ns, err := NewSchemaParser().FirstNodeSchema(src)
	assert.Nil(t, err)
	t.Log(ns.PrimaryPredict())
}

func TestParseSchemaComplex(t *testing.T) {
	schemaSrc := `
    type Crew{
        id(isIndex: Boolean = true, isPrimaryKey: Boolean = true):                   Int
        name(termIndex: Boolean = true):               String
        gender(isIndex: Boolean = true):               Int
        credit_id(termIndex: Boolean = true):          String
        job(termIndex: Boolean = true):                String
        department(termIndex: Boolean = true):         String
    }

    type Cast{
        id(isIndex: Boolean = true, isPrimaryKey: Boolean = true):                    Int
        cast_id(isIndex: Boolean = true):               Int
        character(termIndex: Boolean = true):           String
        credit_id(termIndex: Boolean = true):           String
        gender(isIndex: Boolean = true):                Int
        name(termIndex: Boolean = true):                String
        order(isIndex: Boolean = true):                 Int
    }

    type Company{
        id(isIndex: Boolean = true, isPrimaryKey: Boolean = true):                   Int
        name(termIndex: Boolean = true):               String
    }

    type Cinema{
        id(isIndex: Boolean = true, isPrimaryKey: Boolean = true):                    Int
        popularity(isIndex: Boolean = true):            Float
        release_date(isIndex: Boolean = true):           DateTime
        runtime(isIndex: Boolean = true):                Int
        title(termIndex: Boolean = true):               String
        companies:                                      [Int]
        crew:                                           [Int]
        cast:                                           [Int]
    }

    type EdgeCinemaCompany{
        Cinema(companies: [Int]): EdgeFrom
        Company(id: Int): EdgeTo
    }

    type EdgeCinemaCrew{
        Cinema(crews: [Int]): EdgeFrom
        Crew(id: Int): EdgeTo
    }

    type EdgeCinemaCast{
        Cinema(casts: [Int]): EdgeFrom
        Cast(id: Int): EdgeTo
    }

    `
	css, err := NewSchemaParser().Parse(schemaSrc)
	assert.Nil(t, err)
	assert.Len(t, css, 10, spew.Sdump(css))

	nodes, _ := SplitByType(css...)
	for _, cs := range nodes {
		builder := NewSchemaMetaBuilder(cs)
		schemas := builder.Build()
		t.Log("\n", schemas.String())
	}
}
