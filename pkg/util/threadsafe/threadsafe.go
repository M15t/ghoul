/*
Package threadsafe provides an interface for thread-safe functions that can be safely used between multiple goroutines.

See the tsmap and tsslice packages for examples of how to use this interface.
*/
package threadsafe

import (
	"sync"
)

// Locker is an interface that implements the Lock and Unlock methods.
type Locker sync.Locker

// RLocker is an interface that implements the RLock and RUnlock methods.
type RLocker interface {
	RLock()
	RUnlock()
}

// SimpleSafeSlice for simple use
type SimpleSafeSlice[T any] struct {
	mu sync.RWMutex
	v  []T
}

// NewSimpleSlice creates a new SimpleSafeSlice with initial values
func NewSimpleSlice[T any](v []T) *SimpleSafeSlice[T] {
	return &SimpleSafeSlice[T]{v: v}
}

// Len returns the length of the slice
func (s *SimpleSafeSlice[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.v)
}

// Append appends elements to the slice
func (s *SimpleSafeSlice[T]) Append(x ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.v = append(s.v, x...)
}

// All returns all elements in the slice
func (s *SimpleSafeSlice[T]) All() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.v
}

// Get retrieves the element at the specified index
func (s *SimpleSafeSlice[T]) Get(i int) T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.v[i]
}

// RemoveAt removes the element at the specified index
func (s *SimpleSafeSlice[T]) RemoveAt(i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i < 0 || i >= len(s.v) {
		return
	}
	s.v = append(s.v[:i], s.v[i+1:]...)
}

// Filter returns a new slice containing only the elements for which the predicate returns true
func (s *SimpleSafeSlice[T]) Filter(predicate func(T) bool) []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []T
	for _, v := range s.v {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map returns a new slice resulting from applying the function to each element of the slice
func (s *SimpleSafeSlice[T]) Map(transform func(T) T) []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]T, len(s.v))
	for i, v := range s.v {
		result[i] = transform(v)
	}
	return result
}

// Reduce applies the function to each element of the slice, returning a single result
func (s *SimpleSafeSlice[T]) Reduce(initial T, reducer func(T, T) T) T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := initial
	for _, v := range s.v {
		result = reducer(result, v)
	}
	return result
}
