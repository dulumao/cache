package cache

import (
	"errors"
	"time"

	"github.com/dulumao/cache/drivers"
	"github.com/dulumao/cache/drivers/boltdb"
	"github.com/dulumao/cache/drivers/cache2go"
	"github.com/dulumao/cache/drivers/redis"
	"github.com/dulumao/cache/options"
)

const (
	// 永远存在
	Forever = 0
	// 1 分钟
	OneMinutes = 60 * time.Second
	// 2 分钟
	TwoMinutes = 120 * time.Second
	// 3 分钟
	ThreeMinutes = 180 * time.Second
	// 5 分钟
	FiveMinutes = 300 * time.Second
	// 10 分钟
	TenMinutes = 600 * time.Second
	// 半小时
	HalfHour = 1800 * time.Second
	// 1 小时
	OneHour = 3600 * time.Second
	// 2 小时
	TwoHour = 7200 * time.Second
	// 3 小时
	ThreeHour = 10800 * time.Second
	// 12 小时(半天)
	HalfDay = 43200 * time.Second
	// 24 小时(1 天)
	OneDay = 86400 * time.Second
	// 2 天
	TwoDay = 172800 * time.Second
	// 3 天
	ThreeDay = 259200 * time.Second
	// 7 天(一周)
	OneWeek = 604800 * time.Second
)

type ICache interface {
	Get(key string) (data interface{}, err error)
	MustGet(key string) (data interface{})
	Set(key string, data interface{}, ttl time.Duration) error
	NotFoundAdd(key string, data interface{}, ttl time.Duration) error
	Exists(key string) bool
	Delete(key string) bool
	Clear() error
	Count() int
	Close() error
}

type Cache struct {
	ICache
}

func (c *Cache) GetAndSet(key string, fn func() (interface{}, error), ttl time.Duration) (interface{}, error) {
	if c.Exists(key) {
		return c.Get(key)
	}

	var err error
	var data interface{}

	data, err = fn()

	if err != nil {
		return nil, err
	}

	err = c.Set(key, data, ttl)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func Use(driver string, ops *options.Options) (*Cache, error) {
	if ops.Name == "" {
		ops.Name = "cache"
	}

	var c = new(Cache)
	var err error

	if driver == drivers.DriverCache2go {
		c.ICache, err = cache2go.New(ops)

		return c, err
	}

	if driver == drivers.DriverRedis {
		c.ICache, err = redis.New(ops)

		return c, err
	}

	if driver == drivers.DriverBoltdb {
		c.ICache, err = boltdb.New(ops)

		return c, err
	}

	return nil, errors.New("driver error")
}
