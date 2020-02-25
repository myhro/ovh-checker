package auth

import (
	"fmt"
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

type TokensTestSuite struct {
	suite.Suite

	handler Handler
	mini    *miniredis.Miniredis
	router  *gin.Engine
}

func TestTokensTestSuite(t *testing.T) {
	suite.Run(t, new(TokensTestSuite))
}

func (s *TokensTestSuite) SetupTest() {
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
}

func (s *TokensTestSuite) TearDownTest() {
	s.mini.Close()
}

func (s *TokensTestSuite) TestCacheError() {
	s.router.GET("/", s.handler.Tokens)

	s.mini.Close()

	w := tests.Get(s.router, "/")

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *TokensTestSuite) TestSingleToken() {
	id := 1
	tk := token.NewAuthToken(id, s.handler.Cache)
	err := tk.Save()
	assert.NoError(s.T(), err)

	tests.SetGinContext(s.router, map[string]interface{}{
		"auth_id": id,
	})
	s.router.GET("/", s.handler.Tokens)

	w := tests.Get(s.router, "/")

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Regexp(s.T(), fmt.Sprintf(`"id":"%v"`, tk.ID), w.Body.String())
}
