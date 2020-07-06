package cache

import (
	"errors"
	"sort"
	"sync"
)

type Index struct {
	c    ICache
	keys []string
	l    sync.Mutex
}

func NewIndex(c ICache) *Index {
	return &Index{
		keys: make([]string, 0),
		c:    c,
	}
}

func (i *Index) Add(keys ...string) *Index {
	i.l.Lock()
	defer i.l.Unlock()

	for _, k := range keys {
		i.keys = append(i.keys, k)
	}

	return i
}

func (i *Index) Delete(key string) {
	i.l.Lock()
	defer i.l.Unlock()

	var keys []string

	for idx, _ := range i.keys {
		if i.keys[idx] != key {
			keys = append(keys, i.keys[idx])
		}
	}
}

func (i *Index) Sort() {
	i.l.Lock()
	defer i.l.Unlock()

	sort.Strings(i.keys)
}

func (i *Index) FindIndex(key string) int {
	i.l.Lock()
	defer i.l.Unlock()

	for idx, k := range i.keys {
		if k == key {
			return idx
		}
	}

	return -1
}

func (i *Index) Insert(idx int, key string) error {
	i.l.Lock()
	defer i.l.Unlock()

	if idx > len(i.keys) {
		return errors.New("长度出错")
	}

	var keys []string

	keys = append(keys, i.keys[0:idx]...)
	keys = append(keys, key)
	keys = append(keys, i.keys[idx:]...)

	i.keys = keys

	return nil
}

func (i *Index) Before(key, newKey string) error {
	var idx = i.FindIndex(key)

	if idx != -1 {
		return i.Insert(idx, newKey)
	}

	return errors.New("key 不存在")
}

func (i *Index) After(key, newKey string) error {
	var idx = i.FindIndex(key)

	if idx != -1 {
		return i.Insert(idx+1, newKey)
	}

	return errors.New("key 不存在")
}

func (i *Index) Replace(key, newKey string) error {
	var idx = i.FindIndex(key)

	i.l.Lock()
	defer i.l.Unlock()

	var keys []string

	if idx != -1 {
		keys = append(keys, i.keys[0:idx]...)
		keys = append(keys, newKey)
		keys = append(keys, i.keys[idx+1:]...)

		i.keys = keys
	}

	return errors.New("key 不存在")
}

func (i *Index) GetCacheByKeys(keys ...string) ([]interface{}, []string) {
	var caches []interface{}
	var notFoundKeys []string

	for _, k := range keys {
		if c, err := i.c.Get(k); err == nil {
			caches = append(caches, c)
		} else {
			notFoundKeys = append(notFoundKeys, k)
		}
	}

	return caches, notFoundKeys
}

func (i *Index) GetCaches() ([]interface{}, []string) {
	return i.GetCacheByKeys(i.keys...)
}
