package fifo

import (
	"container/list"
	"sync"
	"time"

	"go-cache/internal"

	"go-cache/util"
)

type Cache struct {
	capacity int64 // 最大容量
	length   int64 // 已使用容量
	mu       sync.RWMutex
	list     *list.List
	store    map[string]*list.Element
}

const (
	DefaultCapacity int64 = 1 << (10 * 2)
)

func New() *Cache {
	c := &Cache{
		capacity: DefaultCapacity,
		length:   0,
		store:    make(map[string]*list.Element),
		list:     list.New(),
	}

	return c
}

func (c *Cache) SetMaxMemory(size string) bool {
	userSize, err := util.ParseBytes(size)
	if err != nil {
		return false
	}
	c.capacity = int64(userSize)
	return true
}

func (c *Cache) Set(key string, value interface{}, expire time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.store[key]; ok {

		c.list.MoveToBack(e)
		en := e.Value.(*internal.Entry)
		c.length = c.length - util.CalcLen(en.Value()) + util.CalcLen(value)
		en.SetValue(value)

		return
	}

	en := internal.NewEntry(key, value, 0)
	if expire > 0 {
		en.SetExpiration(time.Now().Add(expire).UnixMicro())
	}
	e := c.list.PushBack(en)
	c.store[key] = e

	c.length += en.Len()
	if c.capacity > 0 && c.length > c.capacity {
		c.delOldest()
	}
}

func (c *Cache) delOldest() {
	c.delete(c.list.Front())
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if e, ok := c.store[key]; ok {
		return e.Value.(*internal.Entry).Value(), ok
	}

	return nil, false
}

func (c *Cache) Del(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.store[key]; ok {
		c.delete(e)
	}
	return true
}

// 检测⼀个值 是否存在
func (c *Cache) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.store[key]
	return ok
}

// 清空所有值
func (c *Cache) Flush() bool {
	c.length = 0
	c.store = make(map[string]*list.Element)
	c.list.Init()
	return true
}

// 返回所有的key 多少
func (c *Cache) Keys() int64 {
	return int64(len(c.store))
}

func (c *Cache) delete(e *list.Element) {
	if e == nil {
		return
	}

	c.list.Remove(e)
	en := e.Value.(*internal.Entry)

	c.length = c.length - en.Len()
	delete(c.store, en.Key())
}
