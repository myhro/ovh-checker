package storage

import (
	"errors"
	"os"

	"github.com/gin-contrib/sessions/cookie"
)

// CookieStore is the default cookie store interface
type CookieStore interface {
	cookie.Store
}

// NewCookieStore creates a new cookie store
func NewCookieStore() (cookie.Store, error) {
	secret := os.Getenv("COOKIE_STORE_SECRET")
	nonProdSecret := "non-production-secret"
	if secret == "" {
		secret = nonProdSecret
	}

	store := cookie.NewStore([]byte(secret))
	if secret == nonProdSecret {
		return store, errors.New("COOKIE_STORE_SECRET not found, using non-production secret")
	}

	return store, nil
}
