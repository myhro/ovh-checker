package auth

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
)

const (
	incorrectEmailTokenError = "incorrect email or token"
	invalidPairError         = "invalid email/token pair"
)

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

func tokenKey(id int, token string) string {
	return fmt.Sprintf("user:%v:token:%v", id, token)

}

func tokenSetKey(id int) string {
	return fmt.Sprintf("user:%v:token-set", id)
}

func (h *Handler) checkTokenAuth(c *gin.Context) {
	email, token, err := parseTokenAuth(c)
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

	key := tokenSetKey(id)
	exists, err := h.Cache.SIsMember(key, token).Result()
	if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return
	} else if !exists {
		errors.UnauthorizedWithMessage(c, incorrectEmailTokenError)
		return
	}

	key = tokenKey(id, token)
	_, err = h.Cache.HSet(key, "last_used", now()).Result()
	if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return
	}

	c.Set("auth_id", id)
	c.Set("email", email)
}