package auth

import (
	"time"
)

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}
