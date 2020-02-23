package storage

import (
	"database/sql"
	"os"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	// Load Postgres driver
	_ "github.com/lib/pq"
)

// DB is an interface that allows database handlers to be mocked
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

// NewCache creates a Redis client
func NewCache() (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	opts := &redis.Options{
		Addr: addr,
	}
	cache := redis.NewClient(opts)
	_, err := cache.Ping().Result()
	if err != nil {
		return nil, err
	}

	return cache, nil
}

// NewDB creates a Postgres database handler
func NewDB() (*sqlx.DB, error) {
	conn := os.Getenv("POSTGRES_CONN")
	if conn == "" {
		conn = "dbname=ovh sslmode=disable"
	}

	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
