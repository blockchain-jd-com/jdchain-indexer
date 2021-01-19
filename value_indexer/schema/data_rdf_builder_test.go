package schema

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"testing"
)

func newFakeCache() *fakeCache {
	return &fakeCache{
		store: make(map[string]string),
	}
}

type fakeCache struct {
	store map[string]string
}

func (fc *fakeCache) SetUidCache(key, value string) {
	fc.store[key] = value
}

func (fc *fakeCache) GetUidInCache(key string) (string, bool) {
	v, ok := fc.store[key]
	return v, ok
}

func TestNewRDFBuilder(t *testing.T) {
	src := `
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

    type Movie{
        id(isIndex: Boolean = true, isPrimaryKey: Boolean = true):                    Int
        popularity(isIndex: Boolean = true):            Float
        release_date(isIndex: Boolean = true):          DateTime
        runtime(isIndex: Boolean = true):               Int
        title(termIndex: Boolean = true):               String
        companies:                                      [Int]
        crew:                                           [Int]
        casts:                                           [Int]
    }


    type EdgeMovieCompany{
        Movie(companies: [Int]): EdgeFrom
        Company(id: Int): EdgeTo
    }

    type EdgeMovieCrew{
        Movie(crew: [Int]): EdgeFrom
        Crew(id: Int): EdgeTo
    }

    type EdgeMovieCast{
        Movie(casts: [Int]): EdgeFrom
        Cast(id: Int): EdgeTo
    }


    `
	css, err := NewSchemaParser().Parse(src)
	assert.Nil(t, err)
	assert.Len(t, css, 7)

	t.Log("---> company: ")
	companyBuilder := NewRDFBuilder(css.FindNodeSchema("company"))
	assert.NotNil(t, companyBuilder)
	src = `
    {
        "production_companies": [
            {
                "name": "Ingenious Film Partners",
                "id": 289
            },
            {
                "name": "Lightstorm Entertainment",
                "id": 574
            }
        ]
    }
    `
	cache := newFakeCache()

	for _, company := range gjson.Parse(src).Get("production_companies").Array() {
		result, success := companyBuilder.Build(company.Raw, cache)
		assert.True(t, success)
		cache.SetUidCache(result.QueryKey, result.UID)
		t.Log("\n", result.Mutations.Assembly())
	}

	t.Log("---> cast: ")
	castBuilder := NewRDFBuilder(css.FindNodeSchema("cast"))
	assert.NotNil(t, castBuilder)
	src = `
    [
        {
            "cast_id": 9,
            "character": "Moat",
            "credit_id": "52fe48009251416c750ac9e5",
            "gender": 1,
            "id": 30484,
            "name": "CCH Pounder",
            "order": 7
        },
        {
            "cast_id": 9,
            "character": "Moat",
            "credit_id": "52fe48009251416c750ac9e5",
            "gender": 1,
            "id": 30485,
            "name": "CCH Pounder",
            "order": 7
        }
    ]
    `
	gjsonResults := gjson.Parse(src).Array()
	for _, gr := range gjsonResults {
		result, success := castBuilder.Build(gr.Raw, cache)
		assert.True(t, success)
		t.Log("\n", result.Mutations.Assembly())
		cache.SetUidCache(result.QueryKey, result.UID)
	}

	t.Log("---> crew: ")
	crewBuilder := NewRDFBuilder(css.FindNodeSchema("crew"))
	assert.NotNil(t, crewBuilder)
	src = `
    [
        {
            "credit_id": "52fe48009251416c750ac9c3",
            "department": "Directing",
            "gender": 2,
            "id": 2710,
            "job": "Director",
            "name": "James Cameron"
        },
        {
            "credit_id": "52fe48009251416c750ac9c3",
            "department": "Directing",
            "gender": 2,
            "id": 2711,
            "job": "Director",
            "name": "James Cameron"
        }
    ]`
	gjsonResults = gjson.Parse(src).Array()
	for _, gr := range gjsonResults {
		result, success := crewBuilder.Build(gr.Raw, cache)
		assert.True(t, success)
		t.Log("\n", result.Mutations.Assembly())
		cache.SetUidCache(result.QueryKey, result.UID)
	}

	t.Log("---> cinema: ")
	src = `
    {
        "id": "19995",
        "popularity": "150.437577",
        "release_date": "2009-12-10",
        "runtime": "162",
        "title": "Avatar",
        "companies": [289, 574],
        "casts": [30484, 30485],
        "crew": [2710,2711]
    }
    `
	relativeSchemas := css.FindRelativeSchemas("movie")
	assert.Len(t, relativeSchemas, 4, spew.Sdump(relativeSchemas))

	cinemaBuilder := NewRDFBuilder(relativeSchemas...)
	assert.NotNil(t, cinemaBuilder)

	result, success := cinemaBuilder.Build(src, cache)
	assert.True(t, success)
	t.Log("\n", result.Mutations.Assembly())

}
