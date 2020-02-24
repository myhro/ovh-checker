package auth

import (
	"fmt"
	"sort"

	"github.com/myhro/ovh-checker/storage"
)

const (
	authStoragePrefix    = "auth"
	sessionStoragePrefix = "session"
)

type tokenDetails map[string]string

func tokenSetKey(keyPrefix string, id int) string {
	return fmt.Sprintf("user:%v:%v-set", id, keyPrefix)
}

func tokenKey(keyPrefix string, id int, token string) string {
	return fmt.Sprintf("user:%v:%v:%v", id, keyPrefix, token)
}

func (h *Handler) addToken(keyPrefix string, id int, token, client, ip string) error {
	details := map[string]interface{}{
		"id":         token,
		"client":     client,
		"ip":         ip,
		"created_at": storage.NowString(),
	}

	tx := h.Cache.TxPipeline()
	key := tokenSetKey(keyPrefix, id)
	tx.SAdd(key, token)
	key = tokenKey(keyPrefix, id, token)
	tx.HMSet(key, details)

	_, err := tx.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) countTokens(keyPrefix string, id int) (int64, error) {
	key := tokenSetKey(keyPrefix, id)
	count, err := h.Cache.SCard(key)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (h *Handler) deleteToken(keyPrefix string, id int, token string) error {
	tx := h.Cache.TxPipeline()
	key := tokenSetKey(keyPrefix, id)
	tx.SRem(key, token)
	key = tokenKey(keyPrefix, id, token)
	tx.Del(key)

	_, err := tx.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) getTokens(id int) (map[string][]tokenDetails, error) {
	result := map[string][]tokenDetails{}
	prefixes := []string{authStoragePrefix, sessionStoragePrefix}
	for _, p := range prefixes {
		key := tokenSetKey(p, id)
		members, err := h.Cache.SMembers(key)
		if err != nil {
			return nil, err
		}

		list := make([]tokenDetails, 0)
		for _, token := range members {
			key := tokenKey(p, id, token)
			details, err := h.Cache.HGetAll(key)
			if err != nil {
				return nil, err
			}
			list = append(list, details)
		}
		sort.Slice(list, func(i, j int) bool {
			return list[i]["created_at"] > list[j]["created_at"]
		})
		result[p] = list
	}

	return result, nil
}

func (h *Handler) updateTokenLastUsed(keyPrefix string, id int, token string) error {
	key := tokenKey(keyPrefix, id, token)
	_, err := h.Cache.HSet(key, "last_used", storage.NowString())
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) validToken(keyPrefix string, id int, token string) (bool, error) {
	key := tokenSetKey(keyPrefix, id)
	valid, err := h.Cache.SIsMember(key, token)
	if err != nil {
		return false, err
	}
	return valid, nil
}
