package queue

import "container/heap"

type Comparer[T any] interface {
	Compare(other T) int
}

type innerHeap[T Comparer[T]] []T

func (h innerHeap[T]) Len() int           { return len(h) }
func (h innerHeap[T]) Less(i, j int) bool { return h[i].Compare(h[j]) < 0 }
func (h innerHeap[T]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *innerHeap[T]) Push(x any) {
	*h = append(*h, x.(T))
}

func (h *innerHeap[T]) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

type PriorityQueue[T Comparer[T]] struct {
	h innerHeap[T]
}

func NewPriorityQueue[T Comparer[T]](capacity int) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		h: make(innerHeap[T], 0, capacity),
	}
	heap.Init(&pq.h)
	return pq
}

func (pq *PriorityQueue[T]) Push(item T) {
	heap.Push(&pq.h, item)
}

func (pq *PriorityQueue[T]) Pop() (T, bool) {
	if len(pq.h) == 0 {
		var zero T
		return zero, false
	}
	return heap.Pop(&pq.h).(T), true
}

func (pq *PriorityQueue[T]) Peek() (T, bool) {
	if len(pq.h) == 0 {
		var zero T
		return zero, false
	}
	return pq.h[0], true
}

func (pq *PriorityQueue[T]) Len() int {
	return len(pq.h)
}
