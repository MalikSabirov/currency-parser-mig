package main

import (
	"currency-parser-mig/internal/database"
	"currency-parser-mig/internal/parser"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// init db
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to init db: %v", err)
	}
	defer db.Close()

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	parser.ParseCurrencies()

	log.Printf("Server running on port %s", port)
}
