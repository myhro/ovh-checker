package auth

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
	"github.com/myhro/ovh-checker/api/token"
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
		id := session.Get("auth_id").(int)

		tk, err := token.LoadSessionToken(id, sessionID.(string), h.Cache)
		if err == token.ErrNoToken {
			errors.UnauthorizedWithMessage(c, invalidSessionError)
			return
		} else if err != nil {
			log.Print(err)
			errors.InternalServerError(c)
			return
		}

		err = tk.UpdateLastUsed()
		if err != nil {
			log.Print(err)
			errors.InternalServerError(c)
			return
		}

		c.Set("auth_id", id)
		c.Set("token", tk)
		return
	} else if hasTokenAuth(c) {
		h.checkTokenAuth(c)
		return
	}

	errors.UnauthorizedWithMessage(c, unsupportedAuthError)
}
