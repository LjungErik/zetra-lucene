package collection

import "container/list"

type Stack[T any] struct {
	list *list.List
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		list: list.New(),
	}
}

func (q *Stack[T]) Push(item T) {
	q.list.PushFront(item)
}

func (q *Stack[T]) Pop() (T, bool) {
	var zero T

	front := q.list.Front()
	if front == nil {
		return zero, false
	}

	q.list.Remove(front)
	return front.Value.(T), true
}

func (q *Stack[T]) Len() int {
	return q.list.Len()
}

func (q *Stack[T]) IsEmpty() bool {
	return q.list.Len() == 0
}

func (q *Stack[T]) Peak() T {
	var zero T
	front := q.list.Front()
	if front == nil {
		return zero
	}

	return front.Value.(T)
}
