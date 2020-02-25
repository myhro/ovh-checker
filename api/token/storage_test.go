package token

import (
	"log"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/myhro/ovh-checker/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TokenStorageTestSuite struct {
	suite.Suite

	storage *Storage
	mini    *miniredis.Miniredis
}

func TestTokenStorageTestSuite(t *testing.T) {
	suite.Run(t, new(TokenStorageTestSuite))
}

func (s *TokenStorageTestSuite) SetupTest() {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	s.mini = mr

	opts := &redis.Options{
		Addr: s.mini.Addr(),
	}
	cache := &storage.Redis{
		Client: redis.NewClient(opts),
	}

	s.storage = &Storage{
		Cache: cache,
	}
}

func (s *TokenStorageTestSuite) TearDownTest() {
	s.mini.Close()
}

func (s *TokenStorageTestSuite) TestListAll() {
	id := 1

	authPrefix := prefixes[Auth]
	authToken := NewAuthToken(id, s.storage.Cache)

	sessionPrefix := prefixes[Session]
	sessionToken := NewSessionToken(id, s.storage.Cache)

	hash, err := s.storage.ListAll(id)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), hash[authPrefix], 0)
	assert.Len(s.T(), hash[sessionPrefix], 0)

	err = authToken.Save()
	assert.NoError(s.T(), err)
	err = sessionToken.Save()
	assert.NoError(s.T(), err)

	hash, err = s.storage.ListAll(id)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), hash[authPrefix], 1)
	assert.Len(s.T(), hash[sessionPrefix], 1)
}

func (s *TokenStorageTestSuite) TestListAllError() {
	authPrefix := prefixes[Auth]
	sessionPrefix := prefixes[Session]

	s.mini.Close()

	hash, err := s.storage.ListAll(1)
	assert.Error(s.T(), err)
	assert.Len(s.T(), hash[authPrefix], 0)
	assert.Len(s.T(), hash[sessionPrefix], 0)
}

func (s *TokenStorageTestSuite) TestLoad() {
	token := NewAuthToken(1, s.storage.Cache)
	token.Client = "token-storage-test"
	token.IP = "127.0.0.1"

	err := token.Save()
	assert.NoError(s.T(), err)

	loaded, err := s.storage.Load(token.Type, token.UserID, token.ID)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), token.Key, loaded.Key)
	assert.Equal(s.T(), token.SetKey, loaded.SetKey)
	assert.Equal(s.T(), token.Type, loaded.Type)
	assert.Equal(s.T(), token.UserID, loaded.UserID)

	assert.Equal(s.T(), token.ID, loaded.ID)
	assert.Equal(s.T(), token.Client, loaded.Client)
	assert.Equal(s.T(), token.IP, loaded.IP)
	assert.Equal(s.T(), token.CreatedAt, loaded.CreatedAt)
	assert.Equal(s.T(), token.LastUsedAt, loaded.LastUsedAt)
}

func (s *TokenStorageTestSuite) TestLoadError() {
	token := NewAuthToken(1, s.storage.Cache)

	s.mini.Close()

	_, err := s.storage.Load(token.Type, token.UserID, token.ID)
	assert.Error(s.T(), err)
	assert.NotEqual(s.T(), ErrNoToken, err)
}

func (s *TokenStorageTestSuite) TestLoadNonExistentToken() {
	token := NewAuthToken(1, s.storage.Cache)

	_, err := s.storage.Load(token.Type, token.UserID, token.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrNoToken, err)
}
