package sqlsuite

import (
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
	s.db = NewDB()
	s.mig = NewMigrate()
	s.queries = NewQueries("auth")

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

	var id int
	err = s.db.Get(&id, s.queries["check-password"], email, "p4ssw0rd")
	assert.Error(s.T(), err)
}
