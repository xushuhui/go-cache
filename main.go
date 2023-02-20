package main

import (
	"godis/fifo"
	"godis/internal"
	"godis/local"
	"godis/lru"
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
