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

type TokenEndpointTestSuite struct {
	suite.Suite

	handler Handler
	mini    *miniredis.Miniredis
	router  *gin.Engine
}

func TestTokenEndpointTestSuite(t *testing.T) {
	suite.Run(t, new(TokenEndpointTestSuite))
}

func (s *TokenEndpointTestSuite) SetupTest() {
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
	s.router.POST("/", s.handler.Token)
}

func (s *TokenEndpointTestSuite) TearDownTest() {
	s.mini.Close()
}

func (s *TokenEndpointTestSuite) TestCacheError() {
	s.mini.Close()

	w := tests.Post(s.router, "/", "")

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *TokenEndpointTestSuite) TestTokenCreated() {
	w := tests.Post(s.router, "/", "")

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Regexp(s.T(), `"token":".+"`, w.Body.String())
}
