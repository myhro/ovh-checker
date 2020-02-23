package auth

import (
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func (h *Handler) addToken(id int, token, client, ip string) error {
	details := map[string]interface{}{
		"id":         token,
		"client":     client,
		"ip":         ip,
		"created_at": now(),
	}

	tx := h.Cache.TxPipeline()
	key := tokenSetKey(id)
	tx.SAdd(key, token)
	key = tokenKey(id, token)
	tx.HMSet(key, details)

	_, err := tx.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) getTokens(c *gin.Context) ([]map[string]string, error) {
	id := c.GetInt("auth_id")
	key := tokenSetKey(id)

	members, err := h.Cache.SMembers(key)
	if err != nil {
		return nil, err
	}

	var list []map[string]string
	for _, token := range members {
		key := tokenKey(id, token)
		details, err := h.Cache.HGetAll(key)
		if err != nil {
			return nil, err
		}
		list = append(list, details)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i]["created_at"] > list[j]["created_at"]
	})

	return list, nil
}
