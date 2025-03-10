package api

import (
	"currency-parser-mig/pkg/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLatestExchangeRates(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		const query = `
			SELECT c.id, c.currency_code, c.buy_rate, c.sell_rate, c.timestamp
			FROM currencies c
			WHERE c.timestamp = (
				SELECT MAX(timestamp) 
				FROM currencies 
				WHERE currency_code = c.currency_code
			)
		`

		rates, err := fetchCurrencyRates(db, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error on get data"})
			return
		}

		c.JSON(http.StatusOK, rates)
	}
}

func fetchCurrencyRates(db *sql.DB, query string) ([]models.Currency, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates []models.Currency
	for rows.Next() {
		var rate models.Currency
		if err := rows.Scan(&rate.ID, &rate.CurrencyCode, &rate.BuyRate, &rate.SellRate, &rate.Timestamp); err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}

	return rates, nil
}
