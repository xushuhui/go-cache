package lru

import (
	"container/list"
	"sync"
	"time"

	"godis/util"
)

type Cache struct {
	capacity int64 // 最大容量
	length   int64 // 已使用容量

	mu    sync.RWMutex
	list  *list.List
	store map[string]*list.Element
}
type entry struct {
	value      interface{}
	expiration int64 // 单位毫秒
}

func (e *entry) Len() int64 {
	return util.CalcLen(e.value)
}

func (c *Cache) SetMaxMemory(size string) bool {
	panic("implement me")
}

func (c *Cache) Set(key string, value interface{}, expire time.Duration) {
	if e, ok := c.store[key]; ok {
		c.list.MoveToBack(e)
		en := e.Value.(*entry)
		c.length = c.length - util.CalcLen(en.value) + util.CalcLen(value)
		en.value = value
		return
	}

	en := &entry{value, 0}
	e := c.list.PushBack(en)
	c.store[key] = e

	c.length += en.Len()
	if c.capacity > 0 && c.length > c.capacity {
		// c.DelOldest()
	}
}

func (c Cache) Get(key string) (interface{}, bool) {
	panic("implement me")
}

func (c Cache) Del(key string) bool {
	panic("implement me")
}

func (c Cache) Exists(key string) bool {
	panic("implement me")
}

func (c Cache) Flush() bool {
	panic("implement me")
}

func (c Cache) Keys() int64 {
	panic("implement me")
}
