package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/darrenjon/restaurant-ordering-system/internal/config"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Manager handles database connections and operations
type Manager struct {
	db *sql.DB
}

// NewManager creates a new database manager
func NewManager(cfg *config.DatabaseConfig) (*Manager, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Successfully connected to the database")

	return &Manager{db: db}, nil
}

// Close closes the database connection
func (m *Manager) Close() error {
	return m.db.Close()
}

// GetDB returns the underlying database connection
func (m *Manager) GetDB() *sql.DB {
	return m.db
}

// ExecuteQuery executes a query that returns rows
func (m *Manager) ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return m.db.Query(query, args...)
}

// ExecuteQueryRow executes a query that is expected to return at most one row
func (m *Manager) ExecuteQueryRow(query string, args ...interface{}) *sql.Row {
	return m.db.QueryRow(query, args...)
}

// ExecuteCommand executes a query that doesn't return rows (e.g., INSERT, UPDATE, DELETE)
func (m *Manager) ExecuteCommand(query string, args ...interface{}) (sql.Result, error) {
	return m.db.Exec(query, args...)
}

// BeginTx starts a transaction
func (m *Manager) BeginTx() (*sql.Tx, error) {
	return m.db.Begin()
}
