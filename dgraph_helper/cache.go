package dgraph_helper

import (
	"fmt"
	"github.com/bluele/gcache"
)

func NewUidLruCache(uidQuery RemoteUIDQuery, size int) *UidLruCache {
	lru := &UidLruCache{
		uidRemoteQuery: uidQuery,
	}
	if size <= 0 {
		lru.db = gcache.New(1000 * 1000).LRU().Build()
	} else {
		lru.db = gcache.New(size).LRU().Build()
	}
	return lru
}

type RemoteUIDQuery interface {
	QueryUID(predict, value string) (uid string, exists bool, e error)
}

type UidLruCache struct {
	uidRemoteQuery RemoteUIDQuery
	db             gcache.Cache
}

func (lru *UidLruCache) updateUidFromRemote(predict, value string) (uid string, exists bool, e error) {
	uid, exists, e = lru.uidRemoteQuery.QueryUID(predict, value)
	if e != nil {
		return
	}
	if exists == false {
		return
	}
	lru.UpdateUid(lru.FormatCacheKey(predict, value), uid)
	return
}

func (lru *UidLruCache) FormatCacheKey(predict, value string) string {
	return fmt.Sprintf("%s||%s", predict, value)
}

func (lru *UidLruCache) UpdateUid(key, uid string) {
	lru.db.Set(key, uid)
}

func (lru *UidLruCache) QueryUid(predict, value string) (uid string, exists bool, e error) {
	cacheKey := lru.FormatCacheKey(predict, value)
	v, err := lru.db.Get(cacheKey)
	if err != nil {
		if err == gcache.KeyNotFoundError {
			//zlog.Debugf("cache [%s] not found in local cache", cacheKey)
			return lru.updateUidFromRemote(predict, value)
		} else {
			e = err
			return
		}
	}
	exists = true
	uid = v.(string)
	return
}
