package lru

import (
	"container/list"
	"errors"
	"sync"
	"time"

	"github.com/dulumao/cache/options"
)

type Cache struct {
	size    int
	c       *list.List
	hashMap map[string]*list.Element
	lock    *sync.Mutex
}

type node struct {
	Key   string
	Value interface{}
}

func New(ops *options.Options) (*Cache, error) {
	return &Cache{
		size:    ops.Lru.Size,
		c:       list.New(),
		hashMap: make(map[string]*list.Element),
		lock:    new(sync.Mutex),
	}, nil
}

func (c *Cache) Get(key string) (data interface{}, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if element, ok := c.hashMap[key]; ok {
		c.c.MoveToFront(element)

		n, ok := element.Value.(*node)

		if !ok {
			return nil, errors.New("element type error")
		}

		return n.Value, nil
	}

	return nil, errors.New("element empty")
}

func (c *Cache) MustGet(key string) interface{} {
	d, err := c.Get(key)

	if err != nil {
		panic(err)
	}

	return d
}

func (c *Cache) Set(key string, data interface{}, _ time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if element, isFound := c.hashMap[key]; isFound {
		c.c.MoveToFront(element)

		n, ok := element.Value.(*node)

		if !ok {
			return errors.New("element type error")
		}

		n.Value = data

		return nil
	}

	var newElement = c.c.PushFront(&node{key, data})

	c.hashMap[key] = newElement

	if c.c.Len() > c.size {
		lastElement := c.c.Back()

		if lastElement == nil {
			return nil
		}

		n, ok := lastElement.Value.(*node)

		if !ok {
			return errors.New("element type error")
		}

		delete(c.hashMap, n.Key)

		c.c.Remove(lastElement)
	}

	return nil
}

func (c *Cache) NotFoundAdd(key string, data interface{}, ttl time.Duration) error {
	if !c.Exists(key) {
		return c.Set(key, data, ttl)
	}

	return nil
}

func (c *Cache) Exists(key string) bool {
	_, err := c.Get(key)

	if err != nil {
		return false
	}

	return true
}

func (c *Cache) Delete(key string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.hashMap == nil {
		return false
	}

	if element, ok := c.hashMap[key]; ok {
		n, ok := element.Value.(*node)

		if !ok {
			return false
		}

		delete(c.hashMap, n.Key)

		c.c.Remove(element)

		return true
	}

	return false
}

func (c *Cache) Clear() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.c = list.New()
	c.hashMap = make(map[string]*list.Element)

	return nil
}

func (c *Cache) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.c.Len()
}

func (c *Cache) Close() error {
	return nil
}
