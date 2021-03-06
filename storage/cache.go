package storage

import (
	"os"

	"github.com/go-redis/redis"
)

// Cache is an interface that allows cache operations to be wrapped and mocked
type Cache interface {
	HGet(key, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	HSet(key, field string, value interface{}) (bool, error)
	SCard(key string) (int64, error)
	SIsMember(key string, member interface{}) (bool, error)
	SMembers(key string) ([]string, error)
	TxPipeline() redis.Pipeliner
	Z(score float64, member interface{}) redis.Z
	ZRangeByScore(key string, opt redis.ZRangeBy) ([]string, error)
}

// NewCache creates a wrapped Redis client
func NewCache() *Redis {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	opts := &redis.Options{
		Addr: addr,
	}
	cache := &Redis{
		Client: redis.NewClient(opts),
	}

	return cache
}
