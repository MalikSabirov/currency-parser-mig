package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// connect to server
	psqlServerInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword)

	serverDb, err := sql.Open("postgres", psqlServerInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL server: %v", err)
	}
	defer serverDb.Close()

	var exists bool
	err = serverDb.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if database exists: %v", err)
	}

	if !exists {
		log.Printf("Database '%s' does not exist, creating it now...", dbName)
		_, err = serverDb.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %v", err)
		}
		log.Printf("Database '%s' created successfully", dbName)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS currencies (
			id SERIAL PRIMARY KEY,
			currency_code VARCHAR(10) NOT NULL,
			buy_rate DECIMAL(10,4) NOT NULL,
			sell_rate DECIMAL(10,4) NOT NULL,
			timestamp TIMESTAMP NOT NULL
		)
	`); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	if _, err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_currency_timestamp ON currencies (currency_code, timestamp)
	`); err != nil {
		return nil, fmt.Errorf("failed to create index: %v", err)
	}

	log.Println("PostgreSQL database connection established")
	return db, nil
}
