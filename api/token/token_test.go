package token

import (
	"log"
	"strings"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/myhro/ovh-checker/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TokenTestSuite struct {
	suite.Suite

	storage *TokenStorage
	mini    *miniredis.Miniredis
}

func TestTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}

func (s *TokenTestSuite) SetupTest() {
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

	s.storage = &TokenStorage{
		Cache: cache,
	}
}

func (s *TokenTestSuite) TearDownTest() {
	s.mini.Close()
}

func (s *TokenTestSuite) TestAuthToken() {
	token := NewAuthToken(1)

	assert.Len(s.T(), token.ID, 36)
	assert.Equal(s.T(), 4, strings.Count(token.ID, "-"))
}

func (s *TokenTestSuite) TestCount() {
	token := NewAuthToken(1)
	token.Storage = s.storage

	count, err := token.Count()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(0), count)

	err = token.Save()
	assert.NoError(s.T(), err)

	count, err = token.Count()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), count)
}

func (s *TokenTestSuite) TestCountError() {
	token := NewAuthToken(1)
	token.Storage = s.storage

	s.mini.Close()

	count, err := token.Count()
	assert.Error(s.T(), err)
	assert.Equal(s.T(), int64(0), count)
}

func (s *TokenTestSuite) TestDelete() {
	token := NewAuthToken(1)
	token.Storage = s.storage

	err := token.Save()
	assert.NoError(s.T(), err)

	count, err := token.Count()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), count)

	err = token.Delete()
	assert.NoError(s.T(), err)

	count, err = token.Count()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(0), count)
}

func (s *TokenTestSuite) TestDeleteError() {
	token := NewAuthToken(1)
	token.Storage = s.storage

	s.mini.Close()

	err := token.Delete()
	assert.Error(s.T(), err)
}

func (s *TokenTestSuite) TestField() {
	token := NewAuthToken(1)

	table := []struct {
		in  string
		out string
	}{
		{
			in:  token.field("ID"),
			out: "id",
		},
		{
			in:  token.field("Client"),
			out: "client",
		},
		{
			in:  token.field("IP"),
			out: "ip",
		},
		{
			in:  token.field("CreatedAt"),
			out: "created_at",
		},
		{
			in:  token.field("LastUsedAt"),
			out: "last_used_at",
		},
		{
			in:  token.field("UserID"),
			out: "",
		},
		{
			in:  token.field("NonExistentStructField"),
			out: "",
		},
	}

	for _, tt := range table {
		assert.Equal(s.T(), tt.out, tt.in)
	}
}

func (s *TokenTestSuite) TestInvalid() {
	token := NewAuthToken(1)
	token.Storage = s.storage

	valid, err := token.Valid()
	assert.NoError(s.T(), err)
	assert.False(s.T(), valid)
}

func (s *TokenTestSuite) TestKeys() {
	authToken := NewAuthToken(1)
	sessionToken := NewSessionToken(1)

	table := []struct {
		in  string
		out string
	}{
		{
			in:  authToken.Key,
			out: "user:1:auth:" + authToken.ID,
		},
		{
			in:  authToken.SetKey,
			out: "user:1:auth-set",
		},
		{
			in:  sessionToken.Key,
			out: "user:1:session:" + sessionToken.ID,
		},
		{
			in:  sessionToken.SetKey,
			out: "user:1:session-set",
		},
	}

	for _, tt := range table {
		assert.Equal(s.T(), tt.out, tt.in)
	}
}

func (s *TokenTestSuite) TestLastUsed() {
	token := NewAuthToken(1)
	token.Storage = s.storage

	lastUsed := token.LastUsedAt

	err := token.Save()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), lastUsed, token.LastUsedAt)

	err = token.UpdateLastUsed()
	assert.NoError(s.T(), err)
	assert.NotEqual(s.T(), lastUsed, token.LastUsedAt)
}

func (s *TokenTestSuite) TestLastUsedError() {
	token := NewAuthToken(1)
	token.Storage = s.storage

	err := token.Save()
	assert.NoError(s.T(), err)

	s.mini.Close()

	lastUsed := token.LastUsedAt

	err = token.UpdateLastUsed()
	assert.Error(s.T(), err)
	assert.Equal(s.T(), lastUsed, token.LastUsedAt)
}

func (s *TokenTestSuite) TestSave() {
	token := NewAuthToken(1)
	token.Storage = s.storage

	err := token.Save()
	assert.NoError(s.T(), err)
}

func (s *TokenTestSuite) TestSaveError() {
	token := NewAuthToken(1)
	token.Storage = s.storage

	s.mini.Close()

	err := token.Save()
	assert.Error(s.T(), err)
}

func (s *TokenTestSuite) TestSessionToken() {
	token := NewSessionToken(1)

	assert.Len(s.T(), token.ID, 32)
	assert.Equal(s.T(), 0, strings.Count(token.ID, "-"))
}

func (s *TokenTestSuite) TestSet() {
	token := NewAuthToken(1)
	token.Storage = s.storage

	set, err := token.Set()
	assert.NoError(s.T(), err)
	assert.Len(s.T(), set, 0)

	err = token.Save()
	assert.NoError(s.T(), err)

	set, err = token.Set()
	assert.NoError(s.T(), err)
	assert.Len(s.T(), set, 1)
}

func (s *TokenTestSuite) TestSetError() {
	token := NewAuthToken(1)
	token.Storage = s.storage

	s.mini.Close()

	set, err := token.Set()
	assert.Error(s.T(), err)
	assert.Len(s.T(), set, 0)
}

func (s *TokenTestSuite) TestValid() {
	token := NewAuthToken(1)
	token.Storage = s.storage

	err := token.Save()
	assert.NoError(s.T(), err)

	valid, err := token.Valid()
	assert.NoError(s.T(), err)
	assert.True(s.T(), valid)
}

func (s *TokenTestSuite) TestValidError() {
	token := NewAuthToken(1)
	token.Storage = s.storage

	s.mini.Close()

	valid, err := token.Valid()
	assert.Error(s.T(), err)
	assert.False(s.T(), valid)
}
