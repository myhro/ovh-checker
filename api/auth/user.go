package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
)

// User returns information about the current user
func (h *Handler) User(c *gin.Context) {
	tokens, err := h.getTokens(c)
	if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return
	}

	body := gin.H{
		"email":  c.GetString("email"),
		"tokens": tokens,
	}

	c.JSON(http.StatusOK, body)
}
