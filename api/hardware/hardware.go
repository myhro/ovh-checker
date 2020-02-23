package hardware

import (
	"github.com/myhro/ovh-checker/storage"
	"github.com/nleof/goyesql"
)

// Handler holds objects to be reused between requests, like a database connection
type Handler struct {
	DB      storage.DB
	Queries goyesql.Queries
}

// NewHandler creates a new Handler
func NewHandler(db storage.DB) (*Handler, error) {
	queries, err := goyesql.ParseFile("sql/hardware.sql")
	if err != nil {
		return nil, err
	}

	handler := Handler{
		DB:      db,
		Queries: queries,
	}

	return &handler, nil
}
