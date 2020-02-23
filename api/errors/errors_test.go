package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ErrorsTestSuite struct {
	suite.Suite

	router *gin.Engine
}

func TestErrorsTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorsTestSuite))
}

func (s *ErrorsTestSuite) SetupSuite() {
	gin.SetMode(gin.ReleaseMode)
}

func (s *ErrorsTestSuite) TestBadRequestWithMessage() {
	msg := "bad request"
	handler := func(c *gin.Context) {
		BadRequestWithMessage(c, msg)
	}

	s.router = gin.New()
	s.router.GET("/", handler)
	w := tests.Get(s.router, "/")

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Regexp(s.T(), "application/json", w.HeaderMap["Content-Type"][0])
	assert.Regexp(s.T(), "error.*"+msg, w.Body.String())
}

func (s *ErrorsTestSuite) TestInternalServerError() {
	handler := func(c *gin.Context) {
		InternalServerError(c)
	}

	s.router = gin.New()
	s.router.GET("/", handler)
	w := tests.Get(s.router, "/")

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Regexp(s.T(), "text/plain", w.HeaderMap["Content-Type"][0])
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *ErrorsTestSuite) TestNew() {
	msg := "an error occurred"
	assert.Equal(s.T(), New(msg), errors.New(msg))
}

func (s *ErrorsTestSuite) TestUnauthorizedWithMessage() {
	msg := "unauthorized"
	handler := func(c *gin.Context) {
		UnauthorizedWithMessage(c, msg)
	}

	s.router = gin.New()
	s.router.GET("/", handler)
	w := tests.Get(s.router, "/")

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Regexp(s.T(), "application/json", w.HeaderMap["Content-Type"][0])
	assert.Regexp(s.T(), "error.*"+msg, w.Body.String())
}

func (s *ErrorsTestSuite) TestValidationMessage() {
	table := []struct {
		in  error
		out string
	}{
		{
			in:  errors.New("strconv.ParseInt: parsing \"xyz\": invalid syntax"),
			out: "parsing \"xyz\": invalid syntax",
		},
		{
			in:  errors.New("Key: 'OffersRequest.Last' Error:Field validation for 'Last' failed on the 'required' tag"),
			out: "field validation for 'last' failed on the 'required' tag",
		},
		{
			in:  errors.New("Key: 'OffersRequest.Country' Error:Field validation for 'Country' failed on the 'required' tag\nKey: 'OffersRequest.Last' Error:Field validation for 'Last' failed on the 'required' tag"),
			out: "field validation for 'country' failed on the 'required' tag",
		},
		{
			in:  errors.New("invalid character '}' looking for beginning of object key string"),
			out: "invalid json",
		},
		{
			in:  errors.New("EOF"),
			out: "invalid json",
		},
		{
			in:  errors.New("not mapped failure"),
			out: "unknown error",
		},
	}

	for _, tt := range table {
		assert.Equal(s.T(), tt.out, ValidationMessage(tt.in))
	}
}
