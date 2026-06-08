package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"sales-module/internal/config"
)

// DB is the global database connection pool
var DB *sql.DB

// Connect initializes the database connection pool
func Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBSSLMode,
	)
	if cfg.DBPassword != "" {
		connStr += fmt.Sprintf(" password=%s", cfg.DBPassword)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	log.Println("Database connected successfully")
	return db, nil
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
