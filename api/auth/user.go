package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
)

// User returns information about the current user
func (h *Handler) User(c *gin.Context) {
	id := c.GetInt("auth_id")

	key := tokenSetKey(id)
	count, err := h.Cache.SCard(key).Result()
	if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return
	}

	body := gin.H{
		"email":  c.GetString("email"),
		"tokens": count,
	}

	c.JSON(http.StatusOK, body)
}
