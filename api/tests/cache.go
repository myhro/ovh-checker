package tests

import (
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/mock"
)

// MockedCache is a dummy cache handler
type MockedCache struct {
	mock.Mock
}

// HGet mocks the same redis.Client method
func (m *MockedCache) HGet(key, field string) (string, error) {
	args := m.Called(key, field)
	return args.String(0), args.Error(1)
}

// HGetAll mocks the same redis.Client method
func (m *MockedCache) HGetAll(key string) (map[string]string, error) {
	args := m.Called(key)
	return args[0].(map[string]string), args.Error(1)
}

// HSet mocks the same redis.Client method
func (m *MockedCache) HSet(key, field string, value interface{}) (bool, error) {
	args := m.Called(key, field, value)
	return args.Bool(0), args.Error(1)
}

// SCard mocks the same redis.Client method
func (m *MockedCache) SCard(key string) (int64, error) {
	args := m.Called(key)
	return int64(args.Int(0)), args.Error(1)
}

// SIsMember mocks the same redis.Client method
func (m *MockedCache) SIsMember(key string, member interface{}) (bool, error) {
	args := m.Called(key, member)
	return args.Bool(0), args.Error(1)
}

// SMembers mocks the same redis.Client method
func (m *MockedCache) SMembers(key string) ([]string, error) {
	args := m.Called(key)
	return args[0].([]string), args.Error(1)
}

// TxPipeline mocks the same redis.Client method
func (m *MockedCache) TxPipeline() redis.Pipeliner {
	return nil
}

// Z mocks the redis.Z struct
func (m *MockedCache) Z(score float64, member interface{}) redis.Z {
	return redis.Z{}
}

func (m *MockedCache) ZRangeByScore(key string, opt redis.ZRangeBy) ([]string, error) {
	args := m.Called(key, opt)
	return args[0].([]string), args.Error(1)
}
