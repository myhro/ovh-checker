package auth

import (
	"github.com/go-redis/redis"
	"github.com/myhro/ovh-checker/storage"
	"github.com/nleof/goyesql"
)

// Handler holds objects to be reused between requests, like a database connection
type Handler struct {
	Cache   *redis.Client
	DB      storage.DB
	Queries goyesql.Queries
}

// NewHandler creates a new Handler
func NewHandler(cache *redis.Client, db storage.DB) (*Handler, error) {
	queries, err := goyesql.ParseFile("sql/auth.sql")
	if err != nil {
		return nil, err
	}

	handler := Handler{
		Cache:   cache,
		DB:      db,
		Queries: queries,
	}

	return &handler, nil
}
