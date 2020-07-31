package cache2go

import (
	"time"

	"github.com/muesli/cache2go"

	"github.com/dulumao/cache/options"
)

type Cache struct {
	c *cache2go.CacheTable
}

func New(ops *options.Options) (*Cache, error) {
	return &Cache{
		c: cache2go.Cache(ops.Name),
	}, nil
}

func (c *Cache) Get(key string) (data interface{}, err error) {
	d, err := c.c.Value(key)

	if err != nil {
		return nil, err
	}

	return d.Data(), nil
}

func (c *Cache) MustGet(key string) interface{} {
	d, err := c.Get(key)

	if err != nil {
		panic(err)
	}

	return d
}

func (c *Cache) Set(key string, data interface{}, ttl time.Duration) error {
	c.c.Add(key, ttl, data)

	return nil
}

func (c *Cache) NotFoundAdd(key string, data interface{}, ttl time.Duration) error {
	c.c.NotFoundAdd(key, ttl, data)

	return nil
}

func (c *Cache) Exists(key string) bool {
	return c.c.Exists(key)
}

func (c *Cache) Delete(key string) bool {
	c.c.Delete(key)

	return true
}

func (c *Cache) Clear() error {
	c.c.Flush()

	return nil
}

func (c *Cache) Count() int {
	return c.c.Count()
}

func (c *Cache) Close() error {
	return nil
}
