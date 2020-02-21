package hardware

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HardwareTestSuite struct {
	suite.Suite

	handler Handler
	router  *gin.Engine
}

func TestHardwareTestSuite(t *testing.T) {
	suite.Run(t, new(HardwareTestSuite))
}

func (s *HardwareTestSuite) SetupSuite() {
	log.SetOutput(ioutil.Discard)

	s.handler = Handler{}

	gin.SetMode(gin.ReleaseMode)
	s.router = gin.New()
	s.router.GET("/hardware/offers", s.handler.Offers)
}

func (s *HardwareTestSuite) TestAllParameters() {
	db := &tests.MockedDatabase{}
	db.On("Select", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.handler.DB = db

	w := tests.Get(s.router, "/hardware/offers?country=fr&first=1&last=1")

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), "[]", strings.TrimSpace(w.Body.String()))
}

func (s *HardwareTestSuite) TestDatabaseError() {
	db := &tests.MockedDatabase{}
	db.On("Select", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("dummy"))
	s.handler.DB = db

	w := tests.Get(s.router, "/hardware/offers?country=fr&first=1&last=1")

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *HardwareTestSuite) TestMissingAllParameters() {
	w := tests.Get(s.router, "/hardware/offers")

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Regexp(s.T(), "field validation .* failed", w.Body.String())
}

func (s *HardwareTestSuite) TestMissingSomeParameters() {
	w := tests.Get(s.router, "/hardware/offers?country=fr")

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Regexp(s.T(), "field validation .* failed", w.Body.String())
}

func (s *HardwareTestSuite) TestNonIntParameter() {
	w := tests.Get(s.router, "/hardware/offers?country=fr&first=1&last=xyz")

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Regexp(s.T(), "parsing .* invalid syntax", w.Body.String())
}
