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

	mini    *miniredis.Miniredis
	storage *Storage
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

	s.storage = &Storage{
		Cache: cache,
	}
}

func (s *TokenTestSuite) TearDownTest() {
	s.mini.Close()
}

func (s *TokenTestSuite) TestAuthToken() {
	token := s.storage.NewAuthToken(1)

	assert.Len(s.T(), token.ID, 36)
	assert.Equal(s.T(), 4, strings.Count(token.ID, "-"))
}

func (s *TokenTestSuite) TestCount() {
	token := s.storage.NewAuthToken(1)

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
	token := s.storage.NewAuthToken(1)

	s.mini.Close()

	count, err := token.Count()
	assert.Error(s.T(), err)
	assert.Equal(s.T(), int64(0), count)
}

func (s *TokenTestSuite) TestDelete() {
	token := s.storage.NewAuthToken(1)

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
	token := s.storage.NewAuthToken(1)

	s.mini.Close()

	err := token.Delete()
	assert.Error(s.T(), err)
}

func (s *TokenTestSuite) TestField() {
	token := s.storage.NewAuthToken(1)

	table := []struct {
		in  string
		out string
	}{
		{
			in:  token.dbField("ID"),
			out: "id",
		},
		{
			in:  token.dbField("Client"),
			out: "client",
		},
		{
			in:  token.dbField("IP"),
			out: "ip",
		},
		{
			in:  token.dbField("CreatedAt"),
			out: "created_at",
		},
		{
			in:  token.dbField("LastUsedAt"),
			out: "last_used_at",
		},
		{
			in:  token.dbField("UserID"),
			out: "user_id",
		},
		{
			in:  token.dbField("Type"),
			out: "",
		},
		{
			in:  token.dbField("NonExistentStructField"),
			out: "",
		},
	}

	for _, tt := range table {
		assert.Equal(s.T(), tt.out, tt.in)
	}
}

func (s *TokenTestSuite) TestInvalid() {
	token := s.storage.NewAuthToken(1)

	valid, err := token.Valid()
	assert.NoError(s.T(), err)
	assert.False(s.T(), valid)
}

func (s *TokenTestSuite) TestKeys() {
	authToken := s.storage.NewAuthToken(1)
	sessionToken := s.storage.NewSessionToken(1)

	table := []struct {
		in  string
		out string
	}{
		{
			in:  authToken.Key,
			out: "token-auth:" + authToken.ID,
		},
		{
			in:  authToken.SetKey,
			out: "user:1:tokenset-auth",
		},
		{
			in:  sessionToken.Key,
			out: "token-session:" + sessionToken.ID,
		},
		{
			in:  sessionToken.SetKey,
			out: "user:1:tokenset-session",
		},
	}

	for _, tt := range table {
		assert.Equal(s.T(), tt.out, tt.in)
	}
}

func (s *TokenTestSuite) TestLastUsed() {
	token := s.storage.NewAuthToken(1)

	lastUsed := token.LastUsedAt

	err := token.Save()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), lastUsed, token.LastUsedAt)

	err = token.UpdateLastUsed()
	assert.NoError(s.T(), err)
	assert.NotEqual(s.T(), lastUsed, token.LastUsedAt)
}

func (s *TokenTestSuite) TestLastUsedError() {
	token := s.storage.NewAuthToken(1)

	err := token.Save()
	assert.NoError(s.T(), err)

	s.mini.Close()

	lastUsed := token.LastUsedAt

	err = token.UpdateLastUsed()
	assert.Error(s.T(), err)
	assert.Equal(s.T(), lastUsed, token.LastUsedAt)
}

func (s *TokenTestSuite) TestSave() {
	token := s.storage.NewAuthToken(1)

	err := token.Save()
	assert.NoError(s.T(), err)
}

func (s *TokenTestSuite) TestSaveError() {
	token := s.storage.NewAuthToken(1)

	s.mini.Close()

	err := token.Save()
	assert.Error(s.T(), err)
}

func (s *TokenTestSuite) TestSessionToken() {
	token := s.storage.NewSessionToken(1)

	assert.Len(s.T(), token.ID, 32)
	assert.Equal(s.T(), 0, strings.Count(token.ID, "-"))
}

func (s *TokenTestSuite) TestSet() {
	token := s.storage.NewAuthToken(1)

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
	token := s.storage.NewAuthToken(1)

	s.mini.Close()

	set, err := token.Set()
	assert.Error(s.T(), err)
	assert.Len(s.T(), set, 0)
}

func (s *TokenTestSuite) TestValid() {
	token := s.storage.NewAuthToken(1)

	err := token.Save()
	assert.NoError(s.T(), err)

	valid, err := token.Valid()
	assert.NoError(s.T(), err)
	assert.True(s.T(), valid)
}

func (s *TokenTestSuite) TestValidError() {
	token := s.storage.NewAuthToken(1)

	s.mini.Close()

	valid, err := token.Valid()
	assert.Error(s.T(), err)
	assert.False(s.T(), valid)
}
