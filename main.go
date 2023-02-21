package main

import (
	"go-cache/fifo"
	"go-cache/internal"
	"go-cache/lfu"
	"go-cache/local"
	"go-cache/lru"
)

func main() {
	c := NewCache("local")
	c.SetMaxMemory("1MB")
}

func NewCache(policy string) internal.Cache {
	switch policy {
	case "local":
		return local.New()
	case "fifo":
		return fifo.New()
	case "lru":
		return lru.New()
	case "lfu":
		return lfu.New()

	}
	return local.New()
}
