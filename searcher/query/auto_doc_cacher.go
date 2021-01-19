package query

import (
	"github.com/dgraph-io/dgo/v200"
	"github.com/go-ego/riot"
)

type AutoDocCacher struct {
	cacher        *DocCacher
	ledgerMonitor *LedgerMonitor
}

func NewAutoDocCacher(engine *riot.Engine, cacheSize int, dgClient *dgo.Dgraph) *AutoDocCacher {
	adc := &AutoDocCacher{
		cacher: NewDocCacher(engine, cacheSize, dgClient),
	}
	adc.ledgerMonitor = NewLedgerMonitor(adc.cacher, dgClient)
	adc.ledgerMonitor.Run(0)
	return adc
}
