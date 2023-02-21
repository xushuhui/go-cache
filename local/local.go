package local

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"go-cache/internal"
	"go-cache/util"
)

type Cache struct {
	capacity      int64 // 最大容量
	length        int64 // 已使用容量
	store         map[string]*internal.Entry
	mu            sync.RWMutex
	waitDel       chan string   //
	clearDuration time.Duration // 定时检查并删除过期缓存
}

const (
	DefaultCapacity         int64         = 1 << (10 * 2)
	DefaultExpiration       time.Duration = 0
	DefaultDelChannelLength               = 100
)

func New() *Cache {
	c := &Cache{
		capacity: DefaultCapacity,
		length:   0,
		store:    make(map[string]*internal.Entry),
		waitDel:  make(chan string, DefaultDelChannelLength),
	}
	go c.work()
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

// Set 设置⼀个缓存项，并且在expire时间之后过期
func (c *Cache) Set(key string, value interface{}, expire time.Duration) {
	if c.length >= c.capacity {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if en, ok := c.store[key]; ok {
		c.length = c.length - util.CalcLen(en.Value()) + util.CalcLen(value)

		en.SetValue(value)
		return
	}

	e := internal.NewEntry(key, value, int64(DefaultExpiration))
	if c.length+e.Len() > c.capacity {
		return
	}
	if expire > 0 {
		e.SetExpiration(time.Now().Add(expire).UnixMicro())
	}

	c.store[key] = e
	c.length = c.length + e.Len()

}

// set 优化版本
func (c *Cache) SetE(key string, val interface{}, expire time.Duration) error {
	if c.length >= c.capacity {
		return errors.New("cache is out of capacity")
	}
	e := internal.NewEntry(key, val, int64(DefaultExpiration))

	if c.length+e.Len() > c.capacity {
		return errors.New("cache is out of capacity")
	}
	if expire > 0 {
		e.SetExpiration(time.Now().Add(expire).UnixMicro())
	}
	c.mu.Lock()
	c.store[key] = e
	c.length = c.length + e.Len()
	c.mu.Unlock()
	return nil
}

// Get 获取⼀个值
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	val, ok := c.get(key)
	if !ok {
		c.mu.RUnlock()
		return nil, false
	}
	if val.Expired() {
		c.waitDel <- key
		c.mu.RUnlock()
		return nil, false
	}
	c.mu.RUnlock()
	return val.Value(), ok
}
func (c *Cache) get(key string) (*internal.Entry, bool) {
	val, ok := c.store[key]
	return val, ok
}

// Del 删除⼀个值
func (c *Cache) Del(key string) bool {
	if v, ok := c.get(key); ok {

		c.mu.Lock()
		deleted := c.delete(key, v)
		c.mu.Unlock()
		return deleted
	}

	return false
}

func (c *Cache) delete(key string, val *internal.Entry) bool {
	delete(c.store, key)
	c.length = c.length - val.Len()
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
	c.store = make(map[string]*internal.Entry)
	return true
}

// 返回所有的key 多少
func (c *Cache) Keys() int64 {
	return int64(len(c.store))
}

func (c *Cache) expireClear() {
	for k, v := range c.store {
		if v.Expired() {
			fmt.Printf("del expire key:%s \n", k)
			c.delete(k, v)
		}
	}
}

func (c *Cache) work() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			c.expireClear()
		case k := <-c.waitDel:
			fmt.Printf("del wait key:%s \n", k)
			c.delete(k, c.store[k])
		}
	}
}
