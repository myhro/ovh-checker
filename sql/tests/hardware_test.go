package sqlsuite

import (
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/myhro/ovh-checker/models/hardware"
	"github.com/nleof/goyesql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HardwareTestSuite struct {
	suite.Suite

	db      *sqlx.DB
	mig     *migrate.Migrate
	queries goyesql.Queries
}

func TestHardwareTestSuite(t *testing.T) {
	suite.Run(t, new(HardwareTestSuite))
}

func (s *HardwareTestSuite) SetupSuite() {
	s.db = newDB()
	s.mig = newMigrate()
	s.queries = newQueries("hardware")

	s.mig.Up()
}

func (s *HardwareTestSuite) TearDownSuite() {
	s.mig.Down()
}

func (s *HardwareTestSuite) TestSingleLatestOffer() {
	loadOffers("ks-1-unavailable.json")
	loadOffers("ks-1-eu.json")

	offers := []hardware.LatestOffers{}
	err := s.db.Select(&offers, s.queries["latest-offers"], "fr", 1, 1)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(offers))

	if len(offers) < 1 {
		return
	}

	assert.Equal(s.T(), "KS-1", offers[0].Server)
	assert.Equal(s.T(), "France", offers[0].Country)
	assert.NotNil(s.T(), offers[0].UpdatedAt)
}
