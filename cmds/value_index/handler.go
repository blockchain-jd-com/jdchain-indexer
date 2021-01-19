package value_index

import (
	"git.jd.com/jd-blockchain/explorer/response"
	"git.jd.com/jd-blockchain/explorer/value_indexer/sql"
	"git.jd.com/jd-blockchain/explorer/value_indexer/worker"
	"github.com/gin-gonic/gin"
	"github.com/ssor/zlog"
	"net/http"
)

func initRouter() *gin.Engine {
	router := gin.Default()

	router.POST("schema/querysql", doQuerySql)
	router.POST("schema/query", doQuery)
	router.GET("/schema/list", listSchemas)
	router.GET("/schema/start/:id", startSchema)
	router.GET("/schema/stop/:id", stopSchema)
	router.POST("/schema", addSchema)
	router.DELETE("/schema/:id", deleteSchema)
	router.PUT("/schema/:id", updateSchema)

	return router
}

func stopSchema(context *gin.Context) {
	id := context.Param("id")
	if len(id) <= 0 {
		context.JSON(http.StatusOK, response.NewFailedResponse("id required"))
		return
	}
	err := schemaCenter.Stop(id)
	if err != nil {
		context.JSON(http.StatusOK, response.NewFailedResponse(err.Error()))
		return
	}
	context.JSON(http.StatusOK, response.NewSuccessResponse(nil))
}

func startSchema(context *gin.Context) {
	id := context.Param("id")
	if len(id) <= 0 {
		context.JSON(http.StatusOK, response.NewFailedResponse("id required"))
		return
	}
	err := schemaCenter.Start(id)
	if err != nil {
		context.JSON(http.StatusOK, response.NewFailedResponse(err.Error()))
		return
	}
	context.JSON(http.StatusOK, response.NewSuccessResponse(nil))
}

func deleteSchema(context *gin.Context) {
	id := context.Param("id")
	if len(id) <= 0 {
		context.JSON(http.StatusOK, response.NewFailedResponse("id required"))
		return
	}
	err := schemaCenter.Delete(id)
	if err != nil {
		context.JSON(http.StatusOK, response.NewFailedResponse(err.Error()))
		return
	}

	context.JSON(http.StatusOK, response.NewSuccessResponse(nil))
}

func addSchema(context *gin.Context) {
	bs, err := context.GetRawData()
	if err != nil {
		zlog.Errorf("get raw input failed: %s", err)
		context.JSON(http.StatusOK, response.NewFailedResponse("cannot fetch input data"))
		return
	}
	info, err := worker.ParseSchemaInfo(string(bs))
	if err != nil {
		zlog.Errorf("parse input data failed: %s", err)
		context.JSON(http.StatusOK, response.NewFailedResponse("input data invalid"))
		return
	}

	err = schemaCenter.Add(info)
	if err != nil {
		zlog.Errorf("add schema failed: %s", err)
		context.JSON(http.StatusOK, response.NewFailedResponse(err.Error()))
		return
	}

	context.JSON(http.StatusOK, response.NewSuccessResponse(nil))
}

func updateSchema(context *gin.Context) {
	id := context.Param("id")
	if len(id) <= 0 {
		context.JSON(http.StatusOK, response.NewFailedResponse("id required"))
		return
	}
	bs, err := context.GetRawData()
	if err != nil {
		zlog.Errorf("get raw input failed: %s", err)
		context.JSON(http.StatusOK, response.NewFailedResponse("cannot fetch input data"))
		return
	}
	info, err := worker.ParseSchemaInfo(string(bs))
	if err != nil {
		zlog.Errorf("parse input data failed: %s", err)
		context.JSON(http.StatusOK, response.NewFailedResponse("input data invalid"))
		return
	}

	if id != info.ID {
		context.JSON(http.StatusOK, response.NewFailedResponse("invalid id"))
	}

	err = schemaCenter.Update(info)
	if err != nil {
		zlog.Errorf("update schema failed: %s", err)
		context.JSON(http.StatusOK, response.NewFailedResponse(err.Error()))
		return
	}

	context.JSON(http.StatusOK, response.NewSuccessResponse(nil))
}

func listSchemas(context *gin.Context) {
	context.JSON(http.StatusOK, response.NewSuccessResponse(schemaCenter.IndexStatusList))
}

func doQuery(context *gin.Context) {
	raw, err := context.GetRawData()
	if err != nil {
		zlog.Warnf("get data failed: %s", err)
		context.JSON(http.StatusOK, response.NewFailedResponse(err.Error()))
		return
	}
	result, err := dgraphHelper.QueryObj(string(raw))
	if err != nil {
		zlog.Warnf("do query failed: %s", err)
		context.JSON(http.StatusOK, response.NewFailedResponse(err.Error()))
		return
	}
	//zlog.Info(string(pretty.Pretty(result)))
	context.Data(http.StatusOK, "application/json; charset=utf-8", result)
}

func doQuerySql(context *gin.Context) {
	raw, err := context.GetRawData()
	if err != nil {
		zlog.Warnf("get data failed: %s", err)
		context.JSON(http.StatusOK, response.NewFailedResponse(err.Error()))
		return
	}
	from, tree, err := sql.NewSqlParser().Parse(string(raw))
	if err != nil {
		zlog.Warnf("parse sql failed: %s", err)
		context.JSON(http.StatusOK, response.NewFailedResponse("invalid sql: "+err.Error()))
		return
	}

	nodeSchema := schemaCenter.Find(from)
	if nodeSchema == nil {
		zlog.Warnf("schema %s in sql does not exist", from)
		context.JSON(http.StatusOK, response.NewFailedResponse("schema not found"))
		return
	}
	converter := sql.NewConverter(nodeSchema)
	graphql, ok := converter.Do(from, tree)
	if ok == false {
		zlog.Warnf("convert sql to graphql failed: %s", string(raw))
		context.JSON(http.StatusOK, response.NewFailedResponse("sql unsupported"))
		return
	}
	zlog.Infof("sql( %s ) -> graphQL( %s )", string(raw), graphql)
	result, err := dgraphHelper.QueryObj(graphql)
	if err != nil {
		zlog.Warnf("do query failed: %s", err)
		context.JSON(http.StatusOK, response.NewFailedResponse(err.Error()))
		return
	}
	//zlog.Info(string(pretty.Pretty(result)))
	context.Data(http.StatusOK, "application/json; charset=utf-8", result)
}
