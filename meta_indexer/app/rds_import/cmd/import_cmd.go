package cmd

import (
	"errors"
	"log"
	"os"
	"sync/atomic"
	"time"

	"git.jd.com/jd-blockchain/explorer/adaptor"
	"git.jd.com/jd-blockchain/explorer/event"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/app/rds_import/tasks"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/app/rds_import/types"
	"git.jd.com/jd-blockchain/explorer/worker"
	"github.com/mkideal/cli"
	"github.com/ssor/zlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"

	gomysql "github.com/go-sql-driver/mysql"
)

var Import = &cli.Command{
	Name: "ledger-import",
	Desc: "fetch data from ledger, and import to mysql",
	Argv: func() interface{} { return new(ImportArgs) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*ImportArgs)
		startServer(argv)
		return nil
	},
}

var totalTask int64
var complteTask int64
var errorTask int64

type ImportArgs struct {
	cli.Helper
	ApiHost   string `cli:"*ledger-host" usage:"api server host" dft:"http://127.0.0.1:8080"`
	Ledger    string `cli:"*ledger" usage:"ledger to search in" dft:""`
	DSN       string `cli:"*dsn" usage:"database dsn"`
	BlockFrom int64  `cli:"from" usage:"start from block" dft:"0"`
	BlockTo   int64  `cli:"to" usage:"stop at block" dft:"-1"`
}

type ImportHandler struct {
	Ledger string
	DB     *gorm.DB
}

func (handler *ImportHandler) ID() string {
	return "import handler"
}

func (handler *ImportHandler) EventReceived(e event.Event) bool {
	switch e.GetName() {
	case event.EventWorkerNoTask:
		zlog.Infof("no task in worker now")
	case event.EventWorkerTaskComplete:
		completeTasks := e.GetData().([]worker.Task)
		processCompleteTasks(completeTasks, handler.DB)
	default:
		zlog.Infof("no handler for event [%s]", e.GetName())
	}
	return true
}

func startServer(cmd *ImportArgs) {

	from, to := initBlockRange(cmd)
	if from > to {
		panic(errors.New("block from is bigger than to"))
	}

	db, err := initDB(cmd.DSN)
	if err != nil {
		panic(err)
	}

	dataWorker := worker.NewConcurrentWorker(cmd.Ledger, 8)
	handler := ImportHandler{
		Ledger: cmd.Ledger,
		DB:     db,
	}
	dataWorker.AddListeners(&handler)

	zlog.Infof("Begin Import Task. From Height: %d, To Height: %d", from, to)

	addAllTasks(dataWorker, cmd.ApiHost, cmd.Ledger, from, to)

	completeCh := make(chan interface{}, 1)
	go waitComplete(completeCh)
	<-completeCh

	zlog.Infof("All Import Task Done. Total Tasks: %d, Error Tasks: %d", atomic.LoadInt64(&totalTask), atomic.LoadInt64(&errorTask))
}

func addAllTasks(dataWorker *worker.ConcurrentWorker, apiHost, ledger string, fromHeight, toHeight int64) {
	addBlockTasks(dataWorker, apiHost, ledger, fromHeight, toHeight)
	addTxTasks(dataWorker, apiHost, ledger, fromHeight, toHeight)
	addContractTasks(dataWorker, apiHost, ledger)
	addUserTasks(dataWorker, apiHost, ledger)
	addDataAccountTasks(dataWorker, apiHost, ledger)
	addEventAccountTasks(dataWorker, apiHost, ledger)
}

func addEventAccountTasks(dataWorker *worker.ConcurrentWorker, apiHost, ledger string) {
	eventAccountCount, err := adaptor.GetTotalEventAccountCountInLedgerFromServer(apiHost, ledger)
	if err == nil {
		eventTasks := tasks.NewEventAccountTasks(apiHost, ledger, eventAccountCount, 10)
		for _, eventTask := range eventTasks {
			//事件任务
			addTask(dataWorker, eventTask)
		}
	}
}

func addDataAccountTasks(dataWorker *worker.ConcurrentWorker, apiHost, ledger string) {
	dataAccountCount, err := adaptor.GetTotalAccountCountInLedgerFromServer(apiHost, ledger)
	if err == nil {
		dataAccountTasks := tasks.NewDataAccountTasks(apiHost, ledger, dataAccountCount, 10)
		for _, dataAccountTask := range dataAccountTasks {
			//数据账户任务
			addTask(dataWorker, dataAccountTask)
		}
	}
}

func addUserTasks(dataWorker *worker.ConcurrentWorker, apiHost, ledger string) {
	userCount, err := adaptor.GetTotalUserCountInLedgerFromServer(apiHost, ledger)
	if err == nil {
		userTasks := tasks.NewUserTasks(apiHost, ledger, userCount, 10)
		for _, userTask := range userTasks {
			//用户任务
			addTask(dataWorker, userTask)
		}
	}
}

func addContractTasks(dataWorker *worker.ConcurrentWorker, apiHost, ledger string) {
	contractCount, err := adaptor.GetTotalContractCountInLedgerFromServer(apiHost, ledger)
	if err == nil {
		contractTasks := tasks.NewContractTasks(apiHost, ledger, contractCount, 10)
		for _, contractTask := range contractTasks {
			//合约任务
			addTask(dataWorker, contractTask)
		}
	}
}

func addTxTasks(dataWorker *worker.ConcurrentWorker, apiHost, ledger string, fromHeight, toHeight int64) {
	for i := fromHeight; i <= toHeight; i++ {
		//交易任务
		txTask := tasks.NewTxTask(apiHost, ledger, i)
		addTask(dataWorker, txTask)
	}
}

func addBlockTasks(dataWorker *worker.ConcurrentWorker, apiHost, ledger string, fromHeight, toHeight int64) {
	for i := fromHeight; i <= toHeight; i++ {
		//区块信息任务
		blockTask := tasks.NewBlockTask(apiHost, ledger, i)
		addTask(dataWorker, blockTask)
	}
}

func initBlockRange(cmd *ImportArgs) (from, to int64) {
	from = cmd.BlockFrom
	to = cmd.BlockTo

	if from <= -1 {
		from = 0
	}

	if to <= -1 {
		ledger, err := adaptor.GetLedgerDetailFromServer(cmd.ApiHost, cmd.Ledger)
		if err != nil {
			panic(err)
		}
		to = ledger.Height
	}
	return
}

func initDB(dsn string) (*gorm.DB, error) {
	cfg, err := gomysql.ParseDSN(dsn)
	if err != nil {
		panic(err)
	}

	cfg.Params = map[string]string{}
	cfg.Params["charset"] = "utf8mb4"
	cfg.Params["parseTime"] = "True"
	cfg.Params["loc"] = "Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             2 * time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	originDb, _ := db.DB()
	originDb.SetConnMaxLifetime(time.Minute * 3)
	originDb.SetMaxOpenConns(10)
	originDb.SetMaxIdleConns(10)

	return db, nil
}

func addTask(dataWorker *worker.ConcurrentWorker, task types.Task) {
	for {
		if dataWorker.AddTask(task) {
			zlog.Infof("add task: %s to worker", task.ID())
			atomic.AddInt64(&totalTask, 1)
			return
		}
		time.Sleep(time.Second)
	}
}

func waitComplete(ch chan interface{}) {
	for {
		total := atomic.LoadInt64(&totalTask)
		complete := atomic.LoadInt64(&complteTask)

		if complete >= total {
			ch <- struct{}{}
			return
		}

		time.Sleep(time.Second)
	}
}

func processCompleteTasks(completeTasks []worker.Task, db *gorm.DB) {
	var blocks []*types.Blocks
	var txs []*types.Txs
	var contracts []*types.Contracts
	var users []*types.Users
	var eventAccounts []*types.EventAccounts
	var events []*types.EventAccountEvents
	var dataAccounts []*types.DataAccounts
	var kvs []*types.DataAccountKVS

	for _, task := range completeTasks {

		t := task.(types.Task)
		result := "success"
		if t.Status() != nil {
			result = t.Status().Error()
			atomic.AddInt64(&errorTask, 1)
		}

		zlog.Infof("task: %s complete, result is: %s ", t.ID(), result)

		if blockTask, ok := task.(*tasks.BlockTask); ok {
			blocks = append(blocks, blockTask.Data())
		}

		if txTask, ok := task.(*tasks.TxTask); ok {
			txs = append(txs, txTask.Data()...)
		}

		if contractTask, ok := task.(*tasks.ContractTask); ok {
			contracts = append(contracts, contractTask.Data()...)
		}

		if userTask, ok := task.(*tasks.UserTask); ok {
			users = append(users, userTask.Data()...)
		}

		if eventTask, ok := task.(*tasks.EventAccountTask); ok {
			eventAccounts = append(eventAccounts, eventTask.Accounts()...)
			events = append(events, eventTask.Events()...)
		}

		if dataAccountTask, ok := task.(*tasks.DataAccountTask); ok {
			dataAccounts = append(dataAccounts, dataAccountTask.Accounts()...)
			kvs = append(kvs, dataAccountTask.KVS()...)
		}
	}

	if len(blocks) > 0 {
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(blocks)
	}

	if len(txs) > 0 {
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(txs)
	}

	if len(contracts) > 0 {
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(contracts)
	}

	if len(users) > 0 {
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(users)
	}

	if len(eventAccounts) > 0 {
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(eventAccounts)
	}

	if len(events) > 0 {
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(events)
	}

	if len(dataAccounts) > 0 {
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(dataAccounts)
	}

	if len(kvs) > 0 {
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(kvs)
	}

	atomic.AddInt64(&complteTask, int64(len(completeTasks)))

}
