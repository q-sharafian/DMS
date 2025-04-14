package graph

import (
	e "DMS/internal/error"
	l "DMS/internal/logger"
	"fmt"
	"strings"
	"sync"
)

// memoryStorage implements Storage interface using in-memory map
type memoryStorage struct {
	data   map[string]bool
	mu     sync.RWMutex
	logger l.Logger
}

// It uses memory/ram storage
func NewMemoryStorage(logger l.Logger) storage {
	logger.Infof("Created an instance of in-memory storage")
	return &memoryStorage{
		data:   make(map[string]bool),
		logger: logger,
	}
}

func (s *memoryStorage) makeKey(pair Edge) string {
	return fmt.Sprintf("%s:%s", pair.Start, pair.End)
}
func (s *memoryStorage) Get(key Edge) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, exists := s.data[s.makeKey(key)]
	if !exists {
		return false, e.ErrNotFound
	}
	return val, nil
}

func (s *memoryStorage) Set(key Edge, value bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[s.makeKey(key)] = value
	return nil
}

func (s *memoryStorage) Delete(key Edge) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, s.makeKey(key))
}

func (s *memoryStorage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]bool)
	return nil
}

func (s *memoryStorage) Size() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data), nil
}

func (s *memoryStorage) DeleteByPrefix(start Vertex) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key := range s.data {
		if s.getStartVertex(key).Equals(start) {
			delete(s.data, key)
		}
	}
	return nil
}

// Return start vertex of the given edge
func (s *memoryStorage) getStartVertex(edge string) Vertex {
	i := strings.Index(edge, ":")
	if i == -1 {
		s.logger.Warnf("Invalid key format: %s", edge)
		return nil
	}
	return Vertex(edge[:i])
}
