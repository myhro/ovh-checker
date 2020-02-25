package auth

import (
	"strings"
)

const headerPrefix = "Token "

func parseAuthHeader(header string) string {
	list := strings.Split(header, headerPrefix)
	id := strings.TrimSpace(list[1])
	return id
}

func validAuthHeader(header string) bool {
	if header == "" {
		return false
	}
	if !strings.HasPrefix(header, headerPrefix) {
		return false
	}
	return true
}
