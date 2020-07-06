package redis

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/dulumao/cache/options"
)

type Cache struct {
	name string
	c    *redis.Client
	ctx  context.Context
}

func New(ops *options.Options) (*Cache, error) {
	c := redis.NewClient(ops.Redis)

	pong, err := c.Ping().Result()

	if err != nil {
		return nil, err
	}

	_ = pong

	return &Cache{
		c:    c,
		name: ops.Name,
		ctx:  context.Background(),
	}, nil
}

func (c *Cache) getKey(key string) string {
	return strings.Join([]string{c.name, key}, "_")
}

func (c *Cache) Get(key string) (data interface{}, err error) {
	cmd, err := c.c.Get(c.getKey(key)).Result()

	return cmd, err
}

func (c *Cache) MustGet(key string) interface{} {
	d, err := c.Get(key)

	if err != nil {
		panic(err)
	}

	return d
}

func (c *Cache) Set(key string, data interface{}, ttl time.Duration) error {
	_, err := c.c.SetNX(c.getKey(key), data, ttl).Result()

	return err
}

func (c *Cache) NotFoundAdd(key string, data interface{}, ttl time.Duration) error {
	if !c.Exists(key) {
		return c.Set(c.getKey(key), data, ttl)
	}

	return errors.New("key has found")
}

func (c *Cache) Exists(key string) bool {
	d, err := c.c.Exists(c.getKey(key)).Result()

	if err != nil {
		panic(err)
	}

	if d == 1 {
		return true
	}
	return false
}

func (c *Cache) Delete(key string) bool {
	_, err := c.c.Del(c.getKey(key)).Result()

	if err == nil {
		return true
	}

	return false
}

func (c *Cache) Clear() error {
	keys, err := c.c.Keys("cache*").Result()

	if err != nil {
		return err
	}

	for _, k := range keys {
		c.c.Del(k).Result()
	}

	return nil
}

func (c *Cache) Count() int {
	return 0
}

func (c *Cache) Close() error {
	return c.c.Close()
}
