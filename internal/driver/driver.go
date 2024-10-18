package driver

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConnection = &DB{}

const (
	maxOpenDbConn    = 25
	maxIdleDbConn    = 25
	maxDbLifetime    = 5 * time.Minute
	maxDbIdleTime    = 5 * time.Minute
	connMaxRetryTime = 5 * time.Second
)

// ConnectSQL creates database pool for Postgres
func ConnectSQL(dbConnectionString string) (*DB, error) {
	// No need to modify the connection string for lib/pq
	newDatabase, err := NewDatabase(dbConnectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Configure connection pool
	newDatabase.SetMaxOpenConns(maxOpenDbConn)
	newDatabase.SetMaxIdleConns(maxIdleDbConn)
	newDatabase.SetConnMaxLifetime(maxDbLifetime)
	newDatabase.SetConnMaxIdleTime(maxDbIdleTime)

	dbConnection.SQL = newDatabase

	// Test the database connection
	ctx, cancel := context.WithTimeout(context.Background(), connMaxRetryTime)
	defer cancel()

	err = newDatabase.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return dbConnection, nil
}

// NewDatabase creates a new database for the application
func NewDatabase(dbConnectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	return db, nil
}
