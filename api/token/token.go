package token

import (
	"fmt"
	"reflect"
	"time"

	"github.com/myhro/ovh-checker/storage"
	"github.com/satori/go.uuid"
)

// Type is the token type
type Type int

const (
	// Auth persistent token
	Auth Type = iota
	// Session expiring token
	Session
)

var prefixes = []string{
	"auth",
	"session",
}

// Token holds all the token information
type Token struct {
	Storage *Storage `json:"-"`
	Key     string   `json:"-"`
	SetKey  string   `json:"-"`
	Type    Type     `json:"-"`
	UserID  int      `json:"-"`

	ID         string    `json:"id"`
	Client     string    `json:"client"`
	IP         string    `json:"ip"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}

// LoadAuthToken loads an existing Auth token from storage
func LoadAuthToken(userID int, tokenID string, cache storage.Cache) (*Token, error) {
	return loadToken(Auth, userID, tokenID, cache)
}

// LoadSessionToken loads an existing Session token from storage
func LoadSessionToken(userID int, sessionID string, cache storage.Cache) (*Token, error) {
	return loadToken(Session, userID, sessionID, cache)
}

func loadToken(tt Type, userID int, tokenID string, cache storage.Cache) (*Token, error) {
	ts := &Storage{
		Cache: cache,
	}

	token, err := ts.Load(tt, userID, tokenID)
	if err != nil {
		return nil, err
	}
	token.Storage = ts

	return token, nil
}

// NewAuthToken creates a new Auth token
func NewAuthToken(userID int, cache storage.Cache) *Token {
	tokenID := uuid.NewV4().String()
	return newToken(Auth, userID, tokenID, cache)
}

// NewSessionToken creates a new Session token
func NewSessionToken(userID int, cache storage.Cache) *Token {
	tokenID := fmt.Sprintf("%x", uuid.NewV4().Bytes())
	return newToken(Session, userID, tokenID, cache)
}

func newToken(tt Type, userID int, tokenID string, cache storage.Cache) *Token {
	ts := &Storage{
		Cache: cache,
	}

	token := &Token{
		Storage: ts,
		Type:    tt,
		UserID:  userID,

		ID:        tokenID,
		CreatedAt: storage.Now(),
	}
	token.keys()

	return token
}

// Count returns how many tokens are part of its set
func (t *Token) Count() (int64, error) {
	count, err := t.Storage.Cache.SCard(t.SetKey)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Delete removes a token from storage
func (t *Token) Delete() error {
	tx := t.Storage.Cache.TxPipeline()
	tx.SRem(t.SetKey, t.ID)
	tx.Del(t.Key)
	_, err := tx.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (t *Token) field(f string) string {
	field, ok := reflect.TypeOf(t).Elem().FieldByName(f)
	if ok {
		v := field.Tag.Get("json")
		if v != "-" {
			return v
		}
	}
	return ""
}

func (t *Token) keys() {
	prefix := prefixes[t.Type]
	t.Key = fmt.Sprintf("user:%v:%v:%v", t.UserID, prefix, t.ID)
	t.SetKey = fmt.Sprintf("user:%v:%v-set", t.UserID, prefix)
}

// Save adds a token to storage
func (t *Token) Save() error {
	details := map[string]interface{}{
		t.field("ID"):         t.ID,
		t.field("Client"):     t.Client,
		t.field("IP"):         t.IP,
		t.field("CreatedAt"):  storage.TimeFormat(t.CreatedAt),
		t.field("LastUsedAt"): storage.TimeFormat(t.LastUsedAt),
	}

	tx := t.Storage.Cache.TxPipeline()
	tx.SAdd(t.SetKey, t.ID)
	tx.HMSet(t.Key, details)
	_, err := tx.Exec()
	if err != nil {
		return err
	}

	return nil
}

// Set returns the tokens which are part of its set
func (t *Token) Set() ([]string, error) {
	members, err := t.Storage.Cache.SMembers(t.SetKey)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// UpdateLastUsed updates the LastUsedAt token information
func (t *Token) UpdateLastUsed() error {
	now := storage.Now()
	_, err := t.Storage.Cache.HSet(t.Key, t.field("LastUsedAt"), storage.TimeFormat(now))
	if err != nil {
		return err
	}
	t.LastUsedAt = now
	return nil
}

// Valid returns whether a token is valid or not
func (t *Token) Valid() (bool, error) {
	valid, err := t.Storage.Cache.SIsMember(t.SetKey, t.ID)
	if err != nil {
		return false, err
	}
	return valid, nil
}
