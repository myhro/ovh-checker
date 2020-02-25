package auth

import (
	"database/sql"
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
	"github.com/myhro/ovh-checker/api/token"
	"github.com/myhro/ovh-checker/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const validCreds = `{"email": "user@example.com", "password": "123"}`

type LoginTestSuite struct {
	suite.Suite

	handler Handler
	mini    *miniredis.Miniredis
	router  *gin.Engine
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}

func (s *LoginTestSuite) SetupTest() {
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
	cache := &storage.Redis{
		Client: redis.NewClient(opts),
	}

	s.handler.TokenStorage = &token.Storage{
		Cache: cache,
	}

	store := cookie.NewStore([]byte("login-test"))

	gin.SetMode(gin.ReleaseMode)
	s.router = gin.New()
	s.router.Use(sessions.Sessions("session", store))
	s.router.POST("/", s.handler.Login)
}

func (s *LoginTestSuite) TearDownTest() {
	s.mini.Close()
}

func (s *LoginTestSuite) TestCacheError() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.handler.DB = db

	s.mini.Close()

	w := tests.Post(s.router, "/", validCreds)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *LoginTestSuite) TestDatabaseError() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("database error"))
	s.handler.DB = db

	w := tests.Post(s.router, "/", validCreds)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *LoginTestSuite) TestMissingAllParameters() {
	w := tests.Post(s.router, "/", "{}")

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Regexp(s.T(), "field validation .*'email'.* failed", w.Body.String())
}

func (s *LoginTestSuite) TestMissingPassword() {
	w := tests.Post(s.router, "/", `{"email": "user@example.com"}`)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Regexp(s.T(), "field validation .*'password'.* failed", w.Body.String())
}

func (s *LoginTestSuite) TestNonExistentUser() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows)
	s.handler.DB = db

	w := tests.Post(s.router, "/", validCreds)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Regexp(s.T(), incorrectEmailPasswordError, w.Body.String())
}

func (s *LoginTestSuite) TestValidUser() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.handler.DB = db

	w := tests.Post(s.router, "/", validCreds)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Regexp(s.T(), successfulLogin, w.Body.String())
}
