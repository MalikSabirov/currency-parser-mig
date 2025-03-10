package api

import (
	"currency-parser-mig/pkg/models"
	"database/sql"
	"fmt"
	"net/http"
	"time"

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

func GetAverageExchangeRates(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		startDateStr, endDateStr := c.Query("start_date"), c.Query("end_date")

		// validate start_date and end_date
		if startDateStr == "" || endDateStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Both start_date and end_date must be provided"})
			return
		}

		// parse dates
		startDate, err := parseDate(startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid start_date format: %s. Expected format: DD-MM-YYYY", startDateStr)})
			return
		}

		endDate, err := parseDate(endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid end_date format: %s. Expected format: DD-MM-YYYY", endDateStr)})
			return
		}

		// T+1
		endDate = endDate.AddDate(0, 0, 1)

		rates, err := fetchAverageRates(db, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve currency averages"})
			return
		}

		// send to client
		c.JSON(http.StatusOK, gin.H{
			"period": gin.H{
				"from": startDateStr,
				"to":   endDateStr,
			},
			"data": rates,
		})
	}
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("02.01.2006", dateStr)
}

func fetchAverageRates(db *sql.DB, start, end time.Time) ([]models.AverageCurrency, error) {
	query := `
		SELECT 
			currency_code, 
			AVG(buy_rate) AS avg_buy, 
			AVG(sell_rate) AS avg_sell
		FROM currencies
		WHERE timestamp BETWEEN $1 AND $2
		GROUP BY currency_code
	`

	rows, err := db.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.AverageCurrency
	for rows.Next() {
		var record models.AverageCurrency
		if err := rows.Scan(&record.CurrencyCode, &record.AverageBuy, &record.AverageSell); err != nil {
			return nil, err
		}
		results = append(results, record)
	}

	return results, nil
}
