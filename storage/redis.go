package storage

import (
	"github.com/go-redis/redis"
)

// Redis wraps redis.Client to avoid the need to call 'Result()' every time
type Redis struct {
	Client *redis.Client
}

// HGetAll wraps the same redis.Client method
func (r *Redis) HGetAll(key string) (map[string]string, error) {
	return r.Client.HGetAll(key).Result()
}

// HSet wraps the same redis.Client method
func (r *Redis) HSet(key, field string, value interface{}) (bool, error) {
	return r.Client.HSet(key, field, value).Result()
}

// Ping wraps the same redis.Client method
func (r *Redis) Ping() (string, error) {
	return r.Client.Ping().Result()
}

// SAdd wraps the same redis.Client method
func (r *Redis) SAdd(key string, members ...interface{}) (int64, error) {
	return r.Client.SAdd(key, members...).Result()
}

// SCard wraps the same redis.Client method
func (r *Redis) SCard(key string) (int64, error) {
	return r.Client.SCard(key).Result()
}

// SIsMember wraps the same redis.Client method
func (r *Redis) SIsMember(key string, member interface{}) (bool, error) {
	return r.Client.SIsMember(key, member).Result()
}

// SMembers wraps the same redis.Client method
func (r *Redis) SMembers(key string) ([]string, error) {
	return r.Client.SMembers(key).Result()
}

// TxPipeline wraps the same redis.Client method
func (r *Redis) TxPipeline() redis.Pipeliner {
	return r.Client.TxPipeline()
}
