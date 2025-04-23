// Simple stack implementation

package stack

type Stack[T any] struct {
	data []T
}

// New creates a new stack.
func New[T any]() *Stack[T] {
	return &Stack[T]{data: []T{}}
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(value T) {
	s.data = append(s.data, value)
}

// Pop removes and returns the top element of the stack.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.data) == 0 {
		var zeroValue T
		return zeroValue, false
	}
	value := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return value, true
}

// Peek returns the top element of the stack without removing it.
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.data) == 0 {
		var zeroValue T
		return zeroValue, false
	}
	return s.data[len(s.data)-1], true
}

// IsEmpty checks if the stack is empty.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.data) == 0
}

// Size returns the number of elements in the stack.
func (s *Stack[T]) Size() int {
	return len(s.data)
}

// Clear removes all elements from the stack.
func (s *Stack[T]) Clear() {
	s.data = []T{}
}
