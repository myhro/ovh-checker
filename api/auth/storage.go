package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func (h *Handler) newToken(c *gin.Context) (string, error) {
	id := c.GetInt("auth_id")
	token := uuid.NewV4().String()
	details := map[string]interface{}{
		"id":         token,
		"client":     c.GetHeader("User-Agent"),
		"ip":         c.ClientIP(),
		"created_at": now(),
	}

	tx := h.Cache.TxPipeline()
	key := tokenSetKey(id)
	tx.SAdd(key, token)
	key = tokenKey(id, token)
	tx.HMSet(key, details)
	_, err := tx.Exec()
	if err != nil {
		return "", err
	}

	return token, nil
}
