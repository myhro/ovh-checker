package sqlsuite

import (
	"database/sql"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/icrowley/fake"
	"github.com/jmoiron/sqlx"
	"github.com/nleof/goyesql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite

	db      *sqlx.DB
	mig     *migrate.Migrate
	queries goyesql.Queries
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) SetupSuite() {
	s.db = newDB()
	s.mig = newMigrate()
	s.queries = newQueries("auth")

	s.mig.Up()
}

func (s *AuthTestSuite) TearDownSuite() {
	s.mig.Down()
}

func (s *AuthTestSuite) TestAddUser() {
	_, err := s.db.Exec(s.queries["add-user"], fake.EmailAddress(), fake.SimplePassword())
	assert.NoError(s.T(), err)
}

func (s *AuthTestSuite) TestCheckPassword() {
	email := fake.EmailAddress()
	password := fake.SimplePassword()

	_, err := s.db.Exec(s.queries["add-user"], email, password)
	assert.NoError(s.T(), err)

	var id int
	err = s.db.Get(&id, s.queries["check-password"], email, password)
	assert.NoError(s.T(), err)
	assert.NotEqual(s.T(), 0, id)
}

func (s *AuthTestSuite) TestCheckWrongPassword() {
	email := fake.EmailAddress()

	_, err := s.db.Exec(s.queries["add-user"], email, fake.SimplePassword())
	assert.NoError(s.T(), err)

	err = s.db.Get(new(int), s.queries["check-password"], email, "p4ssw0rd")
	assert.Equal(s.T(), sql.ErrNoRows, err)
}

func (s *AuthTestSuite) TestUserDoesntExists() {
	err := s.db.Get(new(int), s.queries["user-exists"], fake.EmailAddress())
	assert.Equal(s.T(), sql.ErrNoRows, err)
}

func (s *AuthTestSuite) TestUserEmail() {
	email := fake.EmailAddress()

	_, err := s.db.Exec(s.queries["add-user"], email, fake.SimplePassword())
	assert.NoError(s.T(), err)

	var id int
	err = s.db.Get(&id, s.queries["user-exists"], email)
	assert.NoError(s.T(), err)

	var dbEmail string
	err = s.db.Get(&dbEmail, s.queries["user-email"], id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), email, dbEmail)
}

func (s *AuthTestSuite) TestUserEmailNotFound() {
	var dbEmail string
	err := s.db.Get(&dbEmail, s.queries["user-email"], 0)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), sql.ErrNoRows, err)
}

func (s *AuthTestSuite) TestUserExists() {
	email := fake.EmailAddress()

	_, err := s.db.Exec(s.queries["add-user"], email, fake.SimplePassword())
	assert.NoError(s.T(), err)

	var id int
	err = s.db.Get(&id, s.queries["user-exists"], email)
	assert.NoError(s.T(), err)
	assert.NotEqual(s.T(), 0, id)
}
