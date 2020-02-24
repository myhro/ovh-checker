package auth

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/myhro/ovh-checker/api/tests"
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
	s.handler.Cache = &storage.Redis{
		Client: redis.NewClient(opts),
	}

	gin.SetMode(gin.ReleaseMode)
	s.router = gin.New()
	s.router.POST("/", s.handler.Logout)
}

func (s *LogoutTestSuite) TearDownTest() {
	s.mini.Close()
}

func (s *LogoutTestSuite) SetGinContext(ctxs map[string]interface{}) {
	rec := httptest.NewRecorder()
	_, s.router = gin.CreateTestContext(rec)
	s.router.Use(func(c *gin.Context) {
		for k, v := range ctxs {
			c.Set(k, v)
		}
	})
	s.router.POST("/", s.handler.Logout)
}

func (s *LogoutTestSuite) TestCacheError() {
	s.mini.Close()

	w := tests.Post(s.router, "/", "")

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *LogoutTestSuite) TestMultipleTokens() {
	id := 0
	token1 := "xyz"
	token2 := "abc"
	s.handler.addToken(id, token1, "logout-test", "127.0.0.1")
	s.handler.addToken(id, token2, "logout-test", "127.0.0.1")

	list, err := s.handler.getTokens(id)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), list, 2)

	s.SetGinContext(map[string]interface{}{
		"auth_id": id,
		"token":   token1,
	})
	w := tests.Post(s.router, "/", "")

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Regexp(s.T(), logoutMessage, w.Body.String())

	list, err = s.handler.getTokens(id)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), list, 1)
}

func (s *LogoutTestSuite) TestSingleToken() {
	id := 0
	token := "xyz"
	s.handler.addToken(id, token, "logout-test", "127.0.0.1")

	list, err := s.handler.getTokens(id)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), list, 1)

	s.SetGinContext(map[string]interface{}{
		"auth_id": id,
		"token":   token,
	})
	w := tests.Post(s.router, "/", "")

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Regexp(s.T(), logoutMessage, w.Body.String())

	list, err = s.handler.getTokens(id)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), list, 0)
}
