package storage

import (
	"time"
)

// Now returns the current UTC time in RFC 3339 format
func Now() string {
	return time.Now().UTC().Format(time.RFC3339)
}
