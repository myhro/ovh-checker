package notification

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type NotificationEndpointTestSuite struct {
	suite.Suite

	handler Handler
	router  *gin.Engine
}

func TestNotificationEndpointTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationEndpointTestSuite))
}

func (s *NotificationEndpointTestSuite) SetupSuite() {
	log.SetOutput(ioutil.Discard)

	s.handler = Handler{}

	gin.SetMode(gin.ReleaseMode)
	s.router = gin.New()
	s.router.POST("/", s.handler.Notification)
}

func (s *NotificationEndpointTestSuite) TestDatabaseError() {
	db := &tests.MockedDatabase{}
	db.On("Exec", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
	s.handler.DB = db

	w := tests.Post(s.router, "/", `{"country":"fr","server":"KS-1"}`)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *NotificationEndpointTestSuite) TestMissingAllParameters() {
	w := tests.Post(s.router, "/", "{}")

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Regexp(s.T(), "field validation .*'country'.* failed", w.Body.String())
}

func (s *NotificationEndpointTestSuite) TestMissingSomeParameters() {
	w := tests.Post(s.router, "/", `{"country":"fr"}`)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Regexp(s.T(), "field validation .*'server'.* failed", w.Body.String())
}

func (s *NotificationEndpointTestSuite) TestNotificationCreated() {
	db := &tests.MockedDatabase{}
	db.On("Exec", mock.Anything, mock.Anything).Return(nil, nil)
	s.handler.DB = db

	w := tests.Post(s.router, "/", `{"country":"fr","server":"KS-1"}`)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Regexp(s.T(), successfullyCreated, w.Body.String())
}
