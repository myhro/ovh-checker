package auth

import (
	"database/sql"
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

const bodyEqualPassword = `
		{
			"email": "user@example.com",
			"password1": "123",
			"password2": "123"
		}
	`

var errDB = errors.New("database error")

type SignupTestSuite struct {
	suite.Suite

	handler Handler
	router  *gin.Engine
}

func TestSignupTestSuite(t *testing.T) {
	suite.Run(t, new(SignupTestSuite))
}

func (s *SignupTestSuite) SetupSuite() {
	log.SetOutput(ioutil.Discard)

	s.handler = Handler{}

	gin.SetMode(gin.ReleaseMode)
	s.router = gin.New()
	s.router.POST("/auth/signup", s.handler.Signup)
}

func (s *SignupTestSuite) TestAddUser() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows)
	db.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	s.handler.DB = db

	w := tests.Post(s.router, "/auth/signup", bodyEqualPassword)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Regexp(s.T(), userCreated, w.Body.String())
}

func (s *SignupTestSuite) TestDatabaseErrorAddUser() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows)
	db.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(nil, errDB)
	s.handler.DB = db

	w := tests.Post(s.router, "/auth/signup", bodyEqualPassword)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *SignupTestSuite) TestDatabaseErrorUserExists() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(errDB)
	s.handler.DB = db

	w := tests.Post(s.router, "/auth/signup", bodyEqualPassword)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(s.T(), "Internal Server Error", w.Body.String())
}

func (s *SignupTestSuite) TestMissingAllParameters() {
	w := tests.Post(s.router, "/auth/signup", "{}")

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Regexp(s.T(), "field validation .*'email'.* failed", w.Body.String())
}

func (s *SignupTestSuite) TestMissingSomeParameters() {
	body := `
		{
			"email": "user@example.com"
		}
	`
	w := tests.Post(s.router, "/auth/signup", body)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Regexp(s.T(), "field validation .*'password1'.* failed", w.Body.String())
}

func (s *SignupTestSuite) TestPasswordsDoesntMatch() {
	body := `
		{
			"email": "user@example.com",
			"password1": "123",
			"password2": "456"
		}
	`
	w := tests.Post(s.router, "/auth/signup", body)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Regexp(s.T(), passwordError, w.Body.String())
}

func (s *SignupTestSuite) TestUserExists() {
	db := &tests.MockedDatabase{}
	db.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.handler.DB = db

	w := tests.Post(s.router, "/auth/signup", bodyEqualPassword)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	assert.Regexp(s.T(), emailError, w.Body.String())
}
