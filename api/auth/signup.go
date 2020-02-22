package auth

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
)

const (
	emailError    = "email is invalid or already taken"
	passwordError = "passwords doesn't match"
	userCreated   = "user was successfully created"
)

type signupRequest struct {
	Email     string `form:"email" binding:"required"`
	Password1 string `form:"password1" binding:"required"`
	Password2 string `form:"password2" binding:"required"`
}

// Signup allows a user to sign up
func (h *Handler) Signup(c *gin.Context) {
	req := signupRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		msg := errors.ValidationMessage(err)
		errors.BadRequestWithMessage(c, msg)
		return
	}

	if req.Password1 != req.Password2 {
		errors.BadRequestWithMessage(c, passwordError)
		return
	}

	err = h.DB.Get(new(bool), h.Queries["user-exists"], req.Email)
	if err == nil {
		errors.BadRequestWithMessage(c, emailError)
		return
	} else if err != sql.ErrNoRows {
		log.Print(err)
		errors.InternalServerError(c)
		return
	}

	_, err = h.DB.Exec(h.Queries["add-user"], req.Email, req.Password1)
	if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return
	}

	resp := gin.H{
		"email":   req.Email,
		"message": userCreated,
	}

	c.JSON(http.StatusOK, resp)
}
