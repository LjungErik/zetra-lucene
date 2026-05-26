package queue

import "container/list"

type Queue[T any] struct {
	list *list.List
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		list: list.New(),
	}
}

func (q *Queue[T]) Push(item T) {
	q.list.PushBack(item)
}

func (q *Queue[T]) Pop() (T, bool) {
	var zero T

	front := q.list.Front()
	if front == nil {
		return zero, false
	}

	q.list.Remove(front)
	return front.Value.(T), true
}

func (q *Queue[T]) Len() int {
	return q.list.Len()
}

func (q *Queue[T]) IsEmpty() bool {
	return q.list.Len() == 0
}

func (q *Queue[T]) Peak() T {
	var zero T
	front := q.list.Front()
	if front == nil {
		return zero
	}

	return front.Value.(T)
}
