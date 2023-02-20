package local

import (
	"testing"
	"time"
)

func TestCache_Del(t *testing.T) {
	c := New()

	c.Set("k1", "v1", 3*time.Second)
	c.Del("k1")
	t.Log(c.Exists("k1"))
}

func TestCache_Exists(t *testing.T) {
	c := New()

	c.Set("k1", "v1", 3*time.Second)
	t.Log(c.Exists("k1"))
	t.Log(c.Exists("k2"))
}

func TestCache_Keys(t *testing.T) {
	c := New()
	t.Log(c.Keys())
	c.Set("k1", "v1", 3*time.Second)
	t.Log(c.Keys())
}

func TestCache_Flush(t *testing.T) {
	c := New()

	c.Set("k1", "v1", 3*time.Second)
	c.Flush()
	// t.Log(c.Get("k1"))
}

func TestCache_Get(t *testing.T) {
	c := New()

	c.Set("k3", "v3", 1*time.Minute)
	c.Set("k2", "v2", 3*time.Second)

	c.Set("k1", "v1", 2*time.Second)
	c.Set("k0", "v0", 1*time.Second)
	// t.Log(c.Get("k1"))

	time.Sleep(3 * time.Second)
	t.Log(c.Get("k0"))
	t.Log(c.Get("k1"))
	t.Log(c.Get("k2"))
	t.Log(c.Get("k3"))
	// t.Log(c)
}

func TestCache_Set(t *testing.T) {
	c := New()
	c.Set("k1", "v1", 3*time.Second)
	// t.Log(c.Get("k1"))
}

func TestCache_SetMaxMemory(t *testing.T) {
	c := New()
	c.SetMaxMemory("1b")
	t.Log(c)
	c.Set("k1", "v1", 3*time.Second)
	t.Log(c.Get("k1"))
}
