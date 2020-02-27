package storage

import (
	"github.com/lib/pq"
)

// ErrUniqueViolation checks whether the Postgres error is a unique constraint violation
func ErrUniqueViolation(err error) bool {
	pgerr, ok := err.(*pq.Error)
	if ok && pgerr.Code.Name() == "unique_violation" {
		return true
	}
	return false
}
