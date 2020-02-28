package errors

import (
	"strings"
)

func jsonError(msg string) bool {
	if strings.HasPrefix(msg, "json: ") {
		return true
	} else if strings.HasPrefix(msg, "invalid character") {
		return true
	} else if msg == "EOF" {
		return true
	}
	return false
}

// ValidationMessage turns a validation error into a more user-friendly message
func ValidationMessage(e error) string {
	errorList := []string{
		"Error:",
		"strconv.ParseInt:",
	}

	firstLine := strings.Split(e.Error(), "\n")[0]
	msg := firstLine
	for _, elem := range errorList {
		if strings.Contains(msg, elem) {
			msg = strings.Split(msg, elem)[1]
			break
		}
	}

	if jsonError(msg) {
		msg = "invalid JSON"
	} else if msg == firstLine {
		msg = "unknown error"
	}
	msg = strings.ToLower(msg)

	return strings.TrimSpace(msg)
}
