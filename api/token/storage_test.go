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

	authToken := s.storage.NewAuthToken(id)
	sessionToken := s.storage.NewSessionToken(id)

	hash, err := s.storage.ListAll(id)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), hash[AuthPrefix], 0)
	assert.Len(s.T(), hash[SessionPrefix], 0)

	err = authToken.Save()
	assert.NoError(s.T(), err)
	err = sessionToken.Save()
	assert.NoError(s.T(), err)

	hash, err = s.storage.ListAll(id)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), hash[AuthPrefix], 1)
	assert.Len(s.T(), hash[SessionPrefix], 1)
}

func (s *TokenStorageTestSuite) TestListAllError() {
	s.mini.Close()

	hash, err := s.storage.ListAll(1)
	assert.Error(s.T(), err)
	assert.Len(s.T(), hash[AuthPrefix], 0)
	assert.Len(s.T(), hash[SessionPrefix], 0)
}

func (s *TokenStorageTestSuite) TestLoad() {
	token := s.storage.NewAuthToken(1)
	token.Client = "token-storage-test"
	token.IP = "127.0.0.1"

	err := token.Save()
	assert.NoError(s.T(), err)

	loaded, err := s.storage.LoadAuthToken(token.ID)
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
	token := s.storage.NewSessionToken(1)

	s.mini.Close()

	_, err := s.storage.LoadSessionToken(token.ID)
	assert.Error(s.T(), err)
	assert.NotEqual(s.T(), ErrNoToken, err)
}

func (s *TokenStorageTestSuite) TestLoadNonExistentToken() {
	_, err := s.storage.LoadAuthToken("xyz")
	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrNoToken, err)
}
