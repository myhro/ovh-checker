package auth

import (
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
	"github.com/stretchr/testify/suite"
)

type LogoutTestSuite struct {
	suite.Suite

	handler Handler
	mini    *miniredis.Miniredis
	router  *gin.Engine
}

func TestLogoutTestSuite(t *testing.T) {
	suite.Run(t, new(LogoutTestSuite))
}

func (s *LogoutTestSuite) SetupTest() {
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

func (s *LogoutTestSuite) TearDownTest() {
	s.mini.Close()
}

func (s *LogoutTestSuite) TestCacheError() {
	tk := s.handler.TokenStorage.NewAuthToken(1)
	err := tk.Save()
	assert.NoError(s.T(), err)

	s.mini.Close()

	tests.SetGinContext(s.router, map[string]interface{}{
		"token": tk,
	})
	s.router.POST("/", s.handler.Logout)

	w := tests.Post(s.router, "/", "")

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *LogoutTestSuite) TestSingleToken() {
	id := 1
	tk := s.handler.TokenStorage.NewAuthToken(id)
	err := tk.Save()
	assert.NoError(s.T(), err)

	tests.SetGinContext(s.router, map[string]interface{}{
		"token": tk,
	})
	s.router.POST("/", s.handler.Logout)

	res, err := s.handler.TokenStorage.ListAll(id)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), res["auth"], 1)

	w := tests.Post(s.router, "/", "")

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Regexp(s.T(), logoutMessage, w.Body.String())

	res, err = s.handler.TokenStorage.ListAll(id)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), res["auth"], 0)
}
