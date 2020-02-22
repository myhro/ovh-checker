package errors

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// BadRequest returns a regular Bad Request response
func BadRequest(c *gin.Context) {
	c.String(http.StatusBadRequest, "Bad Request")
}

// BadRequestWithMessage returns a Bad Request response with a custom error message
func BadRequestWithMessage(c *gin.Context, msg string) {
	body := gin.H{
		"error": msg,
	}

	c.JSON(http.StatusBadRequest, body)
}

// InternalServerError returns a regular Internal Server Error response
func InternalServerError(c *gin.Context) {
	c.String(http.StatusInternalServerError, "Internal Server Error")
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
