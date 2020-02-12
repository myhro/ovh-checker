package postgres

import (
	"os"

	"github.com/jmoiron/sqlx"
	// Load Postgres driver
	_ "github.com/lib/pq"
)

// New creates a Postgres database handler
func New() (*sqlx.DB, error) {
	conn := os.Getenv("POSTGRES_CONN")
	if conn == "" {
		conn = "dbname=ovh sslmode=disable"
	}

	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
