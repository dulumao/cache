package boltdb

import (
	"bytes"
	"encoding/gob"
	"errors"
	"time"

	"go.etcd.io/bbolt"

	"github.com/dulumao/cache/options"
)

type Cache struct {
	c      *bbolt.DB
	file   string
	bucket []byte
}

type Item struct {
	Data      interface{}
	ExpiredAt time.Duration
	CreatedAt time.Time
}

func New(ops *options.Options) (*Cache, error) {
	var err error
	var db *bbolt.DB
	var bucket = []byte(ops.Name)

	db, err = bbolt.Open(ops.Bboltdb.Path, ops.Bboltdb.Mode, ops.Bboltdb.Options)

	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	})

	return &Cache{
		c:      db,
		bucket: bucket,
	}, nil
}

func (c *Cache) Get(key string) (data interface{}, err error) {
	var d Item

	err = c.c.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(c.bucket)
		bs := b.Get([]byte(key))

		if len(bs) == 0 {
			return errors.New(key + " error")
		}

		buf := bytes.NewBuffer(bs)
		dec := gob.NewDecoder(buf)

		if err := dec.Decode(&d); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if d.CreatedAt.Add(d.ExpiredAt).Before(time.Now()) {
		return nil, errors.New(key + " has Expired")
	}

	return d.Data, err
}

func (c *Cache) MustGet(key string) interface{} {
	d, err := c.Get(key)

	if err != nil {
		panic(err)
	}

	return d
}

func (c *Cache) Set(key string, data interface{}, ttl time.Duration) error {
	var buf = new(bytes.Buffer)
	var enc = gob.NewEncoder(buf)

	var item = Item{
		Data:      data,
		ExpiredAt: ttl,
		CreatedAt: time.Now(),
	}

	if err := enc.Encode(&item); err != nil {
		return err
	}

	return c.c.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(c.bucket)

		return b.Put([]byte(key), buf.Bytes())
	})
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
	err := c.c.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(c.bucket).Delete([]byte(key))
	})

	if err != nil {
		return false
	}

	return true
}

func (c *Cache) Clear() error {
	return c.c.Update(func(tx *bbolt.Tx) error {
		var err error

		err = tx.DeleteBucket(c.bucket)

		if err != nil {
			return err
		}

		// 防止清空后，直接Bucket -> get 出错
		_, err = tx.CreateBucket(c.bucket)

		return err
	})
}

func (c *Cache) Count() int {
	var count = 0

	//c.c.View(func(tx *bbolt.Tx) error {
	//	b := tx.Bucket(c.bucket)
	//
	//	b.ForEach(func(k, v []byte) error {
	//		count++
	//
	//		return nil
	//	})
	//	return nil
	//})

	err := c.c.View(func(tx *bbolt.Tx) error {
		count = tx.Bucket(c.bucket).Stats().KeyN

		return nil
	})

	if err != nil {
		return count
	}

	return count
}

func (c *Cache) Close() error {
	return c.c.Close()
}
