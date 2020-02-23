package errors

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// BadRequestWithMessage returns a Bad Request response with a custom error message
func BadRequestWithMessage(c *gin.Context, msg string) {
	body := gin.H{
		"error": msg,
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, body)
}

// New returns a standard error
func New(e string) error {
	return errors.New(e)
}

// InternalServerError returns a regular Internal Server Error response
func InternalServerError(c *gin.Context) {
	c.String(http.StatusInternalServerError, "Internal Server Error")
	c.Abort()
}

// UnauthorizedWithMessage returns an Unauthorized response with a custom error message
func UnauthorizedWithMessage(c *gin.Context, msg string) {
	body := gin.H{
		"error": msg,
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, body)
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
	if msg == "EOF" || strings.HasPrefix(msg, "invalid character") {
		msg = "invalid JSON"
	} else if msg == firstLine {
		msg = "unknown error"
	}
	msg = strings.ToLower(msg)

	return strings.TrimSpace(msg)
}
