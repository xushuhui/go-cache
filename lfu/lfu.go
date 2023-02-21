package lfu

import (
	"container/heap"
	"sync"
	"time"

	"go-cache/util"
)

type Cache struct {
	capacity int64 // 最大容量
	length   int64 // 已使用容量
	store    map[string]*entry
	mu       sync.RWMutex
	queue    *queue
}

const (
	DefaultCapacity int64 = 1 << (10 * 2)
)

func New() *Cache {
	q := make(queue, 0, 1024)
	return &Cache{
		capacity: DefaultCapacity,
		length:   0,
		store:    make(map[string]*entry),
		queue:    &q,
	}
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
	if en, ok := c.store[key]; ok {
		c.length = c.length - util.CalcLen(en.value) + util.CalcLen(value)

		c.queue.update(en, value, en.weight+1)
		return
	}
	en := &entry{key: key, value: value}
	if expire > 0 {
		en.expiration = time.Now().Add(expire).UnixMicro()
	}
	heap.Push(c.queue, en)
	c.store[key] = en
	c.length = c.length + en.Len()
	if c.capacity > 0 && c.length > c.capacity {
		c.removeElement(heap.Pop(c.queue))
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if e, ok := c.store[key]; ok {
		c.queue.update(e, e.value, e.weight+1)
		return e.value, true
	}

	return nil, false
}

func (c *Cache) Del(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.store[key]; ok {
		heap.Remove(c.queue, e.index)
		c.removeElement(e)
	}
	return true
}

func (c *Cache) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if e, ok := c.store[key]; ok {
		c.queue.update(e, e.value, e.weight+1)
		return true
	}
	return false
}

func (c *Cache) Flush() bool {
	c.length = 0
	c.store = make(map[string]*entry)
	c.queue = nil
	return true
}

func (c *Cache) Keys() int64 {
	return int64(len(c.store))
}

func (c *Cache) removeElement(x interface{}) {
	if x == nil {
		return
	}
	en := x.(*entry)
	delete(c.store, en.key)
	c.length = c.length - en.Len()
}
