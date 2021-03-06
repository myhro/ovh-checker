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
	"github.com/myhro/ovh-checker/api/token"
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
	cache := &storage.Redis{
		Client: redis.NewClient(opts),
	}

	s.handler.TokenStorage = &token.Storage{
		Cache: cache,
	}

	gin.SetMode(gin.ReleaseMode)
	s.router = gin.New()
}

func (s *UserTestSuite) TearDownTest() {
	s.mini.Close()
}

func (s *UserTestSuite) TestCacheError() {
	tk := s.handler.TokenStorage.NewAuthToken(1)
	err := tk.Save()
	assert.NoError(s.T(), err)

	tests.SetGinContext(s.router, map[string]interface{}{
		"token": tk,
	})
	s.router.GET("/", s.handler.User)

	s.mini.Close()

	w := tests.Get(s.router, "/")

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *UserTestSuite) TestDatabaseError() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("database error"))
	s.handler.DB = db

	tk := s.handler.TokenStorage.NewAuthToken(1)
	err := tk.Save()
	assert.NoError(s.T(), err)

	tests.SetGinContext(s.router, map[string]interface{}{
		"token": tk,
	})
	s.router.GET("/", s.handler.User)

	w := tests.Get(s.router, "/")

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *UserTestSuite) TestSingleToken() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.handler.DB = db

	tk := s.handler.TokenStorage.NewAuthToken(1)
	err := tk.Save()
	assert.NoError(s.T(), err)

	tests.SetGinContext(s.router, map[string]interface{}{
		"token": tk,
	})
	s.router.GET("/", s.handler.User)

	w := tests.Get(s.router, "/")

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Regexp(s.T(), `"tokens":1`, w.Body.String())
}
