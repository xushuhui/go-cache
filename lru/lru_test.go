package lru

import (
	"testing"
	"time"

	tis "github.com/matryer/is"
)

func TestCache_Del(t *testing.T) {
	is := tis.New(t)
	c := New()
	c.Set("k1", "v1", 3*time.Second)
	is.Equal(c.Del("k1"), true)
	is.Equal(c.Exists("k1"), false)
}

func TestCache_Exists(t *testing.T) {
	is := tis.New(t)
	c := New()
	c.Set("k1", "v1", 3*time.Second)
	is.Equal(c.Exists("k1"), true)
	is.Equal(c.Exists("k2"), false)
}

func TestCache_Keys(t *testing.T) {
	is := tis.New(t)
	c := New()

	c.Set("k1", "v1", 3*time.Second)
	is.Equal(c.Keys(), int64(1))
}

func TestCache_Flush(t *testing.T) {
	is := tis.New(t)
	c := New()

	c.Set("k1", "v1", 3*time.Second)
	is.Equal(c.Keys(), int64(1))
	c.Flush()
	is.Equal(c.Keys(), int64(0))
}

func TestCache_Get(t *testing.T) {
	is := tis.New(t)
	c := New()
	c.Set("string", "string", 1*time.Minute)
	v, ok := c.Get("string")
	is.Equal(v, "string")
	is.Equal(ok, true)
	c.Set("k0", "v0", 1*time.Second)
}

func TestCache_Set(t *testing.T) {
	is := tis.New(t)
	c := New()
	c.Set("string", "string", 1*time.Minute)
	v, ok := c.Get("string")
	is.Equal(v, "string")
	is.Equal(ok, true)
}

func TestCache_SetMaxMemory(t *testing.T) {
	is := tis.New(t)
	c := New()
	c.SetMaxMemory("1kb")

	is.Equal(c.capacity, int64(1024))
}
