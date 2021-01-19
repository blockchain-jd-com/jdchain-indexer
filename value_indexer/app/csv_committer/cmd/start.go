package cmd

import (
	"encoding/csv"
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/value_indexer/schema"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/cli"
	"github.com/ssor/zlog"
	"github.com/tidwall/gjson"
	"io"
	"math"
	"net/http"
	"os"
	"path"
)

var Root = &cli.Command{
	Name: "data",
	Desc: "fetch data from ledger, and generate RDF mutations, and then commit to dgraph",
	Argv: func() interface{} { return new(Server) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*Server)
		return startLedgerServer(argv.DgraphHost, argv.Port)
	},
}

type Server struct {
	cli.Helper
	Port       int    `cli:"p,port" usage:"server listening port" dft:"8888"`
	DgraphHost string `cli:"dgraph" usage:"dgraph server host" dft:"127.0.0.1:9080"`
}

var (
	dgraphHelper     *dgraph_helper.Helper
	definitionsCache []*schema.NodeSchema
	schemas          schema.CommonSchemas
)

func prepareSchemas() schema.CommonSchemas {
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
        runtime(isIndex: Boolean = true):               Float
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
	css, _ := schema.NewSchemaParser().Parse(src)
	return css
}

func startLedgerServer(dgraphHost string, port int) error {
	dgraphHelper = dgraph_helper.NewHelper(dgraphHost)

	schemas = prepareSchemas()
	err := commitSchema(schemas)
	if err != nil {
		return err
	}

	//ops := dgraph_live_client.NewBatchMutaionOptions(batch, concurrent)
	//liveLoader = dgraph_live_client.NewMemoryLoader(zero, dgraphHost, ops)
	router := gin.Default()

	router.POST("/schema", updateSchema)

	router.POST("/data/index", indexInputData)
	router.POST("/data/query", doQuery)
	router.GET("/search", handleSearch)

	return router.Run(fmt.Sprintf("0.0.0.0:%d", port))
}

func doQuery(context *gin.Context) {
	raw, err := context.GetRawData()
	if err != nil {
		zlog.Warnf("get data failed: %s", err)
		context.JSON(http.StatusBadRequest, err)
		return
	}
	result, err := dgraphHelper.QueryObj(string(raw))
	if err != nil {
		zlog.Warnf("do query failed: %s", err)
		context.JSON(http.StatusBadRequest, err)
		return
	}
	context.Data(http.StatusOK, "application/json", result)
}
func handleSearch(context *gin.Context) {
	keyword := context.Query("kw")
	if len(keyword) <= 0 {
		context.JSON(http.StatusBadRequest, "no keyword found")
		return
	}
	builder := schema.NewQLBuilder()
	//ql := definitionsCache[0].GenerateGraphQL(keyword)
	ql := builder.Build(definitionsCache[0], keyword)
	context.Data(http.StatusOK, "json", []byte(ql))
}

var (
	uidCache = NewStatusCache("abc", 1000*1000)
)

func indexCompanyInfo(inputFile string, startLine, endLine int) error {
	companyBuilder := schema.NewRDFBuilder(schemas.FindNodeSchema("company"))

	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}

	r := csv.NewReader(in)
	r.Comma = ','
	//r.LazyQuotes = true
	max := endLine

	if max <= 0 {
		max = math.MaxInt64
	}
	var buildResults []schema.RDFBuildResult
	min := startLine - 1
	index := 0
	for {
		if index > max {
			break
		}

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if index == 0 {
			fmt.Println("csv info:")
			for i, f := range record {
				fmt.Println(i, " -> ", f)
			}
		}
		if len(record) != 20 {
			for i, f := range record {
				fmt.Println(i, " -> ", f)
			}
			return fmt.Errorf("there should be 20 fields each record, but it actually has [%d]", len(record))
		}

		if index < min {
			index++
			continue
		}

		results := gjson.Parse(record[9]).Array()
		for _, company := range results {
			result, success := companyBuilder.Build(company.Raw, uidCache)
			if success == false {
				return fmt.Errorf("build company rdf failed, movie id [%s]", record[3])
			}
			buildResults = append(buildResults, result)

		}

		index++
	}

	return commitBuildResult(buildResults)
}

func indexCastInfo(inputFile string, startLine, endLine int) error {
	companyBuilder := schema.NewRDFBuilder(schemas.FindNodeSchema("cast"))

	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}

	r := csv.NewReader(in)
	r.Comma = ','
	//r.LazyQuotes = true
	max := endLine

	if max <= 0 {
		max = math.MaxInt64
	}
	var buildResults []schema.RDFBuildResult
	min := startLine - 1
	index := 0
	for {
		if index > max {
			break
		}

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if index == 0 {
			fmt.Println("csv info:")
			for i, f := range record {
				fmt.Println(i, " -> ", f)
			}
		}
		if len(record) != 4 {
			for i, f := range record {
				fmt.Println(i, " -> ", f)
			}
			return fmt.Errorf("there should be 4 fields each record, but it actually has [%d]", len(record))
		}

		if index < min {
			index++
			continue
		}

		results := gjson.Parse(record[2]).Array()
		for _, company := range results {
			result, success := companyBuilder.Build(company.Raw, uidCache)
			if success == false {
				return fmt.Errorf("build cast rdf failed, movie id [%s]", record[0])
			}
			buildResults = append(buildResults, result)
		}

		index++
	}

	return commitBuildResult(buildResults)
}

func indexCrewInfo(inputFile string, startLine, endLine int) error {
	companyBuilder := schema.NewRDFBuilder(schemas.FindNodeSchema("crew"))

	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}

	r := csv.NewReader(in)
	r.Comma = ','
	//r.LazyQuotes = true
	max := endLine

	if max <= 0 {
		max = math.MaxInt64
	}
	var buildResults []schema.RDFBuildResult
	min := startLine - 1
	index := 0
	for {
		if index > max {
			break
		}

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if index == 0 {
			fmt.Println("csv info:")
			for i, f := range record {
				fmt.Println(i, " -> ", f)
			}
		}
		if len(record) != 4 {
			for i, f := range record {
				fmt.Println(i, " -> ", f)
			}
			return fmt.Errorf("there should be 4 fields each record, but it actually has [%d]", len(record))
		}

		if index < min {
			index++
			continue
		}

		results := gjson.Parse(record[3]).Array()
		for _, company := range results {
			result, success := companyBuilder.Build(company.Raw, uidCache)
			if success == false {
				return fmt.Errorf("build crew rdf failed, movie id [%s]", record[0])
			}
			buildResults = append(buildResults, result)
		}

		index++
	}

	return commitBuildResult(buildResults)
}

func indexMovieInfo(inputFile string, startLine, endLine int) error {
	companyBuilder := schema.NewRDFBuilder(schemas.FindRelativeSchemas("movie")...)
	zlog.Infof("index movie, load data from file %s", inputFile)

	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}

	r := csv.NewReader(in)
	r.Comma = ','
	//r.LazyQuotes = true
	max := endLine

	if max <= 0 {
		max = math.MaxInt64
	}
	var buildResults []schema.RDFBuildResult
	min := startLine - 1
	index := 0
	for {
		if index > max {
			break
		}

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if index == 0 {
			fmt.Println("csv info:")
			for i, f := range record {
				fmt.Println(i, " -> ", f)
			}
		}
		if len(record) != 2 {
			for i, f := range record {
				fmt.Println(i, " -> ", f)
			}
			return fmt.Errorf("there should be 2 fields each record, but it actually has [%d]", len(record))
		}

		if index < min {
			index++
			continue
		}

		result, success := companyBuilder.Build(record[1], uidCache)
		if success == false {
			return fmt.Errorf("build movie rdf failed, movie id [%s]", record[0])
		}
		//fmt.Println("movie rdf:")
		//fmt.Println(result.mutations.Assembly())
		buildResults = append(buildResults, result)
		//results := gjson.Parse(record[1]).Array()
		//for _, movie := range results {
		//}

		index++
	}
	//return nil
	return commitBuildResult(buildResults)
}

func commitBuildResult(buildResults []schema.RDFBuildResult) error {
	groups := splitBuildResults(buildResults, 1000)
	for _, group := range groups {
		var mutations dgraph_helper.Mutations
		for _, result := range group {
			mutations = mutations.Add(result.Mutations...)
		}
		raw := []byte(mutations.Assembly())
		uids, e := dgraphHelper.MutationRdfs(raw)
		if e != nil {
			return fmt.Errorf("commit rdfs to dgraph failed: %s", e)
		}
		uidCache.UpdateUidsInCache(uids)
		//fmt.Println(uids)
	}
	return nil
}

func splitBuildResults(buildResults []schema.RDFBuildResult, size int) (groups [][]schema.RDFBuildResult) {
	max := len(buildResults)
	if max <= size {
		groups = append(groups, buildResults[:])
		return
	}

	for index := 0; ; index += size {
		if max > (index + size) {
			group := buildResults[index : index+size]
			groups = append(groups, group)
		} else {
			groups = append(groups, buildResults[index:])
			return
		}
	}
}

func indexInputData(context *gin.Context) {
	var input struct {
		Credits    string `json:"credits"`
		Movies     string `json:"movies"`
		MoviesInfo string `json:"movies_info"`
	}
	err := context.Bind(&input)
	if err != nil {
		zlog.Warnf("get input para failed: %s", err)
		context.JSON(http.StatusBadRequest, err)
		return
	}
	endLine := 0
	err = indexCompanyInfo(path.Join("data", input.Movies), 2, endLine)
	if err != nil {
		zlog.Warnf("index companies failed: %s", err)
		return
	}
	zlog.Success("index companies success")

	err = indexCastInfo(path.Join("data", input.Credits), 2, endLine)
	if err != nil {
		zlog.Warnf("index cast failed: %s", err)
		return
	}
	zlog.Success("index cast success")

	err = indexCrewInfo(path.Join("data", input.Credits), 2, endLine)
	if err != nil {
		zlog.Warnf("index crew failed: %s", err)
		return
	}
	zlog.Success("index crew success")

	err = indexMovieInfo(path.Join("data", input.MoviesInfo), 2, endLine)
	if err != nil {
		zlog.Warnf("index movie failed: %s", err)
		return
	}
	zlog.Success("index movie success")
	return
}

func updateSchema(context *gin.Context) {
	raw, err := context.GetRawData()
	if err != nil {
		zlog.Warnf("get data failed: %s", err)
		context.JSON(http.StatusBadRequest, err)
		return
	}
	definitions, err := schema.NewSchemaParser().Parse(string(raw))
	if err != nil {
		zlog.Warnf("parse schema failed: %s", err)
		context.JSON(http.StatusBadRequest, err)
		return
	}
	if len(definitions) <= 0 {
		err = fmt.Errorf("no definitions found from: %s", string(raw))
		zlog.Warnf("no definitions found from: %s", err)
		context.JSON(http.StatusBadRequest, err)
		return
	}
	nodes, _ := schema.SplitByType(definitions...)
	definitionsCache = nodes
	builder := schema.NewSchemaMetaBuilder(nodes[0])
	schemas := builder.Build()
	//schemas := definitions[0].ToRDF()
	fmt.Println("schema: ")
	fmt.Println(schemas.String())

	err = dgraphHelper.Alter(schemas)
	if err != nil {
		zlog.Failedf("alter schema failed: %s", err)
		context.JSON(http.StatusInternalServerError, err)
		return
	}
	zlog.Success("update schema success")
	context.JSON(http.StatusOK, "OK")
}

func commitSchema(schemas schema.CommonSchemas) error {
	nodes, _ := schema.SplitByType(schemas...)
	for _, node := range nodes {
		builder := schema.NewSchemaMetaBuilder(node)
		schemas := builder.Build()
		fmt.Println("schema: ")
		fmt.Println(schemas.String())

		err := dgraphHelper.Alter(schemas)
		if err != nil {
			zlog.Failedf("alter schema failed: %s", err)
			return err
		}
		zlog.Successf("commit schema [%s] success", node.Name)
	}
	return nil
}
