package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const ConnectionTimeout = 5 * time.Second

var schemaMu sync.Mutex
var schemaReady bool

func InitDB() error {

	connStr := os.Getenv("POSTGRES_URL")

	if connStr == "" {
		return errors.New("POSTGRES_URL is not configured")
	}
	// pq needs connect_timeout to bound the PostgreSQL handshake as well as dialing.
	timeout := strconv.Itoa(int(ConnectionTimeout / time.Second))
	if strings.HasPrefix(connStr, "postgres://") || strings.HasPrefix(connStr, "postgresql://") {
		parsed, err := url.Parse(connStr)
		if err != nil {
			return errors.New("POSTGRES_URL is invalid")
		}
		query := parsed.Query()
		query.Set("connect_timeout", timeout)
		parsed.RawQuery = query.Encode()
		connStr = parsed.String()
	} else {
		connStr += " connect_timeout=" + timeout
	}
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return fmt.Errorf("DB configuration failed: %w", err)
	}

	// configuring the connections pool manager
	db.SetMaxOpenConns(10)                //Max connections if requests are more other have to wait
	db.SetMaxIdleConns(5)                 //Max Idle connections which are ready to be used
	db.SetConnMaxLifetime(time.Hour * 24) //After 24 hours the connection will be replaced with new one

	// Assigning the db instance to global variable
	DB = db
	schemaReady = false

	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()
	return EnsureReady(ctx)
}

// EnsureReady retries failed startup initialization on subsequent requests.
// Keep the pool even after a failed ping so it can reconnect when PostgreSQL returns.
func EnsureReady(ctx context.Context) error {
	if DB == nil {
		return errors.New("database is not configured")
	}
	if err := DB.PingContext(ctx); err != nil {
		return fmt.Errorf("DB ping failed: %w", err)
	}
	schemaMu.Lock()
	defer schemaMu.Unlock()
	if schemaReady {
		return nil
	}
	if err := createTables(ctx); err != nil {
		return err
	}
	schemaReady = true
	return nil
}

func createTables(ctx context.Context) error {

	// Creating the required tables

	// users table
	usersTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	name TEXT,
	email TEXT UNIQUE NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW()
	);
	`

	_, err := DB.ExecContext(ctx, usersTableQuery)

	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// urls table
	urlsTableQuery := `
	CREATE TABLE IF NOT EXISTS urls (
	id SERIAL PRIMARY KEY,
	original_url TEXT NOT NULL,
	shortcode VARCHAR(10) UNIQUE NOT NULL, 
	clicks INT DEFAULT(0),
	expired_at TIMESTAMPTZ,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	user_id TEXT REFERENCES users(id) ON DELETE CASCADE
	);
	`

	_, err = DB.ExecContext(ctx, urlsTableQuery)

	if err != nil {
		return fmt.Errorf("failed to create urls table: %w", err)
	}

	return nil

}
