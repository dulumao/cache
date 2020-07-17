package options

import (
	"os"

	"github.com/go-redis/redis/v7"
	"go.etcd.io/bbolt"
)

type Options struct {
	Name    string
	Redis   *redis.Options
	Bboltdb *Bboltdb
	Lru     *Lru
}

type Bboltdb struct {
	Path    string
	Mode    os.FileMode
	Options *bbolt.Options
}

type Lru struct {
	Size int
}
