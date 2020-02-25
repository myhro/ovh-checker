package token

import (
	"errors"
	"fmt"

	"github.com/myhro/ovh-checker/storage"
	"github.com/satori/go.uuid"
)

// ErrNoToken is returned when a token isn't found
var ErrNoToken = errors.New("non-existent token")

// Storage holds the underlying token storage
type Storage struct {
	Cache storage.Cache
}

// List returns the list which the token is part of
func (s *Storage) List(token *Token) ([]Token, error) {
	set, err := token.Set()
	if err != nil {
		return nil, err
	}

	list := make([]Token, 0)
	for _, tokenID := range set {
		t, err := s.load(token.Type, token.UserID, tokenID)
		if err != nil {
			return nil, err
		}
		list = append(list, *t)
	}

	return list, nil
}

// ListAll returns all token lists
func (s *Storage) ListAll(userID int) (map[string][]Token, error) {
	authToken := s.NewAuthToken(userID)
	sessionToken := s.NewSessionToken(userID)

	authList, err := s.List(authToken)
	if err != nil {
		return nil, err
	}
	sessionList, err := s.List(sessionToken)
	if err != nil {
		return nil, err
	}

	result := map[string][]Token{}

	result[AuthPrefix] = authList
	result[SessionPrefix] = sessionList

	return result, nil
}

// Load loads a token from storage
func (s *Storage) load(tt Type, userID int, tokenID string) (*Token, error) {
	t := &Token{
		Cache:  s.Cache,
		ID:     tokenID,
		Type:   tt,
		UserID: userID,
	}
	t.keys()

	details, err := s.Cache.HGetAll(t.Key)
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

// LoadAuthToken loads an existing Auth token from storage
func (s *Storage) LoadAuthToken(userID int, tokenID string) (*Token, error) {
	return s.load(Auth, userID, tokenID)
}

// LoadSessionToken loads an existing Session token from storage
func (s *Storage) LoadSessionToken(userID int, sessionID string) (*Token, error) {
	return s.load(Session, userID, sessionID)
}

// NewAuthToken creates a new Auth token
func (s *Storage) NewAuthToken(userID int) *Token {
	tokenID := uuid.NewV4().String()
	return s.newToken(Auth, userID, tokenID)
}

// NewSessionToken creates a new Session token
func (s *Storage) NewSessionToken(userID int) *Token {
	tokenID := fmt.Sprintf("%x", uuid.NewV4().Bytes())
	return s.newToken(Session, userID, tokenID)
}

func (s *Storage) newToken(tt Type, userID int, tokenID string) *Token {
	token := &Token{
		Cache:     s.Cache,
		CreatedAt: storage.Now(),
		ID:        tokenID,
		Type:      tt,
		UserID:    userID,
	}
	token.keys()

	return token
}
