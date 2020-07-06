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

	c, _ := Use(drivers.DriverBoltdb, &options.Options{
		Bboltdb: &options.Bboltdb{
			Path: "cache.db",
			Mode: 0666,
		},
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
