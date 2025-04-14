package db

import (
	l "DMS/internal/logger"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConnDetails struct {
	Addr     string
	Password string
	DB       int
	// Maximum time a key-value would be kept in the cache. (In seconds)
	// Zero means the key-value will never expire.
	Expire time.Duration
}

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
	logger l.Logger
	// Maximum time a key-value would be kept in the cache. (In seconds)
	// Zero means the key-value will never expire.
	expire time.Duration
}

func NewRedisConn(conn *RedisConnDetails, logger l.Logger) *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conn.Addr,
		Password: conn.Password,
		DB:       conn.DB,
	})
	logger.Infof("Created an instance of Redis database \"%s\" ", conn.Addr)
	return &RedisStorage{
		client: rdb,
		ctx:    context.Background(),
		logger: logger,
		expire: conn.Expire,
	}
}

// If the key doesn't exists, returns ("", redis.Nil)
func (s *RedisStorage) Get(key string) (string, error) {
	val, err := s.client.Get(s.ctx, key).Result()
	return val, err
}

func (s *RedisStorage) Set(key, value string) error {
	result := s.client.Set(s.ctx, key, value, s.expire)
	return result.Err()
}

func (s *RedisStorage) Delete(key string) error {
	result := s.client.Del(s.ctx, key)
	return result.Err()
}

// Clear key-values in the cahce that their keys match the pattern.
func (s *RedisStorage) Clear(pattern string) error {
	// pattern := fmt.Sprintf("%s:*", s.prefix)
	iter := s.client.Scan(s.ctx, 0, pattern, 0).Iterator()
	if iter.Err() != nil {
		return iter.Err()
	}
	for iter.Next(s.ctx) {
		s.client.Del(s.ctx, iter.Val())
	}
	return nil
}

// Returns the number of keys that match the pattern.
func (s *RedisStorage) Size(pattern string) (int, error) {
	keys, err := s.client.Keys(s.ctx, pattern).Result()
	return len(keys), err
}

func (s *RedisStorage) Scan(pattern string) (*redis.ScanIterator, error) {
	result := s.client.Scan(s.ctx, 0, pattern, 0)
	return result.Iterator(), result.Err()
}
