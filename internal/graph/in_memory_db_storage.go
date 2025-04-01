package graph

import (
	"DMS/internal/dal"
	l "DMS/internal/logger"
	"fmt"
	"io"

	"github.com/redis/go-redis/v9"
)

type inMemoryDBStorage struct {
	client dal.InMemoryDAL
	// The prefix appended to each created key in the database
	prefix []byte
	logger l.Logger
}

// It uses in-memoey databases like Redis.
func NewInMemoryDBStorage(client dal.InMemoryDAL, prefix []byte, logger l.Logger) storage {
	return &inMemoryDBStorage{
		client: client,
		prefix: prefix,
		logger: logger,
	}
}

func (s *inMemoryDBStorage) makeKey(pair Edge) string {
	return fmt.Sprintf("%s:%s:%s", s.prefix, pair.Start, pair.End)
}

func (s *inMemoryDBStorage) Get(key Edge) (bool, bool) {
	val, err := s.client.Get(s.makeKey(key))
	if err == redis.Nil {
		return false, false
	}
	if err != nil {
		return false, false // Handle error appropriately in production
	}
	return val == "1", true
}

func (s *inMemoryDBStorage) Set(key Edge, value bool) error {
	val := "0"
	if value {
		val = "1"
	}
	err := s.client.Set(s.makeKey(key), val)
	if err != nil {
		err = s.client.Set(s.makeKey(key), val)
		if err != nil {
			s.logger.Warnf("Error setting key: %s, value: %s: %s", s.makeKey(key), val, err.Error())
			return err
		}
	}

	size, err2 := s.Size()
	if err2 == nil {
		err2 = fmt.Errorf("")
	}
	s.logger.Debugf("Set key: %s, value: %s, number of keys created so far: %d (%s)",
		s.makeKey(key), val, size, err2.Error())
	return nil
}

func (s *inMemoryDBStorage) Delete(key Edge) {
	s.client.Delete(s.makeKey(key))
}

func (s *inMemoryDBStorage) Clear() error {
	pattern := fmt.Sprintf("%s:*", s.prefix)
	iter, err := s.client.Scan(pattern)
	if err != nil {
		return err
	}
	for {
		val, err2 := iter.Next()
		if err2 != nil {
			if err2 == io.EOF {
				break
			}
			return fmt.Errorf("raised error during iteration action in clearing cache: %s", err2.Error())
		}
		s.client.Delete(val)
	}
	return nil
}

func (s *inMemoryDBStorage) Size() (int, error) {
	return s.client.Size(fmt.Sprintf("%s:*", s.prefix))
}

func (s *inMemoryDBStorage) DeleteByPrefix(start Vertex) error {
	pattern := fmt.Sprintf("%s:%d:*", s.prefix, start)
	iter, err := s.client.Scan(pattern)
	if err != nil {
		return err
	}
	for {
		val, err2 := iter.Next()
		if err2 != nil {
			if err2 == io.EOF {
				break
			}
			return fmt.Errorf("raised error during iteration action in deleting cache: %s", err2.Error())
		}
		s.client.Delete(val)
	}
	return nil
}
