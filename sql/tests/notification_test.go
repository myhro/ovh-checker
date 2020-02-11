package sqlsuite

import (
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/nleof/goyesql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NotificationTestSuite struct {
	suite.Suite

	db      *sqlx.DB
	mig     *migrate.Migrate
	queries goyesql.Queries
}

func TestNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationTestSuite))
}

func (s *NotificationTestSuite) SetupTest() {
	s.db = NewDB()
	s.mig = NewMigrate()
	s.queries = NewQueries("notification")

	s.mig.Up()
}

func (s *NotificationTestSuite) TearDownTest() {
	s.mig.Down()
}

type Notification struct {
	ID       int
	Email    string
	Server   string
	Cores    int
	Threads  int
	Memory   int
	Storage  string
	Country  string
	Hardware string
}

func (s *NotificationTestSuite) TestAddNotification() {
	email := AddRandomUser()
	res, err := s.db.Exec(s.queries["add-notification"], email, "KS-1", "ca", false)
	assert.NoError(s.T(), err)

	rows, err := res.RowsAffected()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), rows)
}

func (s *NotificationTestSuite) TestPendingNotification() {
	email := AddRandomUser()
	_, err := s.db.Exec(s.queries["add-notification"], email, "KS-1", "fr", false)
	assert.NoError(s.T(), err)

	LoadOffers("ks-1-eu.json")

	res := []Notification{}
	err = s.db.Select(&res, s.queries["pending-notifications"])
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(res))
}

func (s *NotificationTestSuite) TestMarkedAsSentNotification() {
	email := AddRandomUser()
	_, err := s.db.Exec(s.queries["add-notification"], email, "KS-1", "fr", false)
	assert.NoError(s.T(), err)

	LoadOffers("ks-1-eu.json")

	res1 := []Notification{}
	err = s.db.Select(&res1, s.queries["pending-notifications"])
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(res1))

	_, err = s.db.Exec(s.queries["mark-as-sent"], time.Now(), 1)
	assert.NoError(s.T(), err)

	res2 := []Notification{}
	err = s.db.Select(&res2, s.queries["pending-notifications"])
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, len(res2))
}

func (s *NotificationTestSuite) TestRecurrentNotification() {
	email := AddRandomUser()
	_, err := s.db.Exec(s.queries["add-notification"], email, "KS-1", "fr", true)
	assert.NoError(s.T(), err)

	LoadOffers("ks-1-unavailable.json")

	res1 := []Notification{}
	err = s.db.Select(&res1, s.queries["pending-notifications"])
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, len(res1))

	hour_ago := time.Now().Add(-1 * time.Hour)
	_, err = s.db.Exec(s.queries["mark-as-sent"], hour_ago, 1)
	assert.NoError(s.T(), err)

	res2 := []Notification{}
	err = s.db.Select(&res2, s.queries["pending-notifications"])
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, len(res2))

	LoadOffers("ks-1-eu.json")

	res3 := []Notification{}
	err = s.db.Select(&res3, s.queries["pending-notifications"])
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(res3))
}
