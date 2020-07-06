package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/muesli/cache2go"
)

type testCache struct {
	Cached string
}

func (tc *testCache) UnmarshalCache(data interface{}) {
	tc.Cached = data.(string)
}

func (tc *testCache) MarshalCache() interface{} {
	return tc.Cached
}

func TestInstance(t *testing.T) {
	Instance().Flush()
	//Instance().SetAddedItemCallback(func(entry *cache2go.CacheItem) {
	//	fmt.Println("Added Callback 1:", entry.Key(), entry.Data(), entry.CreatedOn())
	//})
	//Instance().AddAddedItemCallback(func(entry *cache2go.CacheItem) {
	//	fmt.Println("Added Callback 2:", entry.Key(), entry.Data(), entry.CreatedOn())
	//})
	//
	Instance().SetDataLoader(func(key interface{}, args ...interface{}) *cache2go.CacheItem {
		fmt.Println("This is a test with key " + key.(string))
		spew.Dump(args)

		return nil
	})
	//
	//Instance().SetAboutToDeleteItemCallback(func(entry *cache2go.CacheItem) {
	//	fmt.Println("Deleting:", entry.Key(), entry.Data(), entry.CreatedOn())
	//})

	//Instance().Add("k1", 10*time.Second, "ahaha")
	//Instance().Add("k1", 0, "ahaha")
	//
	//data, err := Instance().Value("k1","1111","2222")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//spew.Dump(data.Data())
	//
	//Instance().Delete("k1")

	var tc testCache
	tc.Cached = "hahah"

	MarshalCache("k1", &tc, 0)

	err := UnmarshalCache("k1", &tc)

	if err != nil {
		panic(err)
	}

	spew.Dump(tc.Cached)
}

func TestSaveFile(t *testing.T) {
	Instance().Flush()
	Instance().Add("t1", 0, "t1-value")
	Instance().Add("t2", 0, "t2-value")
	Instance().Add("t3", 0, "t3-value")

	SaveFile("cache.db")

	Instance().Flush()

	LoadFile("cache.db")

	spew.Dump(Instance().Value("t3"))

}

func TestReplace(t *testing.T) {
	Instance().Flush()
	Instance().Add("t1", 0, "1111")
	Instance().Add("t1", 0, "2222")

	spew.Dump(Instance().Value("t1"))
}

func TestCache(t *testing.T) {
	Instance().Flush()
	Instance().Add("e", 0, "e-5555")
	Instance().Add("a", time.Duration(2)*time.Second, "a-1111")
	Instance().Add("c", 0, "c-3333")
	Instance().Add("b", 0, "b-2222")
	Instance().Add("d", 0, "d-4444")

	time.Sleep(3 * time.Second)

	i := NewCacheIndex()

	i.Add("a")
	i.Add("b")
	i.Add("c")
	i.Add("d")

	//i.Insert(0, "111")
	//i.Before("b", "111")
	//i.After("b", "222")

	//i.Replace("e", "333")

	spew.Dump(i.GetCacheByKeys("a", "b"))
	//spew.Dump(i.GetCaches())
}

func TestTable(t *testing.T) {
	Instance().Flush()
	t1 := NewCacheTable().Create("t1", 2*time.Second)
	t1.Add("k1", "1111")
	t1.Add("k2", "2222")
	t1.Add("k3", "3333")

	t1.Before("k3", &TableData{
		Key:  "k3-before",
		Data: "k3-before-1111",
	})

	t1.After("k3", &TableData{
		Key:  "k3-after",
		Data: "k3-after-2222",
	})

	t1.Save()

	time.Sleep(3 * time.Second)

	var isFound = NewCacheTable().TableIf("t1", func(table *Table) {
		println(table.Get("k1").Data.(string))
		println(table.Get("k3").Data.(string))
		println(table.Get("k3-before").Data.(string))

		table.Replace("k2", "k2-replace")
		table.Save(true)
	})

	fmt.Printf("表是否存在: %#v\n", isFound)

	NewCacheTable().TableIf("t1", func(table *Table) {
		spew.Dump(table.Get("k2").Data.(string))
	})
}
