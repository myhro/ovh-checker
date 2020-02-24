package storage

import (
	"time"
)

// Now returns the current UTC time rounded to second
func Now() time.Time {
	return time.Now().UTC().Round(time.Second)
}

// NowString returns the current UTC time formatted with TimeFormat()
func NowString() string {
	return TimeFormat(Now())
}

// ParseTime parses a previously stored time
func ParseTime(t string) time.Time {
	res, _ := time.Parse(time.RFC3339, t)
	return res
}

// TimeFormat formats a time.Time object in RFC 3339 format
func TimeFormat(t time.Time) string {
	return t.Format(time.RFC3339)
}
