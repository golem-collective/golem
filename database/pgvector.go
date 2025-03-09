package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB() {
	var err error
	// Use environment variables for sensitive information
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)   // Maximum number of open connections to the database
	db.SetMaxIdleConns(25)   // Maximum number of idle connections in the pool
	db.SetConnMaxLifetime(0) // Connection lifetime (0 means no limit)

	// Check if the database is reachable
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Database connection established successfully")
}

// CloseDB closes the database connection
func CloseDB() {
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	} else {
		log.Println("Database connection closed")
	}
}

// CreateAgentsTable creates the agents table if it does not exist
func CreateAgentsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS agents (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL,
		context TEXT
	);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating agents table: %w", err)
	}
	log.Println("Agents table created or already exists")
	return nil
}

// CreateChatHistoryTable creates the chat_history table if it does not exist
func CreateChatHistoryTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS chat_history (
		id SERIAL PRIMARY KEY,
		agent_id INTEGER NOT NULL,
		role VARCHAR(50) NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (agent_id) REFERENCES agents(id)
	);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating chat_history table: %w", err)
	}
	log.Println("Chat history table created or already exists")
	return nil
}

// Exec executes a query without returning any rows
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}
