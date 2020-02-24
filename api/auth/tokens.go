package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
)

// Tokens returns the list of tokens for the current user
func (h *Handler) Tokens(c *gin.Context) {
	id := c.GetInt("auth_id")

	tokens, err := h.getTokens(id)
	if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, tokens)
}
