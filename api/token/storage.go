package token

import (
	"errors"

	"github.com/myhro/ovh-checker/storage"
)

// ErrNoToken is returned when a token isn't found
var ErrNoToken = errors.New("non-existent token")

// Storage holds the underlying token storage
type Storage struct {
	Cache storage.Cache
}

// List returns the list which the token is part of
func (ts *Storage) List(token *Token) ([]Token, error) {
	token.Storage = ts

	set, err := token.Set()
	if err != nil {
		return nil, err
	}

	list := make([]Token, 0)
	for _, tokenID := range set {
		t, err := ts.Load(token.Type, token.UserID, tokenID)
		if err != nil {
			return nil, err
		}
		list = append(list, *t)
	}

	return list, nil
}

// ListAll returns all token lists
func (ts *Storage) ListAll(userID int) (map[string][]Token, error) {
	authToken := NewAuthToken(userID, ts.Cache)
	sessionToken := NewSessionToken(userID, ts.Cache)

	authList, err := ts.List(authToken)
	if err != nil {
		return nil, err
	}
	sessionList, err := ts.List(sessionToken)
	if err != nil {
		return nil, err
	}

	result := map[string][]Token{}

	authPrefix := prefixes[Auth]
	sessionPrefix := prefixes[Session]

	result[authPrefix] = authList
	result[sessionPrefix] = sessionList

	return result, nil
}

// Load loads a token from storage
func (ts *Storage) Load(tt Type, userID int, tokenID string) (*Token, error) {
	t := &Token{
		ID:     tokenID,
		Type:   tt,
		UserID: userID,
	}
	t.keys()

	details, err := ts.Cache.HGetAll(t.Key)
	if err != nil {
		return nil, err
	} else if len(details) == 0 {
		return nil, ErrNoToken
	}

	t.ID = details[t.field("ID")]
	t.Client = details[t.field("Client")]
	t.IP = details[t.field("IP")]

	t.CreatedAt = storage.ParseTime(details[t.field("CreatedAt")])
	t.LastUsedAt = storage.ParseTime(details[t.field("LastUsedAt")])

	return t, nil
}
