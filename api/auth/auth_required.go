package auth

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
	"github.com/myhro/ovh-checker/api/token"
)

const (
	invalidAuthError     = "invalid authorization token"
	invalidSessionError  = "invalid or expired session"
	unsupportedAuthError = "no supported authentication method supplied"
)

// AuthRequired ensures that only authenticated requests will reach next handlers
func (h *Handler) AuthRequired(c *gin.Context) {
	var err error
	var tk *token.Token
	header := c.GetHeader("Authorization")
	session := sessions.Default(c)
	sessionID := session.Get("session_id")

	if sessionID != nil {
		tk, err = h.TokenStorage.LoadSessionToken(sessionID.(string))
		if err == token.ErrNoToken {
			errors.UnauthorizedWithMessage(c, invalidSessionError)
			return
		}
	} else if validAuthHeader(header) {
		tokenID := parseAuthHeader(header)
		tk, err = h.TokenStorage.LoadAuthToken(tokenID)
		if err == token.ErrNoToken {
			errors.UnauthorizedWithMessage(c, invalidAuthError)
			return
		}
	}

	if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return
	} else if tk == nil {
		errors.UnauthorizedWithMessage(c, unsupportedAuthError)
		return
	}

	err = tk.UpdateLastUsed()
	if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return
	}

	c.Set("auth_id", tk.UserID)
	c.Set("token", tk)
}
