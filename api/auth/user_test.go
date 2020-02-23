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
	s.handler.Cache = redis.NewClient(opts)

	gin.SetMode(gin.ReleaseMode)
	s.router = gin.New()
	s.router.GET("/", s.handler.AuthRequired, s.handler.User)
}

func (s *UserTestSuite) TearDownTest() {
	s.mini.Close()
}

func (s *UserTestSuite) TestProperRequest() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.handler.DB = db

	email := "user@example.com"
	token := "xyz"
	headers := map[string]string{
		"Authorization": tests.AuthHeader(email, token),
	}

	s.handler.addToken(0, token, "user-test", "127.0.0.1")

	w := tests.GetWithHeaders(s.router, "/", headers)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Regexp(s.T(), "email.*"+email, w.Body.String())
	assert.Regexp(s.T(), "token.*"+token, w.Body.String())
}
