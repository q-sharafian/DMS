package dal

import (
	"DMS/internal/db"
	l "DMS/internal/logger"
	"context"
	"io"

	"github.com/redis/go-redis/v9"
)

type InMemoryIterator interface {
	// Returns the next key-value pair. If there are no more key-value pairs, returns `io.EOF` error.
	Next() (string, error)
}

type InMemoryDAL interface {
	// If both returned string and error be nil, means there's not such key
	Get(key string) (*string, error)
	Set(key, value string) error
	Delete(key string) error
	// Clear the key-values in im-memory cache that their keys match the pattern.
	Clear(pattern string) error
	// Returns the number of keys that match the pattern.
	Size(pattern string) (int, error)
	// Returns the keys that match the pattern
	Scan(pattern string) (InMemoryIterator, error)
	// Try deleting multiple times if couldn't delete key-value.
	// If tryTimes be zero, try to delete an entity and if couldn't, return error.
	// If try times be one, try to delete entity and if couldn't, tries one more time to
	// delete and return error if couldn't delete.
	DeleteWithTry(key string, tryTimes int) error
}

type redisInMemoeyDAL struct {
	db     *db.RedisStorage
	logger l.Logger
}

func (r *redisInMemoeyDAL) Clear(pattern string) error {
	return r.db.Clear(pattern)
}

func (r *redisInMemoeyDAL) Delete(key string) error {
	return r.db.Delete(key)
}

func (r *redisInMemoeyDAL) Get(key string) (*string, error) {
	val, err := r.db.Get(key)
	if err == redis.Nil {
		return nil, nil
	}
	return &val, err
}

func (r *redisInMemoeyDAL) Set(key string, value string) error {
	return r.db.Set(key, value)
}

func (r *redisInMemoeyDAL) Size(pattern string) (int, error) {
	return r.db.Size(pattern)
}

func (r *redisInMemoeyDAL) Scan(pattern string) (InMemoryIterator, error) {
	iter, err := r.db.Scan(pattern)
	if err != nil {
		return nil, err
	}
	return &redisInMemoryIterator{iter, context.Background()}, nil
}

func (r *redisInMemoeyDAL) DeleteWithTry(key string, tryTimes int) error {
	err := r.Delete(key)
	if tryTimes <= 0 || err == nil {
		return err
	}

	for i := 0; i < tryTimes; i++ {
		err = r.Delete(key)
		if err == nil {
			return nil
		}
	}
	return err
}
func NewRedisInMemoeyDAL(connDetails *db.RedisConnDetails, logger l.Logger) InMemoryDAL {
	redisClient := db.NewRedisConn(connDetails, logger)
	logger.Infof("Created an instance of Redis in-memory database")
	return &redisInMemoeyDAL{redisClient, logger}
}

type redisInMemoryIterator struct {
	iter *redis.ScanIterator
	ctx  context.Context
}

func (r *redisInMemoryIterator) Next() (string, error) {
	if r.iter.Next(r.ctx) {
		if r.iter.Err() != nil {
			return "", r.iter.Err()
		}
		return r.iter.Val(), nil
	}
	return "", io.EOF
}
