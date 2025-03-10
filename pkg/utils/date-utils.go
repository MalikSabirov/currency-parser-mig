package utils

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func ReplaceRussianMonth(dateStr string) (string, error) {
	russianMonth, err := extractMonth(dateStr)

	if err != nil {
		return "", fmt.Errorf("failed to replace Russian month: %s", err)
	}

	months := map[string]string{
		"января":   "January",
		"февраля":  "February",
		"марта":    "March",
		"апреля":   "April",
		"мая":      "May",
		"июня":     "June",
		"июля":     "July",
		"августа":  "August",
		"сентября": "September",
		"октября":  "October",
		"ноября":   "November",
		"декабря":  "December",
	}

	if enMonth, exists := months[russianMonth]; exists {
		return strings.Replace(dateStr, russianMonth, enMonth, 1), nil
	}
	return dateStr, nil
}

const ETALON_DATE = "02 January 2006 15:04"

func extractMonth(dateStr string) (string, error) {
	parts := strings.Split(dateStr, " ")

	if len(parts) < 2 {
		return "", fmt.Errorf("incorrect date format")
	}

	return parts[1], nil
}

func ParseOrFallback(dateStr string) time.Time {
	loc, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		log.Printf("failed to load timezone: %v", err)
		return time.Now().UTC() // Фолбэк
	}

	if parsedTime, err := time.ParseInLocation(ETALON_DATE, dateStr, loc); err == nil {
		return parsedTime.UTC()
	} else {
		log.Printf("failed to parse date: %v", err)
	}

	return time.Now().UTC()
}
