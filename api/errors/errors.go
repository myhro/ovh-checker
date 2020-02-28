package errors

import (
	"errors"
	"net/http"

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
