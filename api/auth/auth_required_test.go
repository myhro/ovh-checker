package auth

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/myhro/ovh-checker/api/tests"
	"github.com/myhro/ovh-checker/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var validTokenHeader = tests.AuthHeader("user@example.com", "xyz")

type AuthRequiredTestSuite struct {
	suite.Suite

	handler Handler
	mini    *miniredis.Miniredis
	router  *gin.Engine
}

func TestAuthRequiredTestSuite(t *testing.T) {
	suite.Run(t, new(AuthRequiredTestSuite))
}

func (s *AuthRequiredTestSuite) SetupTest() {
	log.SetOutput(ioutil.Discard)

	s.handler = Handler{}

	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	s.mini = mr
	opts := &redis.Options{
		Addr: s.mini.Addr(),
	}
	s.handler.Cache = &storage.Redis{
		Client: redis.NewClient(opts),
	}

	store := cookie.NewStore([]byte("login-test"))

	gin.SetMode(gin.ReleaseMode)
	s.router = gin.New()
	s.router.Use(sessions.Sessions("session", store))
	s.router.GET("/", s.handler.AuthRequired)
}

func (s *AuthRequiredTestSuite) TearDownTest() {
	s.mini.Close()
}

func (s *AuthRequiredTestSuite) TestCacheError() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.handler.DB = db

	s.mini.Close()

	headers := map[string]string{
		"Authorization": validTokenHeader,
	}
	w := tests.GetWithHeaders(s.router, "/", headers)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *AuthRequiredTestSuite) TestCacheErrorTokenLastUsed() {
	cache := &tests.MockedCache{}
	cache.On("SIsMember", mock.Anything, mock.Anything).Return(true, nil)
	cache.On("HSet", mock.Anything, mock.Anything, mock.Anything).Return(false, errors.New("cache error"))
	s.handler.Cache = cache

	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.handler.DB = db

	headers := map[string]string{
		"Authorization": validTokenHeader,
	}
	w := tests.GetWithHeaders(s.router, "/", headers)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *AuthRequiredTestSuite) TestDatabaseError() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("database error"))
	s.handler.DB = db

	headers := map[string]string{
		"Authorization": validTokenHeader,
	}
	w := tests.GetWithHeaders(s.router, "/", headers)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *AuthRequiredTestSuite) TestExistingUserTokenNotInRedis() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.handler.DB = db

	headers := map[string]string{
		"Authorization": validTokenHeader,
	}
	w := tests.GetWithHeaders(s.router, "/", headers)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Regexp(s.T(), incorrectEmailTokenError, w.Body.String())
}

func (s *AuthRequiredTestSuite) TestExistingUserTokenOk() {
	token := "xyz"
	s.handler.addToken(authStoragePrefix, 0, token, "auth-required-test", "127.0.0.1")

	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.handler.DB = db

	headers := map[string]string{
		"Authorization": validTokenHeader,
	}
	w := tests.GetWithHeaders(s.router, "/", headers)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), "", w.Body.String())
}

func (s *AuthRequiredTestSuite) TestInvalidToken() {
	headers := map[string]string{
		"Authorization": "Token xyz",
	}
	w := tests.GetWithHeaders(s.router, "/", headers)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Regexp(s.T(), invalidPairError, w.Body.String())
}

func (s *AuthRequiredTestSuite) TestMalformedToken() {
	token := base64.StdEncoding.EncodeToString([]byte("xyz"))
	headers := map[string]string{
		"Authorization": "Token " + token,
	}
	w := tests.GetWithHeaders(s.router, "/", headers)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Regexp(s.T(), invalidPairError, w.Body.String())
}

func (s *AuthRequiredTestSuite) TestMissingHeader() {
	w := tests.Get(s.router, "/")

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Regexp(s.T(), unsupportedAuthError, w.Body.String())
}

func (s *AuthRequiredTestSuite) TestNonExistentUser() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows)
	s.handler.DB = db

	headers := map[string]string{
		"Authorization": validTokenHeader,
	}
	w := tests.GetWithHeaders(s.router, "/", headers)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Regexp(s.T(), incorrectEmailTokenError, w.Body.String())
}

func (s *AuthRequiredTestSuite) TestWrongHeader() {
	headers := map[string]string{
		"Authorization": "Basic xyz",
	}
	w := tests.GetWithHeaders(s.router, "/", headers)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Regexp(s.T(), unsupportedAuthError, w.Body.String())
}
