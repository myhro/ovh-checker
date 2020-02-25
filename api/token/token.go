package token

import (
	"fmt"
	"reflect"
	"time"

	"github.com/myhro/ovh-checker/storage"
)

// Type is the token type
type Type int

const (
	// Auth persistent token
	Auth Type = iota
	// Session expiring token
	Session
)

const (
	// AuthPrefix is the prefix for Auth tokens
	AuthPrefix = "auth"
	// SessionPrefix is the prefix for Session tokens
	SessionPrefix = "session"
)

// Token holds all the token information
type Token struct {
	Cache  storage.Cache `json:"-"`
	Key    string        `json:"-"`
	SetKey string        `json:"-"`
	Type   Type          `json:"-"`
	UserID int           `json:"-"`

	ID         string    `json:"id"`
	Client     string    `json:"client"`
	IP         string    `json:"ip"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}

// Count returns how many tokens are part of its set
func (t *Token) Count() (int64, error) {
	count, err := t.Cache.SCard(t.SetKey)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Delete removes a token from storage
func (t *Token) Delete() error {
	tx := t.Cache.TxPipeline()
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
	var prefix string
	switch t.Type {
	case Auth:
		prefix = AuthPrefix
	case Session:
		prefix = SessionPrefix
	}
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

	tx := t.Cache.TxPipeline()
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
	members, err := t.Cache.SMembers(t.SetKey)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// UpdateLastUsed updates the LastUsedAt token information
func (t *Token) UpdateLastUsed() error {
	now := storage.Now()
	_, err := t.Cache.HSet(t.Key, t.field("LastUsedAt"), storage.TimeFormat(now))
	if err != nil {
		return err
	}
	t.LastUsedAt = now
	return nil
}

// Valid returns whether a token is valid or not
func (t *Token) Valid() (bool, error) {
	valid, err := t.Cache.SIsMember(t.SetKey, t.ID)
	if err != nil {
		return false, err
	}
	return valid, nil
}
