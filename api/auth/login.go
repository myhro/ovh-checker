package auth

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
)

const (
	incorrectEmailPasswordError = "incorrect email or password"
	successfulLogin             = "successfully logged in"
)

type loginRequest struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// Login generates a new session
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

	client := c.GetHeader("User-Agent")
	ip := c.ClientIP()
	token, err := h.newSessionToken(id, client, ip)
	if err != nil {
		errors.InternalServerError(c)
		return
	}

	session := sessions.Default(c)
	session.Set("auth_id", id)
	session.Set("session_id", token)
	session.Save()

	body := gin.H{
		"message": successfulLogin,
	}

	c.JSON(http.StatusOK, body)
}
