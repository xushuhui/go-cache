package main

import (
	"go-cache/fifo"
	"go-cache/internal"
	"go-cache/local"
	"go-cache/lru"
)

func main() {
	c := local.New()
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

	}
	return local.New()
}
