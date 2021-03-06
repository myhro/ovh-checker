package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
)

const logoutMessage = "successfully logged out"

// Logout logs out a user removing the token used for that
func (h *Handler) Logout(c *gin.Context) {
	token := getToken(c)
	err := token.Delete()
	if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return
	}

	body := gin.H{
		"message": logoutMessage,
	}

	c.JSON(http.StatusOK, body)
}
