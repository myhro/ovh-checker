package auth

import (
	"github.com/myhro/ovh-checker/api/token"
	"github.com/myhro/ovh-checker/storage"
	"github.com/nleof/goyesql"
)

// Handler holds objects to be reused between requests, like a database connection
type Handler struct {
	DB           storage.DB
	Queries      goyesql.Queries
	TokenStorage *token.Storage
}

// NewHandler creates a new Handler
func NewHandler(cache storage.Cache, db storage.DB) (*Handler, error) {
	queries, err := goyesql.ParseFile("sql/auth.sql")
	if err != nil {
		return nil, err
	}

	ts := &token.Storage{
		Cache: cache,
	}

	handler := Handler{
		DB:           db,
		Queries:      queries,
		TokenStorage: ts,
	}

	return &handler, nil
}
