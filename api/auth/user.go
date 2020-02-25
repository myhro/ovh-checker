package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
)

// User returns information about the current user
func (h *Handler) User(c *gin.Context) {
	tk := getToken(c)
	count, err := tk.Count()
	if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return
	}

	id := c.GetInt("auth_id")
	var email string
	err = h.DB.Get(&email, h.Queries["user-email"], id)
	if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return
	}

	body := gin.H{
		"email":  email,
		"tokens": count,
	}

	c.JSON(http.StatusOK, body)
}
