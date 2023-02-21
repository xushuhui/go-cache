package internal

import (
	"time"

	"go-cache/util"
)

// ⽀持设定过期时间，精度为秒级。
// ⽀持设定最⼤内存，当内存超出时候做出合理的处理。
// ⽀持并发安全。
type Cache interface {
	// size 是⼀个字符串。⽀持以下参数: 1KB，100KB，1MB，2MB，1GB 等
	SetMaxMemory(size string) bool
	// 设置⼀个缓存项，并且在expire时间之后过期
	Set(key string, val interface{}, expire time.Duration)
	// 获取⼀个值
	Get(key string) (interface{}, bool)
	// 删除⼀个值
	Del(key string) bool
	// 检测⼀个值 是否存在
	Exists(key string) bool
	// 情况所有值
	Flush() bool
	// 返回所有的key 多少
	Keys() int64
}

type Entry struct {
	key        string
	value      interface{}
	expiration int64 // 单位毫秒
}

func (v *Entry) Expired() bool {
	if v.expiration == 0 {
		return false
	}
	return time.Now().UnixMicro() > v.expiration
}

func (e *Entry) Len() int64 {
	return util.CalcLen(e.value)
}

func (e *Entry) Value() interface{} {
	return e.value
}

func (e *Entry) Key() string {
	return e.key
}

func (e *Entry) SetValue(value interface{}) {
	e.value = value
}

func NewEntry(key string, value interface{}, expiration int64) *Entry {
	return &Entry{
		key:        key,
		value:      value,
		expiration: expiration,
	}
}

func (v *Entry) SetExpiration(expiration int64) {
	v.expiration = expiration
}
