package cache

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/dulumao/cache/drivers"
	"github.com/dulumao/cache/options"
)

func TestNewIndex(t *testing.T) {
	c, _ := Use(drivers.DriverCache2go, &options.Options{
		Name: "cache",
	})

	if err := c.Set("k4", "v4", OneMinutes); err != nil {
		panic(err)
	}
	if err := c.Set("k1", "v1", OneMinutes); err != nil {
		panic(err)
	}
	if err := c.Set("k3", "v3", OneMinutes); err != nil {
		panic(err)
	}
	if err := c.Set("k2", "v2", OneMinutes); err != nil {
		panic(err)
	}

	i := NewIndex(c)

	i.Add("k1", "k2", "k3", "k4")

	spew.Dump(i.GetCaches())
}
