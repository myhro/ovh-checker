package auth

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
)

const (
	invalidSessionError  = "invalid or expired session"
	unsupportedAuthError = "no supported authentication method supplied"
)

// AuthRequired ensures that only authenticated requests will reach next handlers
func (h *Handler) AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	sessionID := session.Get("session_id")
	if sessionID != nil {
		token := sessionID.(string)
		id := session.Get("auth_id").(int)

		valid, err := h.validToken(sessionStoragePrefix, id, token)
		if err != nil {
			log.Print(err)
			errors.InternalServerError(c)
			return
		} else if !valid {
			errors.UnauthorizedWithMessage(c, invalidSessionError)
			return
		}

		err = h.updateTokenLastUsed(sessionStoragePrefix, id, token)
		if err != nil {
			log.Print(err)
			errors.InternalServerError(c)
			return
		}

		c.Set("auth_id", id)

		return
	} else if hasTokenAuth(c) {
		h.checkTokenAuth(c)
		return
	}

	errors.UnauthorizedWithMessage(c, unsupportedAuthError)
}
