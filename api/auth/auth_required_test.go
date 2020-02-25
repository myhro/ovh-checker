package auth

import (
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

var dummyAuthHeaders = map[string]string{
	"Authorization": "Token xyz",
}

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
	s.router.GET("/", s.handler.AuthRequired)
}

func (s *AuthRequiredTestSuite) TearDownTest() {
	s.mini.Close()
}

func (s *AuthRequiredTestSuite) TestCacheError() {
	s.mini.Close()

	w := tests.GetWithHeaders(s.router, "/", dummyAuthHeaders)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *AuthRequiredTestSuite) TestInvalidSession() {
	s.router.GET("/set", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("auth_id", 1)
		session.Set("session_id", "xyz")
		session.Save()
	})

	w1 := tests.Get(s.router, "/set")
	cookies := w1.HeaderMap.Get("Set-Cookie")

	w2 := tests.GetWithHeaders(s.router, "/", map[string]string{"Cookie": cookies})

	assert.Equal(s.T(), http.StatusUnauthorized, w2.Code)
	assert.Regexp(s.T(), invalidSessionError, w2.Body.String())
}

func (s *AuthRequiredTestSuite) TestInvalidToken() {
	headers := map[string]string{
		"Authorization": "Token xyz",
	}
	w := tests.GetWithHeaders(s.router, "/", headers)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Regexp(s.T(), invalidAuthError, w.Body.String())
}

func (s *AuthRequiredTestSuite) TestMalformedToken() {
	headers := map[string]string{
		"Authorization": "Tokenxyz",
	}
	w := tests.GetWithHeaders(s.router, "/", headers)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Regexp(s.T(), unsupportedAuthError, w.Body.String())
}

func (s *AuthRequiredTestSuite) TestMissingHeader() {
	w := tests.Get(s.router, "/")

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Regexp(s.T(), unsupportedAuthError, w.Body.String())
}

func (s *AuthRequiredTestSuite) TestSessionCacheError() {
	s.router.GET("/set", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("auth_id", 1)
		session.Set("session_id", "xyz")
		session.Save()
	})

	w1 := tests.Get(s.router, "/set")
	cookies := w1.HeaderMap.Get("Set-Cookie")

	s.mini.Close()

	w2 := tests.GetWithHeaders(s.router, "/", map[string]string{"Cookie": cookies})

	assert.Equal(s.T(), http.StatusInternalServerError, w2.Code)
	assert.Equal(s.T(), "Internal Server Error", w2.Body.String())
}

func (s *AuthRequiredTestSuite) TestUpdateLastUsedCacheError() {
	tk := s.handler.TokenStorage.NewAuthToken(1)

	cache := &tests.MockedCache{}
	cache.On("HGetAll", mock.Anything).Return(map[string]string{"id": tk.ID, "user_id": "1"}, nil)
	cache.On("HSet", mock.Anything, mock.Anything, mock.Anything).Return(false, errors.New("cache error"))
	s.handler.TokenStorage.Cache = cache

	headers := map[string]string{
		"Authorization": "Token " + tk.ID,
	}
	w := tests.GetWithHeaders(s.router, "/", headers)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *AuthRequiredTestSuite) TestValidSession() {
	id := 1
	tk := s.handler.TokenStorage.NewSessionToken(id)
	err := tk.Save()
	assert.NoError(s.T(), err)

	s.router.GET("/set", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("auth_id", id)
		session.Set("session_id", tk.ID)
		session.Save()
	})

	w1 := tests.Get(s.router, "/set")
	cookies := w1.HeaderMap.Get("Set-Cookie")

	w2 := tests.GetWithHeaders(s.router, "/", map[string]string{"Cookie": cookies})

	assert.Equal(s.T(), http.StatusOK, w2.Code)
	assert.Equal(s.T(), "", w2.Body.String())
}

func (s *AuthRequiredTestSuite) TestWrongHeader() {
	headers := map[string]string{
		"Authorization": "Basic xyz",
	}
	w := tests.GetWithHeaders(s.router, "/", headers)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Regexp(s.T(), unsupportedAuthError, w.Body.String())
}
