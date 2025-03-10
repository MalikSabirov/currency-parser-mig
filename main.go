package main

import (
	"currency-parser-mig/internal/api"
	"currency-parser-mig/internal/database"
	"currency-parser-mig/internal/parser"
	"log"
	"os"

	"github.com/gin-gonic/gin"
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

	// init gin router
	r := gin.Default()

	r.GET("/currencies/latest", api.GetLatestExchangeRates(db))
	r.GET("/currencies/average", api.GetAverageExchangeRates(db))

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	parser.ParseCurrencies(db)

	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
