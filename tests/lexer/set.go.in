// Simple implementation of a set in Go.

package set

import (
	"fmt"
)

type Set[T comparable] map[T]struct{}

// New creates a new set.
func New[T comparable]() Set[T] {
	return make(Set[T])
}

// Add adds an element to the set.
func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) AddAll(values ...T) {
	for _, value := range values {
		s.Add(value)
	}
}

// Remove removes an element from the set.
func (s Set[T]) Remove(value T) {
	delete(s, value)
}

// Contains checks if the set contains an element.
func (s Set[T]) Contains(value T) bool {
	_, exists := s[value]
	return exists
}

// Size returns the number of elements in the set.
func (s Set[T]) Size() int {
	count := 0
	for range s {
		count++
	}
	return count
}

// Clear removes all elements from the set.
func (s Set[T]) Clear() {
	for key := range s {
		delete(s, key)
	}
}

// Union returns a new set that is the union of two sets.
func (s Set[T]) Union(other Set[T]) Set[T] {
	union := New[T]()
	for key := range s {
		union.Add(key)
	}
	for key := range other {
		union.Add(key)
	}
	return union
}

// Intersection returns a new set that is the intersection of two sets.
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	intersection := New[T]()
	for key := range s {
		if other.Contains(key) {
			intersection.Add(key)
		}
	}
	return intersection
}

// Difference returns a new set that is the difference of two sets.
func (s Set[T]) Difference(other Set[T]) Set[T] {
	difference := New[T]()
	for key := range s {
		if !other.Contains(key) {
			difference.Add(key)
		}
	}
	return difference
}

// IsSubset checks if the set is a subset of another set.
func (s Set[T]) IsSubset(other Set[T]) bool {
	for key := range s {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

// IsSuperset checks if the set is a superset of another set.
func (s Set[T]) IsSuperset(other Set[T]) bool {
	for key := range other {
		if !s.Contains(key) {
			return false
		}
	}
	return true
}

// Elements returns a slice of all elements in the set.
func (s Set[T]) Elements() []T {
	elements := make([]T, 0, s.Size())
	for key := range s {
		elements = append(elements, key)
	}
	return elements
}

// String returns a string representation of the set.
func (s Set[T]) String() string {
	str := "{"
	for key := range s {
		str += fmt.Sprintf("%v ", key)
	}
	str += "}"
	return str
}

// Copy creates a shallow copy of the set.
func (s Set[T]) Copy() Set[T] {
	c := New[T]()
	for key := range s {
		c.Add(key)
	}
	return c
}

// Equal checks if two sets are equal.
func (s Set[T]) Equal(other Set[T]) bool {
	if s.Size() != other.Size() {
		return false
	}
	for key := range s {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

// ForEach applies a function to each element in the set.
func (s Set[T]) ForEach(f func(T)) {
	for key := range s {
		f(key)
	}
}

// Map applies a function to each element in the set and returns a new set.
func (s Set[T]) Map(f func(T) T) Set[T] {
	mapped := New[T]()
	for key := range s {
		mapped.Add(f(key))
	}
	return mapped
}

// Filter returns a new set containing only the elements that satisfy the predicate.
func (s Set[T]) Filter(predicate func(T) bool) Set[T] {
	filtered := New[T]()
	for key := range s {
		if predicate(key) {
			filtered.Add(key)
		}
	}
	return filtered
}

// Reduce applies a function to each element in the set and returns a single value.
func (s Set[T]) Reduce(f func(T, T) T, initial T) T {
	result := initial
	for key := range s {
		result = f(result, key)
	}
	return result
}

// ToSlice converts the set to a slice.
func (s Set[T]) ToSlice() []T {
	slice := make([]T, 0, s.Size())
	for key := range s {
		slice = append(slice, key)
	}
	return slice
}
