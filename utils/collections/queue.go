// Simple queue implementation in Go

package collections

type Queue[T any] struct {
	items []T
}

// NewQueue creates a new queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Enqueue adds an item to the end of the queue
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the item at the front of the queue
func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Peek returns the item at the front of the queue without removing it
func (q *Queue[T]) Peek() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	return q.items[0], true
}

// IsEmpty checks if the queue is empty
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// Size returns the number of items in the queue
func (q *Queue[T]) Size() int {
	return len(q.items)
}

// Clear removes all items from the queue
func (q *Queue[T]) Clear() {
	q.items = []T{}
}

// ToSlice returns the items in the queue as a slice
func (q *Queue[T]) ToSlice() []T {
	return q.items
}
