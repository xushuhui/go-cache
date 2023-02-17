package local

import (
	"sync"
	"time"

	"godis/internal"
	"godis/util"
)

type Cache struct {
	maxBytes  int64
	usedBytes int64
	store     map[string]entry
	mu        sync.RWMutex
}
type entry struct {
	value      interface{}
	expiration int64 // 单位毫秒
}

func (e *entry) Len() int64 {
	return util.CalcLen(e.value)
}

const (
	DefaultMaxBytes   int64         = 1 << (10 * 3)
	DefaultExpiration time.Duration = 0
)

func New() internal.Cache {
	return &Cache{
		maxBytes:  DefaultMaxBytes,
		usedBytes: 0,
		store:     make(map[string]entry),
	}
}

func (c *Cache) SetMaxMemory(size string) bool {
	userSize, err := util.ParseBytes(size)
	if err != nil {
		return false
	}
	c.maxBytes = int64(userSize)
	return true
}

// Set 设置⼀个缓存项，并且在expire时间之后过期
func (c *Cache) Set(key string, val interface{}, expire time.Duration) {
	if c.usedBytes >= c.maxBytes {
		return
	}
	e := entry{
		value:      val,
		expiration: int64(DefaultExpiration),
	}
	if c.usedBytes+e.Len() >= c.maxBytes {
		return
	}
	if expire > 0 {
		e.expiration = time.Now().Add(expire).UnixMicro()
	}
	c.mu.Lock()
	c.store[key] = e
	c.usedBytes = c.usedBytes + e.Len()
	c.mu.Unlock()
}

// Get 获取⼀个值
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	val, ok := c.store[key]
	if !ok {
		c.mu.RUnlock()
		return nil, false
	}
	if val.expiration > 0 {
		if time.Now().UnixMicro() > val.expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}
	c.mu.RUnlock()
	return val.value, ok
}

// Del 删除⼀个值
func (c *Cache) Del(key string) bool {

	if c.Exists(key) {
		c.mu.Lock()
		delete(c.store, key)
		c.mu.Unlock()
		return true
	}

	return false
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
	c.usedBytes = 0
	c.store = make(map[string]entry)
	return true
}

// 返回所有的key 多少
func (c *Cache) Keys() int64 {
	return int64(len(c.store))
}
