package sqlsuite

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/myhro/ovh-checker/postgres"
	"github.com/nleof/goyesql"
)

func NewDB() *sqlx.DB {
	os.Setenv("POSTGRES_CONN", "dbname=ovh_test sslmode=disable")
	defer os.Unsetenv("POSTGRES_CONN")

	db, err := postgres.New()
	if err != nil {
		log.Fatal("NewDB: ", err)
	}

	return db
}

func NewMigrate() *migrate.Migrate {
	mig, err := migrate.New("file://../migrations", "postgres:///ovh_test?sslmode=disable")
	if err != nil {
		log.Fatal("NewMigrate: ", err)
	}
	return mig
}

func NewQueries(file string) goyesql.Queries {
	file = fmt.Sprintf("../%v.sql", file)
	queries, err := goyesql.ParseFile(file)
	if err != nil {
		log.Fatal("NewQueries: ", err)
	}
	return queries
}

func ReadFile(file string) string {
	content, err := ioutil.ReadFile(path.Join("testdata/", file))
	if err != nil {
		log.Fatal("ReadFile: ", err)
	}
	return string(content)
}
