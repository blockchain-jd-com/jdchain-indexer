package cmd

import (
	"fmt"
	"git.jd.com/jd-blockchain/explorer/dgraph_helper"
	"git.jd.com/jd-blockchain/explorer/event"
	"git.jd.com/jd-blockchain/explorer/meta_indexer/links"
	"git.jd.com/jd-blockchain/explorer/worker"
	"github.com/davecgh/go-spew/spew"
	"github.com/mkideal/cli"
	"github.com/ssor/zlog"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var Root = &cli.Command{
	Name: "ledger-rdf",
	Desc: "fetch data from ledger, and generate RDF mutations, or generate an examples of RDFs",
	Argv: func() interface{} { return new(LedgerRDF) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*LedgerRDF)
		startLedgerServer(argv.ApiHost, argv.Ledger, argv.BlockFrom, argv.BlockTo)
		return nil
	},
}

// Root command
type LedgerRDF struct {
	cli.Helper
	ApiHost   string `cli:"ledger-host" usage:"api server host" dft:"http://127.0.0.1:8080"`
	Ledger    string `cli:"l,ledger" usage:"ledger to search in" dft:""`
	BlockFrom int64  `cli:"from" usage:"start from block" dft:"1"`
	BlockTo   int64  `cli:"to" usage:"stop at block" dft:"100"`
}

func startLedgerServer(apiHost, ledger string, from, to int64) {
	dataWorker := worker.NewConcurrentWorker(ledger, 8)
	handler := newLedgerDataHandler(ledger, func(s string) error {
		outputCache.WriteString(s)
		return nil
	})
	dataWorker.AddListeners(handler)
	//
	for i := from; i <= to; i++ {
		task := worker.NewFetchTask(apiHost, ledger, i)
		if dataWorker.AddTask(task) == false {
			time.Sleep(time.Second)
			continue
		}
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	err := saveRDF(outputCache.String())
	if err != nil {
		panic(err)
	}
}

func newLedgerDataHandler(ledger string, rdfSaver func(string) error) *LedgerDataHandler {
	handler := &LedgerDataHandler{
		ledger:   ledger,
		rdfSaver: rdfSaver,
	}
	fakeDB := newFakeDb()
	cache := dgraph_helper.NewUidLruCache(fakeDB, 100)
	cache.UpdateUid("ledger-hash_id||"+ledger, "0x001")
	handler.cache = cache
	return handler
}

type LedgerDataHandler struct {
	ledger   string
	cache    *dgraph_helper.UidLruCache
	rdfSaver func(string) error
}

func (handler *LedgerDataHandler) ID() string {
	return "handler"
}

func (handler *LedgerDataHandler) EventReceived(e event.Event) bool {
	switch e.GetName() {
	case event.EventWorkerNoTask:
		zlog.Infof("no task in worker now")
	case event.EventWorkerTaskComplete:
		zlog.Infof("some task completed")
		tasks := e.GetData().([]worker.Task)
		var ledgerDataList []*worker.LedgerData
		for _, task := range tasks {
			ledgerDataList = append(ledgerDataList, task.(*worker.FetchTask).Data().(*worker.LedgerData))
		}

		for _, ld := range ledgerDataList {
			var builder strings.Builder
			cns := links.ToCommonNode(ld)
			for _, cn := range cns {
				rdfs, err := dgraph_helper.AssembleMutationDatas(handler.cache, cn)
				if err != nil {
					zlog.Errorf("assemble failed: %s", err)
					break
				}
				builder.WriteString(rdfs)
			}
			if err := handler.rdfSaver(builder.String()); err != nil {
				spew.Dump(tasks)
				panic(err)
			}

			nodeLinks := links.ToLinks(ld)
			rdfs, err := dgraph_helper.AssembleMutationDatas(handler.cache, nodeLinks...)
			if err != nil {
				zlog.Errorf("assemble failed: %s", err)
				break
			}
			if err := handler.rdfSaver(rdfs); err != nil {
				spew.Dump(tasks)
				panic(err)
			}
		}

	default:
		zlog.Infof("no handler for event [%s]", e.GetName())
	}
	return true
}

func newFakeDb() *FakeDb {
	return &FakeDb{
		kvs: map[string]string{},
	}
}

type FakeDb struct {
	kvs map[string]string
}

func (db *FakeDb) SetUID(predict, value, uid string) {
	db.kvs[predict+"-"+value] = uid
}

func (db *FakeDb) QueryUID(predict, value string) (uid string, exists bool, e error) {
	//uid, exists = db.kvs[predict+"-"+value]
	//fmt.Printf("get uid from db by %s-%s -> %s\n", predict, value, uid)
	uid = fmt.Sprintf("0x-%s-%s", predict, value)
	exists = true
	return
}

var outputCache strings.Builder

func saveRDF(raw string) error {
	if len(raw) <= 0 {
		return fmt.Errorf("no data to save")
	}

	fileName := fmt.Sprintf("txs_output.rdf")

	err := ioutil.WriteFile(fileName, []byte(raw), os.ModePerm)
	if err != nil {
		zlog.Errorf("write RDF to file failed: %s", err)
		return err
	}
	zlog.Successf("save RDF [%s] success", fileName)
	//zlog.Debug(raw)
	return nil
}
