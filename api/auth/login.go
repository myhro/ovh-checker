package auth

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
)

const incorrectEmailPasswordError = "incorrect email or password"

type loginRequest struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// Login generates a new token and allows a user to log in
func (h *Handler) Login(c *gin.Context) {
	req := loginRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		msg := errors.ValidationMessage(err)
		errors.BadRequestWithMessage(c, msg)
		return
	}

	var id int
	err = h.DB.Get(&id, h.Queries["check-password"], req.Email, req.Password)
	if err != nil && err != sql.ErrNoRows {
		log.Print(err)
		errors.InternalServerError(c)
		return
	} else if err == sql.ErrNoRows {
		errors.UnauthorizedWithMessage(c, incorrectEmailPasswordError)
		return
	}

	token, err := h.newToken(c, id)
	if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return

	}

	body := gin.H{
		"token": token,
	}

	c.JSON(http.StatusOK, body)
}
