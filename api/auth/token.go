package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
	"github.com/myhro/ovh-checker/api/token"
)

func getToken(c *gin.Context) *token.Token {
	tk, _ := c.Get("token")
	return tk.(*token.Token)
}

func (h *Handler) newAuthToken(c *gin.Context, id int) (*token.Token, error) {
	return h.newToken(c, token.Auth, id)
}

func (h *Handler) newSessionToken(c *gin.Context, id int) (*token.Token, error) {
	tk, err := h.newToken(c, token.Session, id)
	if err != nil {
		return nil, err
	}

	err = tk.SetExpiration()
	if err != nil {
		return nil, err
	}

	return tk, nil
}

func (h *Handler) newToken(c *gin.Context, tt token.Type, id int) (*token.Token, error) {
	var tk *token.Token
	switch tt {
	case token.Auth:
		tk = h.TokenStorage.NewAuthToken(id)
	case token.Session:
		tk = h.TokenStorage.NewSessionToken(id)
	}

	tk.Client = c.GetHeader("User-Agent")
	tk.IP = c.ClientIP()
	err := tk.Save()

	if err != nil {
		return nil, err
	}

	return tk, nil
}

// Token creates a new Auth token
func (h *Handler) Token(c *gin.Context) {
	id := c.GetInt("auth_id")
	tk, err := h.newAuthToken(c, id)
	if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return
	}

	body := gin.H{
		"token": tk.ID,
	}

	c.JSON(http.StatusOK, body)
}
