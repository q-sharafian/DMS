package common

import "sync"

// A circular queue/buffer
type CircularQueue[T any] struct {
	data []T
	size int
	head int
	tail int
}

func NewCircularQueue[T any](size int) *CircularQueue[T] {
	return &CircularQueue[T]{
		data: make([]T, size),
		size: size,
		head: 0,
		tail: 0,
	}
}

// An iterator interface
type Iterator[T any] interface {
	// Return true if there is a next value
	Next() (T, bool)
}

type Stack[T any] struct {
	data []T
	// size of the stack. If the stack is empty, size = 0
	size int
	mu   *sync.RWMutex
}

// If mu be nil, a new RWMutex will be created
func NewStack[T any](mu *sync.RWMutex) *Stack[T] {
	if mu != nil {
		return &Stack[T]{
			data: make([]T, 0),
			mu:   mu,
			size: 0,
		}
	}
	return &Stack[T]{
		data: make([]T, 0),
		mu:   &sync.RWMutex{},
		size: 0,
	}
}
func (s *Stack[T]) Push(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.size++
	s.data = append(s.data, value)
}

func (s *Stack[T]) Pop() T {
	s.mu.Lock()
	defer s.mu.Unlock()
	value := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	s.size--
	return value
}
func (s *Stack[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.size
}
func (s *Stack[T]) IsEmpty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.size == 0
}
