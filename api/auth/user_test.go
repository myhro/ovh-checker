package auth

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/myhro/ovh-checker/api/tests"
	"github.com/myhro/ovh-checker/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite

	handler Handler
	mini    *miniredis.Miniredis
	router  *gin.Engine
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) SetupTest() {
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

	gin.SetMode(gin.ReleaseMode)
	s.router = gin.New()
	s.router.GET("/", s.handler.User)
}

func (s *UserTestSuite) TearDownTest() {
	s.mini.Close()
}

func (s *UserTestSuite) TestCacheError() {
	cache := &tests.MockedCache{}
	cache.On("SCard", mock.Anything).Return(0, errors.New("cache error"))
	s.handler.Cache = cache

	w := tests.Get(s.router, "/")

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *UserTestSuite) TestMultipleTokens() {
	token1 := "xyz"
	token2 := "abc"
	s.handler.addToken(0, token1, "user-test", "127.0.0.1")
	s.handler.addToken(0, token2, "user-test", "127.0.0.1")

	w := tests.Get(s.router, "/")

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Regexp(s.T(), "email", w.Body.String())
	assert.Regexp(s.T(), "tokens.*2", w.Body.String())
}

func (s *UserTestSuite) TestSingleToken() {
	token := "xyz"
	s.handler.addToken(0, token, "user-test", "127.0.0.1")

	w := tests.Get(s.router, "/")

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Regexp(s.T(), "email", w.Body.String())
	assert.Regexp(s.T(), "tokens.*1", w.Body.String())
}
