package fifo

import "time"

type Cache struct{}

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
