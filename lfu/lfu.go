package lfu

import (
	"go-cache/internal"
	"sync"
	"time"
)

type Cache struct {
	capacity int64 // 最大容量
	length   int64 // 已使用容量
	store    map[string]*internal.Entry
	mu       sync.RWMutex
	queue    *queue
}

func (c Cache) SetMaxMemory(size string) bool {
	panic("implement me")
}

func (c Cache) Set(key string, val interface{}, expire time.Duration) {
	panic("implement me")
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

const (
	DefaultCapacity         int64         = 1 << (10 * 2)
	DefaultExpiration       time.Duration = 0
	DefaultDelChannelLength               = 100
)

func New() *Cache {
	q := make(queue, 0, 1024)
	return &Cache{
		capacity: DefaultCapacity,
		length:   0,
		store:    make(map[string]*internal.Entry),
		queue:    &q,
	}
}
