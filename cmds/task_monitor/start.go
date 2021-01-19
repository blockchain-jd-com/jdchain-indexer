package task_monitor

import (
	"crypto/md5"
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/response"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/cli"
	"github.com/tidwall/gjson"
	"net/http"
	"strconv"
)

var (
	dgraphHelper *dgraph_helper.Helper
)

var Root = &cli.Command{
	Name: "task",
	Desc: "running tasks info for monitor",
	Argv: func() interface{} { return new(ServerArgs) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*ServerArgs)
		return StartLedgerServer(argv)
	},
}

type ServerArgs struct {
	cli.Helper
	Port       int    `cli:"p,port" usage:"server listening port" dft:"10005"`
	DgraphHost string `cli:"dgraph" usage:"dgraph server host" dft:"127.0.0.1:9080"`
}

func StartLedgerServer(args *ServerArgs) error {
	listeningHost := "0.0.0.0"
	listeningPort := args.Port

	dgraphHelper = dgraph_helper.NewHelper(args.DgraphHost)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/tasks", allTasks)

	err := r.Run(fmt.Sprintf("%s:%d", listeningHost, listeningPort))
	if err != nil {
		return err
	}
	return nil
}

func allTasks(c *gin.Context) {
	type TaskInfo struct {
		ID     string `json:"id"`
		Status int64  `json:"status"`
	}

	type Tasks struct {
		TotalCount int64      `json:"total_count"`
		TaskList   []TaskInfo `json:"tasks"`
	}
	var tasks Tasks

	qlTask := `
    {
      tasks(func: gt(level-task-level,0), first:%d)
      {
        uid
        level-task-name
		level-task-level
		level-task-ledger
		level-task-block
      }
    }
    `

	qlTaskCount := `
    {
      tasks(func: gt(level-task-level, 0))
      {
        count(uid)
      }
    }
    `
	rawTaskCount, err := dgraphHelper.QueryObj(qlTaskCount)
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(err.Error()))
		return
	}
	resultCount := gjson.ParseBytes(rawTaskCount).Get("tasks.0.count")
	if resultCount.Exists() == false {
		c.JSON(http.StatusOK, response.NewFailedResponse("query task count failed"))
		return
	}

	tasks.TotalCount = resultCount.Int()

	count := 2048

	maxQuery := c.Query("max")
	if len(maxQuery) > 0 {
		max, err := strconv.Atoi(maxQuery)
		if err == nil {
			count = max
		}
	}
	rawTasks, err := dgraphHelper.QueryObj(fmt.Sprintf(qlTask, count))
	if err != nil {
		c.JSON(http.StatusOK, response.NewFailedResponse(err.Error()))
		return
	}
	var taskList []TaskInfo
	resultTasks := gjson.ParseBytes(rawTasks).Get("tasks").Array()
	for _, rt := range resultTasks {
		bs := []byte(rt.Get("level-task-content").String())
		status := rt.Get("level-task-level").Int()
		taskList = append(taskList, TaskInfo{
			ID:     fmt.Sprintf("%x", md5.Sum(bs)),
			Status: status,
		})
	}
	tasks.TaskList = taskList

	res := response.NewSuccessResponse(tasks)

	callback := c.Query("callback")
	if len(callback) > 0 {
		raw, _ := res.MarshalJSON()
		content := fmt.Sprintf("%s(%s)", callback, string(raw))
		c.Data(http.StatusOK, "applicaton/json", []byte(content))
	} else {
		c.JSON(http.StatusOK, res)
	}
}
