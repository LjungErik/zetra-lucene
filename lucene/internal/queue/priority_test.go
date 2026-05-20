package queue_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LjungErik/zetra-lucene/lucene/internal/queue"
)

type Item struct {
	Priority int
	Value    string
}

func (a Item) Compare(b Item) int {
	return a.Priority - b.Priority
}

func Test_PriorityQueue_Pop_returns_items_in_priority_order(t *testing.T) {
	t.Parallel()

	pq := queue.NewPriorityQueue[Item](4)

	pq.Push(Item{Priority: 3, Value: "low"})
	pq.Push(Item{Priority: 1, Value: "high"})
	pq.Push(Item{Priority: 4, Value: "lowest"})
	pq.Push(Item{Priority: 2, Value: "medium"})

	expected := []Item{
		{Priority: 1, Value: "high"},
		{Priority: 2, Value: "medium"},
		{Priority: 3, Value: "low"},
		{Priority: 4, Value: "lowest"},
	}

	for _, exp := range expected {
		item, ok := pq.Pop()
		require.True(t, ok)
		assert.Equal(t, exp, item)
	}
}

func Test_PriorityQueue_Pop_on_empty_returns_false(t *testing.T) {
	t.Parallel()

	pq := queue.NewPriorityQueue[Item](0)

	_, ok := pq.Pop()
	assert.False(t, ok)
}

func Test_PriorityQueue_Peek_returns_smallest_without_removing(t *testing.T) {
	t.Parallel()

	pq := queue.NewPriorityQueue[Item](2)

	pq.Push(Item{Priority: 5, Value: "five"})
	pq.Push(Item{Priority: 2, Value: "two"})

	item, ok := pq.Peek()
	require.True(t, ok)
	assert.Equal(t, Item{Priority: 2, Value: "two"}, item)
	assert.Equal(t, 2, pq.Len())
}

func Test_PriorityQueue_Peek_on_empty_returns_false(t *testing.T) {
	t.Parallel()

	pq := queue.NewPriorityQueue[Item](0)

	_, ok := pq.Peek()
	assert.False(t, ok)
}

func Test_PriorityQueue_Len_tracks_size(t *testing.T) {
	t.Parallel()

	pq := queue.NewPriorityQueue[Item](3)

	assert.Equal(t, 0, pq.Len())

	pq.Push(Item{Priority: 1, Value: "a"})
	assert.Equal(t, 1, pq.Len())

	pq.Push(Item{Priority: 2, Value: "b"})
	assert.Equal(t, 2, pq.Len())

	pq.Pop()
	assert.Equal(t, 1, pq.Len())
}

func Test_PriorityQueue_handles_equal_priorities(t *testing.T) {
	t.Parallel()

	pq := queue.NewPriorityQueue[Item](3)

	pq.Push(Item{Priority: 1, Value: "first"})
	pq.Push(Item{Priority: 1, Value: "second"})
	pq.Push(Item{Priority: 1, Value: "third"})

	for range 3 {
		item, ok := pq.Pop()
		require.True(t, ok)
		assert.Equal(t, 1, item.Priority)
	}

	assert.Equal(t, 0, pq.Len())
}

func Test_PriorityQueue_push_after_pop(t *testing.T) {
	t.Parallel()

	pq := queue.NewPriorityQueue[Item](2)

	pq.Push(Item{Priority: 3, Value: "three"})
	pq.Push(Item{Priority: 1, Value: "one"})

	item, ok := pq.Pop()
	require.True(t, ok)
	assert.Equal(t, 1, item.Priority)

	pq.Push(Item{Priority: 2, Value: "two"})

	item, ok = pq.Pop()
	require.True(t, ok)
	assert.Equal(t, Item{Priority: 2, Value: "two"}, item)

	item, ok = pq.Pop()
	require.True(t, ok)
	assert.Equal(t, Item{Priority: 3, Value: "three"}, item)
}
