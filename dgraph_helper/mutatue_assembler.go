package dgraph_helper

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"strings"
)

type LinkMutationData interface {
	SetUidLeft(uid string)
	SetUidRight(uid string)
	LeftQueryBy() (predict, value string, ok bool)
	RightQueryBy() (predict, value string, ok bool)
	MutationData
}

type Queryable interface {
	QueryBy() (predict, value string)
}

type CommonMutationData interface {
	MutationData
}

type MutationData interface {
	Mutations() Mutations
}

type UidCache interface {
	QueryUid(predict, v string) (uid string, exists bool, err error)
}

func AssembleMutationDatas(uidCache UidCache, datas ...MutationData) (rdfs string, err error) {
	var builder strings.Builder
	for _, data := range datas {
		raw, e := AssembleMutationData(data, uidCache)
		if e != nil {
			err = e
			return
		}
		builder.WriteString(raw)
	}
	rdfs = builder.String()
	return
}

func AssembleMutationData(data MutationData, uidCache UidCache) (rdfs string, err error) {
	switch t := data.(type) {
	case LinkMutationData:
		if predict, val, ok := t.LeftQueryBy(); ok {
			uid, exists, e := uidCache.QueryUid(predict, val)
			if e != nil {
				err = e
				return
			}
			if exists == false {
				err = fmt.Errorf("uid for [%s - %s] not found", predict, val)
				return
			}
			t.SetUidLeft(uid)
		}
		if predict, val, ok := t.RightQueryBy(); ok {
			uid, exists, e := uidCache.QueryUid(predict, val)
			if e != nil {
				err = e
				return
			}
			if exists == false {
				err = fmt.Errorf("uid for [%s - %s] not found", predict, val)
				return
			}
			t.SetUidRight(uid)
		}

		rdfs = t.Mutations().Assembly()
	case CommonMutationData:
		rdfs = t.Mutations().Assembly()
	default:
		logger.Warnf("unknown data for MutateAssembler to handle: %s", spew.Sdump(data))
	}
	return
}
