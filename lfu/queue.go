package lfu

import (
	"container/heap"
	"runtime"

	"go-cache/util"
)

type entry struct {
	key        string
	value      interface{}
	weight     int
	index      int
	expiration int64 // 单位毫秒
}

func (e *entry) Len() int64 {
	if runtime.GOARCH == "amd64" {
		return util.CalcLen(e.value) + 8 + 8
	}
	return util.CalcLen(e.value) + 4 + 4
}

type queue []*entry

func (q queue) Len() int {
	return len(q)
}

func (q queue) Less(i, j int) bool {
	return q[i].weight < q[j].weight
}

func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *queue) Push(x interface{}) {
	n := len(*q)
	en := x.(*entry)
	en.index = n
	*q = append(*q, en)
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	en := old[n-1]
	old[n-1] = nil // avoid memory leak
	en.index = -1  // for safety
	*q = old[0 : n-1]
	return en
}

// update modifies the weight and value of an entry in the queue.
func (q *queue) update(en *entry, value interface{}, weight int) {
	en.value = value
	en.weight = weight
	heap.Fix(q, en.index)
}
