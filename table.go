package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/muesli/cache2go"
)

type Table struct {
	Name     string
	Datas    []*TableData
	LifeSpan time.Duration
	l        sync.Mutex
}

type TableData struct {
	Key  interface{}
	Data interface{}
}

func NewCacheTable() *Table {
	return &Table{
	}
}

func (t *Table) Create(name string, lifeSpan time.Duration) *Table {
	t.Name = name
	t.LifeSpan = lifeSpan

	return t
}

func (t *Table) Save(forces ...bool) {
	var isForce = false

	if len(forces) > 0 {
		isForce = forces[0]
	}

	if Instance().Exists(t.Name) {
		if !isForce {
			return
		}
	}

	Add(t.Name, t.LifeSpan, t.Datas)
}

func (t *Table) Table(name string) (*Table, error) {
	var table *cache2go.CacheItem
	var datas []*TableData
	var err error

	table, err = Value(name)

	if err != nil {
		return nil, err
	}

	datas = table.Data().([]*TableData)

	return &Table{
		Name:     table.Key().(string),
		LifeSpan: table.LifeSpan(),
		Datas:    datas,
	}, nil
}

func (t *Table) TableIf(name string, f func(*Table)) bool {
	var table *Table
	var err error

	table, err = t.Table(name)

	if err == nil {
		f(table)

		return true
	}

	return false
}

func (t *Table) Add(key interface{}, data interface{}) {
	t.l.Lock()
	defer t.l.Unlock()

	t.Datas = append(t.Datas, &TableData{
		Key:  key,
		Data: data,
	})
}

func (t *Table) Get(key interface{}) *TableData {
	t.l.Lock()
	defer t.l.Unlock()

	for idx, k := range t.Datas {
		if k.Key == key {
			return t.Datas[idx]
		}
	}

	return nil
}

func (t *Table) Delete(key interface{}) {
	t.l.Lock()
	defer t.l.Unlock()

	var datas []*TableData

	for idx, _ := range t.Datas {
		if t.Datas[idx] != key {
			datas = append(datas, t.Datas[idx])
		}
	}
}

func (t *Table) FindIndex(key interface{}) int {
	t.l.Lock()
	defer t.l.Unlock()

	for idx, d := range t.Datas {
		if d.Key == key {
			return idx
		}
	}

	return -1
}

func (t *Table) Insert(idx int, td *TableData) error {
	t.l.Lock()
	defer t.l.Unlock()

	if idx > len(t.Datas) {
		return errors.New("长度出错")
	}

	var datas []*TableData

	datas = append(datas, t.Datas[0:idx]...)
	datas = append(datas, td)
	datas = append(datas, t.Datas[idx:]...)

	t.Datas = datas

	return nil
}

func (t *Table) Before(key interface{}, newTableData *TableData) error {
	var idx = t.FindIndex(key)

	if idx != -1 {
		return t.Insert(idx, newTableData)
	}

	return errors.New("key 不存在")
}

func (t *Table) After(key interface{}, newTableData *TableData) error {
	var idx = t.FindIndex(key)

	if idx != -1 {
		return t.Insert(idx+1, newTableData)
	}

	return errors.New("key 不存在")
}

func (t *Table) Replace(key interface{}, data interface{}) error {
	var idx = t.FindIndex(key)

	t.l.Lock()
	defer t.l.Unlock()

	var datas []*TableData

	if idx != -1 {
		datas = append(datas, t.Datas[0:idx]...)
		datas = append(datas, &TableData{
			Key:  key,
			Data: data,
		})
		datas = append(datas, t.Datas[idx+1:]...)

		t.Datas = datas
	}

	return errors.New("key 不存在")
}
