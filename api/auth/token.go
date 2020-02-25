package auth

import (
	"database/sql"
	"encoding/base64"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
	"github.com/myhro/ovh-checker/api/token"
)

const (
	incorrectEmailTokenError = "incorrect email or token"
	invalidPairError         = "invalid email/token pair"
)

func getToken(c *gin.Context) *token.Token {
	tk, _ := c.Get("token")
	return tk.(*token.Token)
}

func hasTokenAuth(c *gin.Context) bool {
	header := c.GetHeader("Authorization")
	if header == "" {
		return false
	}
	if !strings.HasPrefix(header, "Token ") {
		return false
	}
	return true
}

func parseTokenAuth(c *gin.Context) (string, string, error) {
	var email, token string

	header := c.GetHeader("Authorization")
	pair := strings.Split(header, "Token ")[1]
	data, err := base64.StdEncoding.DecodeString(pair)
	if err != nil {
		return "", "", err
	}

	pair = string(data)
	if !strings.Contains(pair, ":") {
		return "", "", errors.New(invalidPairError)
	}

	list := strings.Split(pair, ":")
	email = strings.TrimSpace(list[0])
	token = strings.TrimSpace(list[1])

	return email, token, nil
}

func (h *Handler) checkTokenAuth(c *gin.Context) {
	email, tokenID, err := parseTokenAuth(c)
	if err != nil {
		errors.UnauthorizedWithMessage(c, invalidPairError)
		return
	}

	var id int
	err = h.DB.Get(&id, h.Queries["user-exists"], email)
	if err != nil && err != sql.ErrNoRows {
		log.Print(err)
		errors.InternalServerError(c)
		return
	} else if err == sql.ErrNoRows {
		errors.UnauthorizedWithMessage(c, incorrectEmailTokenError)
		return
	}

	tk, err := h.TokenStorage.LoadAuthToken(id, tokenID)
	if err == token.ErrNoToken {
		errors.UnauthorizedWithMessage(c, incorrectEmailTokenError)
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
}
