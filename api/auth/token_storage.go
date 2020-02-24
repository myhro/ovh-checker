package auth

import (
	"sort"

	"github.com/myhro/ovh-checker/storage"
)

func (h *Handler) addToken(id int, token, client, ip string) error {
	details := map[string]interface{}{
		"id":         token,
		"client":     client,
		"ip":         ip,
		"created_at": storage.Now(),
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

func (h *Handler) deleteToken(id int, token string) error {
	tx := h.Cache.TxPipeline()
	key := tokenSetKey(id)
	tx.SRem(key, token)
	key = tokenKey(id, token)
	tx.Del(key)

	_, err := tx.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) getTokens(id int) ([]map[string]string, error) {
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
