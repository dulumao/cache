package cache

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/dulumao/cache/drivers"
	"github.com/dulumao/cache/options"
)

func TestCreate(t *testing.T) {
	var err error

	//c, _ := Use(drivers.DriverBoltdb, &options.Options{
	//	Bboltdb: &options.Bboltdb{
	//		Path: "cache.db",
	//		Mode: 0666,
	//	},
	//})
	c, err := Use(drivers.DriverCache2go, &options.Options{
		Name: "cache",
	})

	err = c.Set("k1", "v1", 1*time.Second)
	c.Set("k2", "v2", 1*time.Second)
	c.Set("k3", "v3", 1*time.Second)

	if err != nil {
		panic(err)
	}

	//time.Sleep(2*time.Second)

	err = c.Clear()

	//if err != nil {
	//	panic(err)
	//}

	spew.Dump(c.Exists("k1"))
	spew.Dump(c.MustGet("k2"))
}

func TestGetAndSet(t *testing.T) {
	var err error

	c, err := Use(drivers.DriverCache2go, &options.Options{
		Name: "cache",
	})

	spew.Dump(c)

	var k4 interface{}

	k4, err = c.GetAndSet("k4", func() (interface{}, error) {
		return "k4 value", nil
	}, 1*time.Second)

	if err != nil {
		panic(err)
	}

	spew.Dump(k4)

	//time.Sleep(2 * time.Second)

	spew.Dump(c.Exists("k4"))
}

func TestLru(t *testing.T) {
	c, err := Use(drivers.DriverLru, &options.Options{
		Lru: &options.Lru{Size: 10},
	})

	err = c.Set("k1", "v1", 1*time.Second)
	c.Set("k2", "v2", 1*time.Second)
	c.Set("k3", "v3", 1*time.Second)

	if err != nil {
		panic(err)
	}

	//time.Sleep(2*time.Second)

	//err = c.Clear()

	//if err != nil {
	//	panic(err)
	//}

	spew.Dump(c.Exists("k1"))
	spew.Dump(c.MustGet("k2"))
}

func TestDelete(t *testing.T) {
	var err error

	c, err := Use(drivers.DriverCache2go, &options.Options{
		Name: "cache",
	})

	if err != nil {
		panic(err)
	}

	var k4 interface{}

	_ = k4

	//c.Set("k4", "asdad", 10*time.Second)

	k4, err = c.GetAndSet("k4", func() (interface{}, error) {
		return "k4 value", nil
	}, 10*time.Second)

	c.Delete("k4")

	d, _ := c.Get("k4")

	spew.Dump(d)
}
