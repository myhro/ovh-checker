package token

import (
	"fmt"
	"reflect"
	"time"

	"github.com/myhro/ovh-checker/storage"
	"github.com/satori/go.uuid"
)

type TokenType int

const (
	Auth TokenType = iota
	Session
)

var prefixes = []string{
	"auth",
	"session",
}

type Token struct {
	Storage *TokenStorage
	Key     string
	SetKey  string
	Type    TokenType
	UserID  int

	ID         string    `cache:"id"`
	Client     string    `cache:"client"`
	IP         string    `cache:"ip"`
	CreatedAt  time.Time `cache:"created_at"`
	LastUsedAt time.Time `cache:"last_used_at"`
}

func NewAuthToken(userID int) *Token {
	id := uuid.NewV4().String()

	token := &Token{
		Type:   Auth,
		UserID: userID,

		ID:        id,
		CreatedAt: storage.Now(),
	}
	token.keys()

	return token
}

func NewSessionToken(userID int) *Token {
	id := fmt.Sprintf("%x", uuid.NewV4().Bytes())

	token := &Token{
		Type:   Session,
		UserID: userID,

		ID:        id,
		CreatedAt: storage.Now(),
	}
	token.keys()

	return token
}

func (t *Token) Count() (int64, error) {
	count, err := t.Storage.Cache.SCard(t.SetKey)
	if err != nil {
		return 0, err
	}
	return count, nil
}

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
		return field.Tag.Get("cache")
	}
	return ""
}

func (t *Token) keys() {
	prefix := prefixes[t.Type]
	t.Key = fmt.Sprintf("user:%v:%v:%v", t.UserID, prefix, t.ID)
	t.SetKey = fmt.Sprintf("user:%v:%v-set", t.UserID, prefix)
}

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

func (t *Token) Set() ([]string, error) {
	members, err := t.Storage.Cache.SMembers(t.SetKey)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (t *Token) UpdateLastUsed() error {
	now := storage.Now()
	_, err := t.Storage.Cache.HSet(t.Key, t.field("LastUsedAt"), storage.TimeFormat(now))
	if err != nil {
		return err
	}
	t.LastUsedAt = now
	return nil
}

func (t *Token) Valid() (bool, error) {
	valid, err := t.Storage.Cache.SIsMember(t.SetKey, t.ID)
	if err != nil {
		return false, err
	}
	return valid, nil
}
