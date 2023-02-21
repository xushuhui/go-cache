package lfu

import (
	"container/heap"
	"fmt"
	"testing"
)

func TestQueue_Pop(t *testing.T) {
	q := &queue{}
	heap.Init(q)

	heap.Push(q, &entry{key: "1", value: "1"})
	heap.Push(q, &entry{key: "2", value: "2"})
	heap.Push(q, &entry{key: "3", value: "3"})
	for len(*q) > 0 {
		fmt.Println(heap.Pop(q))
	}

	// heap.Push(&q,"2")
	// heap.Push(&q,"3")
	// heap.Push(&q,"4")
	// t.Log(q)
	// heap.Pop(&q)
	// t.Log(q)
}
