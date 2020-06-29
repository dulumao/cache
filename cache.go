package cache

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/muesli/cache2go"
)

// 格式 func_name_op _ param: value
// getPublicChat_getstatus|username:1,id:2
// getADV_homepage|started_at:12345678,ended_at:1234567
// getADV_homepage

var single *cache2go.CacheTable
var once sync.Once

type cache struct {
	Items []item
}

type item struct {
	Key      interface{}
	Data     interface{}
	LifeSpan time.Duration
}

func New(names ...string) *cache2go.CacheTable {
	var name = "cache"

	if len(names) > 0 {
		name = names[0]
	}

	return cache2go.Cache(name)
}

func Instance() *cache2go.CacheTable {
	once.Do(func() {
		single = New()
	})

	return single
}

func Close() {
	Instance().Flush()
}

type ICache interface {
	UnmarshalCache(interface{})
	MarshalCache() interface{}
}

func SetLogger(logger *log.Logger) {
	Instance().SetLogger(logger)
}

func UnmarshalCache(key interface{}, cache ICache, args ...interface{}) error {
	data, err := Instance().Value(key, args...)

	if err != nil {
		return err
	}

	cache.UnmarshalCache(data.Data())

	return nil
}

func MarshalCache(key interface{}, cache ICache, lifeSpan ...time.Duration) *cache2go.CacheItem {
	if len(lifeSpan) > 0 {
		return Instance().Add(key, lifeSpan[0], cache.MarshalCache())
	}

	return Instance().Add(key, 0, cache.MarshalCache())
}

func MarshalCacheIfNotFound(key interface{}, cache ICache, lifeSpan ...time.Duration) bool {
	if len(lifeSpan) > 0 {
		return Instance().NotFoundAdd(key, lifeSpan[0], cache.MarshalCache())
	}

	return Instance().NotFoundAdd(key, 0, cache.MarshalCache())
}

func Add(key interface{}, lifeSpan time.Duration, data interface{}) *cache2go.CacheItem {
	return Instance().Add(key, lifeSpan, data)
}

func NotFoundAdd(key interface{}, lifeSpan time.Duration, data interface{}) bool {
	return Instance().NotFoundAdd(key, lifeSpan, data)
}

func Value(key interface{}, args ...interface{}) (*cache2go.CacheItem, error) {
	return Instance().Value(key, args...)
}

func MustValue(key interface{}, args ...interface{}) *cache2go.CacheItem {
	var item *cache2go.CacheItem
	var err error

	if item, err = Value(key, args...); err != nil {
		panic(err)
	}

	return item
}

func Foreach(trans func(key interface{}, item *cache2go.CacheItem)) {
	Instance().Foreach(trans)
}

func Count() int {
	return Instance().Count()
}

func Exists(key interface{}) bool {
	return Instance().Exists(key)
}

func Delete(key interface{}) (*cache2go.CacheItem, error) {
	return Instance().Delete(key)
}

func Save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)

	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()

	var cache cache
	var items []item

	gob.Register(&cache)

	Instance().Foreach(func(key interface{}, cacheItem *cache2go.CacheItem) {
		items = append(items, item{
			Key:  cacheItem.Key(),
			Data: cacheItem.Data(),
		})
	})

	cache.Items = items

	err = enc.Encode(&cache)

	return
}

func SaveFile(filename string) error {
	fp, err := os.Create(filename)

	if err != nil {
		return err
	}

	err = Save(fp)

	if err != nil {
		fp.Close()

		return err
	}

	return fp.Close()
}

func Load(r io.Reader) error {
	var cache cache
	dec := gob.NewDecoder(r)

	err := dec.Decode(&cache)

	if err == nil {
		for _, i := range cache.Items {
			if !Instance().Exists(i.Key) {
				Instance().Add(i.Key, i.LifeSpan, i.Data)
			}
		}
	}

	return err
}

func LoadFile(filename string) error {
	fp, err := os.Open(filename)

	if err != nil {
		return err
	}

	err = Load(fp)

	if err != nil {
		fp.Close()

		return err
	}

	return fp.Close()
}
