package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/dulumao/cache/drivers"
	"github.com/dulumao/cache/options"
)

func TestNewTable(t *testing.T) {
	c, _ := Use(drivers.DriverCache2go, &options.Options{
		Name: "cache",
	})

	t1 := NewTable(c).Create("t1", 2*time.Second)

	t1.Set("k1", "1111")
	t1.Set("k2", "2222")
	t1.Set("k3", "3333")

	//t1.Before("k3", &TableData{
	//	Key:  "k3-before",
	//	Data: "k3-before-1111",
	//})
	//
	//t1.After("k3", &TableData{
	//	Key:  "k3-after",
	//	Data: "k3-after-2222",
	//})

	t1.Save()

	var isFound = NewTable(c).TableIf("t1", func(table *Table) {
		table.Replace("k2", "k2-replace")
		table.Save(true)
	})

	fmt.Printf("表是否存在: %#v\n", isFound)

	NewTable(c).TableIf("t1", func(table *Table) {
		println(table.Get("k1").Data.(string))
		println(table.Get("k3").Data.(string))
		println(table.Get("k2").Data.(string))
	})
}
