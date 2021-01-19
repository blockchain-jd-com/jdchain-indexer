package cmd

import (
	"fmt"
	"github.com/bluele/gcache"
	"strings"
)

func NewStatusCache(ledger string, size int) *StatusCache {
	cache := &StatusCache{
		Ledger: ledger,
		//Height: height,
		cache: gcache.New(size).LRU().Build(),
	}
	return cache
}

type StatusCache struct {
	Ledger string
	cache  gcache.Cache
}

func (cache *StatusCache) UpdateUidInCache(key, uid string) {
	cache.cache.Set(key, uid)
}

func (cache *StatusCache) UpdateUidsInCache(kv map[string]string) {
	for k, v := range kv {
		cache.cache.Set(k, v)
	}
}

func (cache *StatusCache) GetUidInCache(key string) (v string, ok bool) {
	value, err := cache.cache.Get(key)
	if err != nil {
		return
	}
	v = value.(string)
	ok = true
	return
}

func (cache *StatusCache) PrintKeyValue() {
	fmt.Println(strings.Repeat("-", 100))
	kvs := cache.cache.GetALL()
	for k, v := range kvs {
		fmt.Printf("%s -> %s \n", k, v)
	}
	fmt.Println(strings.Repeat("-", 100))
}
