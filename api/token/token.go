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

// SessionSetExpirationKey is the key to the set which holds information about session expiration
const SessionSetExpirationKey = "tokenset-session-expiration"

// Token holds all the token information
type Token struct {
	Cache  storage.Cache `json:"-"`
	Key    string        `json:"-"`
	SetKey string        `json:"-"`
	Type   Type          `json:"-"`
	UserID int           `db:"user_id" json:"-"`

	ID         string    `db:"id" json:"id"`
	Client     string    `db:"client" json:"client"`
	IP         string    `db:"ip" json:"ip"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	ExpiresAt  time.Time `db:"expires_at" json:"expires_at"`
	LastUsedAt time.Time `db:"last_used_at" json:"last_used_at"`
}

// AuthTokenKey returns the key used to store the Auth token
func AuthTokenKey(tokenID string) string {
	return key(Auth, tokenID)
}

// AuthTokenSetKey returns the key used to store Auth token in its set
func AuthTokenSetKey(userID int) string {
	return setKey(Auth, userID)
}

func key(tt Type, tokenID string) string {
	return fmt.Sprintf("token-%v:%v", prefix(tt), tokenID)
}

func prefix(tt Type) string {
	var prefix string
	switch tt {
	case Auth:
		prefix = AuthPrefix
	case Session:
		prefix = SessionPrefix
	}
	return prefix
}

func setKey(tt Type, userID int) string {
	return fmt.Sprintf("user:%v:tokenset-%v", userID, prefix(tt))
}

// SessionTokenKey returns the key used to store the Session token
func SessionTokenKey(tokenID string) string {
	return key(Session, tokenID)
}

// SessionTokenSetKey returns the key used to store Session token in its set
func SessionTokenSetKey(userID int) string {
	return setKey(Session, userID)
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

func (t *Token) dbField(f string) string {
	field, ok := reflect.TypeOf(t).Elem().FieldByName(f)
	if ok {
		v := field.Tag.Get("db")
		if v != "-" {
			return v
		}
	}
	return ""
}

func (t *Token) keys() {
	t.Key = key(t.Type, t.ID)
	t.SetKey = setKey(t.Type, t.UserID)
}

// Save adds a token to storage
func (t *Token) Save() error {
	details := map[string]interface{}{
		t.dbField("ID"):         t.ID,
		t.dbField("UserID"):     t.UserID,
		t.dbField("Client"):     t.Client,
		t.dbField("IP"):         t.IP,
		t.dbField("CreatedAt"):  storage.TimeFormat(t.CreatedAt),
		t.dbField("LastUsedAt"): storage.TimeFormat(t.LastUsedAt),
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

// SetExpiration sets the token expiration date
func (t *Token) SetExpiration() error {
	expires := t.CreatedAt.Add(storage.CookieMaxAge * time.Second)
	pair := t.Cache.Z(float64(expires.Unix()), t.ID)

	tx := t.Cache.TxPipeline()
	tx.HSet(t.Key, t.dbField("ExpiresAt"), storage.TimeFormat(expires))
	tx.ZAdd(SessionSetExpirationKey, pair)
	_, err := tx.Exec()
	if err != nil {
		return err
	}

	t.ExpiresAt = expires

	return nil
}

// UpdateLastUsed updates the LastUsedAt token information
func (t *Token) UpdateLastUsed() error {
	now := storage.Now()
	_, err := t.Cache.HSet(t.Key, t.dbField("LastUsedAt"), storage.TimeFormat(now))
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
