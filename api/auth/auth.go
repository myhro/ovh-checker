package auth

import (
	"github.com/myhro/ovh-checker/database"
	"github.com/nleof/goyesql"
)

// Handler holds objects to be reused between requests, like a database connection
type Handler struct {
	DB      database.DB
	Queries goyesql.Queries
}

// NewHandler creates a new Handler
func NewHandler() (*Handler, error) {
	db, err := database.New()
	if err != nil {
		return nil, err
	}

	queries, err := goyesql.ParseFile("sql/auth.sql")
	if err != nil {
		return nil, err
	}

	handler := Handler{
		DB:      db,
		Queries: queries,
	}

	return &handler, nil
}
