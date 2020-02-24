package auth

import (
	"fmt"

	"github.com/satori/go.uuid"
)

func (h *Handler) newSessionToken(id int, client, ip string) (string, error) {
	token := fmt.Sprintf("%x", uuid.NewV4().Bytes())
	err := h.addToken(sessionStoragePrefix, id, token, client, ip)
	if err != nil {
		return "", err
	}
	return token, nil
}
