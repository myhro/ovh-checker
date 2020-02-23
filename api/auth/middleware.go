package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
)

const unsupportedAuthError = "no supported authentication method supplied"

// AuthRequired ensures that only authenticated requests will reach next handlers
func (h *Handler) AuthRequired(c *gin.Context) {
	if hasTokenAuth(c) {
		h.checkTokenAuth(c)
		return
	}

	errors.UnauthorizedWithMessage(c, unsupportedAuthError)
}
