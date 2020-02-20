package sqlsuite

import (
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/myhro/ovh-checker/models/offer"
	"github.com/nleof/goyesql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OfferTestSuite struct {
	suite.Suite

	db      *sqlx.DB
	mig     *migrate.Migrate
	queries goyesql.Queries
}

func TestOfferTestSuite(t *testing.T) {
	suite.Run(t, new(OfferTestSuite))
}

func (s *OfferTestSuite) SetupSuite() {
	s.db = newDB()
	s.mig = newMigrate()
	s.queries = newQueries("offer")

	s.mig.Up()
}

func (s *OfferTestSuite) TearDownSuite() {
	s.mig.Down()
}

func (s *OfferTestSuite) TestAvailable() {
	_, err := s.db.Exec(s.queries["import-json"], readFile("ks-1-eu.json"))
	assert.NoError(s.T(), err)

	offers := []offer.Available{}
	err = s.db.Select(&offers, s.queries["available"])
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(offers))

	if len(offers) < 1 {
		return
	}

	o := offers[0]
	assert.Equal(s.T(), "France", o.Country)
	assert.Equal(s.T(), "KS-1", o.Server)
	assert.Equal(s.T(), "1801sk12", o.Code)
}

func (s *OfferTestSuite) TestImportJSON() {
	_, err := s.db.Exec(s.queries["import-json"], readFile("ks-1-eu.json"))
	assert.NoError(s.T(), err)
}
